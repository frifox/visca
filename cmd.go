package visca

import (
	"context"
	"fmt"
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
	toMotors    = 0x6
	toConfig    = 0x7e
)

type Cmd interface {
	String() string
	HandleReply([]byte, *Device)

	InitContext()
	context.Context
	Finish()
}
type CmdAppliable interface {
	Apply(*Device) bool
}

//type CmdWaitable interface {
//	WaitReply(*Device)
//}

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

type Raw struct {
	Cmd []byte
}

func (c *Raw) String() string {
	return fmt.Sprintf("Raw{% X}", c.Cmd)
}
func (c *Raw) ViscaCommand() []byte {
	data := []byte{CamID}
	data = append(data, c.Cmd...)
	data = append(data, EOL)
	return data
}
func (c *Raw) HandleReply(data []byte, device *Device) {
	fmt.Printf(">> Raw Reply: [% X]\n", data)
}
