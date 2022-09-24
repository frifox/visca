package visca

import "fmt"

type HighResolution struct {
	On  bool
	cmd uint8
}

func (c *HighResolution) String() string {
	return fmt.Sprintf("HighResolution{%X}", c.cmd)
}

func (c *HighResolution) Apply(device *Device) bool {
	if c.On {
		c.cmd = 0x2
	} else {
		c.cmd = 0x3
	}

	return true
}

func (c *HighResolution) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x52}
	data = append(data, c.cmd)
	data = append(data, EOL)
	return data
}

func (c *HighResolution) HandleReply(data []byte, device *Device) {
	// TODO
}
