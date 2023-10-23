package jsonextend_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/jaksonlin/go-jsonextend"
	"github.com/jaksonlin/go-jsonextend/astbuilder/golang"
	"github.com/jaksonlin/go-jsonextend/interpreter"
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
		t.Log(rs1)
		t.Log("----------------")
		t.Log(rs2)
		t.Log("===================")
		t.Log(string(rs1))
		t.Log("----------------")
		t.Log(string(rs2))
		t.FailNow()
	}
}

func TestMarshalObjCyclic(t *testing.T) {

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
	sample.Something = &sample
	_, err := jsonextend.Marshal(sample)
	if err != golang.ErrorCyclicAccess {
		t.FailNow()
	}
}

func TestMarshalObjEmbed(t *testing.T) {

	type test1 struct {
		Name      string `json:"test1_name"`
		Age       int    `json:"test1_age"`
		IsOK      bool   `json:"test1_ok"`
		Something *test1 `json:"test1_something"`
	}

	type test2 struct {
		test1
		MyName string `json:"test2_name"`
	}

	var data test1 = test1{
		Name:      "test1",
		Age:       10,
		IsOK:      true,
		Something: nil,
	}
	var sample test2 = test2{
		test1:  data,
		MyName: "Annie",
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

func TestMarshalObjPointer(t *testing.T) {

	type test1 struct {
		Name      string `json:"test1_name"`
		Age       int    `json:"test1_age"`
		IsOK      bool   `json:"test1_ok"`
		Something *test1 `json:"test1_something"`
	}

	type test2 struct {
		Data1  *test1
		MyName string `json:"test2_name"`
	}

	var data test1 = test1{
		Name:      "test1",
		Age:       10,
		IsOK:      true,
		Something: nil,
	}
	var sample test2 = test2{
		Data1:  &data,
		MyName: "Annie",
	}
	rs1, err := json.Marshal(sample)
	if err != nil {
		fmt.Println("----------------------------------")
		t.Log(err)
		t.FailNow()
	}

	rs2, err := jsonextend.Marshal(sample)
	if err != nil {
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
		t.Log(err)
		t.FailNow()
	}
	if !bytes.Equal(rs1, rs2) {
		fmt.Println("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")

		t.Log(string(rs1))
		t.Log(string(rs2))
		t.FailNow()
	}
}
func TestMarshalObjStructField(t *testing.T) {

	type test1 struct {
		Name      string `json:"test1_name"`
		Age       int    `json:"test1_age"`
		IsOK      bool   `json:"test1_ok"`
		Something *test1 `json:"test1_something"`
	}

	type test2 struct {
		Data1  test1
		MyName string `json:"test2_name"`
	}

	var data test1 = test1{
		Name:      "test1",
		Age:       10,
		IsOK:      true,
		Something: nil,
	}
	var sample test2 = test2{
		Data1:  data,
		MyName: "Annie",
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

func TestMarshalObjEmbedSameName(t *testing.T) {

	type test1 struct {
		Name      string `json:"test1_name"`
		Age       int    `json:"test1_age"`
		IsOK      bool   `json:"test1_ok"`
		Something *test1 `json:"test1_something"`
	}

	type test2 struct {
		test1
		MyName string `json:"test1_name"`
	}

	var data test1 = test1{
		Name:      "test1",
		Age:       10,
		IsOK:      true,
		Something: nil,
	}
	var sample test2 = test2{
		test1:  data,
		MyName: "Annie",
	}
	rs1, err := json.Marshal(sample)
	if err != nil {
		t.FailNow()
	}
	var checker test2
	err = json.Unmarshal(rs1, &checker)
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

func TestOutput(t *testing.T) {
	template := `{"hello": "world", "name": "this is my ${name}", "age": ${age}}`
	variables := map[string]interface{}{"name": "jakson", "age": 18}

	result, err := jsonextend.Parse(strings.NewReader(template), variables)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(result))
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

func TestCustomizeMarshallerStealSky1(t *testing.T) {
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

func TestCustomizeMarshallerStealSky2(t *testing.T) {
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

func TestCustomizeMarshallerStealSky3(t *testing.T) {
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
