package visca

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

func (d *Device) Writer() {
	fmt.Printf("[Device.Writer] init\n")
	defer fmt.Printf("[Device.Writer] done\n")

	for {
		select {
		case cmd := <-d.write:
			fmt.Printf("[Device.Writer] Received [% X] %T\n", cmd.bytes(), cmd)

			packet := bytes.Buffer{}
			packet.WriteByte(0x81) // camera address
			packet.Write(cmd.bytes())
			packet.WriteByte(0xff) // EOF

			_, err := d.port.Write(packet.Bytes())
			if err != nil {
				fmt.Printf("[Device.Writer] ERR write: %v\n", err)
			} else {
				fmt.Printf("[Device.Writer] Wrote [% X]\n", packet.Bytes())
			}

			time.Sleep(time.Millisecond * 50)

		case <-d.Done():
			fmt.Printf("[Device.Writer] device ctx Done")
			close(d.write)
			return
		}
	}
}

func (d *Device) Reader() {
	fmt.Printf("[Device.Reader] init\n")
	defer fmt.Printf("[Device.Reader] done\n")

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

		log.Printf("[Device.Reader] Read [% X]\n", packet.Bytes())

		// is camera done talking?
		if bytes.Equal(readByte, []byte{0xff}) {
			log.Printf("[Device.Reader] Read Done [% X]\n", packet.Bytes())

			d.read <- packet.Bytes()

			// clear packet buffer
			packet.Reset()
		}

	}
}

func (d *Device) readHandler() {
	fmt.Printf("[Device.readHandler] init\n")
	defer fmt.Printf("[Device.readHandler] done\n")

	for {
		select {
		case packet := <-d.read:
			fmt.Printf("[Device.readHandler] Read [% X]\n", packet)

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
				fmt.Printf("[Camera] Read unhandled [% X]\n", packet)
				continue
			}

			fmt.Printf("[Camera] Handled [% X]\n", packet)

		case <-d.Done():
			close(d.read)
			return
		}
	}
}
