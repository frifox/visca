package visca

import (
	"fmt"
	"math"
)

type Zoom struct {
	CmdContext

	Z float64
	z int8
}

func (c *Zoom) String() string {
	return fmt.Sprintf("%T{z:%d}", *c, c.z)
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

	//c.start = time.Now()
	//c.ackChan = make(chan bool)
	//c.finChan = make(chan bool)

	return true
}

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

func (c *Zoom) HandleReply(data []byte, device *Device) {
	if len(data) != 1 {
		fmt.Printf("[Zoom.HandleReply] BAD REPLY [% X]\n", data)
		return
	}

	p := data[0] & 0xf0

	switch p {
	case 0x40: // ack
	case 0x50: // fin
		c.Finish()
	default:
		fmt.Printf("[Zoom.HandleReply] Unknown [% X]\n", data)
	}
}
