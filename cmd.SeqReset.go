package visca

import "fmt"
import "context"

type SeqReset struct {
	context.Context
	context.CancelFunc
}

func (c *SeqReset) String() string {
	return fmt.Sprintf("[SeqReset] Reset")
}

func (c *SeqReset) InitContext() {
	c.Context, c.CancelFunc = context.WithCancel(context.Background())
}
func (c *SeqReset) Finish() {
	c.CancelFunc()
}

func (c *SeqReset) ControlCommand() []byte {
	return []byte{0x1}
}

func (c *SeqReset) HandleReply(data []byte, device *Device) {
	if c.Context.Err() == nil {
		c.CancelFunc()
	}

	if len(data) != 1 {
		fmt.Printf("[SeqReset.HandleReply] BAD REPLY [% X]\n", data)
		return
	}

	switch data[0] {
	case 0x1:
		//fmt.Printf("[SeqReset.HandleReply] Ok\n")
	default:
		fmt.Printf("[SeqReset.HandleReply] Unknown [% X]\n", data)
	}
}
