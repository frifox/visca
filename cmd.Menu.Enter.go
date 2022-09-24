package visca

import "fmt"

type MenuEnter struct{}

func (c *MenuEnter) String() string {
	return fmt.Sprintf("MenuEnter{}")
}

func (c *MenuEnter) Apply(device *Device) bool {
	return true
}

func (c *MenuEnter) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toConfig, 0x1}
	data = append(data, 0x2, 0x0, 0x1)
	data = append(data, EOL)
	return data
}

func (c *MenuEnter) HandleReply(data []byte, device *Device) {
	// TODO
}
