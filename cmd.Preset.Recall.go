package visca

import (
	"fmt"
)

type PresetRecall struct {
	CmdContext

	ID uint8
	id uint8
}

func (c *PresetRecall) String() string {
	return fmt.Sprintf("%T{id:%d}", *c, c.id)
}

func (c *PresetRecall) Apply(d *Device) bool {
	if c.ID > 0x63 {
		return false
	}

	c.id = c.ID

	return true
}

func (c *PresetRecall) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x3f}
	data = append(data, 0x2, c.id)
	data = append(data, EOL)
	return data
}

func (c *PresetRecall) HandleReply(data []byte, device *Device) {
	c.Finish()

	if len(data) != 1 {
		fmt.Printf("[PresetRecall.HandleReply] BAD REPLY [% X]\n", data)
		return
	}

	p := data[0] & 0xf0

	switch p {
	case 0x40:
		fmt.Printf("[PresetRecall.HandleReply] ACK\n")
	case 0x50:
		fmt.Printf("[PresetRecall.HandleReply] FIN\n")
	default:
		fmt.Printf("[PresetRecall.HandleReply] Unknown [% X]\n", data)
	}
}
