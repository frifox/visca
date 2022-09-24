package visca

import "fmt"

type Power struct {
	On  bool
	cmd uint8
}

func (c *Power) String() string {
	return fmt.Sprintf("Power{%X}", c.cmd)
}

func (c *Power) Apply(device *Device) bool {
	if c.On {
		c.cmd = 0x2
	} else {
		c.cmd = 0x3
	}

	device.State.Power = *c

	return true
}

func (c *Power) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x0}
	data = append(data, c.cmd)
	data = append(data, EOL)
	return data
}

func (c *Power) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[Power.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] & 0xf0 {
	case 0x40:
		fmt.Printf("[Power.HandleReply] ACK\n")
	case 0x50:
		fmt.Printf("[Power.HandleReply] FIN\n")
	default:
		fmt.Printf("[Power.HandleReply] Unknown [% X]\n", data)
	}
}
