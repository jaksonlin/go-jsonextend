package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jaksonlin/go-jsonextend/token"
)

type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   make([]string, 0),
		values: make(map[string]interface{}),
	}
}

// Add a new key-value pair
func (om *OrderedMap) Add(key string, value interface{}) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

// Get the value for a key
func (om *OrderedMap) Get(key string) (interface{}, bool) {
	val, ok := om.values[key]
	return val, ok
}

// Iterate through the key-value pairs in order
func (om *OrderedMap) Iterate(f func(key string, value interface{})) {
	for _, key := range om.keys {
		f(key, om.values[key])
	}
}

func FlattenJsonStruct(workItem reflect.Value) map[string]reflect.Value {
	if workItem.Kind() != reflect.Struct {
		return nil
	}
	s := NewStack[reflect.Value]()
	s.Push(workItem)
	var flattenFields map[string]reflect.Value = make(map[string]reflect.Value)
	for {
		item, err := s.Pop()
		if err != nil {
			break
		}
		var jsonTagFields map[string]bool = make(map[string]bool)
		var noneJsonTagFields []int
		// json field goes first
		for i := item.NumField() - 1; i >= 0; i -= 1 {
			field := item.Type().Field(i)
			// 0. Anonymous, next level to check
			if field.Anonymous {
				s.Push(item.Field(i))
				continue
			}
			// 1. skip none exported
			if !field.IsExported() {
				continue
			}
			jsonTag, ok := field.Tag.Lookup("json")
			// 2. mark location for none json tag fields
			if !ok {
				noneJsonTagFields = append(noneJsonTagFields, i)
				continue
			}
			fieldName, _ := GetFieldNameAndOptions(jsonTag)
			// 3. json tag without json key name, lower the precedent
			if len(fieldName) == 0 {
				// no json key name in json tag, lower the precedent
				// when a field name colllison with json tag name, the field will be dropped
				noneJsonTagFields = append(noneJsonTagFields, i)
				continue
			}
			// 4. check same level key collision
			if createFromHere, ok := jsonTagFields[fieldName]; !ok {
				jsonTagFields[fieldName] = false
				// 4.1 no collision with upper level, create and mark createFromHere, otherwise upper level take precedent, do nothing
				if _, ok := flattenFields[fieldName]; !ok {
					jsonTagFields[fieldName] = true
					flattenFields[fieldName] = item.Field(i)
				}
			} else {
				// 4.2 collision occurs, and we have created value into the `flattenFields`, removed them
				if createFromHere {
					delete(flattenFields, fieldName)
				}
			}
		}
		// 5. check if any of none json tag field would conflict with json tag fields, if not add them in
		for _, noneJsonfieldIndex := range noneJsonTagFields {
			field := item.Type().Field(noneJsonfieldIndex)
			if _, ok := flattenFields[field.Name]; !ok {
				flattenFields[field.Name] = item.Field(noneJsonfieldIndex)
			}
		}
	}

	return flattenFields
}

func IsPrimitiveType(value reflect.Value) bool {
	for value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Bool, reflect.Int, reflect.Float64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.String:
		return true
	default:
		return false
	}
}

var ErrorUnsupportedDataKind = errors.New("unsupported variable data kind")

func EncodePrimitiveValue(v interface{}) ([]byte, error) {
	if v == nil {
		return token.NullBytes, nil
	}
	switch data := v.(type) {
	case string:
		return EncodeToJsonString(data), nil
	case float32, float64:
		return []byte(fmt.Sprintf("%f", v)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return []byte(fmt.Sprintf("%d", v)), nil
	case bool:
		if data {
			return token.TrueBytes, nil
		}
		return token.FalseBytes, nil
	case nil:
		return token.NullBytes, nil
	default:
		return nil, ErrorUnsupportedDataKind
	}
}