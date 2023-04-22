package visca

import (
	"fmt"
)

type InqExposureCompLevel struct {
	CmdContext
	Level int //  00 - 14
}

func (c *InqExposureCompLevel) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqExposureCompLevel) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x4e}
	data = append(data, EOL)
	return data
}

func (c *InqExposureCompLevel) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	p := sonyInt(pp)
	c.Level = int(p)

	device.Inquiry.InqExposureCompLevel = c
}
