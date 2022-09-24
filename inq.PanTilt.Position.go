package visca

import "fmt"

type InqPanTiltPosition struct {
	X []byte
	Y []byte
}

func (c *InqPanTiltPosition) String() string {
	return fmt.Sprintf("InqPanTiltPosition")
}

func (c *InqPanTiltPosition) Apply(device *Device) bool {
	return true
}

func (c *InqPanTiltPosition) ViscaInquiry() []byte {
	data := []byte{CamID, doInquiry, toMotors, 0x12}
	data = append(data, EOL)
	return data
}
func (c *InqPanTiltPosition) HandleReply(data []byte, device *Device) {
	if len(data) != 11 {
		fmt.Printf("[InqPanTiltPosition.HandleReply] len() != (11 [% X]\n", data)
		return
	}

	//fmt.Printf("[InqPanTiltPosition.HandleReply] Handling [% X]\n", data)

	xy := data[2 : len(data)-1]
	//fmt.Printf("[InqPanTiltPosition.HandleReply] Handling [% X]\n", xy)

	c.X = xy[:4]
	c.Y = xy[4:]

	//fmt.Printf("[InqPanTiltPosition.HandleReply] X = [% X] y = [% X]\n", c.X, c.Y)
}
