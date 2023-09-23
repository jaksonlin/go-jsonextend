package jsonextend_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/jaksonlin/go-jsonextend"
)

func TestPoc(t *testing.T) {
	template := `{"hello": "world", "name": "this is my ${name}", "age": ${age}}`
	variables := map[string]interface{}{"name": "jakson", "age": 18}

	result, err := jsonextend.Parse(strings.NewReader(template), variables)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(result), &jsonMap)
	if err != nil {
		fmt.Println(err)
	}
	if len(jsonMap) != 3 {
		t.FailNow()
	}
	if jsonMap["hello"] != "world" {
		t.FailNow()
	}
	if jsonMap["name"] != "this is my jakson" {
		t.FailNow()
	}
	if jsonMap["age"] != 18.0 {
		t.FailNow()
	}
}
