package visca

import "time"

type Command interface {
	apply() bool   // do we need to write to camera?
	bytes() []byte // command bytes for writing
}

type alwaysApply struct{}

func (c *alwaysApply) apply() bool {
	return true
}

// ungrouped commands

type MoveHome struct {
	alwaysApply
}

func (a *MoveHome) bytes() []byte {
	return []byte{0x1, 0x6, 0x4}
}

type Focus struct {
	alwaysApply
	runtime [3]time.Time // [write, fin, fin]
}

func (a *Focus) bytes() []byte {
	return []byte{0x1, 0x4, 0x38, 0x2}
}

type SavePreset struct {
	ID uint8
	alwaysApply
}

func (a *SavePreset) bytes() []byte {
	return []byte{0x1, 0x4, 0x3f, 0x1, a.ID}
}

type CallPreset struct {
	ID uint8
	alwaysApply
	runtime [3]time.Time // [write, ack, fin]
}

func (a *CallPreset) bytes() []byte {
	return []byte{0x1, 0x4, 0x3f, 0x2, a.ID}
}
