package interpreter_test

import (
	"fmt"
	"reflect"
	"testing"
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
	Name1  string
	Name2  []int
	Name3  map[string]int
	Name4  []interface{}
	Name5  []Bro
	Name6  []*Bro
	Name7  Bro
	Name8  *Bro
	Name9  map[string]interface{}
	Name10 map[int]Bro
	Name11 [3]int
	Name12 MyInterface // pointer
	Name13 MyInterface // struct
	// ... and so on for other cases
	Name14 []map[string][]interface{}
	Name15 interface{}
	Name16 *Bro
	Name17 map[string]Bro
	Name18 []int
	Name19 MyType
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
		elm := valueItem.Elem()
		for i := 0; i < elm.NumField(); i++ {
			f := valueItem.Field(i)
			fmt.Println(f.Kind())
		}
	}
	fmt.Println("END")
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
