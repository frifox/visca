package visca

import "fmt"

func sonyInt(in []byte) (out int64) {
	// int64 = values up to 8 sony bytes

	// 0p 0p 0p 0p => pppp
	for i := len(in) - 1; i >= 0; i-- {
		index := len(in) - 1 - i
		val := in[i]
		out += int64(val) << (4 * index)
	}

	return out
}

func sonyGain(val int) (gain int) {
	switch val {
	case 0x0C:
		return 33
	case 0x0B:
		return 30
	case 0x0A:
		return 27
	case 0x09:
		return 24
	case 0x08:
		return 21
	case 0x07:
		return 18
	case 0x06:
		return 15
	case 0x05:
		return 12
	case 0x04:
		return 9
	case 0x03:
		return 6
	case 0x02:
		return 3
	case 0x01:
		return 0
	case 0x00:
		return -3
	default:
		fmt.Printf(">> unknown gain value [%X]\n", val)
	}

	return 0
}

func sonyIris(val int) (f float64) {
	switch val {
	case 0x15:
		return 2.8
	case 0x14:
		return 3.1
	case 0x13:
		return 3.4
	case 0x12:
		return 3.7
	case 0x11:
		return 4.0
	case 0x10:
		return 4.4
	case 0x0f:
		return 4.8
	case 0x0e:
		return 5.2
	case 0x0d:
		return 5.6
	case 0x0c:
		return 6.2
	case 0x0b:
		return 6.8
	case 0x0a:
		return 7.3
	case 0x09:
		return 8.0
	case 0x08:
		return 8.7
	case 0x07:
		return 9.6
	case 0x06:
		return 10.0
	case 0x05:
		return 11.0
	default:
		fmt.Printf(">> unknown iris value [%X]\n", val)
	}

	return 0
}

func sonyShutter(val int, framerate float64) (oneOverX int) {
	switch framerate {
	case 59.94:
		switch val {
		case 0x15:
			return 10000
		case 0x14:
			return 6000
		case 0x13:
			return 4000
		case 0x12:
			return 3000
		case 0x11:
			return 2000
		case 0x10:
			return 1500
		case 0xf:
			return 1000
		case 0xe:
			return 725
		case 0xd:
			return 500
		case 0xc:
			return 350
		case 0xb:
			return 250
		case 0xa:
			return 180
		case 0x9:
			return 125
		case 0x8:
			return 100
		case 0x7:
			return 90
		case 0x6:
			return 60
		case 0x5:
			return 50
		case 0x4:
			return 30
		case 0x3:
			return 15
		case 0x2:
			return 8
		default:
			fmt.Printf(">> unhandled shutter val [%X]\n", val)
		}
	default:
		fmt.Printf(">> unhandled framerate %f\n", framerate)
	}

	return 0
}
