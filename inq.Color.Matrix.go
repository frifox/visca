package visca

import (
	"fmt"
)

type InqColorMatrix struct {
	CmdContext
	Mode string
}

func (c *InqColorMatrix) String() string {
	return fmt.Sprintf("%T{Mode:%s}", *c, c.Mode)
}

func (c *InqColorMatrix) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x3d)
	data = append(data, EOL)
	return data
}

func (c *InqColorMatrix) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]

	switch p {
	case 0x2:
		c.Mode = "Standard"
	case 0x3:
		c.Mode = "Off"
	case 0x4:
		c.Mode = "High Sat"
	case 0x5:
		c.Mode = "FL Light"
	case 0x6:
		c.Mode = "Movie"
	case 0x7:
		c.Mode = "Still"
	case 0x8:
		c.Mode = "Cinema"
	case 0x9:
		c.Mode = "Pro"
	case 0xA:
		c.Mode = "ITU709"
	case 0xB:
		c.Mode = "B/W"

	}

	device.Inquiry.InqColorMatrix = c
}
