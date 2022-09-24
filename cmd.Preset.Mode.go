package visca

import "fmt"

// ONLY ON BRC-X400

// on Recall:
// 	0x0 = Mode1 = Call PTZF + CamSettings to preset
//	0x1 = Mode2 = Call PTZF only

type PresetMode struct {
	ID uint8
	id uint8
}

func (c *PresetMode) String() string {
	return fmt.Sprintf("PresetMode{%X}", c.id)
}

func (c *PresetMode) Apply(d *Device) bool {
	if c.ID > 0x2 {
		return false
	}

	c.id = c.ID

	return true
}
func (c *PresetMode) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toConfig, 0x4}
	data = append(data, 0x3d, c.id)
	data = append(data, EOL)
	return data
}
func (c *PresetMode) HandleReply(data []byte) {
	if len(data) < 2 {
		fmt.Printf("[PresetMode.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] & 0xf0 {
	case 0x40:
		//fmt.Printf("[PresetMode.HandleReply] ACK\n")
	case 0x50:
		fmt.Printf("[PresetMode.HandleReply] FIN\n")
		//default:
		fmt.Printf("[PresetMode.HandleReply] Unknown [% X]\n", data)
	}
}
