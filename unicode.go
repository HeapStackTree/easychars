package easychars

import (
	"unicode/utf8"
)

// Check whether content is valid under UTF-8 rule
func IsValidUTF8(content []byte) bool {
	return utf8.Valid(content)
}

// Check whether content is valid under UTF-16 rule, reference: https://zh.wikipedia.org/wiki/UTF-16
//
// return: isUTF16 bool, BE bool
//
// BE: true if content is valid under UTF-16 BE rule, false if not
// LE: true if content is valid under UTF-16 LE rule, false if not
func isValidUTF16(content []byte) (isUTF16 bool, BE bool, LE bool) {
	BE = isValidUTF16BE(content)
	LE = isValidUTF16LE(content)
	isUTF16 = BE || LE
	return
}

// Check whether content is valid under UTF-16BE rule, reference: https://zh.wikipedia.org/wiki/UTF-16
func isValidUTF16BE(content []byte) bool {
	if len(content) == 0 {
		return false
	}
	if len(content)&0x1 != 0 {
		return false
	}
	// UTF-16 BE BOM: FE FF
	BOM := uint16(content[0])<<8 ^ uint16(content[1])
	if BOM == 0xFEFF {
		return true
	}

	// If no BOM in content, assume it's encoded by UTF-16-BE and check it
	//
	// UTF-16 is valid in the range 0x0000 - 0xFFFF excluding 0xD800 - 0xFFFF
	// with an exception for surrogate pairs, which must be in the range
	// 0xD800-0xDBFF followed by 0xDC00-0xDFFF
	//
	// https://en.wikipedia.org/wiki/UTF-16
	is_surrogate_pairs := false
	for i := 0; i < len(content); i += 2 {
		c := uint16(content[i])<<8 ^ uint16(content[i+1])
		switch is_surrogate_pairs {
		case true:
			is_surrogate_pairs = false
			if c < 0xDC00 || c > 0xDFFF {
				return false
			}
		case false:
			if c >= 0xD800 && c <= 0xFFFF {
				is_surrogate_pairs = true
			}
		}
	}
	return !is_surrogate_pairs
}

// Check whether content is valid under UTF-16LE rule, reference: https://zh.wikipedia.org/wiki/UTF-16
//
// This function assume content is UTF-16LE and then valid it.
func isValidUTF16LE(content []byte) bool {
	if len(content) == 0 {
		return false
	}
	if len(content)&0x1 != 0 {
		return false
	}

	// UTF-16 LE BOM: FE FF
	BOM := uint16(content[1])<<8 ^ uint16(content[0])
	if BOM == 0xFEFF {
		return true
	}

	// If no BOM in content, assume it's encoded by UTF-16-BE and check it
	//
	// UTF-16 is valid in the range 0x0000 - 0xFFFF excluding 0xD800 - 0xFFFF
	// with an exception for surrogate pairs, which must be in the range
	// 0xD800-0xDBFF followed by 0xDC00-0xDFFF
	//
	// https://en.wikipedia.org/wiki/UTF-16
	is_surrogate_pairs := false
	for i := 0; i < len(content); i += 2 {
		c := uint16(content[i+1])<<8 ^ uint16(content[i])
		switch is_surrogate_pairs {
		case true:
			is_surrogate_pairs = false
			if c < 0xDC00 || c > 0xDFFF {
				return false
			}
		case false:
			if c >= 0xD800 && c <= 0xFFFF {
				is_surrogate_pairs = true
			}
		}
	}
	return !is_surrogate_pairs
}

// Convert unicode to utf-8 encode []byte
func unicodeRuneToUtf8(unicode rune) (utf8codes []byte) {
	if unicode <= 0x7F {
		utf8codes = append(utf8codes, byte(unicode))
		return
	} else if unicode <= 0x7FF {
		var c1 byte = 192
		var c2 byte = 128
		for k := 0; k < 11; k++ {
			if k < 6 {
				c2 |= byte((unicode % 64) & (1 << k))
			} else {
				c1 |= byte((unicode >> 6) & (1 << (k - 6)))
			}
		}
		utf8codes = append(utf8codes, c1)
		utf8codes = append(utf8codes, c2)
		return
	} else if unicode <= 0xFFFF {
		var c1 byte = 224
		var c2 byte = 128
		var c3 byte = 128
		for k := 0; k < 16; k++ {
			if k < 6 {
				c3 |= byte((unicode % 64) & (1 << k))
			} else if k < 12 {
				c2 |= byte((unicode >> 6) & (1 << (k - 6)))
			} else {
				c1 |= byte((unicode >> 12) & (1 << (k - 12)))
			}
		}
		utf8codes = append(utf8codes, c1)
		utf8codes = append(utf8codes, c2)
		utf8codes = append(utf8codes, c3)
		return
	} else if unicode <= 0x10FFFF { // last unicode point which can be represented by utf-8
		var c1 byte = 240
		var c2 byte = 128
		var c3 byte = 128
		var c4 byte = 128
		for k := 0; k < 21; k++ {
			if k < 6 {
				c4 |= byte((unicode % 64) & (1 << k))
			} else if k < 12 {
				c3 |= byte((unicode >> 6) & (1 << (k - 6)))
			} else if k < 18 {
				c2 |= byte((unicode >> 12) & (1 << (k - 12)))
			} else {
				c1 |= byte((unicode >> 18) & (1 << (k - 18)))
			}
		}
		utf8codes = append(utf8codes, c1)
		utf8codes = append(utf8codes, c2)
		utf8codes = append(utf8codes, c3)
		utf8codes = append(utf8codes, c4)
		return
	} else {
		// can't be represented by utf-8
		return
	}
}
