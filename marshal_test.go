package jsonextend_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/jaksonlin/go-jsonextend"
)

func TestMarshalObj(t *testing.T) {

	type test1 struct {
		Name      string `json:"test1_name"`
		Age       int    `json:"test1_age"`
		IsOK      bool   `json:"test1_ok"`
		Something *test1 `json:"test1_something"`
	}

	var sample test1 = test1{
		Name:      "test1",
		Age:       10,
		IsOK:      true,
		Something: nil,
	}
	rs1, err := json.Marshal(sample)
	if err != nil {
		t.FailNow()
	}
	rs2, err := jsonextend.Marshal(sample)
	if err != nil {
		t.FailNow()
	}
	if !bytes.Equal(rs1, rs2) {
		t.FailNow()
	}
}
