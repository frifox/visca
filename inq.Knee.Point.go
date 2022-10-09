package visca

import (
	"fmt"
)

type InqKneePoint struct {
	CmdContext
	Level int // 0 - 12/0xc
}

func (c *InqKneePoint) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqKneePoint) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x6e)
	data = append(data, EOL)
	return data
}

func (c *InqKneePoint) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	p := sonyInt(pp)
	c.Level = int(p)

	device.Inquiry.InqKneePoint = c
}
