package visca

import (
	"fmt"
)

type InqColorLevel struct {
	CmdContext
	Level int // 0 - 15
}

func (c *InqColorLevel) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqColorLevel) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x49}
	data = append(data, EOL)
	return data
}

func (c *InqColorLevel) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 00 0p
	if len(data) != 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[4]
	c.Level = int(p)

	device.Inquiry.InqColorLevel = c
}
