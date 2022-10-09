package visca

import (
	"fmt"
)

type InqColorChromaSuppress struct {
	CmdContext
	Level int // 0:off, 1:weak - 3:strong
}

func (c *InqColorChromaSuppress) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqColorChromaSuppress) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x5f}
	data = append(data, EOL)
	return data
}

func (c *InqColorChromaSuppress) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]
	c.Level = int(p)

	device.Inquiry.InqColorChromaSuppress = c
}
