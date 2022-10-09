package visca

import (
	"fmt"
)

type InqExposureGainPoint struct {
	CmdContext
	On bool
}

func (c *InqExposureGainPoint) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureGainPoint) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toUnknown, 0xc}
	data = append(data, EOL)
	return data
}

func (c *InqExposureGainPoint) HandleReply(data []byte, device *Device) {
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
		fmt.Printf(">> unexpected gainpoint [%X]\n", p)
	}

	device.Inquiry.ExposureGainPoint = c
}
