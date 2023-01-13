package easychars

// Check whether content is valid under GBK rule, referce: https://zh.wikipedia.org/wiki/GBK
func isValidGBK(content []byte) bool {
	nByte := 1 // the number of bytes that current character use, max 2 bytes in GBK
	for _, b := range content {
		switch nByte {
		case 1:
			if b <= 0x7F { // character is ascii
				continue
			}
			if b >= 0x81 && b <= 0xFE { // may be a GBK encoded character, depending on second byte
				nByte = 2
			} else { // not a valid GBK encoded character
				return false
			}
		case 2:
			nByte = 1
			if b < 0x40 || b > 0xFE || b == 0x7F { // not a valid GBK encoded character under these circumustance
				return false
			}
		}
	}
	return nByte == 1
}

// Check whether content is valid under GB18030 rule, referce: https://zh.wikipedia.org/wiki/GB_18030
func isValidGB18030(content []byte) bool {
	nByte := 1 // the number of bytes that current character use, max 4 bytes in GB18030
	for _, b := range content {
		switch nByte {
		case 1:
			if b <= 0x7F { // character is ascii
				continue
			}
			if b >= 0x81 && b <= 0xFE { // may be a GB18030 encoded character, depending on second byte
				nByte = 2
			} else { // not a valid GBK encoded character
				return false
			}
		case 2:
			if b >= 0x40 && b <= 0xFE && b != 0x7F { // is a valid 2 byte GB18030(GBK) encoded character
				nByte = 1
				continue
			}
			if b >= 0x30 && b <= 0x39 { // may be a valid 4 byte GB18030 encoded character, depending on the third and fourth byte
				nByte = 3
				continue
			} else {
				return false
			}
		case 3:
			if b >= 0x81 && b <= 0xFE { // may be a valid 4 byte GB18030 encoded character, depending on the fourth byte
				nByte = 4
				continue
			} else {
				return false
			}
		case 4:
			if b >= 0x30 && b <= 0x39 { // a valid 4 byte GB18030 encoded character
				nByte = 1
				continue
			} else {
				return false
			}
		}
	}
	return nByte == 1
}

// Check whether content is valid under Big5 rule, referce: https://zh.wikipedia.org/wiki/Big5
func isValidBig5(content []byte) bool {
	nByte := 1 // Big5 use ascii && 2 byte encoded character
	for _, b := range content {
		switch nByte {
		case 1:
			if b <= 0x7F {
				continue
			}
			if b >= 0x81 && b <= 0xFE {
				nByte = 2
			} else {
				return false
			}
		case 2:
			nByte = 1
			if !(b >= 0x40 && b <= 0x7E || b >= 0xA1 && b <= 0xFE) {
				return false
			}
		}
	}
	return nByte == 1
}
