package visca

import (
	"fmt"
)

type InqDetailLimit struct {
	CmdContext
	Level int // 0 - 7
}

func (c *InqDetailLimit) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqDetailLimit) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x42}
	data = append(data, 0x6)
	data = append(data, EOL)
	return data
}

func (c *InqDetailLimit) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	c.Level = int(p)

	device.Inquiry.InqDetailLimit = c
}
