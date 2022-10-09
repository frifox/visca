package visca

import (
	"context"
	"time"
)

// first/last bytes
const (
	CamID = 0x81
	EOL   = 0xff
)

// byte 2
const (
	doCommand = 0x1
	doInquiry = 0x9
)

// byte 3
const (
	toInterface = 0x0
	toCamera    = 0x4
	toCamera2   = 0x5
	toMotors    = 0x6
	toConfig    = 0x7e
)

type Cmd interface {
	String() string
	HandleReply([]byte, *Device)

	context.Context
	InitContext()
	Finish()
}
type CmdAppliable interface {
	Apply(*Device) bool
}

type ViscaCommand interface {
	ViscaCommand() []byte
}
type ViscaInquiry interface {
	ViscaInquiry() []byte
}
type ViscaReply interface {
	ViscaReply() []byte
}
type DeviceSettingCommand interface {
	DeviceSettingCommand() []byte
}
type ControlCommand interface {
	ControlCommand() []byte
}
type ControlCommandReply interface {
	ControlCommandReply() []byte
}

//type ViscaInquiry CmdType
//type ViscaReply CmdType
//type DeviceSettingCommand CmdType
//type ControlCommand CmdType
//type ControlCommandReply CmdType

type CmdContext struct {
	context.Context
	context.CancelFunc
}

func (c *CmdContext) InitContext() {
	c.Context, c.CancelFunc = context.WithTimeout(context.Background(), time.Millisecond*100)
}
func (c *CmdContext) Finish() {
	if c.Err() == nil {
		c.CancelFunc()
	}
}
