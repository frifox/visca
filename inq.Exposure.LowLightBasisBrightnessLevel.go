package visca

import (
	"fmt"
)

type InqExposureLowLightBasisBrightnessInfo struct {
	CmdContext
	Level int // 4 - A
}

func (c *InqExposureLowLightBasisBrightnessInfo) String() string {
	return fmt.Sprintf("%T{Level:%d}", *c, c.Level)
}

func (c *InqExposureLowLightBasisBrightnessInfo) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera2, 0x39}
	data = append(data, EOL)
	return data
}

func (c *InqExposureLowLightBasisBrightnessInfo) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]
	c.Level = int(p)

	device.Inquiry.InqExposureLowLightBasisBrightnessInfo = c
}
