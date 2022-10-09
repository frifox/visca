package visca

import (
	"fmt"
)

type InqExposureLowLightBasisBrightness struct {
	CmdContext
	On bool
}

func (c *InqExposureLowLightBasisBrightness) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureLowLightBasisBrightness) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x39}
	data = append(data, EOL)
	return data
}

func (c *InqExposureLowLightBasisBrightness) HandleReply(data []byte, device *Device) {
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

	device.Inquiry.InqExposureLowLightBasisBrightness = c
}
