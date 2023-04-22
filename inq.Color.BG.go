package visca

import (
	"fmt"
)

type InqColorBG struct {
	CmdContext
	Shift int // -99 - 99
}

func (c *InqColorBG) String() string {
	return fmt.Sprintf("%T{Shift:%d}", *c, c.Shift)
}

func (c *InqColorBG) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x7f)
	data = append(data, EOL)
	return data
}

func (c *InqColorBG) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	val := sonyInt(pp)

	// 0x0 - 0xC6 >> -99 - 99
	c.Shift = int(val) - 0x63

	device.Inquiry.InqColorBG = c
}
