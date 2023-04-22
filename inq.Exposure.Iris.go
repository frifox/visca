package visca

import (
	"fmt"
)

type InqExposureIris struct {
	CmdContext
	F float64
}

func (c *InqExposureIris) String() string {
	return fmt.Sprintf("%T{F:%f}", *c, c.F)
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
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	p := sonyInt(pp)
	c.F = sonyIris(int(p))

	device.Inquiry.InqExposureIris = c
}
