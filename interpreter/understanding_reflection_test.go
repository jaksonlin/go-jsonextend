package interpreter_test

import (
	"encoding/json"
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
	Name1  string                 `json:"name1"`  //checked
	Name2  []int                  `json:"name2"`  //checked
	Name3  map[string]int         `json:"name3"`  //checked
	Name4  []interface{}          `json:"name4"`  //checked
	Name5  []Bro                  `json:"name5"`  //checked
	Name6  []*Bro                 `json:"name6"`  //checked
	Name7  Bro                    `json:"name7"`  //checked
	Name8  *Bro                   `json:"name8"`  //checked
	Name9  map[string]interface{} `json:"name9"`  //checked
	Name10 map[int]Bro            `json:"name10"` //checked
	Name11 [3]int                 `json:"name11"` //checked
	Name12 MyInterface            `json:"name12"` // pointer
	Name13 MyInterface            `json:"name13"` // struct
	// ... and so on for other cases
	Name14 []map[string][]interface{} `json:"name14"`
	Name15 interface{}                `json:"name15"` //checked covert to map[string]interface{}
	Name16 *Bro                       `json:"name16"` //checked
	Name17 map[string]Bro             `json:"name17"` //checked
	Name18 []int                      `json:"name18"` //checked
	Name19 MyType                     `json:"name19"` //checked
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
	err = json.Unmarshal(b, &cv) // go cannot unmarshal this format
	if err == nil {
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

func TestAppendNil(t *testing.T) {
	t1 := &[]interface{}{1, true, "string", 4.5}

	v := reflect.ValueOf(t1).Elem()

	nilValue := reflect.Zero(v.Type().Elem())

	newSlice := reflect.Append(v, nilValue)
	v.Set(newSlice)

	fmt.Println(t1)
}

func TestNestedPointerCase(t *testing.T) {
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

}

func TestMapint(t *testing.T) {
	type someStruct struct {
		Name15 map[int]string `json:"name15"`
	}

	// JSON data as a byte slice
	jsonData := []byte(`{"name15": {"1": "John", "2": "Doe"}}`)

	// Create an instance of someStruct
	var data someStruct

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.FailNow()
	}

	// Print the value in the Name15 field
	fmt.Println(data.Name15) // Prints: map[1:John 2:Doe]

}

func TestNilInterfaceint(t *testing.T) {
	type someStruct struct {
		Name15 *interface{}
	}

	// JSON data as a byte slice
	jsonData := someStruct{nil}

	// Unmarshal the JSON data into the struct
	rs, err := json.Marshal(jsonData)
	if err != nil {
		t.FailNow()
	}

	// Create an instance of someStruct
	var data someStruct
	err = json.Unmarshal(rs, &data)
	if err != nil {
		t.FailNow()
	}
	// Print the value in the Name15 field
	fmt.Println(data.Name15) // Prints: map[1:John 2:Doe]

}

func TestInterfaceReceiver(t *testing.T) {
	var myMap map[string]interface{} = map[string]interface{}{"Hello": "World"}
	var myInterface interface{} = myMap
	val := reflect.ValueOf(myInterface)
	fmt.Println(val.Kind()) // prints map
	if val.Kind() != reflect.Map {
		t.FailNow()
	}
	var nestedInterface interface{} = myInterface
	nestedInterfaceCheck := reflect.ValueOf(nestedInterface)
	if nestedInterfaceCheck.Kind() != reflect.Map {
		t.FailNow()
	}
	// this is the only way you can get reflect.Interface
	ptrToInterface := reflect.ValueOf(&myInterface)
	if ptrToInterface.Kind() != reflect.Pointer {
		t.FailNow()
	}
	if ptrToInterface.Elem().Kind() != reflect.Interface {
		t.FailNow()
	}

}
