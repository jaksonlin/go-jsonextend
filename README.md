# go-jsonextend

## Description

A simple Go json parser that support defining variables in json file.

## Usage


### json unmarshaller

```go

type SomeStruct struct {
    Field1 string `json:"field1"`
    Field2 bool   `json:"field2"`
    Field3 int
    Field4 interface{}
}

testExample := SomeStruct{
    Field1: "hello",
    Field2: true,
    Field3: 100,
    Field4: nil,
}

data, _ := json.Marshal(testExample)
var out SomeStruct
err := jsonextend.Unmarshal(bytes.NewReader(data), nil, &out)
if err != nil {
    t.FailNow()
}

```

### json unmarshaller with json.Unmarshaller interface

```go
type SomeStructWithUnmarshaller struct {
    Field1 string
    Field2 int
    Field3 interface{}
}

func (s *SomeStructWithUnmarshaller) UnmarshalJSON(payload []byte) error {
    s.Field1 = "hello world"
    s.Field2 = 100
    s.Field3 = payload
    return nil
}

func TestUnmarshal4(t *testing.T) {
    someData := "Field3 value"
    data, _ := json.Marshal(someData)
    var validator SomeStructWithUnmarshaller
    err := jsonextend.Unmarshal(bytes.NewReader(data), nil, &validator)
    if err != nil {
        t.FailNow()
    }
    if validator.Field1 != "hello world" {
        t.FailNow()
    }
    if validator.Field2 != 100 {
        t.FailNow()
    }

    if !bytes.Equal(validator.Field3.([]byte), data) {
        t.FailNow()
    }

}

```

### Variable in string

when a string contains `${variable}` pattern, it is considered as a `string with variable`, and the variable value will be replaced in the result. a string can contain multiple variables.

if the variable value is not found, it will be kept as it is, giving the chance to be intrepreted later.

```go
type SomeStruct struct {
    Field1 string
}
testExample := `
{
    "Field1": "hello ${var1}",
}`

variables := map[string]interface{}{
    "var1":      "world!",
}

var out SomeStruct
err := jsonextend.Unmarshal(strings.NewReader(testExample), variables, &out)
if err != nil {
    t.FailNow()
}
if out.Field1 != "hello world!" {
    t.FailNow()
}

```

### Variable as field value

we can use variable as field value, as long as the value is json compatible (json.Marshal won't fail).

in this case when the variable value is not found, it will report error.

```go
type SomeStruct struct {
    Field3 interface{}
}
testExample := `
{
    "Field3":${var3}
}`

variables := map[string]interface{}{
    "var3":      []int{1, 2, 3},
}

var out SomeStruct
err := jsonextend.Unmarshal(strings.NewReader(testExample), variables, &out)
if err != nil {
    t.FailNow()
}

for i, v := range out.Field3.([]int) {
    if v != i+1 {
           t.FailNow()
    }
}
```

### Variable as field value router

combine the use of [Variable in string](#variable-in-string) and [Variable as field value](#variable-as-field-value), you can use variable as field value router.

in below case, the value of variable "var2" is used as the field name of the result; and the value of variable "var2Value" is used as the field value. this can help you dynamically set the field name and value.

```go
type SomeStruct struct {
    Field1 string
    Field2 int
}
testExample := `
{
    "Field1": "hello ${var1}",
    "${var2}": ${var2Value},
}`

variables := map[string]interface{}{
    "var1":      "world!",
    "var2":      "Field2",
    "var2Value": 100,
}

var out SomeStruct
err := jsonextend.Unmarshal(strings.NewReader(testExample), variables, &out)
if err != nil {
    t.FailNow()
}
if out.Field1 != "hello world!" {
    t.FailNow()
}
if out.Field2 != 100 {
    t.FailNow()
}

```

this is not limited in struct/map, you can use it in any json compatible value.

```go
dataTemplate := `[1,true,"hello", null, ${var1}, ${var2}]`
var someItemSlice []interface{}

variables := map[string]interface{}{
    "var1": []int{1, 2, 3},
    "var2": map[string]interface{}{"baby": "shark"},
}

err := jsonextend.Unmarshal(strings.NewReader(dataTemplate), variables, &someItemSlice)
if err != nil {
    t.FailNow()
}
if someItemSlice[0] != 1.0 {
    t.FailNow()
}
if someItemSlice[1] != true {
    t.FailNow()
}
if someItemSlice[2] != "hello" {
    t.FailNow()
}
if someItemSlice[3] != nil {
    t.FailNow()
}
for i, v := range someItemSlice[4].([]int) {
    if v != i+1 {
        t.FailNow()
    }
}
if someItemSlice[5].(map[string]interface{})["baby"] != "shark" {
    t.FailNow()
}
```

### json template engine

```go

package main

import (
    "fmt"
    "strings"

    "github.com/jaksonlin/go-jsonextend"
)

func main() {
    template := `{"hello": "world", "name": "this is my ${name}", "age": ${age}}`
    variables := map[string]interface{}{"name": "jakson", "age": 18}

    result, err := jsonextend.Parse(strings.NewReader(template), variables)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(result))

}


```

this will output

``` json
{
    "hello" : "world",
    "name" : "this is my jakson",
    "age" : 18
}
```

## Advance feature - JsonFlexMarshal

use customize tag, one can marshal it with the extension syntax for downstream to interpret, this will be useful when you need to templating some field.

to output the json extension format template, add the `jsonext` tag on the field and specified its key and value's variable represetative.

when giving the property `k`, the key field of the json output will be replaced by the k's variable name,quoted in `"${...}"` as `STRING WITH VARIABLE`

### JsonFlexMarshal - `k`

``` golang
func TestCustomizeMarshaller1(t *testing.T) {
    type MyDataStruct struct {
        Name string `jsonext:"k=var1"`
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
    if string(data) != `{"${var1}":"hello"}` {
        t.FailNow()
    }
}
```


when both `json` tag and `jsonext` tag is applied the MarshalIntoTemplate will use the `k`'s variable name if it is set.

``` golang
func TestCustomizeMarshaller3Ext(t *testing.T) {
    type MyDataStruct struct {
        Name string `json:"myfield" jsonext:"k=var1"`
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
    if string(data) != `{"${var1}":"hello"}` {
        t.FailNow()
    }
}
```

### JsonFlexMarshal - `v`

when giving the property `v`, the value field of the json will be replaced by the variable name `${...}`. Note the differences is that 
in this case it is a variable mark in the json extension syntax (this data format can only be parsed by go-jsonextend as this is not the standard json format)

the reason not to quote the `${...}` with double quotation mark is that: we can later interpret it with any valid json data type: number/string/bool/null...

``` golang
func TestCustomizeMarshaller2(t *testing.T) {
    type MyDataStruct struct {
        Name string `jsonext:"v=var1"`
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
    if string(data) != `{"Name":${var1}}` {
        t.FailNow()
    }
}
```

### JsonFlexMarshal - `k` & `v`

it is also valid to give both properties, regardless of what the data type it is.

``` golang
func TestCustomizeMarshaller3(t *testing.T) {
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



```

### JsonFlexMarshal - Steal the Sky

With the Marshal with variable support, you can easliy implement your json's unmarshal by replacing the corresponding field with variable,
rather than wrighting the entire UnmarshalJson function.

This can make the templating of json easy to manage.

``` golang

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
```

this  feature can be applied even though the field is an object/array.

``` golang
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

    data, err := interpreter.MarshalWithVariable(item, map[string]interface{}{"var1": "hello", "var2": "world"})
    if err != nil {
        t.FailNow()
    }

    fmt.Println(data)
    fmt.Println(data)
    if string(data) != `{"hello":"world"}` {
        t.FailNow()
    }
}

func TestCustomizeMarshallerOnStruct2(t *testing.T) {

    type MyDataStruct struct {
        Name []int `jsonext:"k=var1,v=var2"`
    }
    item := &MyDataStruct{
        Name: []int{1, 2, 3, 4, 5},
    }

    data, err := interpreter.MarshalWithVariable(item, map[string]interface{}{"var1": "hello", "var2": "world"})
    if err != nil {
        t.FailNow()
    }

    fmt.Println(data)
    fmt.Println(data)
    if string(data) != `{"hello":"world"}` {
        t.FailNow()
    }
}


```

