package visca

import (
	"fmt"
)

type InqExposureComp struct {
	CmdContext
	On bool
}

func (c *InqExposureComp) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureComp) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0xe3}
	data = append(data, EOL)
	return data
}

func (c *InqExposureComp) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	val := data[1]
	switch val {
	case 0x2:
		c.On = true
	case 0x3:
		c.On = false
	default:
		fmt.Printf(">> ")
	}
	p := data[1:2]
	c.On = int(sonyInt(p))

	device.Inquiry.ExposureAESpeed = c
}
