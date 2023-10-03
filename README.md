# go-jsonextend

## Description

A simple Go json parser that support defining variables in json file.

## Usage

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
    fmt.Println(result)

}


```

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