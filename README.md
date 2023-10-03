# go-jsonextend

## Description

A simple Go json parser that support defining variables in json file.

## Usage

as a json template processor, it can be used in the following ways:

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

as a json unmarshaller

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

as a dyanmaic json unmarshaller, you can use variable as a route table to route the value to the right field.

```go
type SomeStruct struct {
    Field1 string
    Field2 int
    Field3 interface{}
}
testExample := `
{
    "Field1": "hello ${var1}",
    "${var2}": ${var2Value},
    "Field3":${var3}
}`

variables := map[string]interface{}{
    "var1":      "world!",
    "var2":      "Field2",
    "var2Value": 100,
    "var3":      []int{1, 2, 3},
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
for i, v := range out.Field3.([]int) {
    if v != i+1 {
   	t.FailNow()
}
}
```