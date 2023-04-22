package visca

import "fmt"

type Raw struct {
	CmdContext
	Cmd []byte
}

func (c *Raw) String() string {
	return fmt.Sprintf("%T{% X}", *c, c.Cmd)
}
func (c *Raw) ViscaCommand() []byte {
	data := []byte{CamID}
	data = append(data, c.Cmd...)
	data = append(data, EOL)
	return data
}
func (c *Raw) HandleReply(data []byte, device *Device) {
	c.Finish()
	//fmt.Printf(">> Raw Reply: [% X]\n", data)
}
