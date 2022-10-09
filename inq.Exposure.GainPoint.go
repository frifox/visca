package visca

import (
	"fmt"
)

type InqExposureGainLimit struct {
	CmdContext
	DB int
}

func (c *InqExposureGainLimit) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureGainLimit) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x2c}
	data = append(data, EOL)
	return data
}

func (c *InqExposureGainLimit) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	c.DB = int(data[1])
}
