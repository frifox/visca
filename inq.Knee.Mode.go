package visca

import (
	"fmt"
)

type InqKneeMode struct {
	CmdContext
	Mode string
}

func (c *InqKneeMode) String() string {
	return fmt.Sprintf("%T{Mode:%s}", *c, c.Mode)
}

func (c *InqKneeMode) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x42}
	data = append(data, 0x1)
	data = append(data, EOL)
	return data
}

func (c *InqKneeMode) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	switch p {
	case 0x0:
		c.Mode = "Auto"
	case 0x1:
		c.Mode = "Manual"
	}

	device.Inquiry.InqKneeMode = c
}
