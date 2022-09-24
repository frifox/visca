package visca

type PresetSetCustomSpeed struct {
	ID    uint8
	Speed uint8

	id    uint8
	speed uint8
}

func (c *PresetSetCustomSpeed) Apply(d *Device) bool {
	if c.ID > 0x63 {
		return false
	}
	if c.Speed < 0x1 || c.Speed > 0x19 {
		return false
	}

	c.id = c.ID
	c.speed = c.Speed

	return true
}
func (c *PresetSetCustomSpeed) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toConfig, 0x1}
	data = append(data, 0xb, c.id, c.speed)
	data = append(data, EOL)
	return data
}
