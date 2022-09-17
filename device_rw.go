package visca

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func (d *Device) Writer() {
	fmt.Printf("[Device.Writer] init\n")
	defer fmt.Printf("[Device.Writer] done\n")

	for {
		select {
		case cmd := <-d.write:
			fmt.Printf("[Device.Writer] Received [% X] %T\n", cmd.bytes(), cmd)

			var payload bytes.Buffer
			switch cmd.(type) {
			case *SeqReset, *Raw:
				payload.Write(cmd.bytes())
			default:
				payload.WriteByte(0x81)
				payload.Write(cmd.bytes())
				payload.WriteByte(0xff)
			}

			var header bytes.Buffer

			// Payload Type:
			switch cmd.cmdType().(type) {
			case ViscaCommand:
				header.Write([]byte{0x1, 0x0})
			case ViscaInquiry:
				header.Write([]byte{0x1, 0x10})
			case ViscaReply:
				header.Write([]byte{0x1, 0x11})
			case DeviceSettingCommand:
				header.Write([]byte{0x1, 0x20})
			case ControlCommand:
				header.Write([]byte{0x2, 0x0})
			case ControlCommandReply:
				header.Write([]byte{0x2, 0x1})
			}

			// Payload Length
			length := make([]byte, 2)
			binary.BigEndian.PutUint16(length, uint16(payload.Len()))
			header.Write(length)

			// sequence number
			seq := make([]byte, 4)
			binary.BigEndian.PutUint32(seq, d.writeSeq)
			d.writeSeq++
			header.Write(seq)
			if _, ok := cmd.(SeqReset); ok {
				d.writeSeq = 0
			}

			var message bytes.Buffer
			if strings.HasPrefix(d.Path, "udp://") {
				message.Write(header.Bytes())
				message.Write(payload.Bytes())
			} else {
				message.Write(payload.Bytes())
			}

			_, err := d.port.Write(message.Bytes())
			if err != nil {
				fmt.Printf("[Device.Writer] ERR write: %v\n", err)
			} else {
				fmt.Printf("[Device.Writer] Wrote [% X]\n", message.Bytes())
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
	fmt.Printf("++++++ [Device.Reader] init\n")
	defer fmt.Printf("++++ [Device.Reader] done\n")

	packet := bytes.Buffer{}

	for {
		if d.Err() != nil {
			return
		}

		// read byte from camera
		readByte := make([]byte, 1)
		switch port := d.port.(type) {
		case *net.UDPConn:
			fmt.Printf(">> UDP\n")

			buff := make([]byte, 1024)
			n, _, err := port.ReadFromUDP(buff)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(">>> %d: [% X]\n", n, buff)
		default:
			fmt.Printf(">> Default\n")
			if _, err := d.port.Read(readByte); err != nil {
				d.Close()
				return
			}
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
