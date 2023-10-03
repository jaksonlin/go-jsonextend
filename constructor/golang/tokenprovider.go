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
	value          interface{}
	tokenType      token.TokenType
	isArrayElement bool
	isRoot         bool
}

type tokenProvider struct {
	rootOut      reflect.Value
	workingStack *util.Stack[*workingItem]
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

func (t *tokenProvider) GetNextTokenType() (token.TokenType, error) {

	item, err := t.workingStack.Peek()
	if err != nil {
		return token.TOKEN_DUMMY, err
	}

	// deferences any pointer
	if item.tokenType != token.TOKEN_NULL {
		for item.reflectValue.Kind() == reflect.Pointer {
			item.value = item.reflectValue.Elem()
		}
	}

	switch item.tokenType {
	case token.TOKEN_LEFT_BRACKET:
		len := item.reflectValue.Len()
		// push the end tag
		t.workingStack.Push(&workingItem{tokenType: token.TOKEN_RIGHT_BRACKET, isArrayElement: false, isRoot: false})
		for i := len - 1; i >= 0; i -= 1 {
			element := item.reflectValue.Index(i)
			theTokenType := token.GetTokenTypeByReflection(&element)
			t.workingStack.Push(&workingItem{reflectValue: element, tokenType: theTokenType, isArrayElement: true, isRoot: false})
		}
		t.workingStack.Pop()
		return item.tokenType, nil
	case token.TOKEN_LEFT_BRACE:
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
					t.workingStack.Push(&workingItem{value: field.Name, tokenType: token.TOKEN_STRING, isArrayElement: false, isRoot: false})
				}
			}
		} else {
			for _, key := range item.reflectValue.MapKeys() {
				t.workingStack.Push(&workingItem{item.workValue.MapIndex(key), -1, key.Interface()})
			}
		}
		t.workingStack.Pop()
		return item.tokenType, nil

	default:
		return item.tokenType, nil
	}

	nextTokenType := GetTokenTypeByStartCharacter(nextByte)

	if ShouldUnreadByte(nextTokenType) {
		err = t.dataSource.UnreadByte()
		if err != nil {
			return token.TOKEN_DUMMY, err
		}
	}

	return nextTokenType, nil
}

func (t *tokenProvider) ReadNull() error {
	rs := make([]byte, 4)
	_, err := io.ReadFull(t.dataSource, rs)
	if err != nil {
		return err
	}

	if string(rs) != "null" {
		return ErrorIncorrectValueForState
	}
	return nil

}

func (t *tokenProvider) ReadNumber() (float64, error) {
	lengthOfNumber := 1
	for {
		nextByte, err := t.dataSource.Peek(lengthOfNumber)
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			lengthOfNumber -= 1
			//bare number
			break
		}
		if !isJSONNumberByte(nextByte[lengthOfNumber-1]) {
			lengthOfNumber -= 1 // remove the invalid location
			break
		}
		lengthOfNumber += 1
	}

	result := make([]byte, lengthOfNumber)
	_, err := io.ReadFull(t.dataSource, result)
	if err != nil {
		return 0, err
	}

	f64, err := strconv.ParseFloat(string(result), 64)
	if err != nil {
		return 0, ErrorIncorrectValueForState
	}
	return f64, nil
}

func (t *tokenProvider) ReadString() ([]byte, error) {

	// in order to deal with the multiple slash/escape sequence, we need a flag to check the string state
	isSlashEnclosed := true
	stringLength := 1
	validQuotationCount := 0
	for {
		nextByte, err := t.dataSource.Peek(stringLength)
		if err != nil {
			return nil, err
		}

		// slashes are closed
		if isSlashEnclosed {
			// is double quotation mark, end of string
			if nextByte[stringLength-1] == '"' {
				validQuotationCount += 1
				if validQuotationCount == 2 {
					rs := make([]byte, stringLength)
					_, err := io.ReadFull(t.dataSource, rs)
					if err != nil {
						return nil, err
					}
					return rs, nil
				}
			} else if nextByte[stringLength-1] == 0x5c { // is slash
				isSlashEnclosed = false
			}
		} else {
			isSlashEnclosed = true // flip the slash when there's nextByte, in escape mode skip quotation count check
		}
		stringLength += 1

	}
}

func (t *tokenProvider) ReadVariable() ([]byte, error) {
	variable, err := t.dataSource.ReadBytes('}')
	if err != nil {
		return nil, err
	}
	return variable, nil
}
