package visca

import (
	"math"
)

type Zoom struct {
	Z      float64
	StepsZ float64

	step uint8
	//device *Device
}

func (c *Zoom) apply() bool {
	// 1-7 = zoom, 0 = stop
	str := math.Abs(c.Z)
	step := uint8(math.Ceil(c.StepsZ * str))

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

type ZoomTo struct {
	Step uint16
	step uint16
}

func (c *ZoomTo) cmdType() interface{} {
	return ViscaCommand{}
}

func (c *ZoomTo) apply() bool {
	if c.Step == c.step {
		return false
	}

	c.step = c.Step

	return true
}
func (c *ZoomTo) bytes() []byte {
	var b []byte
	switch c.step {
	case 1:
		b = []byte{0, 0, 0, 0}
	case 2:
		b = []byte{0x0, 0xd, 0xc, 0x1}
	case 3:
		b = []byte{0x1, 0x8, 0x6, 0xc}
	case 4:
		b = []byte{0x2, 0x0, 0x1, 0x5}
	case 5:
		b = []byte{0x2, 0x5, 0x9, 0x4}
	case 6:
		b = []byte{0x2, 0x9, 0xb, 0x7}
	case 7:
		b = []byte{0x2, 0xC, 0xF, 0xB}
	case 8:
		b = []byte{0x2, 0xF, 0xB, 0x0}
	case 9:
		b = []byte{0x3, 0x2, 0x0, 0xC}
	case 10:
		b = []byte{0x3, 0x4, 0x2, 0xD}
	case 11:
		b = []byte{0x3, 0x6, 0x0, 0x8}
	case 12:
		b = []byte{0x3, 0x7, 0xA, 0xA}
	case 13:
		b = []byte{0x3, 0x9, 0x1, 0xC}
	case 14:
		b = []byte{0x3, 0xA, 0x6, 0x6}
	case 15:
		b = []byte{0x3, 0xB, 0x9, 0x0}
	case 16:
		b = []byte{0x3, 0xC, 0x9, 0xC}
	case 17:
		b = []byte{0x3, 0xD, 0x9, 0x1}
	case 18:
		b = []byte{0x3, 0xE, 0x7, 0x2}
	case 19:
		b = []byte{0x3, 0xF, 0x4, 0x0}
	case 20:
		b = []byte{0x4, 0x0, 0x0, 0x0}
	case 30:
		b = []byte{0x5, 0x5, 0x5, 0x6}
	case 40:
		b = []byte{0x6, 0x0, 0x0, 0x0}
	case 60:
		b = []byte{0x6, 0xA, 0xA, 0xB}
	case 80:
		b = []byte{0x7, 0x0, 0x0, 0x0}
	case 100:
		b = []byte{0x7, 0x3, 0x3, 0x4}
	case 120:
		b = []byte{0x7, 0x5, 0x5, 0x6}
	case 140:
		b = []byte{0x7, 0x6, 0xD, 0xC}
	case 160:
		b = []byte{0x7, 0x8, 0x0, 0x0}
	case 180:
		b = []byte{0x7, 0x8, 0xE, 0x4}
	case 200:
		b = []byte{0x7, 0x9, 0x9, 0xA}
	case 220:
		b = []byte{0x7, 0xA, 0x2, 0xF}
	case 240:
		b = []byte{0x7, 0xA, 0xC, 0x0}
	}

	cmd := []byte{0x1, 0x4, 0x47}
	cmd = append(cmd, b...)
	return cmd
}
