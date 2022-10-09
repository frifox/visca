package visca

import (
	"fmt"
)

type InqZoom struct {
	CmdContext
	Z int
}

func (c *InqZoom) String() string {
	return fmt.Sprintf("%T{Z:%d}", *c, c.Z)
}

func (c *InqZoom) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x47}
	data = append(data, EOL)
	return data
}

func (c *InqZoom) HandleReply(data []byte, device *Device) {
	c.Finish()

	// 50 0z 0z 0z 0z
	if len(data) != 5 {
		fmt.Printf(">> bad reply [% X]\n", data)
		return
	}

	zzzz := data[1:5]
	z := sonyInt(zzzz)
	c.Z = int(z)

	/* TODO
	0000 ×1
	1800 ×2
	2340 ×3
	2A40 ×4
	2F00 ×5
	3300 ×6
	3600 ×7
	3880 ×8
	3AC0 ×9
	3CC0 ×10
	3E80 ×11
	4000 ×12
	5580 ×18 (While using Clear Image Zoom)
	6000 ×24 (While using Clear Image Zoom)*1
	*/

	device.Inquiry.InqZoom = c
}
