package golang

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

const maxDepth = 1000

type workingItem struct {
	reflectValue reflect.Value
	tokenType    token.TokenType
	path         []string
	address      uintptr
}

type tokenProvider struct {
	rootOut      reflect.Value
	workingStack *util.Stack[*workingItem]
	visited      map[uintptr][]string // check visited when pop
}

func newRootTokenProvider(out interface{}) (*tokenProvider, error) {
	s := util.NewStack[*workingItem]()
	v := reflect.ValueOf(out)

	theTokenType := token.GetTokenTypeByReflection(&v)
	if theTokenType == token.TOKEN_UNKNOWN {
		return nil, ErrorUnknownData
	}
	var addr uintptr
	if v.CanAddr() {
		addr = getMemoryAddress(v)
	}
	s.Push(&workingItem{reflectValue: v, tokenType: theTokenType, address: addr, path: []string{v.Kind().String()}})

	return &tokenProvider{
		rootOut:      v,
		workingStack: s,
		visited:      make(map[uintptr][]string),
	}, nil
}

func canNilKind(k reflect.Kind) bool {
	return k == reflect.Interface || k == reflect.Ptr || k == reflect.Map || k == reflect.Slice
}

func newContainerWorkingItem(key string, v reflect.Value, parent *workingItem) (*workingItem, error) {

	theTokenType := token.GetTokenTypeByReflection(&v)
	if theTokenType == token.TOKEN_UNKNOWN {
		return nil, ErrorUnknownData
	}
	var addr uintptr = 0
	if v.CanAddr() {
		addr = getMemoryAddress(v)
	}
	itemType := v.Type()
	kind := v.Kind()
	addrable := v.CanAddr()
	var sb strings.Builder
	// Start by writing type and kind
	sb.WriteString(fmt.Sprintf("key:%s(%s#%s)", key, itemType.Name(), kind))
	// Handle different kinds and addressability
	if canNilKind(kind) && v.IsNil() {
		sb.WriteString("@nil->")
		path := append(parent.path, sb.String())
		return &workingItem{reflectValue: v, tokenType: theTokenType, address: addr, path: path}, nil
	}

	if !addrable {
		sb.WriteString(":unaddressable")
		path := append(parent.path, sb.String())
		return &workingItem{reflectValue: v, tokenType: theTokenType, address: addr, path: path}, nil
	}

	if kind == reflect.Slice {
		if v.Len() > 0 {
			sb.WriteString(fmt.Sprintf("@%d:len=%d:cap=%d",
				v.Index(0).UnsafeAddr(),
				v.Len(),
				v.Cap()))
		} else {
			sb.WriteString("@empty")
		}
	} else {
		sb.WriteString(fmt.Sprintf("@%d", v.UnsafeAddr()))
	}

	// Arrow for the next element
	sb.WriteString("->")

	path := append(parent.path, sb.String())
	return &workingItem{reflectValue: v, tokenType: theTokenType, address: addr, path: path}, nil

}
func (t *tokenProvider) detectCyclicAccess(item *workingItem) error {
	if item.address != 0 {
		if paths, ok := t.visited[item.address]; ok {
			currentPath := strings.Join(item.path, "->")
			for _, path := range paths {
				if strings.Contains(currentPath, path) {
					return ErrorCyclicAccess
				} else {
					t.visited[item.address] = append(t.visited[item.address], currentPath)
				}
			}
		} else {
			t.visited[item.address] = make([]string, 0)
			t.visited[item.address] = append(t.visited[item.address], strings.Join(item.path, "->"))
		}
	}
	return nil
}

var _ constructor.TokenProvider = &tokenProvider{}

func (t *tokenProvider) GetNextTokenType() (token.TokenType, error) {

	item, err := t.workingStack.Peek()
	if err != nil {
		return token.TOKEN_DUMMY, err
	}
	if item.tokenType == token.TOKEN_NULL {
		return token.TOKEN_NULL, nil
	}
	if item.reflectValue.Kind() == reflect.Interface && !item.reflectValue.IsNil() {
		item.reflectValue = item.reflectValue.Elem()
	}
	for item.reflectValue.Kind() == reflect.Pointer {
		item.reflectValue = item.reflectValue.Elem()
	}
	if item.reflectValue.Kind() == reflect.Interface && !item.reflectValue.IsNil() {
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
		newItem, err := newContainerWorkingItem(fmt.Sprintf("%d", i), element, item)
		if err != nil {
			return err
		}
		t.workingStack.Push(newItem)
	}
	return nil
}

func (t *tokenProvider) flattenStruct(workItem *workingItem) error {
	allFields := util.FlattenJsonStructForMarshal(workItem.reflectValue)
	for i := 0; i < len(allFields); i += 1 {
		val := allFields[i]
		valueTokenType := token.GetTokenTypeByReflection(&val.FieldValue)
		if valueTokenType == token.TOKEN_UNKNOWN {
			return ErrorInvalidTypeOnExportedField
		}

		if valueTokenType == token.TOKEN_LEFT_BRACE || valueTokenType == token.TOKEN_LEFT_BRACKET {
			// for none primitive type, we need to track the path
			newItem, err := newContainerWorkingItem(val.FieldName, val.FieldValue, workItem)
			if err != nil {
				return err
			}
			t.workingStack.Push(newItem)
		} else {
			t.workingStack.Push(&workingItem{reflectValue: val.FieldValue, tokenType: valueTokenType})
		}
		t.workingStack.Push(&workingItem{reflectValue: reflect.ValueOf(val.FieldName), tokenType: token.TOKEN_STRING})
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
