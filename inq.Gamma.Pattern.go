package visca

import (
	"fmt"
)

type InqGammaPattern struct {
	CmdContext
	Pattern int // 0x1 - 0x200
}

func (c *InqGammaPattern) String() string {
	return fmt.Sprintf("%T{Pattern:%d}", *c, c.Pattern)
}

func (c *InqGammaPattern) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x5b}
	data = append(data, EOL)
	return data
}

func (c *InqGammaPattern) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p 0p 0p
	if len(data) != 4 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[1:4]
	p := sonyInt(pp)
	c.Pattern = int(p)

	device.Inquiry.InqGammaPattern = c
}
