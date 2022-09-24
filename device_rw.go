package visca

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func (d *Device) PanTiltQueueWorker() {
	for {
		// wait for camera to be ready
		//fmt.Printf(">> waiting ptready\n")
		d.PanTiltReady.Wait()
		//fmt.Printf(">> ptready\n")

		// wait for cmd
		//fmt.Printf(">> waiting queue\n")
		cmd := d.PanTiltQueue.Get()
		//fmt.Printf(">> got queue\n")

		// if state is same as before, ignore it
		apply := cmd.Apply(d)
		if !apply {
			//fmt.Printf(">> not applying\n")
			continue
		}

		// send it to camera
		//fmt.Printf(">> doing cmd\n")
		d.PanTiltReady.Add(2)
		d.Do(cmd, true)
		//fmt.Printf(">> cmd do'ed\n")
	}
}

func (d *Device) ZoomQueueWorker() {
	for {
		// wait for camera to be ready
		//fmt.Printf(">> waiting ptready\n")
		d.ZoomReady.Wait()
		//fmt.Printf(">> ptready\n")

		// wait for cmd
		//fmt.Printf(">> waiting queue\n")
		cmd := d.ZoomQueue.Get()
		//fmt.Printf(">> got queue\n")

		// if state is same as before, ignore it
		apply := cmd.Apply(d)
		if !apply {
			//fmt.Printf(">> not applying\n")
			continue
		}

		// send it to camera
		//fmt.Printf(">> doing cmd\n")
		d.ZoomReady.Add(2)
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
				//fmt.Printf("[Device.Writer] Wrote %d bytes [% X]\n", n, data)
			}

			// TODO figure out ready-to-send mechanic
			//time.Sleep(time.Millisecond * 100)

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

				switch {
				// probably response to SeqReset
				case len(payload) == 1:
					if payload[0] == 0x1 {
						//fmt.Printf(">> Ok\n")
					}

				case len(payload) == 3:
					cmd := payload[1 : len(payload)-1]
					switch {
					case cmd[0]&0xf0 == 0x40:
						//fmt.Printf(">> Ack\n")
					case cmd[0]&0xf0 == 0x50:
						//fmt.Printf(">> Fin\n")
					}

				case len(payload) == 4:
					cmd := payload[1 : len(payload)-1]
					switch {
					case cmd[0] == 0x60 && cmd[1] == 0x2:
						fmt.Printf(">> Syntax Error\n")
					case cmd[0] == 0x60 && cmd[1] == 0x3:
						fmt.Printf(">> Command Buffer Full\n")
					case cmd[0]&0xf0 == 0x60 && cmd[1] == 0x4:
						fmt.Printf(">> Command Canceled\n")
					case cmd[0]&0xf0 == 0x60 && cmd[1] == 0x5:
						fmt.Printf(">> No Socket\n")
					case cmd[0]&0xf0 == 0x60 && cmd[1] == 0x41:
						fmt.Printf(">> Command Not Executable\n")
					}
				}

				// find related cmd and launch reply handler
				seqUint32 := binary.BigEndian.Uint32(sequence)
				if cmd, ok := d.writeSeqCmd.Load(seqUint32); ok {
					// cmd found!
					if cmd, yes := cmd.(Cmd); yes {
						// cmd is valid!
						cmd.HandleReply(payload, d)
					} else {
						fmt.Printf(">> Found seq is not Cmd\n")
					}
				} else {
					fmt.Printf(">> Sequence # not found: [%X] {%d}\n", sequence, seqUint32)
				}

			default:
				//switch {
				//// Command ACK
				//case len(packet) == 3 && packet[0] == 0x90 && packet[1]>>4 == 0x4 && packet[2] == 0xff:
				//	//fmt.Printf("[Camera] Read [% X] accepted\n", packet)
				//// Command Completion
				//case len(packet) == 3 && packet[0] == 0x90 && packet[1]>>4 == 0x5 && packet[2] == 0xff:
				//	//fmt.Printf("[Camera] Read [% X] finished\n", packet)
				//// MenuStatus responses
				//case bytes.Equal(packet, []byte{0x90, 0x50, 0x2, 0xFF}):
				//	//fmt.Println("MENU IS ON")
				//	//d.OSDToggle.MenuIsOn = true
				//case bytes.Equal(packet, []byte{0x90, 0x50, 0x3, 0xFF}):
				//	//fmt.Println("MENU IS OFF")
				//	//d.OSDToggle.MenuIsOn = false
				//// unhandled
				//default:
				//	fmt.Printf("[Device.readHandler] Unhandled [% X]\n", packet)
				//	continue
				//}
			}

			//fmt.Printf("[Device.readHandler] Handling [% X]\n", packet)

			//fmt.Printf("[Device.readHandler] Handled [% X]\n", packet)

		case <-d.Done():
			close(d.read)
			return
		}
	}
}
