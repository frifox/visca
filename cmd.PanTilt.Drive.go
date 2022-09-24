package visca

import (
	"fmt"
	"math"
	"time"
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
	X float64
	Y float64

	x int8
	y int8

	ack     bool
	fin     bool
	ackChan chan bool
	finChan chan bool
	start   time.Time
}

func (c *PanTiltDrive) String() string {
	return fmt.Sprintf("PanTiltDrive{x:%d, y:%d}", c.x, c.y)
}

func (c *PanTiltDrive) Apply(device *Device) bool {
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

	//fmt.Printf(">> state updated %s\n", c)
	//device.State.PanTiltDrive = *c
	device.State.PanTiltDrive.X = c.X
	device.State.PanTiltDrive.Y = c.Y
	device.State.PanTiltDrive.x = c.x
	device.State.PanTiltDrive.y = c.y

	c.start = time.Now()
	c.ackChan = make(chan bool)
	c.finChan = make(chan bool)

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

func (c *PanTiltDrive) WaitReply(device *Device) {
	go c.WaitAck(device)
	go c.WaitFin(device)
}

func (c *PanTiltDrive) WaitAck(device *Device) {
	select {
	case <-time.After(time.Millisecond * 100):
		fmt.Printf(">> PanTilt ACK timeout!\n")
	case <-c.ackChan:
		//fmt.Printf(">> PanTilt ACK %d ms\n", time.Now().Sub(c.start).Milliseconds())
	}

	device.PanTiltReady.Done()
}
func (c *PanTiltDrive) WaitFin(device *Device) {
	select {
	case <-time.After(time.Millisecond * 100):
		fmt.Printf(">> PanTilt FIN timeout!\n")
	case <-c.finChan:
		//fmt.Printf(">> PanTilt FIN %d ms\n", time.Now().Sub(c.start).Milliseconds())
	}

	device.PanTiltReady.Done()
}

func (c *PanTiltDrive) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[PanTiltDrive.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] & 0xf0 {
	case 0x40:
		//fmt.Printf("[PanTiltDrive.HandleReply] ACK %s\n", c)
		c.ack = true
		c.ackChan <- true
	case 0x50:
		//fmt.Printf("[PanTiltDrive.HandleReply] FIN %s\n", c)
		c.fin = true
		c.finChan <- true
	default:
		fmt.Printf("[PanTiltDrive.HandleReply] Unknown [% X]\n", data)
	}
}
