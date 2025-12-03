package base64

import "fmt"

const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var ErrInvalidBase64String = fmt.Errorf("invalid base64 string")

func Encode(d []byte) string {
	if len(d) == 0 {
		return ""
	}

	canDivideBy3 := len(d)%3 == 0
	var length int
	if canDivideBy3 {
		length = len(d) / 3 * 4
	} else {
		length = (len(d)/3 + 1) * 4
	}
	res := make([]byte, 0, length)

	for i := 0; i < len(d); i += 3 {
		b1, b2, b3 := d[i], byte(0), byte(0)
		if i+1 < len(d) {
			b2 = d[i+1]
		}
		if i+2 < len(d) {
			b3 = d[i+2]
		}

		// will have 4 groups of 6 bits
		g1 := b1 >> 2            // ex: 101010,01 --> 00,101010
		g2 := b1&0x03<<4 | b2>>4 // ex: 000000,10 1010 --> 10,101000
		g3 := b2&0x0F<<2 | b3>>6 // ex: 0000,0010 101010 --> 10,101000
		g4 := b3 & 0x3F          // ex: 00,101010

		res = append(res, table[g1])
		res = append(res, table[g2])
		if i+1 < len(d) {
			res = append(res, table[g3])
		} else {
			res = append(res, '=')
		}
		if i+2 < len(d) {
			res = append(res, table[g4])
		} else {
			res = append(res, '=')
		}
	}

	return string(res)
}

func Decode(code string) ([]byte, error) {
	if len(code) == 0 {
		return []byte{}, nil
	}

	if len(code)%4 != 0 {
		return nil, ErrInvalidBase64String
	}

	res := make([]byte, 0, len(code)/4*3)

	decodeTable := make(map[byte]byte)
	for i := 0; i < len(table); i++ {
		decodeTable[table[i]] = byte(i)
	}

	for i := 0; i < len(code); i += 4 {
		c1, c2, c3, c4 := code[i], code[i+1], code[i+2], code[i+3]
		var b1, b2, b3 byte

		v1, ok1 := decodeTable[c1]
		v2, ok2 := decodeTable[c2]
		if !ok1 || !ok2 {
			return nil, ErrInvalidBase64String
		}
		b1 = v1<<2 | v2>>4

		res = append(res, b1)

		if c3 != '=' {
			v3, ok3 := decodeTable[c3]
			if !ok3 {
				return nil, ErrInvalidBase64String
			}
			b2 = v2<<4 | v3>>2
			res = append(res, b2)
		}

		if c4 != '=' {
			v4, ok4 := decodeTable[c4]
			if !ok4 {
				return nil, ErrInvalidBase64String
			}
			v3 := decodeTable[c3]
			b3 = v3<<6 | v4
			res = append(res, b3)
		}
	}
	return res, nil
}
