package visca

import (
	"fmt"
)

type InqColorBlueGain struct {
	CmdContext
	Gain int // -128 - 128(?)
}

func (c *InqColorBlueGain) String() string {
	return fmt.Sprintf("%T{Gain:%d}", *c, c.Gain)
}

func (c *InqColorBlueGain) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x44}
	data = append(data, EOL)
	return data
}

func (c *InqColorBlueGain) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	val := int(sonyInt(pp))

	// 0x0 - 0xff; 0x80=0
	c.Gain = val - 0x80

	device.Inquiry.InqColorBlueGain = c
}
