package visca

import (
	"fmt"
)

type InqPower struct {
	CmdContext
	On bool
}

func (c *InqPower) String() string {
	return fmt.Sprintf("%T{On:%t}", *c, c.On)
}

func (c *InqPower) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x0}
	data = append(data, EOL)
	return data
}

func (c *InqPower) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> bad reply [% X]\n", data)
		return
	}

	p := data[1]

	switch p {
	case 0x2:
		c.On = true
		device.State.Power.On = true
	case 0x3:
		c.On = false
		device.State.Power.On = false
	default:
		//fmt.Printf(">> unknown %T value [%X]\n", *c, p)
	}

	device.Inquiry.InqPower = c
}
