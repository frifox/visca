package visca

import (
	"fmt"
)

type InqExposureNDFilter struct {
	CmdContext
	Level int // 0=off, 1=1/4, 2=1/16, 3=1/64
}

func (c *InqExposureNDFilter) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqExposureNDFilter) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x53)
	data = append(data, EOL)
	return data
}

func (c *InqExposureNDFilter) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]
	c.Level = int(p)

	device.Inquiry.InqExposureNDFilter = c
}
