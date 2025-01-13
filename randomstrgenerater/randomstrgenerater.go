package randomstrgenerater

import (
	"fmt"
	"strings"

	uuid "github.com/google/uuid"
)

func pickCharWithJumper(str string, jumper int, expectedLength int) string {
	var result string
	var enough bool
	for i := 0; !enough; i += jumper {
		if i >= len(str) {
			i = i - len(str)
		}
		curChar := string(rune(str[i]))
		result = result + curChar
		if len(result) >= expectedLength {
			enough = true
		}
	}
	return result
}

func GenerateRandomString(length int) (string, error) {
	// UUID v4
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// convert to hex string without hyphens
	uuidWithoutHyphen := strings.ReplaceAll(uuid.String(), "-", "")

	// fmt.Printf("uuidWithoutHyphen: %+v\n", uuidWithoutHyphen)

	// take first n characters based on requested length
	if len(uuidWithoutHyphen) < length {
		return "", fmt.Errorf("requested length %d exceeds UUID string length %d", length, len(uuidWithoutHyphen))
	}

	result := pickCharWithJumper(uuidWithoutHyphen, 1, length)

	return result, nil
}
