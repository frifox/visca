package visca

import (
	"fmt"
)

type InqDetailHVBalance struct {
	CmdContext
	Level int // 5 - 9
}

func (c *InqDetailHVBalance) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqDetailHVBalance) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x42}
	data = append(data, 0x4)
	data = append(data, EOL)
	return data
}

func (c *InqDetailHVBalance) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	c.Level = int(p)

	device.Inquiry.InqDetailHVBalance = c
}
