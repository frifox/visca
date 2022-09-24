package visca

import "fmt"

type PresetSet struct {
	ID uint8
	id uint8
}

func (c *PresetSet) String() string {
	return fmt.Sprintf("PresetSet{ID:%X}", c.id)
}

func (c *PresetSet) Apply(d *Device) bool {
	if c.ID > 0x63 {
		return false
	}

	c.id = c.ID

	return true
}

func (c *PresetSet) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x3f}
	data = append(data, 0x1, c.id)
	data = append(data, EOL)
	return data
}
func (c *PresetSet) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[PresetMode.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] & 0xf0 {
	case 0x40:
		fmt.Printf("[PresetMode.HandleReply] ACK\n")
	case 0x50:
		fmt.Printf("[PresetMode.HandleReply] FIN\n")
	default:
		fmt.Printf("[PresetMode.HandleReply] Unknown [% X]\n", data)
	}
}
