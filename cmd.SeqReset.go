package visca

import "fmt"

type SeqReset struct {
	CmdContext
}

func (c *SeqReset) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *SeqReset) ControlCommand() []byte {
	return []byte{0x1}
}

func (c *SeqReset) HandleReply(data []byte, device *Device) {
	c.Finish()

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
