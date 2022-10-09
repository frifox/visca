package visca

import (
	"fmt"
)

type InqExposureSpotLight struct {
	CmdContext
	On bool
}

func (c *InqExposureSpotLight) String() string {
	return fmt.Sprintf("%T{On:%t}", *c, c.On)
}

func (c *InqExposureSpotLight) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x3a}
	data = append(data, EOL)
	return data
}

func (c *InqExposureSpotLight) HandleReply(data []byte, device *Device) {
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

	device.Inquiry.InqExposureSpotLight = c
}
