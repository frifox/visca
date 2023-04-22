package visca

import (
	"fmt"
)

type InqGammaMode struct {
	CmdContext
	Mode string
}

func (c *InqGammaMode) String() string {
	return fmt.Sprintf("%T{Mode:%s}", *c, c.Mode)
}

func (c *InqGammaMode) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x5b}
	data = append(data, EOL)
	return data
}

func (c *InqGammaMode) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]
	switch p {
	case 0x0:
		c.Mode = "Standard"
	case 0x1:
		c.Mode = "Straight"
	case 0x2:
		c.Mode = "Pattern"
	case 0x8:
		c.Mode = "Movie"
	case 0x9:
		c.Mode = "Still"
	case 0xA:
		c.Mode = "Cine1"
	case 0xB:
		c.Mode = "Cine2"
	case 0xC:
		c.Mode = "Cine3"
	case 0xD:
		c.Mode = "Cine4"
	case 0xE:
		c.Mode = "ITU709"
	}

	device.Inquiry.InqGammaMode = c
}
