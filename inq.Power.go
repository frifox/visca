package visca

import (
	"context"
	"fmt"
)

type InqPower struct {
	On bool

	context.Context
	context.CancelFunc
}

func (c *InqPower) String() string {
	return fmt.Sprintf("InqPower{}")
}

func (c *InqPower) InitContext() {
	c.Context, c.CancelFunc = context.WithCancel(context.Background())
}
func (c *InqPower) Finish() {
	c.CancelFunc()
}

func (c *InqPower) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x0}
	data = append(data, EOL)
	return data
}

func (c *InqPower) HandleReply(data []byte, device *Device) {
	if c.Err() == nil {
		c.CancelFunc()
	}

	if len(data) != 4 {
		fmt.Printf("[InqPower.HandleReply] BAD REPLY [% X]\n", data)
		return
	}

	// [y0 50 0p FF] p: 2=on, 3=standby
	switch data[2] {
	case 0x2:
		device.State.Power.On = true
	case 0x3:
		device.State.Power.On = false
	default:
		fmt.Printf("[InqPower.HandleReply] Unknown [% X]\n", data)
	}
}
