package visca

import (
	"fmt"
	"math"
	"time"
)

type Zoom struct {
	Z float64
	z int8

	ack     bool
	fin     bool
	ackChan chan bool
	finChan chan bool
	start   time.Time
}

func (c *Zoom) String() string {
	return fmt.Sprintf("Zoom{z:%d}", c.z)
}

func (c *Zoom) Apply(device *Device) bool {
	steps := float64(0x8) // 0x0 to 0x7
	if device.Config.ZMaxSpeed > 0 {
		steps = steps * device.Config.ZMaxSpeed
	}

	step := int8(math.Ceil(steps * math.Abs(c.Z)))
	switch {
	case c.Z > 0:
		c.z = step
	case c.Z < 0:
		c.z = -step
	case c.Z == 0:
		c.z = 0
	}

	// no changes?
	if c.z == device.State.Zoom.z {
		//fmt.Printf(">> No change %s\n", c)
		return false
	}

	//fmt.Printf(">> state updated %s\n", c)
	//device.State.Zoom = *c
	device.State.Zoom.Z = c.Z
	device.State.Zoom.z = c.z

	c.start = time.Now()
	c.ackChan = make(chan bool)
	c.finChan = make(chan bool)

	return true
}

//func (c *Zoom) Send(device *Device) bool {
//	if device.State.ZoomLastSeq == nil {
//		return true
//	}
//
//	// send only if last cmd was Ack & Fin
//	if cmd, ok := device.writeSeqCmd.Load(*device.State.ZoomLastSeq); ok {
//		if cmd, ok := cmd.(*Zoom); ok {
//			return cmd.ack && cmd.fin
//		}
//	}
//
//	// otherwise don't send. It'll get auto-sent after last Cmd is fin
//	return false
//}

func (c *Zoom) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x7}

	// speed (8 steps: 0x0 - 0x7)
	switch {
	case c.Z > 0:
		step := byte(c.z)
		data = append(data, 0x20|step-1) // zoom in
	case c.Z < 0:
		step := byte(-c.z)
		data = append(data, 0x30|step-1) // zoom out
	default:
		data = append(data, 0x0) // zoom stop
	}

	data = append(data, EOL)

	return data
}

func (c *Zoom) WaitReply(device *Device) {
	go c.WaitAck(device)
	go c.WaitFin(device)
}

func (c *Zoom) WaitAck(device *Device) {
	select {
	case <-time.After(time.Millisecond * 100):
		fmt.Printf(">> Zoom ACK timeout!\n")
	case <-c.ackChan:
		//fmt.Printf(">> Zoom ACK %d ms\n", time.Now().Sub(c.start).Milliseconds())
	}

	device.ZoomReady.Done()
}
func (c *Zoom) WaitFin(device *Device) {
	select {
	case <-time.After(time.Millisecond * 100):
		fmt.Printf(">> Zoom FIN timeout!\n")
	case <-c.finChan:
		//fmt.Printf(">> Zoom FIN %d ms\n", time.Now().Sub(c.start).Milliseconds())
	}

	device.ZoomReady.Done()
}

func (c *Zoom) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[Zoom.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] & 0xf0 {
	case 0x40:
		//fmt.Printf("[Zoom.HandleReply] ACK\n")
		c.ack = true
		c.ackChan <- true
	case 0x50:
		//fmt.Printf("[Zoom.HandleReply] Fin\n")
		c.fin = true
		c.finChan <- true
	default:
		fmt.Printf("[Zoom.HandleReply] Unknown [% X]\n", data)
	}
}
