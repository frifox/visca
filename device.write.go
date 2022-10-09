package visca

import (
	"fmt"
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
			fmt.Printf("[Device.Writer] device ctx Done\n")
			close(d.write)
			return
		}
	}
}
