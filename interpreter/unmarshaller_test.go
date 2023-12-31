package interpreter_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

// these will be used as cross over validation check.
var myMarshaler = func(v interface{}) ([]byte, error) { return interpreter.Marshal(v) }
var myUnMarshaler = func(data []byte, v interface{}) error {
	return interpreter.Unmarshal(bytes.NewReader(data), nil, v)
}
var jsonMarshaler = func(v interface{}) ([]byte, error) { return json.Marshal(v) }
var jsonUnMarshaler = func(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
var testMarshaler = myMarshaler
var testUnmarshaler = myUnMarshaler

func TestAssignThingsPrimitives(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	t1 := test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))

	var validator test1
	_ = json.Unmarshal(data, &validator)
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(strings.NewReader(string(data)))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 test1
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someTest1.Hello != validator.Hello {
		t.FailNow()
	}
	if someTest1.World != validator.World {
		t.FailNow()
	}
	if someTest1.World2 != validator.World2 {
		t.FailNow()
	}
	if someTest1.Apple != validator.Apple {
		t.FailNow()
	}
	if someTest1.Banana != validator.Banana {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != validator.SomethingNotNil {
		t.FailNow()
	}
}
func TestAssignThingsPrimitivesMapBasic(t *testing.T) {
	var t1 map[string]int = map[string]int{
		"Hello":           1,
		"World":           2,
		"World2":          3,
		"Apple":           4,
		"Banana":          5,
		"Something":       6,
		"SomethingNotNil": 7,
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var validator map[string]int
	err := json.Unmarshal(data, &validator)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(strings.NewReader(string(data)))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 map[string]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for k, v := range someTest1 {
		if validator[k] != v {
			t.FailNow()
		}
	}

}

func TestAssignThingsPrimitivesMapStructKeyBasic(t *testing.T) {
	type someTest struct {
		Number int32
	}
	var t1 map[string]someTest = map[string]someTest{
		"1": someTest{1},
		"2": someTest{1},
		"3": someTest{1},
		"4": someTest{1},
		"5": someTest{1},
		"6": someTest{1},
		"7": someTest{1},
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var validator map[string]someTest
	err := json.Unmarshal(data, &validator)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(strings.NewReader(string(data)))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 map[string]someTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for k, v := range someTest1 {
		if validator[k].Number != v.Number {
			t.FailNow()
		}
	}

}

func TestAssignThingsPrimitivesMapPtrStructKeyBasic(t *testing.T) {
	type someTest struct {
		Number int32
	}
	var t1 map[string]*someTest = map[string]*someTest{
		"1": &someTest{1},
		"2": &someTest{1},
		"3": &someTest{1},
		"4": &someTest{1},
		"5": &someTest{1},
		"6": &someTest{1},
		"7": &someTest{1},
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var validator map[string]someTest
	err := json.Unmarshal(data, &validator)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(strings.NewReader(string(data)))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 map[string]someTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for k, v := range someTest1 {
		if validator[k].Number != v.Number {
			t.FailNow()
		}
	}

}
func TestAssignThingsPrimitivesMapStructKey(t *testing.T) {
	type someTest struct {
		Number int32
	}
	var t1 map[int]someTest = map[int]someTest{
		1: someTest{1},
		2: someTest{1},
		3: someTest{1},
		4: someTest{1},
		5: someTest{1},
		6: someTest{1},
		7: someTest{1},
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var validator map[int]someTest
	err := json.Unmarshal(data, &validator)
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

	var someTest1 map[int]someTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for k, v := range someTest1 {
		if validator[k].Number != v.Number {
			t.FailNow()
		}
	}

}

func TestAssignThingsPrimitivesMapPtrStructKey(t *testing.T) {
	type someTest struct {
		Number int32
	}
	var t1 map[string]*someTest = map[string]*someTest{
		"Hello":           &someTest{1},
		"World":           &someTest{1},
		"World2":          &someTest{1},
		"Apple":           &someTest{1},
		"Banana":          &someTest{1},
		"Something":       &someTest{1},
		"SomethingNotNil": &someTest{1},
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var validator map[string]someTest
	err := json.Unmarshal(data, &validator)
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

	var someTest1 map[string]someTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for k, v := range someTest1 {
		if validator[k].Number != v.Number {
			t.FailNow()
		}
	}

}
func TestAssignThingsPrimitivesMapInterface(t *testing.T) {
	var t1 map[string]interface{} = map[string]interface{}{
		"Hello":           "Peter",
		"World":           101.123,
		"World2":          100,
		"Apple":           true,
		"Banana":          false,
		"Something":       nil,
		"SomethingNotNil": "1234",
	}

	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(data))
	var validator map[string]interface{}
	err = json.Unmarshal(data, &validator)
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
	var someTest1Check map[string]interface{}
	err = json.Unmarshal(data, &someTest1Check)
	if err != nil {
		t.FailNow()
	}
	var someTest1 map[string]interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someTest1["Hello"] != validator["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != validator["World"] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != validator["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != validator["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != validator["Banana"] {
		t.FailNow()
	}
	if validator["Something"] != nil {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}

	if someTest1["SomethingNotNil"] != validator["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesMapInterfaceInNonePointerRoot(t *testing.T) {
	var t1 map[string]interface{} = map[string]interface{}{
		"Hello":           "Peter",
		"World":           101.123,
		"World2":          100,
		"Apple":           true,
		"Banana":          false,
		"Something":       nil,
		"SomethingNotNil": "1234",
	}

	type someRoot struct {
		MapData map[string]interface{}
	}

	testData := someRoot{t1}

	data, _ := testMarshaler(testData)
	fmt.Println(string(data))
	var validator someRoot
	err := json.Unmarshal(data, &validator)
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

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &rootCheck)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := rootCheck.MapData
	if someTest1["Hello"] != validator.MapData["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != validator.MapData["World"] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != validator.MapData["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != validator.MapData["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != validator.MapData["Banana"] {
		t.FailNow()
	}
	if validator.MapData["Something"] != nil {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != validator.MapData["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesMapInterfaceInPointerRoot(t *testing.T) {
	var t1 map[string]interface{} = map[string]interface{}{
		"Hello":           "Peter",
		"World":           101.123,
		"World2":          100,
		"Apple":           true,
		"Banana":          false,
		"Something":       nil,
		"SomethingNotNil": "1234",
	}

	type someRoot struct {
		MapData *map[string]interface{}
	}

	testData := someRoot{&t1}

	data, _ := testMarshaler(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
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

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &rootCheck)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := *rootCheck.MapData
	checkerRootMapData := *checkerRoot.MapData
	if someTest1["Hello"] != checkerRootMapData["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != checkerRootMapData["World"] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != checkerRootMapData["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != checkerRootMapData["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != checkerRootMapData["Banana"] {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != checkerRootMapData["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesMapInterfaceInInterfaceRoot(t *testing.T) {
	var t1 interface{} = map[string]interface{}{
		"Hello":           "Peter",
		"World":           101.123,
		"World2":          100,
		"Apple":           true,
		"Banana":          false,
		"Something":       nil,
		"SomethingNotNil": "1234",
	}

	type someRoot struct {
		MapData interface{}
	}

	testData := someRoot{t1}

	data, _ := testMarshaler(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
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

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &rootCheck)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := (rootCheck.MapData).(map[string]interface{})
	checkerRootMapData := (checkerRoot.MapData).(map[string]interface{})
	if someTest1["Hello"] != checkerRootMapData["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != checkerRootMapData["World"] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != checkerRootMapData["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != checkerRootMapData["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != checkerRootMapData["Banana"] {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != checkerRootMapData["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesMapInterfaceInInterfacePointerRoot(t *testing.T) {
	var t1 interface{} = map[string]interface{}{
		"Hello":           "Peter",
		"World":           101.123,
		"World2":          100,
		"Apple":           true,
		"Banana":          false,
		"Something":       nil,
		"SomethingNotNil": "1234",
	}

	type someRoot struct {
		MapData *interface{}
	}

	testData := someRoot{&t1}

	data, err := testMarshaler(testData)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(data))

	var checkerRoot someRoot
	err = json.Unmarshal(data, &checkerRoot)
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

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &rootCheck)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := (*rootCheck.MapData).(map[string]interface{})
	checkerRootMapData := (*checkerRoot.MapData).(map[string]interface{})
	if someTest1["Hello"] != checkerRootMapData["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != checkerRootMapData["World"] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != checkerRootMapData["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != checkerRootMapData["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != checkerRootMapData["Banana"] {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != checkerRootMapData["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesMapInterfaceInInterfacePointersRoot(t *testing.T) {
	var t1 interface{} = map[string]interface{}{
		"Hello":           "Peter",
		"World":           101.123,
		"World2":          100,
		"Apple":           true,
		"Banana":          false,
		"Something":       nil,
		"SomethingNotNil": "1234",
	}

	type someRoot struct {
		MapData ****interface{}
	}

	t2 := &t1
	t3 := &t2
	t4 := &t3
	testData := someRoot{&t4}

	data, _ := testMarshaler(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
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

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &rootCheck)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := (****rootCheck.MapData).(map[string]interface{})
	checkerRootMapData := (****checkerRoot.MapData).(map[string]interface{})
	if someTest1["Hello"] != checkerRootMapData["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != checkerRootMapData["World"] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != checkerRootMapData["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != checkerRootMapData["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != checkerRootMapData["Banana"] {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != checkerRootMapData["SomethingNotNil"] {
		t.FailNow()
	}
}
func TestAssignThingsPrimitivesMapIntInterface(t *testing.T) {
	var t1 map[int]interface{} = map[int]interface{}{
		1: "Peter",
		2: 101.123,
		3: 100,
		4: true,
		5: false,
		6: nil,
		7: "1234",
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var someTest1Check map[int]interface{}
	err = json.Unmarshal(data, &someTest1Check)
	if err != nil {
		t.FailNow()
	}
	var someTest1 map[int]interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someTest1[1] != someTest1Check[1] {
		t.FailNow()
	}
	if someTest1[2] != someTest1Check[2] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1[3] != someTest1Check[3] {
		t.FailNow()
	}
	if someTest1[4] != someTest1Check[4] {
		t.FailNow()
	}
	if someTest1[5] != someTest1Check[5] {
		t.FailNow()
	}
	if someTest1[6] != nil && someTest1Check[6] != someTest1[6] {
		t.FailNow()
	}
	if someTest1[7] != someTest1Check[7] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesMapBoolInterface(t *testing.T) {
	var t1 map[uint8]interface{} = map[uint8]interface{}{
		1: "Peter",
		2: 101.123,
	}

	data, _ := testMarshaler(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var someTest1Check map[int]interface{}
	err = json.Unmarshal(data, &someTest1Check)
	if err != nil {
		t.FailNow()
	}
	var someTest1 map[uint8]interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someTest1[1] != someTest1Check[1] {
		t.FailNow()
	}
	if someTest1[2] != someTest1Check[2] { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}

}

func TestAssignThingsPrimitivesInNonePointerRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	type testRoot struct {
		Test1Data test1
	}

	t1 := test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	tr := testRoot{t1}
	data, _ := testMarshaler(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := someRoot.Test1Data
	if someTest1.Hello != checkerRoot.Test1Data.Hello {
		t.FailNow()
	}
	if someTest1.World != checkerRoot.Test1Data.World {
		t.FailNow()
	}
	if someTest1.World2 != checkerRoot.Test1Data.World2 {
		t.FailNow()
	}
	if someTest1.Apple != checkerRoot.Test1Data.Apple {
		t.FailNow()
	}
	if someTest1.Banana != checkerRoot.Test1Data.Banana {
		t.FailNow()
	}
	if someTest1.Something != nil && someTest1.Something != checkerRoot.Test1Data.Something {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != checkerRoot.Test1Data.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesInInterfaceRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	type testRoot struct {
		Test1Data interface{}
	}

	t1 := test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	tr := testRoot{t1}
	data, _ := testMarshaler(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := someRoot.Test1Data.(map[string]interface{})
	checkerTest1 := checkerRoot.Test1Data.(map[string]interface{})
	if someTest1["Hello"] != checkerTest1["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != checkerTest1["World"] {
		t.FailNow()
	}
	if someTest1["World2"] != checkerTest1["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != checkerTest1["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != checkerTest1["Banana"] {
		t.FailNow()
	}
	if someTest1["Something"] != nil && someTest1["Something"] != checkerTest1["Something"] {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != checkerTest1["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesInPointerInterfaceRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	type testRoot struct {
		Test1Data *interface{}
	}

	var t1 interface{} = test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	tr := testRoot{&t1}
	data, _ := testMarshaler(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := (*someRoot.Test1Data).(map[string]interface{})
	checkerTest1 := (*checkerRoot.Test1Data).(map[string]interface{})
	if someTest1["Hello"] != checkerTest1["Hello"] {
		t.FailNow()
	}
	if someTest1["World"] != checkerTest1["World"] {
		t.FailNow()
	}
	if someTest1["World2"] != checkerTest1["World2"] {
		t.FailNow()
	}
	if someTest1["Apple"] != checkerTest1["Apple"] {
		t.FailNow()
	}
	if someTest1["Banana"] != checkerTest1["Banana"] {
		t.FailNow()
	}
	if someTest1["Something"] != nil && someTest1["Something"] != checkerTest1["Something"] {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != checkerTest1["SomethingNotNil"] {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesInPointerRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	type testRoot struct {
		Test1Data *test1
	}

	t1 := test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	tr := testRoot{&t1}
	data, _ := testMarshaler(tr)
	fmt.Println(string(data))

	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := someRoot.Test1Data
	if someTest1.Hello != checkerRoot.Test1Data.Hello {
		t.FailNow()
	}
	if someTest1.World != checkerRoot.Test1Data.World {
		t.FailNow()
	}
	if someTest1.World2 != checkerRoot.Test1Data.World2 {
		t.FailNow()
	}
	if someTest1.Apple != checkerRoot.Test1Data.Apple {
		t.FailNow()
	}
	if someTest1.Banana != checkerRoot.Test1Data.Banana {
		t.FailNow()
	}
	if someTest1.Something != nil && someTest1.Something != checkerRoot.Test1Data.Something {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != checkerRoot.Test1Data.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesInNestedPointerRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	type testRoot struct {
		Test1Data **test1
	}

	t1 := &test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	tr := testRoot{&t1}
	data, _ := testMarshaler(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := **someRoot.Test1Data
	checkerTest1 := **checkerRoot.Test1Data
	if someTest1.Hello != checkerTest1.Hello {
		t.FailNow()
	}
	if someTest1.World != checkerTest1.World {
		t.FailNow()
	}
	if someTest1.World2 != checkerTest1.World2 {
		t.FailNow()
	}
	if someTest1.Apple != checkerTest1.Apple {
		t.FailNow()
	}
	if someTest1.Banana != checkerTest1.Banana {
		t.FailNow()
	}
	if someTest1.Something != nil && someTest1.Something != checkerTest1.Something {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != checkerTest1.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivesInNestedStructRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           float64
		World2          int
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
	}

	type test2 struct {
		Test1Data *test1
	}
	type test3 struct {
		Test2Data test2
	}
	type testRoot struct {
		Test3Data **test3
	}

	t1 := &test1{"Peter", 100.123, 100, true, false, nil, "1234"}
	t2 := test2{t1}
	t3 := &test3{t2}
	tr := testRoot{&t3}

	data, err := testMarshaler(tr)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(data))
	var checkerRoot testRoot
	err = json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := *(**someRoot.Test3Data).Test2Data.Test1Data
	checkerTest1 := *(**checkerRoot.Test3Data).Test2Data.Test1Data
	if someTest1.Hello != checkerTest1.Hello {
		t.FailNow()
	}
	if someTest1.World != checkerTest1.World {
		t.FailNow()
	}
	if someTest1.World2 != checkerTest1.World2 {
		t.FailNow()
	}
	if someTest1.Apple != checkerTest1.Apple {
		t.FailNow()
	}
	if someTest1.Banana != checkerTest1.Banana {
		t.FailNow()
	}
	if someTest1.Something != nil && someTest1.Something != checkerTest1.Something {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != checkerTest1.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivePointers(t *testing.T) {
	type test1 struct {
		Hello           *string
		World           *float64
		World2          *int
		Apple           *bool
		Banana          *bool
		Something       *interface{}
		SomethingNotNil *interface{}
	}

	someStr := "Peter"
	someFloat := 100.123
	someInt := 100
	someTrue := true
	someFalse := false

	var someInterface interface{} = nil
	t1 := test1{&someStr, &someFloat, &someInt, &someTrue, &someFalse, &someInterface, new(interface{})}
	*t1.SomethingNotNil = "1234"

	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(data))

	var someRs test1
	err = json.Unmarshal(data, &someRs)
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

	var someTest1 test1
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if *someTest1.Hello != *someRs.Hello {
		t.FailNow()
	}
	if *someTest1.World != *someRs.World {
		t.FailNow()
	}
	if *someTest1.World2 != *someRs.World2 {
		t.FailNow()
	}
	if *someTest1.Apple != *someRs.Apple {
		t.FailNow()
	}
	if *someTest1.Banana != *someRs.Banana {
		t.FailNow()
	}
	if someRs.Something != nil {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if *someTest1.SomethingNotNil != *someRs.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivePointersInPointerRoot(t *testing.T) {
	type test1 struct {
		Hello           *string
		World           *float64
		World2          *int
		Apple           *bool
		Banana          *bool
		Something       *interface{}
		SomethingNotNil *interface{}
	}
	type testRoot struct {
		Test1Data *test1
	}
	someStr := "Peter"
	someFloat := 100.123
	someInt := 100
	someTrue := true
	someFalse := false

	var someInterface interface{} = nil
	t1 := test1{&someStr, &someFloat, &someInt, &someTrue, &someFalse, &someInterface, new(interface{})}
	*t1.SomethingNotNil = "1234"
	var someRoot testRoot = testRoot{&t1}
	data, _ := testMarshaler(someRoot)
	fmt.Println(string(data))

	var someRs testRoot
	err := json.Unmarshal(data, &someRs)
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

	var someTestRoot1 testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTestRoot1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := someTestRoot1.Test1Data
	if *someTest1.Hello != *someRs.Test1Data.Hello {
		t.FailNow()
	}
	if *someTest1.World != *someRs.Test1Data.World {
		t.FailNow()
	}
	if *someTest1.World2 != *someRs.Test1Data.World2 {
		t.FailNow()
	}
	if *someTest1.Apple != *someRs.Test1Data.Apple {
		t.FailNow()
	}
	if *someTest1.Banana != *someRs.Test1Data.Banana {
		t.FailNow()
	}
	if someRs.Test1Data.Something != nil {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if *someTest1.SomethingNotNil != *someRs.Test1Data.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitivePointersInNonePointerRoot(t *testing.T) {
	type test1 struct {
		Hello           *string
		World           *float64
		World2          *int
		Apple           *bool
		Banana          *bool
		Something       *interface{}
		SomethingNotNil *interface{}
	}
	type testRoot struct {
		Test1Data test1
	}
	someStr := "Peter"
	someFloat := 100.123
	someInt := 100
	someTrue := true
	someFalse := false

	var someInterface interface{} = nil
	t1 := test1{&someStr, &someFloat, &someInt, &someTrue, &someFalse, &someInterface, new(interface{})}
	*t1.SomethingNotNil = "1234"
	var someRoot testRoot = testRoot{t1}
	data, _ := testMarshaler(someRoot)
	fmt.Println(string(data))

	var someRs testRoot
	err := json.Unmarshal(data, &someRs)
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

	var someTestRoot1 testRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTestRoot1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	someTest1 := someTestRoot1.Test1Data
	if *someTest1.Hello != *someRs.Test1Data.Hello {
		t.FailNow()
	}
	if *someTest1.World != *someRs.Test1Data.World {
		t.FailNow()
	}
	if *someTest1.World2 != *someRs.Test1Data.World2 {
		t.FailNow()
	}
	if *someTest1.Apple != *someRs.Test1Data.Apple {
		t.FailNow()
	}
	if *someTest1.Banana != *someRs.Test1Data.Banana {
		t.FailNow()
	}
	if someRs.Test1Data.Something != nil {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if *someTest1.SomethingNotNil != *someRs.Test1Data.SomethingNotNil {
		t.FailNow()
	}
}

func TestAssignThingsPrimitiveSlice(t *testing.T) {

	t1 := []int{1, 2, 3, 4, 5}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checker []int
	err := json.Unmarshal(data, &checker)
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

	var someTest1 []int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v != checker[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsStructElementSlice(t *testing.T) {

	type someTest struct {
		Number int
	}
	t1 := []someTest{{1}, {2}, {3}, {4}, {5}}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checker []someTest
	err := json.Unmarshal(data, &checker)
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

	var someTest1 []someTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v.Number != checker[i].Number {
			t.FailNow()
		}
	}
}
func TestAssignThingsStructPtrElementSlice(t *testing.T) {

	type someTest struct {
		Number int
	}
	t1 := []*someTest{&someTest{1}, &someTest{2}, &someTest{3}, &someTest{4}, &someTest{5}}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checker []someTest
	err := json.Unmarshal(data, &checker)
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

	var someTest1 []someTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v.Number != checker[i].Number {
			t.FailNow()
		}
	}
}

func TestAssignThingsStructInterfaceElementSlice(t *testing.T) {

	type someTest struct {
		Number int
	}
	t1 := []interface{}{&someTest{1}, &someTest{2}, &someTest{3}, &someTest{4}, &someTest{5}, nil}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checker []interface{}
	err := json.Unmarshal(data, &checker)
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

	var someTest1 []interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v == nil {
			if checker[i] != nil {
				t.FailNow()
			}
			continue
		}
		m1 := v.(map[string]interface{})
		m2 := checker[i].(map[string]interface{})
		for k1, v1 := range m1 {
			if v1 != m2[k1] {
				t.FailNow()
			}
		}
	}
}
func TestAssignThingsPrimitiveSliceInNonePointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField []int
	}

	test1 := someRoot{[]int{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1.SomeField {
		if v != checkerRoot.SomeField[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsNonePrimitiveSliceInNonePointerRoot(t *testing.T) {
	type test1 struct {
		Hello           string
		World           uint16
		World2          int32
		Apple           bool
		Banana          bool
		Something       interface{}
		SomethingNotNil interface{}
		SomethingArray  interface{}
	}
	type someRoot struct {
		SomeField []test1
	}

	t1 := someRoot{[]test1{
		{"Peter", 12345, 2551, true, false, nil, map[string]int{"hello": 1}, []int{1, 2, 3, 4, 5}},
		{"Peter2", 22345, 3551, true, false, nil, map[string]int{"hello": 2}, []string{"1", "2", "3"}},
		{"Peter3", 32345, 4551, true, false, nil, map[string]int{"hello": 3}, []bool{true, false}},
		{"Peter5", 32345, 4551, true, false, nil, map[string]int{"hello": 3}, []interface{}{"string", 123.3, true, nil}},
	}}

	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}

	fmt.Println(string(data))
	var someRootChecker someRoot
	err = json.Unmarshal(data, &someRootChecker)
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

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someRootChecker.SomeField {
		if v.Hello != someTest1.SomeField[i].Hello {
			t.FailNow()
		}
		if v.World != someTest1.SomeField[i].World {
			t.FailNow()
		}
		if v.World2 != someTest1.SomeField[i].World2 {
			t.FailNow()
		}
		if v.Apple != someTest1.SomeField[i].Apple {
			t.FailNow()
		}
		if v.Banana != someTest1.SomeField[i].Banana {
			t.FailNow()
		}
		if v.Something != someTest1.SomeField[i].Something {
			t.FailNow()
		}
		mapThere := v.SomethingNotNil.(map[string]interface{})
		mapThere2 := someTest1.SomeField[i].SomethingNotNil.(map[string]interface{})
		for k, v := range mapThere {
			if mapThere2[k] != v {
				t.FailNow()
			}
		}
		sliceThere := v.SomethingArray.([]interface{})
		sliceThere2 := someTest1.SomeField[i].SomethingArray.([]interface{})
		for i, v := range sliceThere {
			if v != sliceThere2[i] {
				t.FailNow()
			}
		}
	}
}

func TestAssignThingsPrimitiveSliceInPointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField *[]int
	}

	test1 := someRoot{&[]int{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range *someTest1.SomeField {
		if v != (*checkerRoot.SomeField)[i] {
			t.FailNow()
		}
	}
}
func TestAssignThingsSliceWithInterfaceElement(t *testing.T) {

	t1 := []interface{}{1, true, false, nil}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var validator []interface{}
	_ = json.Unmarshal(data, &validator)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 []interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v != validator[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveSliceInterfaceInNonePointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField []interface{}
	}

	test1 := someRoot{[]interface{}{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var validator someRoot
	_ = json.Unmarshal(data, &validator)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1.SomeField {
		if v != validator.SomeField[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveSliceInterfaceInPointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField *[]interface{}
	}

	test1 := someRoot{&[]interface{}{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var validator someRoot
	_ = json.Unmarshal(data, &validator)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range *someTest1.SomeField {
		if v != (*validator.SomeField)[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveSliceInterfaceInPointerRootAndPointerValue(t *testing.T) {

	type someRoot struct {
		SomeField *[]*interface{}
	}

	var value1 interface{} = 1
	var value2 interface{} = 2
	var value3 interface{} = 3
	var value4 interface{} = 4
	var value5 interface{} = nil

	test1 := someRoot{&[]*interface{}{&value1, &value2, &value3, &value4, &value5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var validator someRoot
	err := json.Unmarshal(data, &validator)
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

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i := 0; i < len(*someTest1.SomeField); i++ {
		v := (*validator.SomeField)[i]
		if v != nil {
			v1 := *v
			v2 := *(*someTest1.SomeField)[i]
			if v1 != v2 {
				t.FailNow()
			}
		} else {
			if (*someTest1.SomeField)[i] != nil {
				t.FailNow()
			}
		}

	}
}

func TestAssignThingsPrimitiveArray(t *testing.T) {

	t1 := [5]int{1, 2, 3, 4, 5}
	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))

	var checkerRoot [5]int
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 [5]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v != checkerRoot[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsArrayWithInterfaceElement(t *testing.T) {

	t1 := [4]interface{}{1, true, false, nil}
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checker [4]interface{}
	_ = json.Unmarshal(data, &checker)

	var checkerRoot [4]interface{}
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 [4]interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v != checkerRoot[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArrayInNonePointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField [5]int
	}

	test1 := someRoot{[5]int{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1.SomeField {
		if v != checkerRoot.SomeField[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArrayInPointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField *[5]int
	}

	test1 := someRoot{&[5]int{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var checkerRoot someRoot

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range *someTest1.SomeField {
		if v != (*checkerRoot.SomeField)[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArrayInterfaceInNonePointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField [5]interface{}
	}

	test1 := someRoot{[5]interface{}{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var mysomeRoot someRoot
	_ = json.Unmarshal(data, &mysomeRoot)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1.SomeField {
		if v != mysomeRoot.SomeField[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArrayInterfaceInPointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField *[5]interface{}
	}

	test1 := someRoot{&[5]interface{}{1, 2, 3, 4, 5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var mysomeRoot someRoot
	_ = json.Unmarshal(data, &mysomeRoot)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range *someTest1.SomeField {
		if v != (*mysomeRoot.SomeField)[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArrayInterfaceInPointerRootAndPointerValue(t *testing.T) {

	type someRoot struct {
		SomeField *[5]*interface{}
	}

	var value1 interface{} = 1
	var value2 interface{} = 2
	var value3 interface{} = 3
	var value4 interface{} = 4
	var value5 interface{} = nil

	test1 := someRoot{&[5]*interface{}{&value1, &value2, &value3, &value4, &value5}}

	data, _ := testMarshaler(test1)
	fmt.Println(string(data))
	var mysomeRoot someRoot
	_ = json.Unmarshal(data, &mysomeRoot)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range *someTest1.SomeField {
		if v != nil {
			if *v != *(*mysomeRoot.SomeField)[i] {
				t.FailNow()
			}
		} else if v == nil {
			if (*mysomeRoot.SomeField)[i] != nil {
				t.FailNow()
			}
		}

	}
}
func TestAssignThingsPrimitiveArraySliceCrossOver(t *testing.T) {

	t1 := [5]int{1, 2, 3, 4, 5}
	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checkerRoot [5]int

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 []int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v != checkerRoot[i] {
			t.FailNow()
		}
	}
}
func TestAssignThingsPrimitiveSliceArrayCrossOver(t *testing.T) {

	t1 := []int{1, 2, 3, 4, 5}
	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checkerRoot []int

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 [5]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if v != checkerRoot[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArrayPointerValues(t *testing.T) {
	v1 := 1
	v2 := 2
	v3 := 3
	v4 := 4
	v5 := 5
	t1 := [5]*int{&v1, &v2, &v3, &v4, &v5}

	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := testMarshaler(t1)
	fmt.Println(string(data))
	var checkerRoot [5]*int

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someTest1 [5]*int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someTest1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i, v := range someTest1 {
		if *v != *checkerRoot[i] {
			t.FailNow()
		}
	}
}

func TestNestedPointerResolve(t *testing.T) {
	var somePtr *********int
	var resultPtr reflect.Value = reflect.ValueOf(&somePtr)
	fmt.Println(resultPtr.Kind())
	resultTye := reflect.TypeOf(somePtr)
	value := 10
	resultValue := reflect.ValueOf(value)

	numberOfPointer := 0
	// get number of pointers
	for resultTye.Kind() == reflect.Pointer {
		resultTye = resultTye.Elem()
		numberOfPointer += 1
	}

	var tmpPtr reflect.Value
	for ; numberOfPointer > 0; numberOfPointer-- {
		tmpPtr = reflect.New(resultValue.Type()) // var tmpPtr *resultValueType
		tmpPtr.Elem().Set(resultValue)           // *tmpPtr = resultValue
		resultValue = tmpPtr
	}
	resultPtr.Elem().Set(resultValue)

	if *********somePtr != 10 {
		t.FailNow()
	}

	content, err := testMarshaler(somePtr)
	if err != nil {
		t.FailNow()
	}
	if string(content) != "10" {
		t.FailNow()
	}
	var someReceiver *********int
	err = json.Unmarshal(content, &someReceiver)
	if err != nil {
		t.FailNow()
	}
	if *********somePtr != 10 {
		t.FailNow()
	}

	// use our own

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(content))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()

	var someReceiver2 *********int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if *********someReceiver2 != 10 {
		t.FailNow()
	}

}

func TestBareStruct(t *testing.T) {
	type someStruct struct {
		Age int16
	}
	var somePtr someStruct

	data, _ := testMarshaler(somePtr)
	var checker someStruct
	err := json.Unmarshal(data, &checker)
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

	var someReceiver2 someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someReceiver2.Age != checker.Age {
		t.FailNow()
	}
}
func TestPtrStruct(t *testing.T) {
	type someStruct struct {
		Age int16
	}
	var somePtr *someStruct

	data, _ := testMarshaler(somePtr)
	var checker someStruct
	err := json.Unmarshal(data, &checker)
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

	var someReceiver2 *someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someReceiver2.Age != checker.Age {
		t.FailNow()
	}
}
func TestCustomizeType(t *testing.T) {
	str := "QUJD"
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	fmt.Println(data) // Output: [65 66 67]

	type MyType uint8
	type someStruct struct {
		Age    MyType
		Ages   []MyType // base64 string
		AgeMap map[MyType][]MyType
		Ages2  [2]MyType
	}
	var somePtr *someStruct = &someStruct{
		12, []MyType{1, 2, 3},
		map[MyType][]MyType{1: []MyType{4, 5, 6}, 2: []MyType{7, 8, 9}, 3: []MyType{10, 11, 12}},
		[2]MyType{13, 14}}

	data, _ = testMarshaler(somePtr)
	var checker someStruct
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

	var someReceiver2 someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someReceiver2.Age != checker.Age {
		t.FailNow()
	}
	for i, v := range somePtr.Ages {
		if checker.Ages[i] != v {
			t.FailNow()
		}
	}
	for i, v := range somePtr.Ages2 {
		if checker.Ages2[i] != v {
			t.FailNow()
		}
	}
	for k, arr1 := range somePtr.AgeMap {
		arr2 := somePtr.AgeMap[k]
		for i, v := range arr1 {
			if arr2[i] != v {
				t.FailNow()
			}
		}
	}
}
func TestMapSliceType(t *testing.T) {

	var someMap []map[string][]interface{} = []map[string][]interface{}{
		map[string][]interface{}{
			"1": []interface{}{1, 2, 3, 4, 5},
		},
		map[string][]interface{}{
			"2": []interface{}{true, false, true, false},
		},
		map[string][]interface{}{
			"3": []interface{}{"hello", "world"},
		},
		map[string][]interface{}{
			"4": []interface{}{nil, "hello", 1, true},
		},
	}

	data, _ := testMarshaler(someMap)
	var checker []map[string][]interface{}
	err := json.Unmarshal(data, &checker)
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

	var someReceiver2 []map[string][]interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for k, v := range someReceiver2 {
		m2 := checker[k]
		for k1, v1 := range v {
			for i, item := range v1 {
				if m2[k1][i] != item {
					t.FailNow()
				}
			}
		}
	}
}

func TestEmptyObject(t *testing.T) {

	type something struct {
		Data map[string]Bro
	}
	var t1 something = something{make(map[string]Bro)}
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}
	var t2 something
	err = json.Unmarshal(data, &t2)
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
	var someReceiver2 something
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(someReceiver2.Data) != 0 {
		t.FailNow()
	}
}
func TestEmptySlice(t *testing.T) {

	type something struct {
		Data []int
	}
	var t1 something = something{make([]int, 0)}
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}
	var t2 something
	err = json.Unmarshal(data, &t2)
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
	var someReceiver2 something
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(someReceiver2.Data) != 0 {
		t.FailNow()
	}
}

func TestEmptyArray(t *testing.T) {

	type something struct {
		Data [0]int
	}
	var t1 something
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}
	var t2 something
	err = json.Unmarshal(data, &t2)
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
	var someReceiver2 something
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(someReceiver2.Data) != 0 {
		t.FailNow()
	}
}

func TestInterfaceNil(t *testing.T) {
	var someinterface interface{}
	data, err := testMarshaler(someinterface)
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
	var someinterface2 interface{}
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2 != nil {
		t.FailNow()
	}
}

func TestPtrSliceNilCaseInMap(t *testing.T) {
	var someinterface map[string]*[]int = map[string]*[]int{
		"a": nil,
		"b": &[]int{0, 1, 2, 3},
	}
	data, err := testMarshaler(someinterface)
	if err != nil {
		t.FailNow()
	}

	var checker map[string]*[]int
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
	var someinterface2 map[string]*[]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2["a"] != nil {
		t.FailNow()
	}
	for k, v := range *checker["b"] {
		if (*someinterface2["b"])[k] != v {
			t.FailNow()
		}
	}
}
func TestPtrMapNilCaseInMap(t *testing.T) {
	var someinterface map[string]*map[int]int = map[string]*map[int]int{
		"a": nil,
		"b": &map[int]int{0: 0, 1: 1, 2: 2, 3: 3},
	}
	data, err := testMarshaler(someinterface)
	if err != nil {
		t.FailNow()
	}

	var checker map[string]*map[int]int
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
	var someinterface2 map[string]*map[int]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2["a"] != nil {
		t.FailNow()
	}
	for k, v := range *checker["b"] {
		if (*someinterface2["b"])[k] != v {
			t.FailNow()
		}
	}
}
func TestPtrSliceNilCaseInSlice(t *testing.T) {
	var someinterface []*[]int = []*[]int{
		nil,
		&[]int{0, 1, 2, 3},
	}
	data, err := testMarshaler(someinterface)
	if err != nil {
		t.FailNow()
	}

	var checker []*[]int
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
	var someinterface2 []*[]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2[0] != nil {
		t.FailNow()
	}
	for k, v := range *checker[1] {
		if (*someinterface2[1])[k] != v {
			t.FailNow()
		}
	}
}

func TestSliceOfPtrToCollections(t *testing.T) {
	type someStruct struct {
		Number int
	}
	var someinterface []*[]someStruct = []*[]someStruct{
		nil,
		&[]someStruct{{1}, {2}},
	}
	data, err := testMarshaler(someinterface)
	if err != nil {
		t.FailNow()
	}

	var checker []*[]someStruct
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
	var someinterface2 []*[]someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2[0] != nil {
		t.FailNow()
	}
	for k, v := range *checker[1] {
		if (*someinterface2[1])[k].Number != v.Number {
			t.FailNow()
		}
	}
}

func TestInterfaceNilOnMapInterface(t *testing.T) {
	var someinterface map[string]interface{} = map[string]interface{}{
		"a": nil,
		"b": &[]int{0, 1, 2, 3},
	}
	data, err := testMarshaler(someinterface)
	if err != nil {
		t.FailNow()
	}

	var checker map[string]*[]int
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
	var someinterface2 map[string]*[]int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2["a"] != nil {
		t.FailNow()
	}
	for k, v := range *checker["b"] {
		if (*someinterface2["b"])[k] != v {
			t.FailNow()
		}
	}
}
func TestEmptyString(t *testing.T) {
	var someinterface string
	data, err := testMarshaler(someinterface)
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
	var someinterface2 string
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2 != "" {
		t.FailNow()
	}
}
func TestString(t *testing.T) {
	var someinterface string = "123"
	data, err := testMarshaler(someinterface)
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
	var someinterface2 string
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2 != "123" {
		t.FailNow()
	}
}

func TestBool(t *testing.T) {
	var someinterface bool = false
	data, err := testMarshaler(someinterface)
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
	var someinterface2 bool
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2 != false {
		t.FailNow()
	}
}

func TestNumber(t *testing.T) {
	var someinterface int = 1234567
	data, err := testMarshaler(someinterface)
	if err != nil {
		t.FailNow()
	}
	var someInt int
	_ = json.Unmarshal(data, &someInt)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var someinterface2 int
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someinterface2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someinterface2 != someInt {
		t.FailNow()
	}
}
func TestFinalExam(t *testing.T) {

	type SomeStruct1 struct {
		Name1  string                 //checked
		Name2  []int                  //checked
		Name3  map[string]int         //checked
		Name4  []interface{}          //checked
		Name5  []Bro                  //checked
		Name6  []*Bro                 //checked
		Name7  Bro                    //checked
		Name8  *Bro                   //checked not fill in, let it nil
		Name9  map[string]interface{} //checked
		Name10 map[int]Bro            //checked
		Name11 [3]int                 //checked

		// ... and so on for other cases
		Name14 []map[string][]interface{}
		Name15 interface{} //checked covert to map[string]interface{}
		Name16 *Bro
		Name17 map[string]Bro
		Name18 []int
		Name19 MyType
	}

	var test1 SomeStruct1 = SomeStruct1{
		Name1: "name1",
		Name2: []int{1, 2, 3},
		Name3: map[string]int{"hello": 123},
		Name4: []interface{}{1, false, 1.23, nil, []int{2, 3, 4}, map[string]int{"world": 223}},
		Name5: []Bro{
			Bro{Name: "Ann", Age: 12}, Bro{Name: "Ken", Age: 13},
		},
		Name6: []*Bro{
			&Bro{Name: "Ann2", Age: 121}, &Bro{Name: "Ken2", Age: 131},
		},
		Name9: map[string]interface{}{
			"First":  1,
			"Second": true,
			"Third":  3.2,
			"Fourth": Bro{Name: "Ann3", Age: 312},
			"Fifth":  &Bro{Name: "Ann2", Age: 421},
			"Sixth":  nil,
		},
		Name10: map[int]Bro{
			11: Bro{Name: "Ann311", Age: 3112},
			12: Bro{Name: "Ann222", Age: 3112222},
		},
		Name11: [3]int{991, 992, 993},

		Name14: []map[string][]interface{}{
			{
				"NestedKey": {1, "string", false, []int{1, 2, 3}},
			},
		},
		Name15: &ImplementingStruct{Field: "Implemented!"},
		Name16: (*Bro)(nil),          // nil pointer to struct
		Name17: make(map[string]Bro), // non-nil empty map
		Name18: []int{},              // non-nil empty slice
		Name19: MyType(123),
	}

	data, err := testMarshaler(test1)
	if err != nil {
		t.FailNow()
	}

	var checker SomeStruct1
	err = json.Unmarshal(data, &checker)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	fd, err := os.OpenFile(`d:\test1.txt`, os.O_CREATE, 0664)
	if err != nil {
		t.FailNow()
	}
	defer fd.Close()
	_, err = fd.Write(data)
	if err != nil {
		t.FailNow()
	}

	err = sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var someReceiver2 SomeStruct1
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if checker.Name1 != someReceiver2.Name1 {
		t.FailNow()
	}
	for i, v := range checker.Name2 {
		if someReceiver2.Name2[i] != v {
			t.FailNow()
		}
	}
	for k, v := range checker.Name3 {
		if someReceiver2.Name3[k] != v {
			t.FailNow()
		}
	}
	//[]interface{}{1, false, 1.23, nil, []int{2, 3, 4}, map[string]int{"world": 223}},
	for i := 0; i < 4; i++ {
		if someReceiver2.Name4[i] != checker.Name4[i] {
			t.FailNow()
		}
	}
	for i, v := range someReceiver2.Name4[4].([]interface{}) {
		if checker.Name4[4].([]interface{})[i] != v {
			t.FailNow()
		}
	}
	for i, v := range someReceiver2.Name4[5].(map[string]interface{}) {
		if checker.Name4[5].(map[string]interface{})[i] != v {
			t.FailNow()
		}
	}
	for i, v := range checker.Name5 {
		myItem := someReceiver2.Name5[i]
		if v.Name != myItem.Name {
			t.FailNow()
		}
		if v.Age != myItem.Age {
			t.FailNow()
		}
	}
	for i, v := range checker.Name6 {
		myItem := someReceiver2.Name6[i]
		if v.Name != myItem.Name {
			t.FailNow()
		}
		if v.Age != myItem.Age {
			t.FailNow()
		}
	}
	if checker.Name7.Name != someReceiver2.Name7.Name {
		t.FailNow()
	}
	if checker.Name7.Age != someReceiver2.Name7.Age {
		t.FailNow()
	}
	if someReceiver2.Name8 != nil && checker.Name8 != someReceiver2.Name8 {
		t.FailNow()
	}

	if checker.Name9["First"] != someReceiver2.Name9["First"] {
		t.FailNow()
	}
	if checker.Name9["Second"] != someReceiver2.Name9["Second"] {
		t.FailNow()
	}
	if checker.Name9["Third"] != someReceiver2.Name9["Third"] {
		t.FailNow()
	}
	if checker.Name9["Sixth"] != someReceiver2.Name9["Sixth"] && someReceiver2.Name9["Sixth"] != nil {
		t.FailNow()
	}
	if checker.Name9["Fourth"].(map[string]interface{})["Name"] != someReceiver2.Name9["Fourth"].(map[string]interface{})["Name"] {
		t.FailNow()
	}
	if checker.Name9["Fourth"].(map[string]interface{})["Age"] != someReceiver2.Name9["Fourth"].(map[string]interface{})["Age"] {
		t.FailNow()
	}
	if checker.Name9["Fifth"].(map[string]interface{})["Name"] != someReceiver2.Name9["Fifth"].(map[string]interface{})["Name"] {
		t.FailNow()
	}
	if checker.Name9["Fifth"].(map[string]interface{})["Age"] != someReceiver2.Name9["Fifth"].(map[string]interface{})["Age"] {
		t.FailNow()
	}
	if checker.Name10[11].Name != someReceiver2.Name10[11].Name {
		t.FailNow()
	}
	if checker.Name10[11].Age != someReceiver2.Name10[11].Age {
		t.FailNow()
	}
	if checker.Name10[12].Name != someReceiver2.Name10[12].Name {
		t.FailNow()
	}
	if checker.Name10[11].Age != someReceiver2.Name10[11].Age {
		t.FailNow()
	}
	for i, v := range checker.Name11 {
		if someReceiver2.Name11[i] != v {
			t.FailNow()
		}
	}
	for i, v := range checker.Name14 {
		m2 := someReceiver2.Name14[i]
		for k1, v1 := range v {
			arr2 := m2[k1]
			for j, e := range v1 {
				if reflect.TypeOf(e).Kind() != reflect.Array && reflect.TypeOf(e).Kind() != reflect.Slice {
					if arr2[j] != e {
						t.FailNow()
					}
				} else {
					for m, n := range e.([]interface{}) {
						if arr2[j].([]interface{})[m] != n {
							t.FailNow()
						}
					}
				}

			}
		}
	}

	if checker.Name15.(map[string]interface{})["Field"] != someReceiver2.Name15.(map[string]interface{})["Field"] {
		t.FailNow()
	}
	if checker.Name16 != someReceiver2.Name16 && someReceiver2.Name16 != nil {
		t.FailNow()
	}
	if len(checker.Name17) != len(someReceiver2.Name17) && len(someReceiver2.Name17) != 0 {
		t.FailNow()
	}
	if len(checker.Name18) != len(someReceiver2.Name18) && len(someReceiver2.Name18) != 0 {
		t.FailNow()
	}
	if checker.Name19 != someReceiver2.Name19 && someReceiver2.Name19 != 123 {
		t.FailNow()
	}

}

func TestVariable(t *testing.T) {
	sample := `{
		"hello1":${myvariable1},
		"hello2":${myvariable2},
		"hello3":${myvariable3},
		"hello4":${myvariable4},
		"hello5":${myvariable5},
		"hello6":${myvariable6},
		"hello7":"hey man: ${myvariable7}",
		"hello8":"hey man! ${myvariable8}",
		"${myvariable9}":${myvariable9},
		"hello10":"${myvariable10}"
		}`
	var variables map[string]interface{} = map[string]interface{}{
		"myvariable1": 1,
		"myvariable2": "123",
		"myvariable3": true,
		"myvariable4": nil,
		"myvariable5": []interface{}{1, 2, 3},
		"myvariable6": map[string]interface{}{"happy": "Cat"},
		"myvariable7": "whats up!",
		"myvariable8": "",
		"myvariable9": "hello9",
	}

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(strings.NewReader(sample))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var someReceiver2 map[string]interface{}
	err = interpreter.UnmarshallAST(node, variables, testMarshaler, testUnmarshaler, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if someReceiver2["hello1"] != 1 {
		t.FailNow()
	}

	if someReceiver2["hello2"] != "123" {
		t.FailNow()
	}

	if someReceiver2["hello3"] != true {
		t.FailNow()
	}

	if someReceiver2["hello4"] != nil {
		t.FailNow()
	}
	for i, v := range someReceiver2["hello5"].([]interface{}) {
		if i+1 != v {
			t.FailNow()
		}
	}
	if someReceiver2["hello6"].(map[string]interface{})["happy"] != "Cat" {
		t.FailNow()
	}
	if someReceiver2["hello7"] != "hey man: whats up!" {
		t.FailNow()
	}
	if someReceiver2["hello8"] != "hey man! " {
		t.FailNow()
	}
	if someReceiver2["hello9"] != "hello9" {
		t.FailNow()
	}
	if someReceiver2["hello10"] != "${myvariable10}" {
		t.FailNow()
	}
}

func TestFieldByTag(t *testing.T) {
	type mydata struct {
		Name1   string `json:"name"`
		Name2   int    `json:"age"`
		Address string
	}

	test1 := mydata{"Ann", 198, "CA Redwood shore"}
	data, _ := testMarshaler(test1)

	var checker mydata
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver mydata
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Name1 != "Ann" {
		t.FailNow()
	}
	if myReceiver.Name2 != 198 {
		t.FailNow()
	}
	if myReceiver.Address != "CA Redwood shore" {
		t.FailNow()
	}

}

func TestEmbeddingFields(t *testing.T) {
	type mydata1 struct {
		Name1 string `json:"name"`
	}
	type database struct {
		Name2 int `json:"age"`
	}
	type mydata2 struct {
		database
		Address string
	}

	type newTest struct {
		mydata1
		mydata2
		Home string
	}

	testdata1 := mydata1{"Ann"}
	testdata2 := mydata2{database: database{198}, Address: "CA Redwood shore"}
	var test1 newTest = newTest{
		mydata1: testdata1,
		mydata2: testdata2,
		Home:    "US",
	}
	data, _ := testMarshaler(test1)

	var checker newTest
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver newTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Name1 != "Ann" {
		t.FailNow()
	}
	if myReceiver.Name2 != 198 {
		t.FailNow()
	}
	if myReceiver.Address != "CA Redwood shore" {
		t.FailNow()
	}
	if myReceiver.Home != "US" {
		t.FailNow()
	}
}

type mybool bool

var _ json.Unmarshaler = (*mybool)(nil)

func (m *mybool) UnmarshalJSON(data []byte) error {
	// your unmarshalling logic here
	if data[0] == 'f' {
		*m = false
	} else {
		*m = true
	}
	return nil
}

func TestCusomUnmarshal(t *testing.T) {

	type mydata struct {
		T1 *mybool
	}

	var someItem mybool = false
	var test1 mydata = mydata{T1: &someItem}
	data, _ := testMarshaler(test1)

	var checker mydata
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver mydata
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if *myReceiver.T1 != false {
		t.FailNow()
	}

}

type myslice []int

var _ json.Unmarshaler = (*myslice)(nil)

func (m *myslice) UnmarshalJSON(data []byte) error {
	// your unmarshalling logic here
	*m = append(*m, 123)
	return nil
}
func TestCusomUnmarshalSlice(t *testing.T) {

	type mydata struct {
		T1 myslice
	}

	var someItem myslice = []int{'1', '2', 3}
	var test1 mydata = mydata{T1: someItem}
	data, _ := testMarshaler(test1)

	var checker myslice = make(myslice, 0)
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver mydata
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

}

type mymap map[string]int

var _ json.Unmarshaler = (mymap)(nil)

func (m mymap) UnmarshalJSON(data []byte) error {
	// your unmarshalling logic here
	m["abc"] = 123
	return nil
}
func TestCusomUnmarshalMap(t *testing.T) {

	type mydata struct {
		T1 *mymap
	}

	var someItem mymap = mymap{"123": 1}
	var test1 mydata = mydata{T1: &someItem}
	data, _ := testMarshaler(test1)

	var checker mymap = make(mymap)
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver mydata
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if (*myReceiver.T1)["abc"] != 123 {
		t.FailNow()
	}

}
func TestCusomUnmarshalMapNonePointer(t *testing.T) {

	type mydata struct {
		T1 mymap
	}

	var someItem mymap = mymap{"123": 1}
	var test1 mydata = mydata{T1: someItem}
	data, _ := testMarshaler(test1)

	var checker mymap = make(mymap)
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver mydata
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.T1["abc"] != 123 {
		t.FailNow()
	}

}

type myUnmarshalstruct struct {
	Name string
	Age  int
}

var _ json.Unmarshaler = (*myUnmarshalstruct)(nil)

func (m *myUnmarshalstruct) UnmarshalJSON(b []byte) error {
	m.Name = string(b)
	m.Age = 123
	return nil
}
func TestCusomUnmarshalStruct(t *testing.T) {

	var test1 myUnmarshalstruct = myUnmarshalstruct{"Kenny", 123}
	data, _ := testMarshaler(test1)

	var checker myUnmarshalstruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver myUnmarshalstruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Name != checker.Name {
		t.FailNow()
	}
	if myReceiver.Age != checker.Age {
		t.FailNow()
	}

}

func TestCusomUnmarshalStructInStructField(t *testing.T) {

	type someStruct struct {
		Field1 myUnmarshalstruct
		Field2 int
	}
	var field1 myUnmarshalstruct = myUnmarshalstruct{"Kenny", 123}
	var test1 someStruct = someStruct{field1, 999}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Field1.Name != checker.Field1.Name {
		t.FailNow()
	}
	if myReceiver.Field1.Age != 123 {
		t.FailNow()
	}
	if myReceiver.Field2 != 999 {
		t.FailNow()
	}

}
func TestCusomUnmarshalStructInStructPointerField(t *testing.T) {

	type someStruct struct {
		Field1 *myUnmarshalstruct
		Field2 int
	}
	var field1 myUnmarshalstruct = myUnmarshalstruct{"Kenny", 123}
	var test1 someStruct = someStruct{&field1, 999}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Field1.Name != checker.Field1.Name {
		t.FailNow()
	}
	if myReceiver.Field1.Age != 123 {
		t.FailNow()
	}
	if myReceiver.Field2 != 999 {
		t.FailNow()
	}

}

type myString string

func (m *myString) UnmarshalJSON(b []byte) error {
	*m = myString(fmt.Sprintf("%s:%s", "abc", string(b)))
	return nil
}
func TestCusomUnmarshalStringInStructField(t *testing.T) {

	type someStruct struct {
		Field1 myString // this won't care if the field is a pointer or not, as long as it has a pointer receiver unmarshaler, it will change the value
	}
	var field1 myString = myString("123")
	var test1 someStruct = someStruct{field1}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Field1 != checker.Field1 {
		t.FailNow()
	}

}
func TestCusomUnmarshalStringInStructPointerField(t *testing.T) {

	type someStruct struct {
		Field1 *myString // this won't care if the field is a pointer or not, as long as it has a pointer receiver unmarshaler, it will change the value
	}
	var field1 myString = myString("123")
	var test1 someStruct = someStruct{&field1}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if *myReceiver.Field1 != *checker.Field1 {
		t.FailNow()
	}

}

type myNumber int

func (m *myNumber) UnmarshalJSON(b []byte) error {
	*m = myNumber(len(b))
	return nil
}
func TestCusomUnmarshalIntInStructField(t *testing.T) {

	type someStruct struct {
		Field1 myNumber // this won't care if the field is a pointer or not, as long as it has a pointer receiver unmarshaler, it will change the value
	}
	var field1 myNumber = myNumber(123)
	var test1 someStruct = someStruct{field1}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Field1 != checker.Field1 {
		t.FailNow()
	}

}
func TestCusomUnmarshalIntInStructPointerField(t *testing.T) {

	type someStruct struct {
		Field1 *myNumber // this won't care if the field is a pointer or not, as long as it has a pointer receiver unmarshaler, it will change the value
	}
	var field1 myNumber = myNumber(123)
	var test1 someStruct = someStruct{&field1}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if *myReceiver.Field1 != *checker.Field1 {
		t.FailNow()
	}

}

type myNil struct {
	Age int
}

func (m *myNil) UnmarshalJSON(b []byte) error {
	m.Age = len(b) * 10
	return nil
}
func TestCusomUnmarshalNullInStructField(t *testing.T) {

	type someStruct struct {
		Field1 *myNil // this won't care if the field is a pointer or not, as long as it has a pointer receiver unmarshaler, it will change the value
	}
	var test1 someStruct = someStruct{nil}
	data, _ := testMarshaler(test1)

	var checker someStruct
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver someStruct
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	// when a value is already nil, the unmarshaler won't help you to convert it to something not nil
	if myReceiver.Field1 != checker.Field1 && myReceiver.Field1 != nil {
		t.FailNow()
	}

}

func TestNameCollision(t *testing.T) {
	type mydata1 struct {
		Name1 string `json:"name1"`
	}
	type database struct {
		mydata1
		Name2 int `json:"name1"`
	}
	type mydata2 struct {
		database
		Address string `json:"name2"`
	}

	type newTest struct {
		mydata2
		Home string `json:"name2"`
	}

	testdata1 := mydata1{"Ann"}
	testdata2 := mydata2{database: database{testdata1, 198}, Address: "CA Redwood shore"}
	var test1 newTest = newTest{
		mydata2: testdata2,
		Home:    "US",
	}
	data, _ := testMarshaler(test1)

	var checker newTest
	_ = json.Unmarshal(data, &checker)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))
	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTBuilder().GetAST()
	var myReceiver newTest
	err = interpreter.UnmarshallAST(node, nil, testMarshaler, testUnmarshaler, &myReceiver)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if myReceiver.Home != checker.Home {
		t.FailNow()
	}
	if myReceiver.Name1 != checker.Name1 {
		t.FailNow()
	}
	if myReceiver.Address != checker.Address {
		t.FailNow()
	}
	if myReceiver.Name2 != checker.Name2 {
		t.FailNow()
	}
}

func TestNameConflic(t *testing.T) {
	type Example struct {
		F2 string
		F1 string `json:"F2"`
	}

	var t1 Example = Example{"hello", "world"}
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}

	var checker Example
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
}

func TestNameConflic2(t *testing.T) {
	type Example2 struct {
		F1 string `json:"F2"`
	}
	type Example struct {
		Example2
		F2 string
	}

	var t1 Example = Example{Example2{"hello"}, "world"}
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}

	var checker Example
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
}

func TestNameConflic3(t *testing.T) {

	type Example struct {
		F1 string `json:",omitempty"`
		F2 string `json:"F1,omitempty"`
	}

	var t1 Example = Example{"hello", "world"}
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}

	var checker Example
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
}

func TestNameConflic4(t *testing.T) {

	type Example struct {
		FX string `json:"F1,omitempty"`
		F1 string
	}

	var t1 Example = Example{"hello", "world"}
	data, err := testMarshaler(t1)
	if err != nil {
		t.FailNow()
	}

	var checker Example
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
}

func TestNameConflic5(t *testing.T) {

	type Example1 struct {
		FX   string `json:"F1,omitempty"`
		FX23 string `json:"F1,omitempty"`
	}
	type Example struct {
		Example1
		Name int
	}

	data := []byte(`{"F1":"Hello", "Name":10}`)
	var checker Example
	err := json.Unmarshal(data, &checker)
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
}

func TestStringOption2(t *testing.T) {

	type Example struct {
		F1 string `json:",string"`
		F2 int    `json:",string"`
		F3 bool   `json:",string"`
	}

	ex := Example{"hello", 123, true}
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
