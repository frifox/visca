package visca

import "time"

type OSDToggle struct {
	MenuIsOn bool
	runtime  [2]time.Time // [write, ack]
}

func (a *OSDToggle) apply() bool {
	return true
}
func (a *OSDToggle) bytes() []byte {
	//return []byte{0x1, 0x4, 0x3f, 0x2, 0x5f}
	return []byte{0x1, 0x6, 0x6, 0x10}
}

type OSDEnter struct {
	alwaysApply
}

func (a *OSDEnter) bytes() []byte {
	//return []byte{0x1, 0x6, 0x6, 0x5}
	return []byte{0x1, 0x7e, 0x1, 0x2, 0x0, 0x1}
}

type OSDReturn struct {
	alwaysApply
}

func (a *OSDReturn) bytes() []byte {
	return []byte{0x1, 0x6, 0x6, 0x4}
}

type OSDLeft struct {
	alwaysApply
}

func (a *OSDLeft) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x1, 0x3}
}

type OSDRight struct {
	alwaysApply
}

func (a *OSDRight) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x2, 0x3}
}

type OSDUp struct {
	alwaysApply
}

func (a *OSDUp) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x3, 0x1}
}

type OSDDown struct {
	alwaysApply
}

func (a *OSDDown) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x3, 0x2}
}

// AskMenuStatus asks camera if menu is open or not
type AskMenuStatus struct {
	alwaysApply
}

func (c *AskMenuStatus) bytes() []byte {
	return []byte{0x9, 0x6, 0x6}
}
