package visca

import (
	"context"
	"encoding/binary"
	"fmt"
	"go.bug.st/serial"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	SonySRGX400 = "Sony SRG-X400"
)

type Device struct {
	Path   string
	Type   string
	Config Config
	State  State

	conn       io.ReadWriter
	remoteAddr *net.UDPAddr

	do          chan Cmd
	read        chan []byte
	write       chan []byte
	writeSeq    uint32
	writeSeqCmd sync.Map

	context.Context
	Close context.CancelFunc

	Booting *sync.WaitGroup

	PanTiltQueue PanTiltQueue
	PanTiltReady sync.WaitGroup
	ZoomQueue    ZoomQueue
	ZoomReady    sync.WaitGroup
}

type Config struct {
	XMaxSpeed float64
	YMaxSpeed float64
	ZMaxSpeed float64

	LocalUDP string
}

type State struct {
	PanTiltDrive        PanTiltDrive
	PanTiltDriveLastSeq *uint32
	Zoom                Zoom
	ZoomLastSeq         *uint32
	ExposureMode        ExposureMode
	Power               Power
}

func (d *Device) Find() (err error) {
	d.Context, d.Close = context.WithCancel(context.Background())
	d.Booting = &sync.WaitGroup{}
	d.Booting.Add(1)

	if strings.HasPrefix(d.Path, "/") {
		mode := serial.Mode{
			BaudRate: 9600,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		}
		if d.conn, err = serial.Open(d.Path, &mode); err != nil {
			return
		}
	}
	if strings.HasPrefix(d.Path, "tcp://") {
		addr := strings.TrimPrefix(d.Path, "tcp://")
		d.conn, err = net.Dial("tcp", addr)
		if err != nil {
			return
		}
	}
	if strings.HasPrefix(d.Path, "udp://") {
		d.remoteAddr, err = net.ResolveUDPAddr("udp4", strings.TrimPrefix(d.Path, "udp://"))
		if err != nil {
			fmt.Printf(">> net.ResolveUDPAddr ERROR %v\n", err)
			return
		}

		ip, _ := net.ResolveUDPAddr("udp", d.Config.LocalUDP)
		d.conn, err = net.ListenUDP("udp4", ip)
		if err != nil {
			fmt.Printf(">> net.ListenUDP ERROR %v\n", err)
			return
		}

		//d.port, err = net.DialUDP("udp", nil, udpAddr)
		//if err != nil {
		//	fmt.Printf(">> net.DialUDP ERROR %v\n", err)
		//	return
		//}
	}
	if strings.HasPrefix(d.Path, "test://") {
		d.conn = &RW{}
	}

	return
}
func (d *Device) Found() bool {
	return d.conn != nil
}

func (d *Device) Run() {
	d.PanTiltQueue = PanTiltQueue{}
	d.PanTiltQueue.Init()
	go d.PanTiltQueueWorker()

	d.ZoomQueue = ZoomQueue{}
	d.ZoomQueue.Init()
	go d.ZoomQueueWorker()

	d.do = make(chan Cmd)
	//go d.DoWorker()
	go d.DoWorker()

	//d.writeSeqCmd = make(map[uint32]Cmd)
	go d.Reader()

	d.read = make(chan []byte)
	go d.readHandler()

	d.write = make(chan []byte)
	go d.Writer()

	d.Booting.Done()
	// sync status

	d.Do(&SeqReset{})
	d.Do(&InqPower{})

	<-d.Done()
}

func (d *Device) Do(cmd Cmd, alreadyApplied ...bool) {
	fmt.Printf("[Do] %s\n", cmd)

	if len(alreadyApplied) == 0 {
		if cmd, ok := cmd.(CmdAppliable); ok {
			okToSend := cmd.Apply(d)
			if !okToSend {
				//fmt.Printf(">> not sending\n")
				return
			}
		}
	}

	d.do <- cmd
}

func (d *Device) DoWorker() {
	for {
		cmd := <-d.do

		fmt.Printf("[Device.DoWordker] Sending %s\n", cmd)

		// retry until ack
		for {
			d.writeSeq++
			if _, ok := cmd.(*SeqReset); ok {
				d.writeSeq = 0
			}

			ack := d.sendAndWaitForAck(cmd, d.writeSeq)
			if ack {
				break
			}

			fmt.Printf(">> retrying %s\n", cmd)
		}
	}
}

func (d *Device) sendAndWaitForAck(cmd Cmd, seq uint32) bool {
	cmd.InitContext()

	// build packet
	var data []byte
	switch cmd := cmd.(type) {
	case ViscaCommand:
		data = append(data, 0x1, 0x0)
		data = append(data, pLen2(cmd.ViscaCommand())...)
		data = append(data, d.pSeq4(seq)...)
		data = append(data, cmd.ViscaCommand()...)
	case ViscaInquiry:
		data = append(data, 0x1, 0x10)
		data = append(data, pLen2(cmd.ViscaInquiry())...)
		data = append(data, d.pSeq4(seq)...)
		data = append(data, cmd.ViscaInquiry()...)
	case ControlCommand:
		data = append(data, 0x2, 0x0)
		data = append(data, pLen2(cmd.ControlCommand())...)
		data = append(data, d.pSeq4(seq)...)
		data = append(data, cmd.ControlCommand()...)
	//case ViscaReply:
	//data = append(data, 0x1, 0x11)
	//case DeviceSettingCommand:
	//data = append(data, 0x1, 0x20)
	//case ControlCommandReply:
	//data = append(data, 0x2, 0x1)
	default:
		fmt.Printf("ERROR unsupported cmd\n")
		return true
	}

	// keep track of all sent Cmds
	d.writeSeqCmd.Store(seq, cmd)
	if _, ok := cmd.(*PanTiltDrive); ok {
		d.State.PanTiltDriveLastSeq = &seq
	}
	if _, ok := cmd.(*Zoom); ok {
		d.State.ZoomLastSeq = &seq
	}

	// send it
	d.write <- data
	//if cmd, ok := cmd.(CmdWaitable); ok {
	//	go cmd.WaitReply(d)
	//}

	// wait for ack
	select {
	case <-time.After(time.Millisecond * 500):
		cmd.Finish()
		return false
	case <-cmd.Done():
		return true
	}
}

func pLen2(payload []byte) []byte {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(payload)))
	return length
}
func (d *Device) pSeq4(seqInt uint32) []byte {
	seq := make([]byte, 4)
	binary.BigEndian.PutUint32(seq, seqInt)

	return seq
}
