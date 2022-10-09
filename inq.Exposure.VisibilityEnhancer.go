package visca

import (
	"fmt"
)

type InqExposureBackLight struct {
	CmdContext
	On bool
}

func (c *InqExposureBackLight) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureBackLight) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x33}
	data = append(data, EOL)
	return data
}

func (c *InqExposureBackLight) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]

	switch p {
	case 0x2:
		c.On = true
	case 0x3:
		c.On = false
	default:
		fmt.Printf(">> unknown %T value [%X]\n", *c, p)
	}

	device.Inquiry.InqExposureBackLight = c
}
