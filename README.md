# go-jsonextend

## Description

A simple Go json parser that support defining variables in json file.

## Usage

```go

package main

import (
    "fmt"
    "github.com/jaksonlin/go-jsonextend"
)

template:= `{"hello": "world", "name": "this is my ${name}", "age": ${age}}`
variables:= {"name": "jakson", "age": 18}

result, err:= jsonextend.Parse(template, variables)
if err != nil {
    fmt.Println(err)
}
fmt.Println(result)

```

