package visca

import (
	"fmt"
)

type InqExposureVisibilityEnhancer struct {
	CmdContext
	On bool
}

func (c *InqExposureVisibilityEnhancer) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureVisibilityEnhancer) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x3d}
	data = append(data, EOL)
	return data
}

func (c *InqExposureVisibilityEnhancer) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]

	switch p {
	case 0x6:
		c.On = true
	case 0x3:
		c.On = false
	default:
		fmt.Printf(">> unknown %T value [%X]\n", *c, p)
	}

	device.Inquiry.InqExposureVisibilityEnhancer = c
}
