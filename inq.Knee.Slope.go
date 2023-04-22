package visca

import (
	"fmt"
)

type InqKneeSlope struct {
	CmdContext
	Level int // 0 - 15
}

func (c *InqKneeSlope) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqKneeSlope) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x6f)
	data = append(data, EOL)
	return data
}

func (c *InqKneeSlope) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	p := sonyInt(pp)
	c.Level = int(p)

	device.Inquiry.InqKneeSlope = c
}
