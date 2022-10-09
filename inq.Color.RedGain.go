package visca

import (
	"fmt"
)

type InqColorWhiteBalanceMode struct {
	CmdContext
	Mode string
}

func (c *InqColorWhiteBalanceMode) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqColorWhiteBalanceMode) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x35}
	data = append(data, EOL)
	return data
}

func (c *InqColorWhiteBalanceMode) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]

	switch p {
	case 0x0:
		c.Mode = "Auto1"
	case 0x1:
		c.Mode = "Indoor"
	case 0x2:
		c.Mode = "Outdoor"
	case 0x3:
		c.Mode = "OnePush WB"
	case 0x4:
		c.Mode = "Auto2"
	case 0x5:
		c.Mode = "Manual"
	}

	device.Inquiry.InqColorWhiteBalanceMode = c
}
