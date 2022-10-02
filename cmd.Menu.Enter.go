package visca

import (
	"context"
	"fmt"
)

type MenuEnter struct {
	context.Context
	context.CancelFunc
}

func (c *MenuEnter) String() string {
	return fmt.Sprintf("MenuEnter{}")
}

func (c *MenuEnter) InitContext() {
	c.Context, c.CancelFunc = context.WithCancel(context.Background())
}
func (c *MenuEnter) Finish() {
	c.CancelFunc()
}

func (c *MenuEnter) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toConfig, 0x1}
	data = append(data, 0x2, 0x0, 0x1)
	data = append(data, EOL)
	return data
}

func (c *MenuEnter) HandleReply(data []byte, device *Device) {
	if c.Err() == nil {
		c.CancelFunc()
	}
}
