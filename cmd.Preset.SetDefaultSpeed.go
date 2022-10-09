package visca

import "fmt"

type PresetSetDefaultSpeed struct {
	CmdContext

	Speed uint8
	speed uint8
}

func (c *PresetSetDefaultSpeed) String() string {
	return fmt.Sprintf("%T{speed:%d}", *c, c.Speed)
}

func (c *PresetSetDefaultSpeed) Apply(d *Device) bool {
	if c.Speed < 0x1 || c.Speed > 0x19 {
		return false
	}

	c.speed = c.Speed

	return true
}
func (c *PresetSetDefaultSpeed) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toConfig, 0x4}

	// TODO 0xAB => []byte{0x0A, 0x0B}
	data = append(data, 0x1c, c.speed&0xf0, c.speed&0x0f)

	data = append(data, EOL)
	return data
}
