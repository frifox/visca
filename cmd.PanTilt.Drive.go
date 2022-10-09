package visca

import (
	"fmt"
	"math"
)

// ptz
const (
	moveLeft  = 0x1
	moveRight = 0x2
	moveUp    = 0x1
	moveDown  = 0x2
	moveNone  = 0x3
)

type PanTiltDrive struct {
	CmdContext

	X float64
	Y float64

	x int8
	y int8
}

func (c *PanTiltDrive) String() string {
	return fmt.Sprintf("%T{x:%d, y:%d}", *c, c.x, c.y)
}

func (c *PanTiltDrive) Apply(device *Device) (needToSend bool) {
	xSteps := float64(0x18) // 0x1 to 0x18
	if device.Config.XMaxSpeed > 0 {
		xSteps = xSteps * device.Config.XMaxSpeed
	}

	ySteps := float64(0x18) // 0x1 to 0x18 TODO 0x1-0x17 when SlowPanTilt=0
	if device.Config.YMaxSpeed > 0 {
		ySteps = ySteps * device.Config.YMaxSpeed
	}

	c.x = int8(math.Ceil(xSteps * math.Abs(c.X)))
	c.y = int8(math.Ceil(ySteps * math.Abs(c.Y)))

	if c.X < 0 {
		c.x = -c.x
	}
	if c.Y < 0 {
		c.y = -c.y
	}

	// no changes?
	if c.x == device.State.PanTiltDrive.x && c.y == device.State.PanTiltDrive.y {
		//fmt.Printf(">> No change %s\n", c)
		return false
	}

	device.State.PanTiltDrive.X = c.X
	device.State.PanTiltDrive.Y = c.Y
	device.State.PanTiltDrive.x = c.x
	device.State.PanTiltDrive.y = c.y

	return true
}

func (c *PanTiltDrive) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toMotors, 0x1}

	// x speed
	if c.x < 0 {
		data = append(data, byte(-c.x))
	} else {
		data = append(data, byte(c.x))
	}

	// y speed (
	if c.y < 0 {
		data = append(data, byte(-c.y))
	} else {
		data = append(data, byte(c.y))
	}

	// x direction
	switch true {
	case c.X > 0:
		data = append(data, moveRight)
	case c.X < 0:
		data = append(data, moveLeft)
	default:
		data = append(data, moveNone)
	}

	// y direction
	switch true {
	case c.Y > 0:
		data = append(data, moveUp)
	case c.Y < 0:
		data = append(data, moveDown)
	default:
		data = append(data, moveNone)
	}

	data = append(data, EOL)

	return data
}

func (c *PanTiltDrive) HandleReply(data []byte, device *Device) {
	if len(data) != 1 {
		fmt.Printf("[PanTiltDrive.HandleReply] BAD REPLY [% X]\n", data)
		return
	}

	p := data[0] & 0xf0

	switch p {
	case 0x40: // ack
	case 0x50: // fin
		c.Finish()
	default:
		fmt.Printf("[PanTiltDrive.HandleReply] Unknown [% X]\n", data)
	}
}
