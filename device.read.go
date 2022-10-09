package visca

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func (d *Device) PanTiltQueueWorker() {
	for {
		// wait for cmd
		//fmt.Printf(">> waiting queue\n")
		cmd := d.PanTiltQueue.Get()
		//fmt.Printf(">> got queue\n")

		// if state is same as before, ignore it
		needToSend := cmd.Apply(d)
		if !needToSend {
			//fmt.Printf(">> not applying\n")
			continue
		}

		// send it to camera
		//fmt.Printf(">> doing cmd\n")
		d.Do(cmd, true)
		//fmt.Printf(">> cmd do'ed\n")
	}
}

func (d *Device) ZoomQueueWorker() {
	for {
		// wait for cmd
		//fmt.Printf(">> waiting queue\n")
		cmd := d.ZoomQueue.Get()
		//fmt.Printf(">> got queue\n")

		// if state is same as before, ignore it
		needToSend := cmd.Apply(d)
		if !needToSend {
			//fmt.Printf(">> not applying\n")
			continue
		}

		// send it to camera
		//fmt.Printf(">> doing cmd\n")
		d.Do(cmd, true)
		//fmt.Printf(">> cmd do'ed\n")
	}
}

func (d *Device) Writer() {
	fmt.Printf("[Device.Writer] init\n")
	defer fmt.Printf("[Device.Writer] done\n")

	for {
		select {
		case data := <-d.write:
			var err error

			// send
			switch port := d.conn.(type) {
			case *net.UDPConn:
				_, err = port.WriteToUDP(data, d.remoteAddr)
			default:
				_, err = port.Write(data)
			}
			if err != nil {
				fmt.Printf("[Device.Writer] [% X] ERR %v\n", data, err)
			} else {
				//fmt.Printf("[Device.Writer] Wrote [% X]\n", data)
			}

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
			// any msg from camera = it's ready to receive new cmd
			//if atomic.CompareAndSwapUint32(&d.writeLockedStatus, 1, 0) {
			//	d.writeLock.Unlock()
			//}

			switch d.conn.(type) {
			case *net.UDPConn:
				// any response = ready to send next cmd

				//payloadType := packet[0:2]
				//payloadLength := packet[2:4]
				sequence := packet[4:8]
				payload := packet[8:]

				var cmdBytes []byte
				if len(payload) > 2 {
					cmdBytes = payload[1 : len(payload)-1]
				}

				// common responses
				switch {
				case bytes.Equal(payload, []byte{0x1}):
					fmt.Printf(">> Ok SeqReset\n")
				case bytes.Equal(payload, []byte{0xf, 0x1}):
					fmt.Printf(">> abnormal seq number\n")
				case bytes.Equal(payload, []byte{0xf, 0x2}):
					fmt.Printf(">> abnormal message type\n")

				case bytes.Equal(cmdBytes, []byte{0x40}) || bytes.Equal(cmdBytes, []byte{0x41}):
					fmt.Printf(">> Ack (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x50}) || bytes.Equal(cmdBytes, []byte{0x51}):
					fmt.Printf(">> Fin (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x60, 0x2}):
					fmt.Printf(">> Syntax Error\n")
				case bytes.Equal(cmdBytes, []byte{0x60, 0x3}):
					fmt.Printf(">> Command Buffer Full\n")
				case bytes.Equal(cmdBytes, []byte{0x60, 0x4}) || bytes.Equal(cmdBytes, []byte{0x61, 0x4}):
					fmt.Printf(">> Command Canceled (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x60, 0x5}) || bytes.Equal(cmdBytes, []byte{0x61, 0x5}):
					fmt.Printf(">> No socket (Socket#%d)\n", cmdBytes[0]&0xf)
				case bytes.Equal(cmdBytes, []byte{0x60, 0x41}) || bytes.Equal(cmdBytes, []byte{0x61, 0x41}):
					fmt.Printf(">> Command not executable (Socket#%d)\n", cmdBytes[0]&0xf)
				}

				// use seq to find cmd and run it's handler
				seqUint32 := binary.BigEndian.Uint32(sequence)
				if cmd, ok := d.writeSeqCmd.Load(seqUint32); ok {
					if cmd, yes := cmd.(Cmd); yes {
						// auto-fin if got an error
						if bytes.Equal(cmdBytes, []byte{0x60, 0x41}) || bytes.Equal(cmdBytes, []byte{0x61, 0x41}) {
							fmt.Printf(">> Command not executable (Socket#%d)\n", cmdBytes[0]&0xf)
							if cmd.Err() == nil {
								cmd.Finish()
							}
						} else {
							cmd.HandleReply(payload, d)
						}
					} else {
						fmt.Printf(">> Found seq is not Cmd\n")
					}
				} else {
					fmt.Printf(">> Sequence # not found: [%X] {%d}\n", sequence, seqUint32)
				}
			default:
				fmt.Printf(">> Not UDP conn\n")
			}

			//fmt.Printf("[Device.readHandler] Handling [% X]\n", packet)

			//fmt.Printf("[Device.readHandler] Handled [% X]\n", packet)

		case <-d.Done():
			close(d.read)
			return
		}
	}
}
