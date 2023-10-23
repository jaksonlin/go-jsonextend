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

	sm, err := tokenizer.NewTokenizerStateMachineFromGoData(sample, nil)
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
	data1, err := json.Marshal(ex)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(data1)
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

func TestStringOptionMarshalWithInterface(t *testing.T) {

	type Example struct {
		F1 interface{} `json:",string"`
		F2 interface{} `json:",string"`
		F3 interface{} `json:",string"`
		F4 interface{} `json:","`
		F5 interface{} `json:",string"`
	}

	ex := Example{"hello", 123, true, []byte("hello"), []byte("hello")}
	var checker Example
	data2, err := json.Marshal(ex)
	if err != nil {
		t.FailNow()
	}

	data, err := testMarshaler(ex)
	if err != nil {
		t.FailNow()
	}

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
}

func TestStringOptionMarshalWithPointer(t *testing.T) {

	type Example struct {
		F1 *string      `json:",string"`
		F2 *int         `json:",string"`
		F3 *bool        `json:",string"`
		F4 *[]byte      `json:",string"`
		F5 *interface{} `json:",string"`
		F6 *[3]string   `json:",string"`
	}

	f1 := "hello"
	f2 := 123
	f3 := true
	f4 := []byte("hello")
	f5 := interface{}(nil)
	f6 := [3]string{}
	ex := Example{&f1, &f2, &f3, &f4, &f5, &f6}
	var checker Example
	data2, err := json.Marshal(ex)
	if err != nil {
		t.FailNow()
	}
	data, err := testMarshaler(ex)
	if err != nil {
		t.FailNow()
	}

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
	if *myReceiver.F1 != *checker.F1 {
		t.FailNow()
	}
	if *myReceiver.F2 != *checker.F2 {
		t.FailNow()
	}
	if *myReceiver.F3 != *checker.F3 {
		t.FailNow()
	}
	if !bytes.Equal(*myReceiver.F4, *checker.F4) {
		t.FailNow()
	}
	if myReceiver.F5 != checker.F5 && myReceiver.F5 != nil && checker.F5 != nil {
		t.FailNow()
	}
	if *myReceiver.F6 != *checker.F6 {
		t.FailNow()
	}
}

func TestStringOptionMarshalWithInterfacePointer(t *testing.T) {

	type Example struct {
		F1 *interface{} `json:",string"`
		F2 *interface{} `json:",string"`
		F3 *interface{} `json:",string"`
		F4 *interface{} `json:",string"`
		F5 *interface{} `json:",string"`
		F6 *interface{} `json:",string"`
	}

	var f1 interface{} = "hello"
	var f2 interface{} = 123
	var f3 interface{} = true
	var f4 interface{} = []byte("hello")
	var f5 interface{} = interface{}(nil)
	var f6 interface{} = [3]string{}
	var ex interface{} = Example{&f1, &f2, &f3, &f4, &f5, &f6}
	var checker Example
	data2, err := json.Marshal(ex)
	if err != nil {
		t.FailNow()
	}
	data, err := testMarshaler(ex)
	if err != nil {
		t.FailNow()
	}

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
	if *myReceiver.F1 != *checker.F1 {
		t.FailNow()
	}
	if *myReceiver.F2 != *checker.F2 {
		t.FailNow()
	}
	if *myReceiver.F3 != *checker.F3 {
		t.FailNow()
	}
	if (*myReceiver.F4).(string) != (*checker.F4).(string) {
		t.FailNow()
	}
	if myReceiver.F5 != checker.F5 && myReceiver.F5 != nil && checker.F5 != nil {
		t.FailNow()
	}
	if len((*myReceiver.F6).([3]string)) != len((*checker.F6).([3]string)) {
		t.FailNow()
	}
}

func TestCustomizeMarshaller(t *testing.T) {
	type MyDataStruct struct {
		Name string `jsonext:"k=var1,v=var2"`
	}
	item := &MyDataStruct{
		Name: "hello",
	}

	data, err := interpreter.MarshalIntoTemplate(item)
	if err != nil {
		t.FailNow()
	}

	fmt.Println(data)
	fmt.Println(data)
	if string(data) != `{"${var1}":${var2}}` {
		t.FailNow()
	}
}

func TestCustomizeMarshallerVariable(t *testing.T) {
	type MyDataStruct struct {
		Name string `json:"myfield" jsonext:"v=var1"`
	}
	item := &MyDataStruct{
		Name: "hello",
	}

	data, err := interpreter.MarshalWithVariable(item, map[string]interface{}{"var1": "my love"})
	if err != nil {
		t.FailNow()
	}

	fmt.Println(data)
	fmt.Println(data)
	if string(data) != `{"myfield":"my love"}` {
		t.FailNow()
	}
}

func TestCustomizeMarshallerKey(t *testing.T) {
	type MyDataStruct struct {
		Name string `json:"myfield" jsonext:"k=var1"`
	}
	item := &MyDataStruct{
		Name: "hello",
	}

	data, err := interpreter.MarshalWithVariable(item, map[string]interface{}{"var1": "my love"})
	if err != nil {
		t.FailNow()
	}

	fmt.Println(data)
	fmt.Println(data)
	if string(data) != `{"my love":"hello"}` {
		t.FailNow()
	}
}

func TestCustomizeMarshallerOnStruct(t *testing.T) {
	type someStruct struct {
		Name2 string
	}
	type MyDataStruct struct {
		Name someStruct `jsonext:"k=var1,v=var2"`
	}
	item := &MyDataStruct{
		Name: someStruct{"ddd"},
	}

	data, err := interpreter.MarshalIntoTemplate(item)
	if err != nil {
		t.FailNow()
	}

	fmt.Println(data)
	fmt.Println(data)
	if string(data) != `{"${var1}":${var2}}` {
		t.FailNow()
	}
}
