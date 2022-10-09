package visca

import (
	"fmt"
)

type InqDetailCrispening struct {
	CmdContext
	Level int // 0 - 7
}

func (c *InqDetailCrispening) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqDetailCrispening) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x42}
	data = append(data, 0x3)
	data = append(data, EOL)
	return data
}

func (c *InqDetailCrispening) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	c.Level = int(p)

	device.Inquiry.InqDetailCrispening = c
}
