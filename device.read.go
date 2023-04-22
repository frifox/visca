package visca

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func (d *Device) Reader() {
	fmt.Printf("[Device.Reader] init\n")
	defer fmt.Printf("[Device.Reader] done\n")

	for {
		if d.Err() != nil {
			return
		}

		data := make([]byte, 1024)
		switch port := d.conn.(type) {
		case *net.UDPConn:
			n, _, err := port.ReadFromUDPAddrPort(data)
			if err != nil {
				log.Fatal(err)
			}
			data = data[:n]
		default:
			n, err := d.conn.Read(data)
			if err != nil {
				d.Close()
				return
			}
			data = data[:n]
		}

		if len(data) == 0 {
			continue
		}

		//fmt.Printf("[Device.Reader] Read [% X]\n", packet.Bytes())
		d.read <- data
	}
}

func (d *Device) readHandler() {
	fmt.Printf("[Device.readHandler] init\n")
	defer fmt.Printf("[Device.readHandler] done\n")

	for {
		select {
		case packet := <-d.read:
			switch d.conn.(type) {
			case *net.UDPConn:
				//payloadType := packet[0:2]
				//payloadLength := packet[2:4]
				sequence := packet[4:8]
				payload := packet[8:]

				// trim prefix & EOD from payload
				cmdBytes := payload
				if len(payload) > 2 && cmdBytes[0] == 0x90 && cmdBytes[len(cmdBytes)-1] == 0xff {
					cmdBytes = payload[1 : len(payload)-1]
				}

				// common responses
				var autoFin bool

				switch {
				case bytes.Equal(payload, []byte{0x1}):
					//fmt.Printf(">> Ok SeqReset\n")
				case bytes.Equal(payload, []byte{0xf, 0x1}):
					//fmt.Printf(">> abnormal seq number\n")
				case bytes.Equal(payload, []byte{0xf, 0x2}):
					//fmt.Printf(">> abnormal message type\n")

				case bytes.Equal(cmdBytes, []byte{0x40}) || bytes.Equal(cmdBytes, []byte{0x41}):
					//fmt.Printf(">> Ack (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x50}) || bytes.Equal(cmdBytes, []byte{0x51}):
					//fmt.Printf(">> Fin (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x60, 0x2}):
					//fmt.Printf(">> Syntax Error\n")
					autoFin = true
				case bytes.Equal(cmdBytes, []byte{0x60, 0x3}):
					//fmt.Printf(">> Command Buffer Full\n")
				case bytes.Equal(cmdBytes, []byte{0x60, 0x4}) || bytes.Equal(cmdBytes, []byte{0x61, 0x4}):
					//fmt.Printf(">> Command Canceled (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x60, 0x5}) || bytes.Equal(cmdBytes, []byte{0x61, 0x5}):
					//fmt.Printf(">> No socket (Socket#%d)\n", cmdBytes[0]&0xf)
					autoFin = true
				case bytes.Equal(cmdBytes, []byte{0x60, 0x41}) || bytes.Equal(cmdBytes, []byte{0x61, 0x41}):
					//fmt.Printf(">> Command not executable (Socket#%d)\n", cmdBytes[0]&0xf)
					autoFin = true
				}

				// use seq to find cmd and run it's handler
				seqUint32 := binary.BigEndian.Uint32(sequence)
				if cmd, ok := d.writeSeqCmd.Load(seqUint32); ok {
					if cmd, yes := cmd.(Cmd); yes {
						//fmt.Printf(">> % X\n", cmdBytes)
						if autoFin {
							cmd.Finish()
						} else {
							cmd.HandleReply(cmdBytes, d)
						}
					} else {
						//fmt.Printf(">> Found seq is not Cmd\n")
					}
				} else {
					//fmt.Printf(">> Sequence # not found: [%X] {%d}\n", sequence, seqUint32)
				}
			default:
				//fmt.Printf(">> Not UDP conn\n")
			}

		case <-d.Done():
			close(d.read)
			return
		}
	}
}
