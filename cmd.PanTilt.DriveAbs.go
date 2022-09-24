package visca

import (
	"fmt"
	"math"
)

type PanTiltDriveAbs struct {
	X      []byte
	Y      []byte
	SpeedX float64
	SpeedY float64

	x      []byte
	y      []byte
	speedX uint8
	speedY uint8
}

func (c *PanTiltDriveAbs) String() string {
	return fmt.Sprintf("PanTiltDriveAbs{x:%X, y:%X}", c.x, c.y)
}

func (c *PanTiltDriveAbs) Apply(device *Device) bool {
	c.x = c.X
	c.y = c.Y

	stepsX := float64(0x18)
	stepsY := float64(0x18) // 0x17 if not in slow mode

	c.speedX = uint8(math.Ceil(stepsX * math.Abs(c.SpeedX)))
	c.speedY = uint8(math.Ceil(stepsY * math.Abs(c.SpeedY)))

	return true
}

func (c *PanTiltDriveAbs) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toMotors, 0x2}

	// speed 0x1 to 0x18
	data = append(data, c.speedX) // pan
	data = append(data, c.speedY) // tilt

	// position
	data = append(data, c.x...)
	data = append(data, c.y...)

	data = append(data, EOL)

	return data
}

func (c *PanTiltDriveAbs) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[PanTiltDriveAbs.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] {
	case 0x41:
		fmt.Printf("[PanTiltDriveAbs.HandleReply] ACK\n")
	case 0x51:
		fmt.Printf("[PanTiltDriveAbs.HandleReply] FIN\n")
	default:
		fmt.Printf("[PanTiltDriveAbs.HandleReply] Unknown [% X]\n", data)
	}
}
