package visca

import (
	"fmt"
)

type InqExposureShutterMin struct {
	CmdContext
	Shutter int
}

func (c *InqExposureShutterMin) String() string {
	return fmt.Sprintf("%T{Shutter:%d}", *c, c.Shutter)
}

func (c *InqExposureShutterMin) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x2a}
	data = append(data, 0x1)
	data = append(data, EOL)
	return data
}

func (c *InqExposureShutterMin) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p 0p
	if len(data) != 3 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[1:3]
	val := int(sonyInt(pp))
	c.Shutter = sonyShutter(val, 59.94) // TODO framerate

	device.Inquiry.InqExposureShutterMin = c
}
