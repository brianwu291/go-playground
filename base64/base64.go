package base64

import "fmt"

const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var ErrInvalidBase64String = fmt.Errorf("invalid base64 string")

func Encode(d []byte) string {
	// separate 3 byte as 1 group
	// for each group, split into 4 group (6 bit a new group)
	// for each new group, make it as a new byte, and access table to get the string
	// concat to result string

	if len(d) == 0 {
		return ""
	}
	// might need add extra "="
	var canDividedBy3 bool
	if len(d)%3 == 0 {
		canDividedBy3 = true
	}
	var resLen int
	if canDividedBy3 {
		resLen = len(d) / 3 * 4
	} else {
		resLen = (len(d)/3 + 1) * 4
	}

	result := make([]byte, 0, resLen)

	for i := 0; i < len(d); i += 3 {
		d1, d2, d3 := d[i], byte(0), byte(0)
		if i+1 < len(d) {
			d2 = d[i+1]
		}
		if i+2 < len(d) {
			d3 = d[i+2]
		}
		one := d1 >> 2
		result = append(result, table[one])

		two := d1&0x03<<4 | d2>>4
		result = append(result, table[two])
		if i+1 < len(d) {
			three := d2&0x0f<<2 | d3>>6
			result = append(result, table[three])
		} else {
			result = append(result, '=')
		}
		if i+2 < len(d) {
			four := d3 & 0x3f
			result = append(result, table[four])
		} else {
			result = append(result, '=')
		}
	}

	return string(result)
}

func buildDecodeTable() map[byte]byte {
	m := make(map[byte]byte, len(table))
	for i := 0; i < len(table); i++ {
		m[table[i]] = byte(i)
	}
	return m
}

var decodeTable = buildDecodeTable()

func Decode(code string) ([]byte, error) {
	if len(code) == 0 {
		return []byte{}, nil
	}

	if len(code)%4 != 0 {
		return nil, ErrInvalidBase64String
	}

	res := make([]byte, 0, len(code)/4*3)

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
