package visca

import (
	"fmt"
)

type InqExposureGainPointPosition struct {
	CmdContext
	Gain int
}

func (c *InqExposureGainPointPosition) String() string {
	return fmt.Sprintf("%T{Gain:%d}", *c, c.Gain)
}

func (c *InqExposureGainPointPosition) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x4c}
	data = append(data, EOL)
	return data
}

func (c *InqExposureGainPointPosition) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p 0p
	if len(data) != 3 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	pp := data[1:3]
	p := sonyInt(pp)
	gain := sonyGain(int(p))
	c.Gain = gain

	device.Inquiry.InqExposureGainPointPosition = c
}
