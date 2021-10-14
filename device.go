package visca

import (
	"context"
	"go.bug.st/serial"
	"io"
	"net"
	"strings"
	"time"
)

type Device struct {
	Path string

	// one-shot commands
	MoveHome   MoveHome
	Focus      Focus
	CallPreset CallPreset
	SavePreset SavePreset

	OSDToggle OSDToggle
	osdIsOpen bool
	OSDEnter  OSDEnter
	OSDReturn OSDReturn
	OSDUp     OSDUp
	OSDRight  OSDRight
	OSDDown   OSDDown
	OSDLeft   OSDLeft

	// info inquiries
	AskMenuStatus AskMenuStatus

	// stateful commands
	Move Move
	Zoom Zoom

	port  io.ReadWriter
	read  chan []byte
	write chan Command

	context.Context
	Close context.CancelFunc
}

type Async struct {
	cmd Command

	id     int
	sentAt time.Time

	accepted   bool
	acceptedAt time.Time

	finished   bool
	finishedAt time.Time

	latency time.Duration
}

func (d *Device) Apply(cmds ...Command) {
	commands := map[Command]bool{
		// one-shot commands
		&d.CallPreset: true,
		&d.SavePreset: true,

		&d.MoveHome: true,
		&d.Focus:    true,

		&d.OSDToggle: true,
		&d.OSDEnter:  true,
		&d.OSDReturn: true,
		&d.OSDUp:     true,
		&d.OSDRight:  true,
		&d.OSDDown:   true,
		&d.OSDLeft:   true,

		// stateful commands
		&d.Move: true,
		&d.Zoom: true,

		// inquiries
		&d.AskMenuStatus: true,
	}

	for _, cmd := range cmds {
		// make sure applied command is found,
		// is allowed to be fired,
		// and actually really needs firing
		if allowed, found := commands[cmd]; found {
			if !allowed {
				//fmt.Printf("NOT ALLOWED\n")
				continue
			}
			if !cmd.apply() {
				//fmt.Printf("NOT APPLIED\n")
				continue
			}

			d.write <- cmd
			// good to go

		} else {
			//fmt.Printf("NOT FOUND\n")
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

	// sync status
	d.Apply(&d.AskMenuStatus)

	<-d.Done()
}
