package visca

import (
	"fmt"
)

type InqExposureGain struct {
	CmdContext
	DB int
}

func (c *InqExposureGain) String() string {
	return fmt.Sprintf("%T{DB:%d}", *c, c.DB)
}

func (c *InqExposureGain) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x4c}
	data = append(data, EOL)
	return data
}

func (c *InqExposureGain) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) != 5 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[3:5]
	p := sonyInt(pp)
	c.DB = sonyGain(int(p))

	device.Inquiry.InqExposureGain = c
}
