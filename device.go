package visca

import (
	"bytes"
	"context"
	"fmt"
	"go.bug.st/serial"
	"io"
	"net"
	"strings"
)

type Device struct {
	Path string

	// one-shot commands
	MoveHome   MoveHome
	Focus      Focus
	CallPreset CallPreset
	SavePreset SavePreset

	OSDToggle OSDToggle
	OSDEnter  OSDEnter
	OSDReturn OSDReturn
	OSDUp     OSDUp
	OSDRight  OSDRight
	OSDDown   OSDDown
	OSDLeft   OSDLeft

	// stateful commands
	Move Move
	Zoom Zoom

	port  io.ReadWriter
	read  chan []byte
	write chan Command

	context.Context
	Close context.CancelFunc
}

func (d *Device) Apply(cmds ...Command) {
	allowed := map[Command]int{
		// one-shot commands
		&d.CallPreset: 0,
		&d.SavePreset: 0,

		&d.MoveHome: 0,
		&d.Focus:    0,

		&d.OSDToggle: 0,
		&d.OSDEnter:  0,
		&d.OSDReturn: 0,
		&d.OSDUp:     0,
		&d.OSDRight:  0,
		&d.OSDDown:   0,
		&d.OSDLeft:   0,

		// stateful commands
		&d.Move: 0,
		&d.Zoom: 0,
	}

	for _, cmd := range cmds {
		if _, ok := allowed[cmd]; ok && cmd.apply() {
			d.write <- cmd
		}
	}
}

func (d *Device) Find() (err error) {
	d.Context, d.Close = context.WithCancel(context.Background())

	if strings.HasPrefix(d.Path, "/") {
		mode := serial.Mode{
			BaudRate: 9600,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		}
		if d.port, err = serial.Open(d.Path, &mode); err != nil {
			return
		}
	}
	if strings.HasPrefix(d.Path, "tcp://") {
		d.port, err = net.Dial("tcp", strings.TrimPrefix(d.Path, "tcp://"))
		if err != nil {
			return
		}
	}
	if strings.HasPrefix(d.Path, "udp://") {
		d.port, err = net.Dial("udp", strings.TrimPrefix(d.Path, "udp://"))
		if err != nil {
			return
		}
	}
	if strings.HasPrefix(d.Path, "test://") {
		d.port = &RW{}
	}

	d.Move.device = d
	d.Zoom.device = d

	// TODO test connectivity

	return
}
func (d *Device) Found() bool {
	return d.port != nil
}

func (d *Device) Run() {
	go d.Reader()

	d.read = make(chan []byte)
	go d.readHandler()

	d.write = make(chan Command)
	go d.Writer()

	<-d.Done()
}

func (d *Device) Reader() {
	packet := bytes.Buffer{}

	for {
		if d.Err() != nil {
			return
		}

		// read byte from camera
		readByte := make([]byte, 1)
		if _, err := d.port.Read(readByte); err != nil {
			d.Close()
			return
		}
		packet.Write(readByte)

		//log.Printf("[Camera] Reading [% X]\n", packet.Bytes())

		// is camera done talking?
		if bytes.Equal(readByte, []byte{0xff}) {
			d.read <- packet.Bytes()

			// clear packet buffer
			packet.Reset()
		}
	}
}

func (d *Device) readHandler() {
	for {
		select {
		case packet := <-d.read:
			fmt.Printf("[Camera] Read [% X]\n", packet)

		case <-d.Done():
			close(d.read)
			return
		}
	}
}

func (d *Device) Writer() {
	for {
		select {
		case cmd := <-d.write:
			fmt.Printf("[%T] [% X] %+v\n", cmd, cmd.bytes(), cmd)

			packet := bytes.Buffer{}
			packet.WriteByte(0x81) // camera address
			packet.Write(cmd.bytes())
			packet.WriteByte(0xff) // EOF

			_, err := d.port.Write(packet.Bytes())
			if err != nil {
				fmt.Printf("ERR write: %v\n", err)
			} else {
				//log.Printf("[Camera] Write len(%d) [% X]\n", n, packet)
			}

		case <-d.Done():
			close(d.write)
			return
		}
	}
}
