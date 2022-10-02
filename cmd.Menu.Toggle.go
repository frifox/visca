package visca

import (
	"context"
	"fmt"
)

type MenuToggle struct {
	context.Context
	context.CancelFunc
}

func (c *MenuToggle) String() string {
	return fmt.Sprintf("MenuToggle{}")
}

func (c *MenuToggle) InitContext() {
	c.Context, c.CancelFunc = context.WithCancel(context.Background())
}
func (c *MenuToggle) Finish() {
	c.CancelFunc()
}

func (c *MenuToggle) ViscaCommand() []byte {
	data := []byte{CamID, doCommand, toMotors, 0x6}
	data = append(data, 0x10)
	data = append(data, EOL)
	return data
}

func (c *MenuToggle) HandleReply(data []byte, device *Device) {
	if c.Err() == nil {
		c.CancelFunc()
	}
}
