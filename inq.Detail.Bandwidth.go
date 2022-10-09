package visca

import (
	"fmt"
)

type InqDetailBandwidth struct {
	CmdContext
	Level int // 0 - 4
}

func (c *InqDetailBandwidth) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqDetailBandwidth) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x42}
	data = append(data, 0x2)
	data = append(data, EOL)
	return data
}

func (c *InqDetailBandwidth) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> bad reply [%X]\n", data)
		return
	}

	p := data[1]

	c.Level = int(p)

	device.Inquiry.InqDetailBandwidth = c
}
