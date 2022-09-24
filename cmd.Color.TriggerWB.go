package visca

import "fmt"

type ColorTriggerWB struct{}

func (c *ColorTriggerWB) String() string {
	return fmt.Sprintf("ColorTriggerWB{}")
}

func (c *ColorTriggerWB) Apply(device *Device) bool {
	return true
}

func (c *ColorTriggerWB) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toCamera, 0x10}
	data = append(data, 0x5)
	data = append(data, EOL)
	return data
}

func (c *ColorTriggerWB) HandleReply(data []byte, device *Device) {
	//
}
