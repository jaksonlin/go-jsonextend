package util

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestGetExtTag(t *testing.T) {
	tag := `jsonext:k=a,v=b`
	// jsonext:k=a,v=b
	// jsonext:k=a
	// jsonext:v=b
	pattern := regexp.MustCompile(`\W(\w+=\w+)`)
	matches := pattern.FindAllStringSubmatch(tag, -1)
	fmt.Println(matches)

}
func TestRemoveQuote(t *testing.T) {
	s := "\"abcdefg\""
	f, err := strconv.Unquote(s)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(f)
}
