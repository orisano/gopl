package main

func IsBlank(b []byte) (bool, int) {
	switch b[0] {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xa0:
		return true, 0
	case 0xe1:
		if len(b) > 2 && b[1] == 0x9a && b[2] == 0x80 {
			return true, 2
		}
	case 0xe2:
		if len(b) <= 2 {
			return false, 0
		}
		x := b[2]

		switch b[1] {
		case 0x80:
			switch {
			case 0x80 <= x && x <= 0x8a:
				return true, 2
			case 0xa8 <= x && x <= 0xa9:
				return true, 2
			case x == 0xaf:
				return true, 2
			}
		case 0x81:
			if x == 0x9f {
				return true, 2
			}
		}
	case 0xe3:
		if len(b) <= 2 {
			return false, 0
		}
		if b[1] == 0x80 && b[2] == 0x80 {
			return true, 2
		}
	}
	return false, 0
}

func CompressBlank(b []byte) []byte {
	r := b[:0]

	skip := false
	for i := 0; i < len(b); i++ {
		space, x := IsBlank(b[i:])
		if space {
			if !skip {
				r = append(r, ' ')
				skip = true
			}
			i += x
		} else {
			skip = false
			r = append(r, b[i])
		}
	}
	return r
}
