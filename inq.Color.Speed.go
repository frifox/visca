package visca

import (
	"fmt"
)

type InqColorSpeed struct {
	CmdContext
	Speed int // 1=slow - 5=fast
}

func (c *InqColorSpeed) String() string {
	return fmt.Sprintf("%T{Speed:%d}", *c, c.Speed)
}

func (c *InqColorSpeed) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x56}
	data = append(data, EOL)
	return data
}

func (c *InqColorSpeed) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0p
	if len(data) != 2 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[1]
	c.Speed = int(p)

	device.Inquiry.InqColorSpeed = c
}
