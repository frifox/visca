package visca

import (
	"fmt"
)

type FocusTrigger struct {
	CmdContext
}

func (c *FocusTrigger) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *FocusTrigger) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x38}
	data = append(data, 0x2)
	data = append(data, EOL)
	return data
}

func (c *FocusTrigger) HandleReply(data []byte, device *Device) {
	c.Finish()

	if len(data) != 1 {
		fmt.Printf("[FocusTrigger.HandleReply] BAD REPLY [% X]\n", data)
		return
	}

	p := data[1] & 0xf0
	switch p {
	case 0x40:
		fmt.Printf("[FocusTrigger.HandleReply] ACK\n")
	case 0x50:
		fmt.Printf("[FocusTrigger.HandleReply] FIN\n")
	default:
		fmt.Printf("[FocusTrigger.HandleReply] Unknown [% X]\n", data)
	}
}
