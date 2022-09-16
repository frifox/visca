package visca

import (
	"bytes"
	"math"
	"time"
)

type Move struct {
	X    float64
	Y    float64
	XMax float64
	YMax float64

	xStep  uint8
	yStep  uint8
	device *Device

	runtime [3]time.Time // [write, fin, fin]
}

func (c *Move) apply() bool {
	xStep := byte(math.Ceil(c.XMax * math.Abs(c.X)))
	yStep := byte(math.Ceil(c.YMax * math.Abs(c.Y)))

	// no changes?
	if c.xStep == xStep && c.yStep == yStep {
		return false
	}

	// save new state
	c.xStep = xStep
	c.yStep = yStep

	return true
}

func (c *Move) bytes() []byte {
	packet := bytes.Buffer{}

	// header
	packet.WriteByte(0x1)
	packet.WriteByte(0x6)
	packet.WriteByte(0x1)

	// speed
	packet.WriteByte(c.xStep)
	packet.WriteByte(c.yStep)

	// x-axis direction
	switch true {
	case c.X > 0:
		packet.WriteByte(0x2) // right
	case c.X < 0:
		packet.WriteByte(0x1) // left
	default:
		packet.WriteByte(0x3) // none
	}

	// y-axis direction
	switch true {
	case c.Y > 0:
		packet.WriteByte(0x1) // up
	case c.Y < 0:
		packet.WriteByte(0x2) // down
	default:
		packet.WriteByte(0x3) // none
	}

	return packet.Bytes()
}
