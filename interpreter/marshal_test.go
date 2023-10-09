package interpreter_test

import (
	"fmt"
	"testing"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
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
	// rs1, err := json.Marshal(sample)
	// if err != nil {
	// 	t.FailNow()
	// }

	sm, err := tokenizer.NewTokenizerStateMachineFromGoData(sample)
	if err != nil {
		t.FailNow()
	}
	err = sm.ProcessData()
	if err != nil {
		t.FailNow()
	}
	if sm.GetASTBuilder().HasOpenElements() {
		t.FailNow()
	}

}
func TestCyclicAccess(t *testing.T) {

	type test1 struct {
		A1  *test1
		Age int
	}

	tr := &test1{}
	ptr := tr
	for i := 0; i < 10; i++ {
		ptr.A1 = &test1{A1: tr, Age: i}
		ptr = ptr.A1
	}
	ptr.A1 = tr

	_, err := interpreter.Marshal(tr)
	if err == nil {
		t.FailNow()
	}

}
func TestCyclicAccess2(t *testing.T) {

	type test1 struct {
		A1  []*test1
		Age int
	}

	tr := &test1{}
	ptr := tr
	rs := []*test1{}
	for i := 0; i < 10; i++ {
		rs = append(rs, &test1{A1: rs, Age: i})
		ptr.A1 = rs
		ptr = ptr.A1[0]
	}
	ptr.A1 = append(ptr.A1, tr)

	_, err := interpreter.Marshal(tr)
	if err == nil {
		t.FailNow()
	}

}
func TestCyclicAccess3(t *testing.T) {

	type test1 struct {
		A1  []*test1
		Age int
	}

	someItem := &test1{nil, 100}
	tr := &test1{nil, 100}
	tr.A1 = make([]*test1, 0)
	for i := 0; i < 10; i++ {
		tr.A1 = append(tr.A1, someItem)
	}

	data, err := interpreter.Marshal(tr)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(data)

}
