package visca

import (
	"fmt"
)

type InqColorOffset struct {
	CmdContext
	Offset int // -7 to +7;
}

func (c *InqColorOffset) String() string {
	return fmt.Sprintf("%T{Offset:%d}", *c, c.Offset)
}

func (c *InqColorOffset) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x2e)
	data = append(data, EOL)
	return data
}

func (c *InqColorOffset) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 00 0p
	if len(data) != 5 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[4]

	// 0x0 - 0xE; 0x7=0
	c.Offset = int(p) - 0x7

	device.Inquiry.InqColorOffset = c
}
