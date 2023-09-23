package util

import "regexp"

var RegStringWithVariable regexp.Regexp = *regexp.MustCompile(`\$\{([a-zA-Z\_]+\w*?)\}`)

func IsSpaces(b byte) bool {
	return b == 0x20 || (b < 0x0E && b > 0x08)
}

func IsNumberStartingCharacter(b byte) bool {
	return (b > 0x2F && b < 0x3A) || b == '-'
}
