package visca

import (
	"math"
)

type Command interface {
	Bytes() []byte
}

type MoveHome struct{}

func (a MoveHome) Bytes() []byte {
	return []byte{0x1, 0x6, 0x4}
}

type Focus struct{}

func (a Focus) Bytes() []byte {
	return []byte{0x1, 0x4, 0x38, 0x2}
}

type SavePreset struct {
	ID uint8
}

func (a SavePreset) Bytes() []byte {
	return []byte{0x1, 0x4, 0x3f, 0x1, a.ID}
}

type CallPreset struct {
	ID uint8
}

func (a CallPreset) Bytes() []byte {
	return []byte{0x1, 0x4, 0x3f, 0x1, a.ID}
}

type Move struct {
	X float64
	Y float64
}

func (a Move) Bytes() []byte {
	// range: 0x0 - 0x18 (0 - 24)
	max := 24.0
	max = max * 0.5 // limit max to 50%!

	xStep := byte(math.Round(max * math.Abs(a.X)))
	yStep := byte(math.Round(max * math.Abs(a.Y)))

	var x byte
	switch true {
	case a.X > 0:
		x = byte(0x2) // right
	case a.X < 0:
		x = byte(0x1) // left
	default:
		x = byte(0x3) // none
	}

	var y byte
	switch true {
	case a.Y > 0:
		y = byte(0x1) // up
	case a.Y < 0:
		y = byte(0x2) // down
	default:
		y = byte(0x3) // none
	}

	return []byte{0x1, 0x6, 0x01, xStep, yStep, x, y}
}

type Zoom struct {
	Z float64
}

func (a *Zoom) Bytes() []byte {
	// 1-7 = zoom, 0 = stop

	steps := 7.0
	step := math.Round(steps * math.Abs(a.Z))

	var z byte
	switch true {
	case a.Z > 0:
		z = byte(step + 0x20)
	case a.Z < 0:
		z = byte(step + 0x30)
	default:
		z = byte(0x0)
	}

	return []byte{0x1, 0x4, 0x7, z}
}

type OSDToggle struct{}

func (a OSDToggle) Bytes() []byte {
	return []byte{0x1, 0x4, 0x3f, 0x2, 0x5f}
}

type OSDEnter struct{}

func (a OSDEnter) Bytes() []byte {
	return []byte{0x1, 0x6, 0x6, 0x5}
}

type OSDReturn struct{}

func (a OSDReturn) Bytes() []byte {
	return []byte{0x1, 0x6, 0x6, 0x4}
}

type OSDLeft struct{}

func (a OSDLeft) Bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x1, 0x3}
}

type OSDRight struct{}

func (a OSDRight) Bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x2, 0x3}
}

type OSDUp struct{}

func (a OSDUp) Bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x3, 0x1}
}

type OSDDown struct{}

func (a OSDDown) Bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x3, 0x2}
}
