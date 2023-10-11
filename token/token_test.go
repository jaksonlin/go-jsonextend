package token_test

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	var something interface{}
	var someInt int = 100
	var somePointers interface{} = &someInt
	for i := 0; i <= 100; i++ {
		somePointers = &somePointers
	}
	something = somePointers
	var somePointers2 interface{} = &something
	for i := 0; i <= 100; i++ {
		somePointers2 = &somePointers2
	}
	something = somePointers2
	fmt.Println(something)
}
