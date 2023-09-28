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
	if someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if someTest1.World != 100.123 {
		t.FailNow()
	}
	if someTest1.World2 != 100 {
		t.FailNow()
	}
	if someTest1.Apple != true {
		t.FailNow()
	}
	if someTest1.Banana != false {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != "1234" {
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

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
	if someTest1["Hello"] != "Peter" {
		t.FailNow()
	}
	if someTest1["World"] != 101.123 { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1["World2"] != 100.0 {
		t.FailNow()
	}
	if someTest1["Apple"] != true {
		t.FailNow()
	}
	if someTest1["Banana"] != false {
		t.FailNow()
	}
	if someTest1["Something"] != nil {
		t.FailNow()
	}
	if someTest1["SomethingNotNil"] != "1234" {
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
	if someTest1[1] != "Peter" {
		t.FailNow()
	}
	if someTest1[2] != 101.123 { // default to ensure int when interface{} cannot tell us what the type is
		t.FailNow()
	}
	if someTest1[3] != 100.0 {
		t.FailNow()
	}
	if someTest1[4] != true {
		t.FailNow()
	}
	if someTest1[5] != false {
		t.FailNow()
	}
	if someTest1[6] != nil {
		t.FailNow()
	}
	if someTest1[7] != "1234" {
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
	if someTest1[1] != "Peter" {
		t.FailNow()
	}
	if someTest1[2] != 101.123 { // default to ensure int when interface{} cannot tell us what the type is
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

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
	if someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if someTest1.World != 100.123 {
		t.FailNow()
	}
	if someTest1.World2 != 100 {
		t.FailNow()
	}
	if someTest1.Apple != true {
		t.FailNow()
	}
	if someTest1.Banana != false {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != "1234" {
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

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
	if someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if someTest1.World != 100.123 {
		t.FailNow()
	}
	if someTest1.World2 != 100 {
		t.FailNow()
	}
	if someTest1.Apple != true {
		t.FailNow()
	}
	if someTest1.Banana != false {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != "1234" {
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

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
	if someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if someTest1.World != 100.123 {
		t.FailNow()
	}
	if someTest1.World2 != 100 {
		t.FailNow()
	}
	if someTest1.Apple != true {
		t.FailNow()
	}
	if someTest1.Banana != false {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != "1234" {
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

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
	if someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if someTest1.World != 100.123 {
		t.FailNow()
	}
	if someTest1.World2 != 100 {
		t.FailNow()
	}
	if someTest1.Apple != true {
		t.FailNow()
	}
	if someTest1.Banana != false {
		t.FailNow()
	}
	if someTest1.Something != nil {
		t.FailNow()
	}
	if someTest1.SomethingNotNil != "1234" {
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
	if *someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if *someTest1.World != 100.123 {
		t.FailNow()
	}
	if *someTest1.World2 != 100 {
		t.FailNow()
	}
	if *someTest1.Apple != true {
		t.FailNow()
	}
	if *someTest1.Banana != false {
		t.FailNow()
	}
	if *someTest1.Something != nil {
		t.FailNow()
	}
	if *someTest1.SomethingNotNil != "1234" {
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
	if *someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if *someTest1.World != 100.123 {
		t.FailNow()
	}
	if *someTest1.World2 != 100 {
		t.FailNow()
	}
	if *someTest1.Apple != true {
		t.FailNow()
	}
	if *someTest1.Banana != false {
		t.FailNow()
	}
	if *someTest1.Something != nil {
		t.FailNow()
	}
	if *someTest1.SomethingNotNil != "1234" {
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
	if *someTest1.Hello != "Peter" {
		t.FailNow()
	}
	if *someTest1.World != 100.123 {
		t.FailNow()
	}
	if *someTest1.World2 != 100 {
		t.FailNow()
	}
	if *someTest1.Apple != true {
		t.FailNow()
	}
	if *someTest1.Banana != false {
		t.FailNow()
	}
	if *someTest1.Something != nil {
		t.FailNow()
	}
	if *someTest1.SomethingNotNil != "1234" {
		t.FailNow()
	}
}

func TestAssignThingsPrimitiveSlice(t *testing.T) {

	t1 := []int{1, 2, 3, 4, 5}
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
		if v != t1[i] {
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
		if v != test1.SomeField[i] {
			t.FailNow()
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
		if v != (*test1.SomeField)[i] {
			t.FailNow()
		}
	}
}
func TestAssignThingsSliceWithInterfaceElement(t *testing.T) {

	t1 := []interface{}{1, true, false, nil}
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

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
		if v != t1[i] {
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
		if v != test1.SomeField[i] {
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
		if v != (*test1.SomeField)[i] {
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
		if *v != *(*test1.SomeField)[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsPrimitiveArray(t *testing.T) {

	t1 := [5]int{1, 2, 3, 4, 5}
	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
		if v != t1[i] {
			t.FailNow()
		}
	}
}

func TestAssignThingsArrayWithInterfaceElement(t *testing.T) {

	t1 := [4]interface{}{1, true, false, nil}
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
		if v != t1[i] {
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
		if v != test1.SomeField[i] {
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
		if v != (*test1.SomeField)[i] {
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
		if v != test1.SomeField[i] {
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
		if v != (*test1.SomeField)[i] {
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
		if *v != *(*test1.SomeField)[i] {
			t.FailNow()
		}
	}
}
func TestAssignThingsPrimitiveArraySliceCrossOver(t *testing.T) {

	t1 := [5]int{1, 2, 3, 4, 5}
	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
		if v != t1[i] {
			t.FailNow()
		}
	}
}
func TestAssignThingsPrimitiveSliceArrayCrossOver(t *testing.T) {

	t1 := []int{1, 2, 3, 4, 5}
	fmt.Println(reflect.TypeOf(t1).Kind())
	data, _ := json.Marshal(t1)
	fmt.Println(string(data))

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
		if v != t1[i] {
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

	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(strings.NewReader(string(data)))
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
		if *v != *t1[i] {
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
