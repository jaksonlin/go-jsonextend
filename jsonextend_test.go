package jsonextend_test

import (
	"bytes"
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

func TestUnmarshal(t *testing.T) {
	type SomeStruct struct {
		Field1 string
		Field2 bool
		Field3 int
		Field4 interface{}
	}
	testExample := SomeStruct{
		Field1: "hello",
		Field2: true,
		Field3: 100,
		Field4: nil,
	}

	data, _ := json.Marshal(testExample)
	var out SomeStruct
	err := jsonextend.Unmarshal(bytes.NewReader(data), nil, &out)
	if err != nil {
		t.FailNow()
	}
	var compared SomeStruct
	_ = json.Unmarshal(data, &compared)
	if out.Field1 != compared.Field1 {
		t.FailNow()
	}
	if out.Field2 != compared.Field2 {
		t.FailNow()
	}
	if out.Field3 != compared.Field3 {
		t.FailNow()
	}
	if out.Field4 != compared.Field4 {
		t.FailNow()
	}
}

func TestUnmarshal2(t *testing.T) {
	type SomeStruct struct {
		Field1 string
		Field2 int
		Field3 interface{}
	}
	testExample := `
	{
		"Field1": "hello ${var1}",
		"${var2}": ${var2Value},
		"Field3":${var3}
	}`

	variables := map[string]interface{}{
		"var1":      "world!",
		"var2":      "Field2",
		"var2Value": 100,
		"var3":      []int{1, 2, 3},
	}

	var out SomeStruct
	err := jsonextend.Unmarshal(strings.NewReader(testExample), variables, &out)
	if err != nil {
		t.FailNow()
	}
	if out.Field1 != "hello world!" {
		t.FailNow()
	}
	if out.Field2 != 100 {
		t.FailNow()
	}
	for i, v := range out.Field3.([]int) {
		if v != i+1 {
			t.FailNow()
		}
	}
}
