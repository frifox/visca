package visca

import (
	"fmt"
)

type InqExposureIris struct {
	CmdContext
	F float64
}

func (c *InqExposureIris) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureIris) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x4b}
	data = append(data, EOL)
	return data
}

func (c *InqExposureIris) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	val := int(sonyInt(pp))

	c.F = sonyIris(val)

	device.Inquiry.ExposureIris = c
}
