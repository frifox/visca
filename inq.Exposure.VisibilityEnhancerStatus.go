package visca

import (
	"fmt"
)

type InqExposureVisibilityEnhancerInfo struct {
	CmdContext
	Level          int // 0=dark - 6=bright
	BrightnessComp int // 0=very dark, 1=dark, 2=standard, 3=bright
	CompLevel      int // 0=low, 1=mid, 2=high
}

func (c *InqExposureVisibilityEnhancerInfo) String() string {
	return fmt.Sprintf("%T{Level:%d,BrightnessComp:%d,CompLevel:%d}", *c, c.Level, c.BrightnessComp, c.CompLevel)
}

func (c *InqExposureVisibilityEnhancerInfo) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x2d}
	data = append(data, EOL)
	return data
}

func (c *InqExposureVisibilityEnhancerInfo) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 0p 0q 0r 00 00 00 00
	if len(data) != 9 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[2]
	q := data[3]
	r := data[4]

	c.Level = int(p)
	c.BrightnessComp = int(q)
	c.CompLevel = int(r)

	device.Inquiry.InqExposureVisibilityEnhancerInfo = c
}
