package visca

import (
	"fmt"
)

type InqExposureShutter struct {
	CmdContext
	Shutter int
}

func (c *InqExposureShutter) String() string {
	return fmt.Sprintf("%T{Shutter:%d}", *c, c.Shutter)
}

func (c *InqExposureShutter) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x4a}
	data = append(data, EOL)
	return data
}

func (c *InqExposureShutter) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	val := int(sonyInt(pp))
	c.Shutter = sonyShutter(val, 59.94) // TODO framerate

	device.Inquiry.InqExposureShutter = c
}
