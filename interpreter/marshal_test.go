package interpreter_test

import (
	"bytes"
	"encoding/json"
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

func TestStringOptionMarshal(t *testing.T) {

	type Example struct {
		F1 string `json:",string"`
		F2 int    `json:",string"`
		F3 bool   `json:",string"`
		F4 []byte `json:","`
		F5 []byte `json:",string"`
	}

	ex := Example{"hello", 123, true, []byte("hello"), []byte("hello")}
	var checker Example
	data, err := testMarshaler(ex)
	if err != nil {
		t.FailNow()
	}

	err = json.Unmarshal(data, &checker)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver Example
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.FailNow()
	}
	if myReceiver.F1 != checker.F1 {
		t.FailNow()
	}
	if myReceiver.F2 != checker.F2 {
		t.FailNow()
	}
	if myReceiver.F3 != checker.F3 {
		t.FailNow()
	}
}

func TestOmitEmptyOptionMarshal(t *testing.T) {

	type Example struct {
		F1 string      `json:",omitempty"`
		F2 int         `json:",omitempty"`
		F3 bool        `json:",omitempty"`
		F4 []byte      `json:",omitempty"`
		F5 interface{} `json:",omitempty"`
		F6 [3]string   `json:",omitempty"`
	}

	ex := Example{}
	var checker Example
	data, err := testMarshaler(ex)
	if err != nil {
		t.FailNow()
	}
	data2, err := json.Marshal(ex)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(data2)
	if !bytes.Equal(data, data2) {
		t.FailNow()
	}
	err = json.Unmarshal(data, &checker)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver Example
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.FailNow()
	}
	if myReceiver.F1 != checker.F1 {
		t.FailNow()
	}
	if myReceiver.F2 != checker.F2 {
		t.FailNow()
	}
	if myReceiver.F3 != checker.F3 {
		t.FailNow()
	}
	if myReceiver.F4 != nil {
		t.FailNow()
	}
	if myReceiver.F5 != nil {
		t.FailNow()
	}
	if myReceiver.F6 != [3]string{} {
		t.FailNow()
	}
}
