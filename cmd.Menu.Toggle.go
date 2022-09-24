package visca

import "fmt"

type MenuToggle struct{}

func (c *MenuToggle) String() string {
	return fmt.Sprintf("MenuToggle{}")
}

func (c *MenuToggle) Apply(device *Device) bool {
	return true
}

func (c *MenuToggle) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toMotors, 0x6}
	data = append(data, 0x10)
	data = append(data, EOL)
	return data
}

func (c *MenuToggle) HandleReply(data []byte, device *Device) {
	// TODO
}
