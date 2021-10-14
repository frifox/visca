package visca

type OSDToggle struct {
	alwaysApply
}

func (a *OSDToggle) bytes() []byte {
	return []byte{0x1, 0x4, 0x3f, 0x2, 0x5f}
}

type OSDEnter struct {
	alwaysApply
}

func (a *OSDEnter) bytes() []byte {
	return []byte{0x1, 0x6, 0x6, 0x5}
}

type OSDReturn struct {
	alwaysApply
}

func (a *OSDReturn) bytes() []byte {
	return []byte{0x1, 0x6, 0x6, 0x4}
}

type OSDLeft struct {
	alwaysApply
}

func (a *OSDLeft) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x1, 0x3}
}

type OSDRight struct {
	alwaysApply
}

func (a *OSDRight) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x2, 0x3}
}

type OSDUp struct {
	alwaysApply
}

func (a *OSDUp) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x3, 0x1}
}

type OSDDown struct {
	alwaysApply
}

func (a *OSDDown) bytes() []byte {
	return []byte{0x1, 0x6, 0x1, 0xe, 0xe, 0x3, 0x2}
}
