package util

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode/utf8"
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

func RepairUTF8(s string) string {
	if utf8.ValidString(s) {
		return s // Already valid UTF-8.
	}

	var repaired []rune
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		repaired = append(repaired, r)
		s = s[size:]
	}

	return string(repaired)
}

// given non-ASCII UTF-8 strings encode it to json string align with RFC 7159
func EncodeToJsonString(input string) []byte {
	var buf bytes.Buffer
	buf.WriteByte('"') // start the JSON string
	for _, r := range input {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 {
				buf.WriteString(fmt.Sprintf(`\u%04x`, r)) //convert the rune value to a 4-character hexadecimal string
			} else {
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteByte('"') // end the JSON string
	return buf.Bytes()
}
