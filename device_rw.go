package visca

import (
	"bytes"
	"fmt"
)

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
			switch true {
			// Command ACK
			case len(packet) == 3 && packet[0] == 0x90 && packet[1]>>4 == 0x4 && packet[2] == 0xff:
				//fmt.Printf("[Camera] Read [% X] accepted\n", packet)
			// Command Completion
			case len(packet) == 3 && packet[0] == 0x90 && packet[1]>>4 == 0x5 && packet[2] == 0xff:
				//fmt.Printf("[Camera] Read [% X] finished\n", packet)
			// MenuStatus responses
			case bytes.Equal(packet, []byte{0x90, 0x50, 0x2, 0xFF}):
				//fmt.Println("MENU IS ON")
				d.OSDToggle.MenuIsOn = true
			case bytes.Equal(packet, []byte{0x90, 0x50, 0x3, 0xFF}):
				//fmt.Println("MENU IS OFF")
				d.OSDToggle.MenuIsOn = false
			// unhandled
			default:
				fmt.Printf("[Camera] Read [% X]\n", packet)
			}

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
			//fmt.Printf("[%T] [% X] %+v\n", cmd, cmd.bytes(), cmd)

			packet := bytes.Buffer{}
			packet.WriteByte(0x81) // camera address
			packet.Write(cmd.bytes())
			packet.WriteByte(0xff) // EOF

			_, err := d.port.Write(packet.Bytes())
			if err != nil {
				fmt.Printf("ERR write: %v\n", err)
			} else {
				//fmt.Printf("\n[Camera] Write [% X]\n", packet.Bytes())
			}

		case <-d.Done():
			close(d.write)
			return
		}
	}
}
