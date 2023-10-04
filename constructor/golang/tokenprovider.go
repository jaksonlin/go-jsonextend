package golang

import (
	"io"
	"reflect"
	"strconv"

	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

type workingItem struct {
	reflectValue   reflect.Value
	tokenType      token.TokenType
	isArrayElement bool
	isRoot         bool
}

type tokenProvider struct {
	rootOut      reflect.Value
	workingStack *util.Stack[*workingItem]
	visited      map[uintptr]bool
}

func newTokenProvider(out interface{}) (*tokenProvider, error) {
	s := util.NewStack[*workingItem]()
	v := reflect.ValueOf(out)

	theTokenType := token.GetTokenTypeByReflection(&v)
	if theTokenType == token.TOKEN_UNKNOWN {
		return nil, ErrorUnknownData
	}

	s.Push(&workingItem{reflectValue: v, tokenType: theTokenType, isArrayElement: false, isRoot: true})

	return &tokenProvider{
		rootOut:      v,
		workingStack: s,
		visited:      make(map[uintptr]bool),
	}, nil
}

var _ constructor.TokenProvider = &tokenProvider{}

func (t *tokenProvider) ReadBool() (bool, error) {
	// in boolean state we will consume until it is the end of boolean.
	data, err := t.dataSource.Peek(1)
	if err != nil {
		return false, err
	}

	numberOfRead := 0
	if data[0] == 't' {
		numberOfRead = 4
	} else if data[0] == 'f' {
		numberOfRead = 5
	} else {
		return false, ErrorIncorrectCharacter
	}

	rs := make([]byte, numberOfRead)

	_, err = io.ReadFull(t.dataSource, rs)
	if err != nil {
		return false, err
	}

	rsBoolean, err := strconv.ParseBool(string(rs))
	if err != nil {
		return false, ErrorIncorrectValueForState
	}
	return rsBoolean, nil
}

func (t *tokenProvider) processArrayItem(item *workingItem) {
	len := item.reflectValue.Len()
	// push the end tag
	t.workingStack.Push(&workingItem{tokenType: token.TOKEN_RIGHT_BRACKET, isArrayElement: false, isRoot: false})
	for i := len - 1; i >= 0; i -= 1 {
		element := item.reflectValue.Index(i)
		theTokenType := token.GetTokenTypeByReflection(&element)
		t.workingStack.Push(&workingItem{reflectValue: element, tokenType: theTokenType, isArrayElement: true, isRoot: false})
	}
}

func (t *tokenProvider) processObjectItem(item *workingItem) error {
	// push the end tag
	t.workingStack.Push(&workingItem{tokenType: token.TOKEN_RIGHT_BRACE, isArrayElement: false, isRoot: false})

	if item.reflectValue.Kind() == reflect.Struct {
		for i := item.reflectValue.NumField() - 1; i >= 0; i -= 1 {
			field := item.reflectValue.Type().Field(i)
			if field.IsExported() {
				element := item.reflectValue.Index(i)
				theTokenType := token.GetTokenTypeByReflection(&element)
				// in stack value first then key
				t.workingStack.Push(&workingItem{reflectValue: element, tokenType: theTokenType, isArrayElement: false, isRoot: false})
				t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(field.Name), tokenType: token.TOKEN_STRING, isArrayElement: false, isRoot: false})
			}
		}
	} else {
		for _, key := range item.reflectValue.MapKeys() {
			mapValue := item.reflectValue.MapIndex(key)
			valueTokenType := token.GetTokenTypeByReflection(&mapValue)
			t.workingStack.Push(&workingItem{reflectValue: mapValue, tokenType: valueTokenType, isArrayElement: false, isRoot: false})
			keyTokenType := token.GetTokenTypeByReflection(&key)
			if keyTokenType == token.TOKEN_NUMBER {
				keyValue, err := convertNumericToString(key)
				if err != nil {
					return err
				}
				t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(keyValue), tokenType: token.TOKEN_STRING, isArrayElement: false, isRoot: false})
			} else {
				t.workingStack.Push(&workingItem{reflectValue: key, tokenType: keyTokenType, isArrayElement: false, isRoot: false})
			}

		}
	}
	return nil

}

func (t *tokenProvider) GetNextTokenType() (token.TokenType, error) {

	item, err := t.workingStack.Peek()
	if err != nil {
		return token.TOKEN_DUMMY, err
	}

	// deferences any pointer
	if item.tokenType != token.TOKEN_NULL {
		for item.reflectValue.Kind() == reflect.Pointer {
			addr := getMemoryAddress(item.reflectValue)
			if _, ok := t.visited[addr]; ok {
				return token.TOKEN_DUMMY, ErrorCyclicAccess
			} else {
				t.visited[addr] = true
			}
			item.reflectValue = item.reflectValue.Elem()
		}
	}
	addr := getMemoryAddress(item.reflectValue)
	if _, ok := t.visited[addr]; ok {
		return token.TOKEN_DUMMY, ErrorCyclicAccess
	} else {
		t.visited[addr] = true
	}
	t.processArrayItem(item)
	switch item.tokenType {
	case token.TOKEN_LEFT_BRACKET:
		t.processArrayItem(item)
		t.workingStack.Pop()
		return item.tokenType, nil
	case token.TOKEN_LEFT_BRACE:
		err := t.processObjectItem(item)
		if err != nil {
			return token.TOKEN_DUMMY, err
		}
		t.workingStack.Pop()
		return item.tokenType, nil
	default:
		// they will be pop when ReadXXX is requested, and we have already marked them visit
		return item.tokenType, nil
	}

}

func (t *tokenProvider) ReadNull() error {

	return nil

}

func (t *tokenProvider) ReadNumber() (float64, error) {

	return f64, nil
}

func (t *tokenProvider) ReadString() ([]byte, error) {

}

func (t *tokenProvider) ReadVariable() ([]byte, error) {

}
