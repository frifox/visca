package visca

import (
	"fmt"
)

type InqDetailLevel struct {
	CmdContext
	Level int // 0 - 14/0xe
}

func (c *InqDetailLevel) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqDetailLevel) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x42}
	data = append(data, EOL)
	return data
}

func (c *InqDetailLevel) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	p := sonyInt(pp)
	c.Level = int(p)

	device.Inquiry.InqDetailLevel = c
}
