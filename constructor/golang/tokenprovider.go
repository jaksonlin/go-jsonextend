package golang

import (
	"reflect"

	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

type workingItem struct {
	reflectValue reflect.Value
	tokenType    token.TokenType
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

	s.Push(&workingItem{reflectValue: v, tokenType: theTokenType})

	return &tokenProvider{
		rootOut:      v,
		workingStack: s,
		visited:      make(map[uintptr]bool),
	}, nil
}

var _ constructor.TokenProvider = &tokenProvider{}

func (t *tokenProvider) ReadBool() (bool, error) {

}

func (t *tokenProvider) processArrayItem(item *workingItem) {
	len := item.reflectValue.Len()
	// push the end tag
	t.workingStack.Push(&workingItem{tokenType: token.TOKEN_RIGHT_BRACKET})
	for i := len - 1; i >= 0; i -= 1 {
		element := item.reflectValue.Index(i)
		theTokenType := token.GetTokenTypeByReflection(&element)
		t.workingStack.Push(&workingItem{reflectValue: element, tokenType: theTokenType})
	}
}

func (t *tokenProvider) flattenStruct(workItem *workingItem) error {
	s := util.NewStack[reflect.Value]()
	s.Push(workItem.reflectValue)
	for {
		item, err := s.Pop()
		if err != nil {
			break
		}
		for i := item.NumField() - 1; i >= 0; i -= 1 {
			field := item.Type().Field(i)
			if field.Anonymous {
				s.Push(item.Field(i))
			}
			if field.IsExported() {
				element := item.Index(i)
				theTokenType := token.GetTokenTypeByReflection(&element)
				// in stack value first then key
				t.workingStack.Push(&workingItem{reflectValue: element, tokenType: theTokenType})
				t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(field.Name), tokenType: token.TOKEN_STRING})
			}
		}
	}

	return nil
}

func (t *tokenProvider) processMapItem(item *workingItem) error {
	for _, key := range item.reflectValue.MapKeys() {
		mapValue := item.reflectValue.MapIndex(key)
		valueTokenType := token.GetTokenTypeByReflection(&mapValue)
		t.workingStack.Push(&workingItem{reflectValue: mapValue, tokenType: valueTokenType})
		keyTokenType := token.GetTokenTypeByReflection(&key)
		if keyTokenType == token.TOKEN_NUMBER {
			keyValue, err := convertNumericToString(key)
			if err != nil {
				return err
			}
			t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(keyValue), tokenType: token.TOKEN_STRING})
		} else if keyTokenType == token.TOKEN_STRING {
			t.workingStack.Push(&workingItem{reflectValue: key, tokenType: token.TOKEN_STRING})
		} else {
			return ErrorInvalidMapKey
		}

	}
	return nil
}

func (t *tokenProvider) processObjectItem(item *workingItem) error {
	// push the end tag
	t.workingStack.Push(&workingItem{tokenType: token.TOKEN_RIGHT_BRACE})

	if item.reflectValue.Kind() == reflect.Struct {
		if err := t.flattenStruct(item); err != nil {
			return err
		}
	} else {
		if err := t.processMapItem(item); err != nil {
			return err
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
