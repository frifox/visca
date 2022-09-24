package visca

import "fmt"

type PresetRecall struct {
	ID uint8
	id uint8
}

func (c *PresetRecall) String() string {
	return fmt.Sprintf("PresetRecall{id:%d}", c.id)
}

func (c *PresetRecall) Apply(d *Device) bool {
	if c.ID > 0x63 {
		return false
	}

	c.id = c.ID

	return true
}
func (c *PresetRecall) Send(d *Device) bool {
	return true
}
func (c *PresetRecall) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x3f}
	data = append(data, 0x2, c.id)
	data = append(data, EOL)
	return data
}

func (c *PresetRecall) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[PresetRecall.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] {
	case 0x41:
		fmt.Printf("[PresetRecall.HandleReply] ACK\n")
	case 0x51:
		fmt.Printf("[PresetRecall.HandleReply] FIN\n")
	default:
		fmt.Printf("[PresetRecall.HandleReply] Unknown [% X]\n", data)
	}
}
