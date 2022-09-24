package visca

import (
	"fmt"
)

type InqZoom struct {
	Z []byte
}

func (c *InqZoom) String() string {
	return fmt.Sprintf("InqZoom{}")
}

func (c *InqZoom) Apply(device *Device) bool {
	return true
}

func (c *InqZoom) ViscaCommand() []byte {
	data := []byte{CamID, doInquiry, toCamera, 0x47}
	data = append(data, EOL)
	return data
}

func (c *InqZoom) HandleReply(data []byte, device *Device) {
	if len(data) < 2 {
		fmt.Printf("[InqZoom.HandleReply] BAD REPLY [% X]\n", data)
		return
	}
	switch data[1] {
	case 0x50:
		c.Z = data[2 : len(data)-1]
		fmt.Printf("[InqZoom.HandleReply] Zoom [%X]\n", c.Z)
	default:
		fmt.Printf("[InqZoom.HandleReply] Unknown [% X]\n", data)
	}
}
