# go-jsonextend

## Description

A simple Go json parser that support defining variables in json file.

## Usage

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

