package visca

import (
	"fmt"
)

type FocusTrigger struct{}

func (c *FocusTrigger) String() string {
	return fmt.Sprintf("FocusTrigger{}")
}

func (c *FocusTrigger) Apply(device *Device) bool {
	return true
}

func (c *FocusTrigger) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x18}
	data = append(data, 0x1)
	data = append(data, EOL)
	return data
}

func (c *FocusTrigger) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[FocusTrigger.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] {
	case 0x41:
		fmt.Printf("[FocusTrigger.HandleReply] ACK\n")
	case 0x51:
		fmt.Printf("[FocusTrigger.HandleReply] FIN\n")
	default:
		fmt.Printf("[FocusTrigger.HandleReply] Unknown [% X]\n", data)
	}
}
