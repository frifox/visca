package visca

import (
	"fmt"
	"github.com/frifox/visca/shared"
)

type Mode struct {
}

func (c *Mode) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *Mode) ViscaCommand() []byte {
	data := []byte{shared.CamID, shared.DoInquiry, shared.ToCamera, 0x47}

	data = append(data, shared.EOL)
	return data
}

func (c *Mode) HandleReply(data []byte, device *Device) {

}
