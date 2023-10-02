package interpreter_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	var validator test1
	_ = json.Unmarshal(data, &validator)
	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 test1
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var validator map[string]interface{}
	err := json.Unmarshal(data, &validator)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()
	var someTest1Check map[string]interface{}
	err = json.Unmarshal(data, &someTest1Check)
	if err != nil {
		t.FailNow()
	}
	var someTest1 map[string]interface{}
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(testData)
	fmt.Println(string(data))
	var validator someRoot
	err := json.Unmarshal(data, &validator)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, &rootCheck)
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

	data, _ := json.Marshal(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, &rootCheck)
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

	data, _ := json.Marshal(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, &rootCheck)
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

	data, _ := json.Marshal(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, &rootCheck)
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

	data, _ := json.Marshal(testData)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var rootCheck someRoot
	err = interpreter.UnmarshallAST(node, nil, &rootCheck)
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

	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()
	var someTest1Check map[int]interface{}
	err = json.Unmarshal(data, &someTest1Check)
	if err != nil {
		t.FailNow()
	}
	var someTest1 map[int]interface{}
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()
	var someTest1Check map[int]interface{}
	err = json.Unmarshal(data, &someTest1Check)
	if err != nil {
		t.FailNow()
	}
	var someTest1 map[uint8]interface{}
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, &someRoot)
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
	data, _ := json.Marshal(tr)
	fmt.Println(string(data))

	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, &someRoot)
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
	data, _ := json.Marshal(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, &someRoot)
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

	data, _ := json.Marshal(tr)
	fmt.Println(string(data))
	var checkerRoot testRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someRoot testRoot
	err = interpreter.UnmarshallAST(node, nil, &someRoot)
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

	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	var someRs test1
	err := json.Unmarshal(data, &someRs)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	node := sm.GetASTConstructor().GetAST()

	var someTest1 test1
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(someRoot)
	fmt.Println(string(data))

	var someRs testRoot
	err := json.Unmarshal(data, &someRs)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	node := sm.GetASTConstructor().GetAST()

	var someTestRoot1 testRoot
	err = interpreter.UnmarshallAST(node, nil, &someTestRoot1)
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
	data, _ := json.Marshal(someRoot)
	fmt.Println(string(data))

	var someRs testRoot
	err := json.Unmarshal(data, &someRs)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	node := sm.GetASTConstructor().GetAST()

	var someTestRoot1 testRoot
	err = interpreter.UnmarshallAST(node, nil, &someTestRoot1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var checker []int
	err := json.Unmarshal(data, &checker)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 []int
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

func TestAssignThingsPrimitiveSliceInNonePointerRoot(t *testing.T) {

	type someRoot struct {
		SomeField []int
	}

	test1 := someRoot{[]int{1, 2, 3, 4, 5}}

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, err := json.Marshal(t1)
	if err != nil {
		t.FailNow()
	}

	fmt.Println(string(data))
	var someRootChecker someRoot
	err = json.Unmarshal(data, &someRootChecker)
	if err != nil {
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var validator []interface{}
	_ = json.Unmarshal(data, &validator)

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 []interface{}
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var validator someRoot
	_ = json.Unmarshal(data, &validator)

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var validator someRoot
	_ = json.Unmarshal(data, &validator)

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var validator someRoot
	err := json.Unmarshal(data, &validator)
	if err != nil {
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	var checkerRoot [5]int
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 [5]int
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var checker [4]interface{}
	_ = json.Unmarshal(data, &checker)

	var checkerRoot [4]interface{}
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 [4]interface{}
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))

	var checkerRoot someRoot
	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var checkerRoot someRoot

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var mysomeRoot someRoot
	_ = json.Unmarshal(data, &mysomeRoot)

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var mysomeRoot someRoot
	_ = json.Unmarshal(data, &mysomeRoot)

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	data, _ := json.Marshal(test1)
	fmt.Println(string(data))
	var mysomeRoot someRoot
	_ = json.Unmarshal(data, &mysomeRoot)

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 someRoot
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var checkerRoot [5]int

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 []int
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var checkerRoot []int

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 [5]int
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))
	var checkerRoot [5]*int

	err := json.Unmarshal(data, &checkerRoot)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(data)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someTest1 [5]*int
	err = interpreter.UnmarshallAST(node, nil, &someTest1)
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

	content, err := json.Marshal(somePtr)
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

	sm := tokenizer.NewTokenizerStateMachine()
	err = sm.ProcessData(strings.NewReader(string(content)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	node := sm.GetASTConstructor().GetAST()

	var someReceiver2 *********int
	err = interpreter.UnmarshallAST(node, nil, &someReceiver2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if *********someReceiver2 != 10 {
		t.FailNow()
	}

}
