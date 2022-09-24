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

		d.conn, err = net.ListenUDP("udp4", &net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 52381, // seems to be hardcoded in camera
		})
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

func (d *Device) Do(cmd Cmd, preApproved ...bool) {
	var data []byte

	// Cmd triggered via Fin is usually pre-approved, no need to check apply/send
	if len(preApproved) == 0 {
		if applied := cmd.Apply(d); !applied {
			return
		}
		fmt.Printf("[Device.Do] Applied %s\n", cmd)

		if cmd, ok := cmd.(CmdSendable); ok {
			if send := cmd.Send(d); !send {
				return
			}
		}
		fmt.Printf("[Device.Do] Sending %s\n", cmd)
	}

	switch cmd := cmd.(type) {
	case ViscaCommand:
		data = append(data, 0x1, 0x0)
		data = append(data, pLen2(cmd.ViscaCommand())...)
		data = append(data, d.pSeq4()...)
		data = append(data, cmd.ViscaCommand()...)
	case ViscaInquiry:
		data = append(data, 0x1, 0x10)
		data = append(data, pLen2(cmd.ViscaInquiry())...)
		data = append(data, d.pSeq4()...)
		data = append(data, cmd.ViscaInquiry()...)
	case ViscaReply:
		//data = append(data, 0x1, 0x11)
	case DeviceSettingCommand:
		//data = append(data, 0x1, 0x20)
	case ControlCommand:
		data = append(data, 0x2, 0x0)
		data = append(data, pLen2(cmd.ControlCommand())...)
		data = append(data, d.pSeq4()...)
		data = append(data, cmd.ControlCommand()...)
	case ControlCommandReply:
		//data = append(data, 0x2, 0x1)
	default:
		fmt.Printf("ERROR unsupported cmd\n")
		return
	}

	// keep track of all sent Cmds
	seq := d.writeSeq - 1
	d.writeSeqCmd.Store(seq, cmd)

	if _, ok := cmd.(*PanTiltDrive); ok {
		d.State.PanTiltDriveLastSeq = &seq
	}
	if _, ok := cmd.(*Zoom); ok {
		d.State.ZoomLastSeq = &seq
	}

	// send it
	d.write <- data
	if cmd, ok := cmd.(CmdWaitable); ok {
		go cmd.WaitReply(d)
	}

	fmt.Printf("[Device.Do] Wrote %s\n", cmd)
}

func pLen2(payload []byte) []byte {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(payload)))
	return length
}
func (d *Device) pSeq4() []byte {
	seq := make([]byte, 4)
	binary.BigEndian.PutUint32(seq, d.writeSeq)
	d.writeSeq++

	return seq
}
