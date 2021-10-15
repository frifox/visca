package visca

import (
	"math"
	"time"
)

type Move struct {
	X float64
	Y float64

	xStep  uint8
	yStep  uint8
	device *Device

	runtime [3]time.Time // [write, fin, fin]
}

func (c *Move) apply() bool {
	xMax := float64(0x18) * 0.5
	xStep := byte(math.Ceil(xMax * math.Abs(c.X)))

	yMax := float64(0x14) * 0.5
	yStep := byte(math.Ceil(yMax * math.Abs(c.Y)))

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
	var xDir, yDir byte

	// x-axis direction
	switch true {
	case c.X > 0:
		xDir = byte(0x2) // right
	case c.X < 0:
		xDir = byte(0x1) // left
	default:
		xDir = byte(0x3) // none
	}

	// y-axis direction
	switch true {
	case c.Y > 0:
		yDir = byte(0x1) // up
	case c.Y < 0:
		yDir = byte(0x2) // down
	default:
		yDir = byte(0x3) // none
	}

	return []byte{0x1, 0x6, 0x01, c.xStep, c.yStep, xDir, yDir}
}
