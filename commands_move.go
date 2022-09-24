package visca

//type Move struct {
//	X      float64
//	Y      float64
//	StepsX float64
//	StepsY float64
//
//	xStep  int8
//	yStep  int8
//	device *Device
//
//	runtime [3]time.Time // [write, fin, fin]
//}
//
//func (c *Move) apply() bool {
//	xStep := int8(math.Ceil(c.StepsX * math.Abs(c.X)))
//	yStep := int8(math.Ceil(c.StepsY * math.Abs(c.Y)))
//
//	if c.X < 0 {
//		xStep = -xStep
//	}
//	if c.Y < 0 {
//		yStep = -yStep
//	}
//
//	// no changes?
//	if c.xStep == xStep && c.yStep == yStep {
//		return false
//	}
//
//	// save new state
//	c.xStep = xStep
//	c.yStep = yStep
//
//	return true
//}
//
//func (c *Move) bytes() []byte {
//	packet := bytes.Buffer{}
//
//	// header
//	packet.Write([]byte{0x1, 0x6, 0x1})
//
//	// X speed
//	xStep := c.xStep
//	if xStep < 0 {
//		xStep = -xStep
//	}
//	packet.WriteByte(byte(xStep))
//
//	// Y speed
//	yStep := c.yStep
//	if yStep < 0 {
//		yStep = -yStep
//	}
//	packet.WriteByte(byte(yStep))
//
//	// X direction
//	switch true {
//	case c.X > 0:
//		packet.WriteByte(0x2) // right
//	case c.X < 0:
//		packet.WriteByte(0x1) // left
//	default:
//		packet.WriteByte(0x3) // none
//	}
//
//	// Y direction
//	switch true {
//	case c.Y > 0:
//		packet.WriteByte(0x1) // up
//	case c.Y < 0:
//		packet.WriteByte(0x2) // down
//	default:
//		packet.WriteByte(0x3) // none
//	}
//
//	return packet.Bytes()
//}
//
//type RampCurve struct {
//	Value uint8
//}
//
//func (c RampCurve) apply() bool {
//	return true
//}
//func (c RampCurve) bytes() []byte {
//	return []byte{0x1, 0x6, 0x31, c.Value}
//}
