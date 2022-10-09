package visca

import (
	"fmt"
)

type InqExposureGain struct {
	CmdContext
	DB int
}

func (c *InqExposureGain) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureGain) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x4c}
	data = append(data, EOL)
	return data
}

func (c *InqExposureGain) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) < 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]

	switch sonyInt(pp) {
	case 0x0C:
		c.DB = 33
	case 0x0B:
		c.DB = 30
	case 0x0A:
		c.DB = 27
	case 0x09:
		c.DB = 24
	case 0x08:
		c.DB = 21
	case 0x07:
		c.DB = 18
	case 0x06:
		c.DB = 15
	case 0x05:
		c.DB = 12
	case 0x04:
		c.DB = 9
	case 0x03:
		c.DB = 6
	case 0x02:
		c.DB = 3
	case 0x01:
		c.DB = 0
	case 0x00:
		c.DB = -3
	}
}
