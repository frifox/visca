package visca

import "fmt"

type InqPanTiltPosition struct {
	CmdContext

	X int
	Y int
}

func (c *InqPanTiltPosition) String() string {
	return fmt.Sprintf("%T{X:%d,Y:%d}", *c, c.X, c.Y)
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
	c.Finish()

	// 50 0p 0p 0p 0p 0p 0t 0t 0t 0t
	if len(data) != 10 {
		//fmt.Printf(">> bad reply [% X]\n", data)
		return
	}

	pppp := data[1:6]
	tttt := data[6:10]

	p := sonyInt(pppp)
	t := sonyInt(tttt)

	c.X = int(p)
	c.Y = int(t)

	device.Inquiry.InqPanTiltPosition = c
}
