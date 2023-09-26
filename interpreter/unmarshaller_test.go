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

type SomeInterfaceV2 interface {
	Method() string
}

type ImplementingStruct struct {
	Field string
}

func (i *ImplementingStruct) Method() string {
	return i.Field
}

type SomeStruct struct {
	Name1  string                 `json:"name1"`
	Name2  []int                  `json:"name2"`
	Name3  map[string]int         `json:"name3"` // coded
	Name4  []interface{}          `json:"name4"`
	Name5  []Bro                  `json:"name5"`
	Name6  []*Bro                 `json:"name6"`
	Name7  Bro                    `json:"name7"`  //coded
	Name8  *Bro                   `json:"name8"`  //coded
	Name9  map[string]interface{} `json:"name9"`  //coded
	Name10 map[int]Bro            `json:"name10"` //go also not support unmarshal, json now allow number as key,
	Name11 [3]int                 `json:"name11"`
	Name12 MyInterface            `json:"name12"` // pointer
	Name13 MyInterface            `json:"name13"` // struct
	// ... and so on for other cases
	Name14 []map[string][]interface{} `json:"name14"`
	Name15 interface{}                `json:"name15"`
	Name16 *Bro                       `json:"name16"` //coded, for nil, the kv pair len==0, no additional code
	Name17 map[string]Bro             `json:"name17"`
	Name18 []int                      `json:"name18"`
	Name19 MyType                     `json:"name19"`
}

var (
	test1 SomeStruct = SomeStruct{
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
			"Thrid":  3.2,
			"Fourth": Bro{Name: "Ann3", Age: 312},
			"Fifth":  &Bro{Name: "Ann2", Age: 421},
		},
		Name10: map[int]Bro{
			11: Bro{Name: "Ann311", Age: 3112},
			12: Bro{Name: "Ann222", Age: 3112222},
		},
		Name11: [3]int{991, 992, 993},
		Name12: NestedStruct{
			Field1: 100,
			Field2: map[string]Bro{
				"Fourth": Bro{Name: "Ann3-122", Age: 1134312},
			}},
		Name13: &Bro{Name: "Ann3-122", Age: 1134312},
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
)

func TestInspaceOfGoStruct(t *testing.T) {
	var anItem interface{} = &test1
	typeItem := reflect.TypeOf(anItem)
	fmt.Println(typeItem)
	valueItem := reflect.ValueOf(anItem)
	fmt.Println(valueItem)

	switch valueItem.Kind() {
	case reflect.Struct:
		for i := 0; i < valueItem.NumField(); i++ {
			f := valueItem.Field(i)
			fmt.Println(f.Kind())
		}
	case reflect.Pointer:
		typeItemElm := typeItem.Elem()
		valueItemElm := valueItem.Elem()
		fmt.Println(valueItemElm.Kind())
		for i := 0; i < typeItemElm.NumField(); i++ {
			f := typeItemElm.Field(i)
			v := valueItemElm.Field(i)
			fmt.Printf("Field Name: %s\n", f.Name)         // Print field name
			fmt.Printf("Field Value: %v\n", v.Interface()) // Print field value
			fmt.Printf("Field Kind: %v\n", v.Kind())       // Print field kind
			fmt.Printf("Field Tag: %v\n", f.Tag)           // Print field tag
			if jsonTag, ok := f.Tag.Lookup("json"); ok {
				fmt.Printf("Field JSON Tag: %s\n", jsonTag)
			}
		}
	}
	fmt.Println("END")
}
func TestInspaceOfGoStruct2(t *testing.T) {
	var anItem interface{} = &test1
	typeItem := reflect.TypeOf(anItem)
	var myAnItem reflect.Value
	var myAnItemType reflect.Type
	if typeItem.Kind() == reflect.Pointer {
		myAnItemType = typeItem.Elem()
		myAnItem = reflect.New(typeItem.Elem())
	} else {
		myAnItemType = typeItem
		myAnItem = reflect.New(typeItem)
	}
	fmt.Println(myAnItemType)
	// new must return ptr
	myAnItemElm := myAnItem.Elem()
	fmt.Println(myAnItemElm)
	for i := 0; i < myAnItemType.NumField(); i++ {
		f := myAnItemType.Field(i)
		v := myAnItemElm.Field(i)
		fmt.Printf("Field Name: %s\n", f.Name)         // Print field name
		fmt.Printf("Field Value: %v\n", v.Interface()) // Print field value
		fmt.Printf("Field Kind: %v\n", v.Kind())       // Print field kind
		if f.Type.Kind() == reflect.Array || f.Type.Kind() == reflect.Slice || f.Type.Kind() == reflect.Pointer {
			fmt.Printf("element type: %v\n", f.Type.Elem().Kind())
		} else if f.Type.Kind() == reflect.Map {
			fmt.Printf("map key type: %v\n", f.Type.Key())
			fmt.Printf("map element type: %v\n", f.Type.Elem().Kind())
		}
		fmt.Printf("Field Tag: %v\n", f.Tag) // Print field tag
		if jsonTag, ok := f.Tag.Lookup("json"); ok {
			fmt.Printf("Field JSON Tag: %s\n", jsonTag)
		}
		fmt.Println("-------")
	}
}

func TestInspaceOfGoStruct3(t *testing.T) {
	b, err := json.Marshal(test1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(b))
	var cv SomeStruct
	err = json.Unmarshal(b, &cv)
	if err != nil {
		t.FailNow()
	}

}

type MyType int

type MyInterface interface {
	SayHi() error
}

type NestedStruct struct {
	Field1 int
	Field2 map[string]Bro
}

func (s NestedStruct) SayHi() error {
	fmt.Println("hi")
	return nil
}

type Bro struct {
	Name string
	Age  int
}

func (s *Bro) SayHi() error {
	fmt.Println("hi")
	return nil
}
func TestGoDefin2(t *testing.T) {
	slice := []int{1, 2, 3}
	sliceType := reflect.TypeOf(slice).Elem()
	fmt.Println(sliceType)
	slice2 := []interface{}{1, 2, []int{1, 2, 3}}
	sliceType2 := reflect.TypeOf(slice2).Elem()
	fmt.Println(sliceType2)
}

func TestAssignThings(t *testing.T) {
	type MyStruct struct {
		Field int
	}

	type Container struct {
		MyField MyStruct
	}

	c := &Container{}

	field := reflect.ValueOf(c).Elem().FieldByName("MyField")
	fmt.Println(field.Type())
	newStruct := reflect.New(field.Type()).Elem()

	newStruct.FieldByName("Field").SetInt(100)
	field.Set(newStruct)
	fmt.Printf("%#v", c)
}

func TestAssignThings2(t *testing.T) {
	type MyStruct struct {
		Field int
	}

	type Container struct {
		MyField *MyStruct
	}

	c := &Container{}

	field := reflect.ValueOf(c).Elem().FieldByName("MyField") // *MyStruct

	fmt.Println(field.IsNil()) // true
	fmt.Println(field.Kind())  // ptr
	fmt.Println(field.Type())  // *interpreter_test.MyStruct

	newStruct := reflect.New(field.Type().Elem()) // not field.Elem().Type(), because field.Elem() is dereferncing, and defernce a nil will panic
	newStruct.Elem().FieldByName("Field").SetInt(100)

	field.Set(newStruct)

	fmt.Printf("Container's MyField: %+v\n", c.MyField) // Prints Container's MyField: {Field:42}

}

func TestAssignThings4(t *testing.T) {
	type MyStruct struct {
		Field int
	}

	type Container struct {
		MyField []MyStruct
	}

	c := &Container{}

	field := reflect.ValueOf(c).Elem().FieldByName("MyField")
	elementType := field.Type().Elem()
	sliceElement := reflect.SliceOf(elementType)
	sliceValue := reflect.MakeSlice(sliceElement, 0, 0)

	for i := 0; i < 5; i++ {
		element := reflect.New(elementType).Elem() // ptr to MyStruct, ok to append
		element.FieldByName("Field").SetInt(int64(i))
		sliceValue = reflect.Append(sliceValue, element)
	}
	field.Set(sliceValue)

	fmt.Printf("%#v", c)
}

func TestAssignThings5(t *testing.T) {
	type MyStruct struct {
		Field int
	}

	type Container struct {
		MyField []*MyStruct
	}

	c := &Container{}

	field := reflect.ValueOf(c)
	if field.Kind() == reflect.Pointer {
		field = field.Elem()
	}
	field = field.FieldByName("MyField")
	myFieldSliceElementType := field.Type().Elem()

	sliceType := reflect.SliceOf(myFieldSliceElementType)
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)

	for i := 0; i < 10; i++ {
		//ptrToStruct := reflect.New(myFieldSliceElementType).Elem() // ptr to ptrMyStruct
		ptrToStruct := reflect.New(myFieldSliceElementType.Elem()) // ptr to MyStruct
		ptrToStruct.Elem().FieldByName("Field").SetInt(int64(i))
		sliceValue = reflect.Append(sliceValue, ptrToStruct)

	}
	field.Set(sliceValue)

	for _, n := range c.MyField {
		fmt.Printf("%#v\n", n)
	}

}

func TestAssignThings11(t *testing.T) {
	type test1 struct {
		Hello     string
		World     float64
		Apple     bool
		Banana    bool
		something interface{}
	}

	t1 := test1{"Peter", 100, true, false, nil}
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
	if someTest1.World != 100 {
		t.FailNow()
	}
	if someTest1.Apple != true {
		t.FailNow()
	}
	if someTest1.Banana != false {
		t.FailNow()
	}
	if someTest1.something != nil {
		t.FailNow()
	}
}
