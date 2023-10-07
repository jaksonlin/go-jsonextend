package golang

import (
	"reflect"
	"strings"

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
	visited      map[uintptr]bool // check visited when pop
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

func (t *tokenProvider) processArrayItem(item *workingItem) error {
	len := item.reflectValue.Len()
	// push the end tag
	t.workingStack.Push(&workingItem{tokenType: token.TOKEN_RIGHT_BRACKET})
	for i := len - 1; i >= 0; i -= 1 {
		element := item.reflectValue.Index(i)
		theTokenType := token.GetTokenTypeByReflection(&element)
		if theTokenType == token.TOKEN_UNKNOWN {
			return ErrorInvalidTypeOnExportedField
		}
		t.workingStack.Push(&workingItem{reflectValue: element, tokenType: theTokenType})
	}
	return nil
}

// for json tag `string`, convert to json for primitive data type
func (t *tokenProvider) parseJsonStringConfig(tokenType token.TokenType, value reflect.Value) error {
	if tokenType == token.TOKEN_BOOLEAN || tokenType == token.TOKEN_NUMBER || tokenType == token.TOKEN_STRING {
		payload, err := util.EncodePrimitiveValue(value.Interface())
		if err != nil {
			return err
		}
		t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(payload), tokenType: token.TOKEN_STRING})
		return nil
	} else {
		return ErrorStringConfigTypeInvalid
	}
}

func (t *tokenProvider) parseJsonFieldInfo(jsonTagFieldName string, defaultFieldName string, inUsedKey map[string]bool) error {
	fieldName := strings.TrimSpace(jsonTagFieldName)
	if len(fieldName) == 0 {
		t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(defaultFieldName), tokenType: token.TOKEN_STRING})
	} else {
		t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(fieldName), tokenType: token.TOKEN_STRING})
	}
	return nil
}

func (t *tokenProvider) parseJsonTag(defaultFieldName, jsonTag string, valueTokenType token.TokenType, value reflect.Value, inUsedKey map[string]bool) error {
	if jsonTag == "-" {
		return nil
	}
	jsonTagConfig := strings.SplitN(jsonTag, ",", 2)
	if len(jsonTagConfig) == 0 {
		return ErrorInvalidJsonTag
	}
	if len(jsonTagConfig) == 1 {
		t.workingStack.Push(&workingItem{reflectValue: value, tokenType: valueTokenType})
	} else {
		fieldOption := strings.TrimSpace(jsonTagConfig[1])
		if fieldOption == "omitempty" && value.IsZero() {
			return nil
		}
		// set the value first
		if fieldOption == "string" {
			err := t.parseJsonStringConfig(valueTokenType, value)
			if err != nil {
				return err
			}
		} else {
			t.workingStack.Push(&workingItem{reflectValue: value, tokenType: valueTokenType})
		}
	}
	return t.parseJsonFieldInfo(jsonTagConfig[0], defaultFieldName, inUsedKey)
}

func (t *tokenProvider) processStructField(field reflect.StructField, element reflect.Value, inUsedKey map[string]bool) error {
	valueTokenType := token.GetTokenTypeByReflection(&element)
	if valueTokenType == token.TOKEN_UNKNOWN {
		return ErrorInvalidTypeOnExportedField
	}

	jsonTag, ok := field.Tag.Lookup("json")
	if !ok {
		inUsedKey[field.Name] = true
		// in stack value first then key
		t.workingStack.Push(&workingItem{reflectValue: element, tokenType: valueTokenType})
		t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(field.Name), tokenType: token.TOKEN_STRING})
	} else {
		err := t.parseJsonTag(field.Name, jsonTag, valueTokenType, element, inUsedKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tokenProvider) flattenStruct(workItem *workingItem) error {
	allFields := util.FlattenJsonStruct(workItem.reflectValue)
	for key, val := range allFields {
		valueTokenType := token.GetTokenTypeByReflection(&val)
		if valueTokenType == token.TOKEN_UNKNOWN {
			return ErrorInvalidTypeOnExportedField
		}

		t.workingStack.Push(&workingItem{reflectValue: val, tokenType: valueTokenType})
		t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(key), tokenType: token.TOKEN_STRING})
	}
	return nil
}

func (t *tokenProvider) processMapItem(item *workingItem) error {
	for _, key := range item.reflectValue.MapKeys() {
		mapValue := item.reflectValue.MapIndex(key)
		valueTokenType := token.GetTokenTypeByReflection(&mapValue)
		if valueTokenType == token.TOKEN_UNKNOWN {
			return ErrorInvalidTypeOnExportedField
		}
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

func (t *tokenProvider) detectCyclicAccess(item *workingItem) error {
	if item.reflectValue.CanAddr() {
		addr := getMemoryAddress(item.reflectValue)
		if _, ok := t.visited[addr]; ok {
			return ErrorCyclicAccess
		} else {
			t.visited[addr] = true
		}
	}
	return nil
}

func (t *tokenProvider) GetNextTokenType() (token.TokenType, error) {

	item, err := t.workingStack.Peek()
	if err != nil {
		return token.TOKEN_DUMMY, err
	}
	if item.tokenType == token.TOKEN_NULL {
		return token.TOKEN_NULL, nil
	}

	for item.reflectValue.Kind() == reflect.Pointer {
		item.reflectValue = item.reflectValue.Elem()
	}

	switch item.tokenType {
	case token.TOKEN_LEFT_BRACKET:
		if err := t.detectCyclicAccess(item); err != nil {
			return token.TOKEN_DUMMY, err
		}
		t.workingStack.Pop()
		t.processArrayItem(item)
		return item.tokenType, nil
	case token.TOKEN_LEFT_BRACE:
		if err := t.detectCyclicAccess(item); err != nil {
			return token.TOKEN_DUMMY, err
		}
		t.workingStack.Pop()
		err := t.processObjectItem(item)
		if err != nil {
			return token.TOKEN_DUMMY, err
		}
		return item.tokenType, nil
	case token.TOKEN_RIGHT_BRACE:
		fallthrough
	case token.TOKEN_RIGHT_BRACKET:
		t.workingStack.Pop()
		return item.tokenType, nil
	default:
		// for primitives, they will be pop when ReadXXX is requested, and we have already marked them visit
		return item.tokenType, nil
	}

}

func (t *tokenProvider) ReadNull() error {
	_, err := t.workingStack.Pop()
	if err != nil {
		return err
	}

	return nil
}
func (t *tokenProvider) ReadBool() (bool, error) {
	item, err := t.workingStack.Pop()
	if err != nil {
		return false, err
	}
	val := item.reflectValue.Bool()
	return val, nil
}

func (t *tokenProvider) ReadString() ([]byte, error) {
	item, err := t.workingStack.Pop()
	if err != nil {
		return nil, err
	}
	val := item.reflectValue.String()
	// use go standard
	v := util.EncodeToJsonString(val)
	return v, nil
}

func (t *tokenProvider) ReadNumber() (float64, error) {
	item, err := t.workingStack.Pop()
	if err != nil {
		return 0.0, err
	}
	val, err := convertNumberBaseOnKind(item.reflectValue)
	if err != nil {
		return 0.0, err
	}
	return val, nil
}

func (t *tokenProvider) ReadVariable() ([]byte, error) {
	// no golang datatype corresponding to variable now, maybe we can extend this later through tag or plugin
	return nil, nil
}
