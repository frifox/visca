package visca

import (
	"fmt"
)

type InqKneeSetting struct {
	CmdContext
	On bool
}

func (c *InqKneeSetting) String() string {
	return fmt.Sprintf("%T{On:%t}", *c, c.On)
}

func (c *InqKneeSetting) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toConfig, 0x1}
	data = append(data, 0x6d)
	data = append(data, EOL)
	return data
}

func (c *InqKneeSetting) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	switch p {
	case 0x2:
		c.On = true
	case 0x3:
		c.On = false
	}

	device.Inquiry.InqKneeSetting = c
}
