package util

import (
	"bytes"
	"regexp"
)

var RegStringWithVariable regexp.Regexp = *regexp.MustCompile(`\$\{([a-zA-Z\_]+\w*?)\}`)

func IsSpaces(b byte) bool {
	return b == 0x20 || (b < 0x0E && b > 0x08)
}

func IsNumberStartingCharacter(b byte) bool {
	return (b > 0x2F && b < 0x3A) || b == '-'
}

func RemoveBytes(b []byte, b2remove []byte) []byte {
	parts := bytes.Split(b, b2remove)
	if len(parts) == 1 {
		return b
	}
	var rebuilt []byte
	for _, part := range parts {
		rebuilt = append(rebuilt, part...)
	}
	return rebuilt
}
