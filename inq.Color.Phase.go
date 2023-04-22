package visca

import (
	"fmt"
)

type InqColorPhase struct {
	CmdContext
	Phase int // -7 to +7
}

func (c *InqColorPhase) String() string {
	return fmt.Sprintf("%T{Phase:%d}", *c, c.Phase)
}

func (c *InqColorPhase) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x49}
	data = append(data, EOL)
	return data
}

func (c *InqColorPhase) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 00 0p
	if len(data) != 5 {
		//fmt.Printf(">> BAD REPLY\n")
		return
	}

	p := data[4]
	c.Phase = int(p) - 0x7

	device.Inquiry.InqColorPhase = c
}
