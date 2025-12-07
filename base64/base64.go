package base64

import (
	"fmt"
)

const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var ErrInvalidBase64String = fmt.Errorf("invalid base64 string")

func Encode(d []byte) string {
	// split d 3 byte a group, for loop
	// each group, every 6 bit a group as a char from table
	// if empty, append "="

	if len(d) == 0 {
		return ""
	}
	var resLen int
	if len(d)%3 == 0 {
		resLen = len(d) / 3 * 4
	} else {
		resLen = (len(d)/3 + 1) * 4
	}

	result := make([]byte, 0, resLen)

	for i := 0; i < len(d); i += 3 {
		b1, b2, b3 := d[i], byte(0), byte(0)
		if i+1 < len(d) {
			b2 = d[i+1]
		}
		if i+2 < len(d) {
			b3 = d[i+2]
		}

		c1 := table[b1>>2]
		result = append(result, c1)

		c2 := table[b1&0x03<<4|b2>>4]
		result = append(result, c2)

		if i+1 < len(d) {
			c3 := table[b2&0x0f<<2|b3>>6]
			result = append(result, c3)
		} else {
			result = append(result, '=')
		}

		if i+2 < len(d) {
			c4 := table[b3&0x3f]
			result = append(result, c4)
		} else {
			result = append(result, '=')
		}
	}

	return string(result)
}

func buildDecodeTable() map[byte]byte {
	res := make(map[byte]byte, len(table))
	for i := 0; i < len(table); i += 1 {
		res[table[i]] = byte(i)
	}
	return res
}

var decodedTable = buildDecodeTable()

func Decode(code string) ([]byte, error) {
	// each chat represent 6 bit of source data
	// split with 4 byte a group
	// each group, every 8 bit, map back to byte

	if code == "" {
		return nil, nil
	}

	if len(code)%4 != 0 { // no need to use rune here since we will handle it later in loop
		return nil, ErrInvalidBase64String
	}

	result := make([]byte, 0, len(code)/4*3)
	for i := 0; i < len(code); i += 4 { // here item access from decodedTable can handle input code contains Not ASCII text
		c1, c2, c3, c4 := code[i], code[i+1], code[i+2], code[i+3]
		v1, ok1 := decodedTable[c1]
		v2, ok2 := decodedTable[c2]
		if !ok1 || !ok2 {
			return nil, ErrInvalidBase64String
		}
		b1 := v1<<2 | v2>>4
		result = append(result, b1)

		if c3 != '=' {
			v3, ok := decodedTable[c3]
			if !ok {
				return nil, ErrInvalidBase64String
			}
			b2 := v2<<4 | v3>>2
			result = append(result, b2)
		}

		if c4 != '=' {
			v4, ok := decodedTable[c4]
			if !ok {
				return nil, ErrInvalidBase64String
			}
			v3 := decodedTable[c3]
			b3 := v3<<6 | v4
			result = append(result, b3)
		}
	}

	return result, nil
}
