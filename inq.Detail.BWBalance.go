package visca

import (
	"fmt"
)

type InqDetailBWBalance struct {
	CmdContext
	Level int // 0 - 4
}

func (c *InqDetailBWBalance) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqDetailBWBalance) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x42}
	data = append(data, 0x5)
	data = append(data, EOL)
	return data
}

func (c *InqDetailBWBalance) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	c.Level = int(p)

	device.Inquiry.InqDetailBWBalance = c
}
