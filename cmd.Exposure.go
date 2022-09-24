package visca

import "fmt"

type ExposureMode struct {
	Mode string // [auto, manual, shutter, iris]
	mode uint8  // [0x0,  0x3,    0xa,     0xb]
}

func (c *ExposureMode) Apply(d *Device) bool {
	switch c.Mode {
	case "auto":
		c.mode = 0x0
	case "manual":
		c.mode = 0x3
	case "shutter":
		c.mode = 0xA
	case "iris":
		c.mode = 0xb
	default:
		fmt.Printf("[ExposureMode] unuspported mode: %s\n", c.Mode)
		return false
	}

	// no changes?
	if c.mode == d.State.ExposureMode.mode {
		return false
	}

	d.State.ExposureMode.mode = c.mode

	return true
}
func (c *ExposureMode) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x39}
	data = append(data, c.mode)
	data = append(data, EOL)
	return data
}

func (c *ExposureMode) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[ExposureMode.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] {
	case 0x41:
		fmt.Printf("[ExposureMode.HandleReply] ACK\n")
	case 0x51:
		fmt.Printf("[ExposureMode.HandleReply] FIN\n")
	default:
		fmt.Printf("[ExposureMode.HandleReply] Unknown [% X]\n", data)
	}
}
