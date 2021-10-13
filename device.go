package visca

import (
	"bytes"
	"context"
	"fmt"
	"go.bug.st/serial"
	"io"
	"log"
	"net"
	"strings"
)

type Device struct {
	Path string

	port   io.ReadWriter
	reader chan []byte
	writer chan Command

	Run  context.Context
	Quit context.CancelFunc
}

func (d *Device) Find() (err error) {
	d.Run, d.Quit = context.WithCancel(context.Background())

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

	d.reader = make(chan []byte)
	d.writer = make(chan Command)

	// TODO test connectivity

	return
}
func (d *Device) Found() bool {
	return d.port != nil
}

func (d *Device) Reader() {
	packet := bytes.Buffer{}

	for {
		// read byte from camera
		readByte := make([]byte, 1)
		if _, err := d.port.Read(readByte); err != nil {
			d.Quit()
			close(d.reader)
			return
		}
		packet.Write(readByte)

		//log.Printf("[Camera] Reading [% X]\n", packet.Bytes())

		// is camera done talking?
		if bytes.Equal(readByte, []byte{0xff}) {
			d.reader <- packet.Bytes()

			// clear packet buffer
			packet.Reset()
		}
	}
}

func (d *Device) readHandler() {
	for {
		select {
		case packet := <-d.reader:
			log.Printf("[Camera] Read [% X]\n", packet)
		case <-d.Run.Done():
			return
		}
	}
}

func (d *Device) Writer() {
	for {
		select {
		case cmd := <-d.writer:
			fmt.Printf("[%T] %+v\n", cmd, cmd)

			packet := []byte{0x81} // camera address
			packet = append(packet, cmd.Bytes()...)
			packet = append(packet, byte(0xff)) // EOF

			n, err := d.port.Write(packet)
			if err != nil {
				fmt.Printf("ERR write: %v\n", err)
			} else {
				log.Printf("[Camera] Write len(%d) [% X]\n", n, packet)
			}

		case <-d.Run.Done():
			close(d.writer)
			return
		}
	}
}

func (d *Device) Do(cmd Command) {
	d.writer <- cmd
}
