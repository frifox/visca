package visca

import (
	"math"
)

type Zoom struct {
	Z float64

	step   uint8
	device *Device
}

func (c *Zoom) apply() bool {
	// 1-7 = zoom, 0 = stop
	steps := 7.0

	str := math.Abs(c.Z)
	step := uint8(math.Ceil(steps * str))

	if c.step == step {
		return false
	}

	// save new state
	c.step = step

	return true
}
func (c *Zoom) bytes() []byte {
	var step byte
	switch true {
	case c.Z > 0:
		step = c.step + 0x20 // zoom in
	case c.Z < 0:
		step = c.step + 0x30 // zom out
	default:
		step = byte(0x0) // stop zoom
	}

	return []byte{0x1, 0x4, 0x7, step}
}
