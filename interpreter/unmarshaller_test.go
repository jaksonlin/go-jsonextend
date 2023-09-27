package interpreter_test

import (
	"encoding/json"
	"fmt"
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

func TestAssignThingsPrimitiveArray(t *testing.T) {

	t1 := [5]int{1, 2, 3, 4, 5}
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
