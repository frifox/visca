package visca

import (
	"fmt"
)

type InqExposureIris struct {
	CmdContext
	F float64
}

func (c *InqExposureIris) String() string {
	return fmt.Sprintf("%T{}", *c)
}

func (c *InqExposureIris) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x4b}
	data = append(data, EOL)
	return data
}

func (c *InqExposureIris) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 00 00 0p 0p
	if len(data) < 5 {
		fmt.Printf(">> BAD REPLY\n")
		return
	}

	/* pp:
	15 F2.8 (Open)
	14 F3.1
	13 F3.4
	12 F3.7
	11 F4.0
	10 F4.4
	0F F4.8
	0E F5.2
	0D F5.6
	0C F6.2
	0B F6.8
	0A F7.3
	09 F8.0
	08 F8.7
	07 F9.6
	06 F10
	05 F11
	*/

	val := ppToP(data[3:5])
	switch val[0] {
	case 0x15:
		c.F = 2.8
	case 0x14:
		c.F = 3.1
	case 0x13:
		c.F = 3.4
	case 0x12:
		c.F = 3.7
	case 0x11:
		c.F = 4.0
	case 0x10:
		c.F = 4.4
	case 0x0f:
		c.F = 4.8
	case 0x0e:
		c.F = 5.2
	case 0x0d:
		c.F = 5.6
	case 0x0c:
		c.F = 6.2
	case 0x0b:
		c.F = 6.8
	case 0x0a:
		c.F = 7.3
	case 0x09:
		c.F = 8.0
	case 0x08:
		c.F = 8.7
	case 0x07:
		c.F = 9.6
	case 0x06:
		c.F = 10.0
	case 0x05:
		c.F = 11.0
	}
	fmt.Printf(">> val = % X\n", val)
}
