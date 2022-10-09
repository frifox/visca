package visca

import (
	"fmt"
)

type InqExposureMode struct {
	CmdContext
	Mode string
}

func (c *InqExposureMode) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureMode) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x39}
	data = append(data, EOL)
	return data
}

func (c *InqExposureMode) HandleReply(data []byte, device *Device) {
	c.Finish()

	if len(data) < 4 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	cmd := data[1:len(data)-1]

	switch cmd[1] {
	case 0x0:
		c.Mode = "Auto"
	case 0x3:
		c.Mode = "Manual"
	case 0xa:
		c.Mode = "Shutter Priority"
	case 0xb:
		c.Mode = "Iris Priority"
	case 0xe:
		c.Mode = "Gain Priority"
	}
}
