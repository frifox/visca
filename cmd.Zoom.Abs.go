package visca

import (
	"fmt"
)

type ZoomAbs struct {
	Z []byte
	z []byte
}

func (c *ZoomAbs) String() string {
	return fmt.Sprintf("ZoomAbs{z:%d}", c.z)
}

func (c *ZoomAbs) Apply(device *Device) bool {
	c.z = c.Z
	return true
}

func (c *ZoomAbs) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x47}
	data = append(data, c.z...)
	data = append(data, EOL)
	return data
}

func (c *ZoomAbs) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[ZoomAbs.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] & 0xf0 {
	case 0x40:
		//fmt.Printf("[ZoomAbs.HandleReply] ACK\n")
	case 0x50:
		//fmt.Printf("[ZoomAbs.HandleReply] FIN\n")
	default:
		fmt.Printf("[ZoomAbs.HandleReply] Unknown [% X]\n", data)
	}
}
