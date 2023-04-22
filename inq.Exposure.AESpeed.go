package visca

import (
	"fmt"
)

type InqExposureAESpeed struct {
	CmdContext
	Speed int
}

func (c *InqExposureAESpeed) String() string {
	return fmt.Sprintf("%T{Speed:%d}", *c, c.Speed)
}

func (c *InqExposureAESpeed) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x5d}
	data = append(data, EOL)
	return data
}

func (c *InqExposureAESpeed) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]
	c.Speed = int(p)

	device.Inquiry.InqExposureAESpeed = c
}
