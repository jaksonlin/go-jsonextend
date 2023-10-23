package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

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

type JSONStructField struct {
	FieldName    string
	FieldValue   reflect.Value
	FieldJsonTag *JsonTagOptions
	ExtendTag    *JsonExtendOptions
}

func (j *JSONStructField) ShouldDropFieldIfSetValue(valueToUnmarshal reflect.Value) bool {
	if j.FieldJsonTag == nil {
		return false
	}
	if j.FieldJsonTag.Omitempty && isEmptyValue(valueToUnmarshal) {
		return true
	}
	return false
}

type JsonTagOptions struct {
	Omitempty    bool
	StringEncode bool
	fieldName    string // will not be used outside, cause the correct field name will be set in the `fieldName` of `JSONStructField`
}

type JsonExtendOptions struct {
	FieldVariableKeyName   string
	FieldVariableValueName string
}

func GetFieldNameAndOptions(jsonTag string) *JsonTagOptions {
	if jsonTag == "-" {
		return nil
	}
	rs := strings.Split(jsonTag, ",")
	ret := &JsonTagOptions{}
	for _, item := range rs {
		if item == "omitempty" {
			ret.Omitempty = true
		} else if item == "string" {
			ret.StringEncode = true
		} else {
			ret.fieldName = item
		}
	}
	return ret
}

var extendTagPattern = regexp.MustCompile(`(\w+=\w+)`)

func getExtensionTags(field reflect.StructField) *JsonExtendOptions {

	tag, ok := field.Tag.Lookup("jsonext")
	if !ok {
		return nil
	}
	matches := extendTagPattern.FindAllStringSubmatch(tag, -1)
	if len(matches) == 0 {
		return nil
	}
	ret := &JsonExtendOptions{}
	for _, match := range matches {
		kv := strings.Split(match[1], "=")
		if len(kv) != 2 {
			continue
		}
		if kv[0] == "k" {
			ret.FieldVariableKeyName = kv[1]
		} else if kv[0] == "v" {
			ret.FieldVariableValueName = kv[1]
		}

	}
	if len(ret.FieldVariableKeyName) == 0 && len(ret.FieldVariableValueName) == 0 {
		return nil
	}
	return ret

}

// this is for Unmarshal, the main differences relies on the fact that when unmarshalling, we cannot use workItem's value
// and omityempty field to drop field that will be returned from here, as the values has not yet been unmarshal into the workItem
func FlattenJsonStructForUnmarshal(workItem reflect.Value) map[string]*JSONStructField {
	if workItem.Kind() != reflect.Struct {
		return nil
	}
	s := NewStack[reflect.Value]()
	s.Push(workItem)
	var flattenFields map[string]*JSONStructField = make(map[string]*JSONStructField)
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
			if jsonTag == "-" {
				continue
			}
			// 3. json tag without json key name, lower the precedent
			if strings.HasPrefix(jsonTag, ",") {
				// json tag without json key name, lower the precedent
				// when a field name colllison with json tag name, the field will be dropped
				noneJsonTagFields = append(noneJsonTagFields, i)
				continue
			}
			tagConfig := GetFieldNameAndOptions(jsonTag)
			extendTag := getExtensionTags(field)
			fieldValue := item.Field(i)

			// 4. check same level key collision
			if createFromHere, ok := jsonTagFields[tagConfig.fieldName]; !ok {
				jsonTagFields[tagConfig.fieldName] = false
				// 4.1 no collision with upper level, create and mark createFromHere, otherwise upper level take precedent, do nothing
				if _, ok := flattenFields[tagConfig.fieldName]; !ok {
					jsonTagFields[tagConfig.fieldName] = true
					flattenFields[tagConfig.fieldName] = &JSONStructField{
						FieldName:    tagConfig.fieldName,
						FieldValue:   fieldValue,
						FieldJsonTag: tagConfig,
						ExtendTag:    extendTag,
					}
				}
			} else {
				// 4.2 collision occurs, and we have created value into the `flattenFields`, removed them
				if createFromHere {
					delete(flattenFields, tagConfig.fieldName)
				}
			}
		}
		// 5. check if any of none json tag field would conflict with json tag fields, if not add them in
		for _, noneJsonfieldIndex := range noneJsonTagFields {
			field := item.Type().Field(noneJsonfieldIndex)
			if _, ok := flattenFields[field.Name]; !ok {
				fieldValue := item.Field(noneJsonfieldIndex)
				result := &JSONStructField{
					FieldName:    field.Name,
					FieldValue:   fieldValue,
					FieldJsonTag: nil,
				}
				jsonTag, ok := field.Tag.Lookup("json")

				if ok {
					// no json tag field name but have option
					tagConfig := GetFieldNameAndOptions(jsonTag)
					result.FieldJsonTag = tagConfig
				}
				extendTag := getExtensionTags(field)
				result.ExtendTag = extendTag
				flattenFields[field.Name] = result
			}
		}
	}

	return flattenFields
}

func shouldDropField(value reflect.Value, tagConfig *JsonTagOptions) bool {
	if tagConfig == nil {
		return false
	}
	if tagConfig.Omitempty && isEmptyValue(value) {
		return true
	}
	return false
}

func isEmptyValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.String:
		return value.String() == ""
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() == 0
	case reflect.Struct:
		return false
	default:
		return false
	}
}

// FlattenJsonStruct flatten a struct into a list of JSONStructField, will drop field base on json tag settings, as the workItem already have all the information set.
func FlattenJsonStructForMarshal(workItem reflect.Value) []*JSONStructField {
	if workItem.Kind() != reflect.Struct {
		return nil
	}
	s := NewStack[reflect.Value]()
	s.Push(workItem)
	var flattenFieldsState map[string]int = make(map[string]int)
	var flattenFields []*JSONStructField = make([]*JSONStructField, 0)

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
			// 2.1 drop when the field is marked with `-`
			if jsonTag == "-" {
				continue
			}
			// 2.2 for json tag without json key name, lower the precedent
			if strings.HasPrefix(jsonTag, ",") {
				// no json key name in json tag, lower the precedent
				// when a field name colllison with json tag name, the field will be dropped
				noneJsonTagFields = append(noneJsonTagFields, i)
				continue
			}
			// 2.3 check for omitempty
			tagConfig := GetFieldNameAndOptions(jsonTag)
			extendTag := getExtensionTags(field)
			fieldValue := item.Field(i)
			if shouldDropField(fieldValue, tagConfig) {
				continue
			}
			// 3. check same level key collision
			if createFromHere, ok := jsonTagFields[tagConfig.fieldName]; !ok {
				jsonTagFields[tagConfig.fieldName] = false
				// 4.1 no collision with upper level, create and mark createFromHere, otherwise upper level take precedent, do nothing
				if _, ok := flattenFieldsState[tagConfig.fieldName]; !ok {
					jsonTagFields[tagConfig.fieldName] = true
					flattenFieldsState[tagConfig.fieldName] = len(flattenFields)
					flattenFields = append(flattenFields, &JSONStructField{
						FieldName:    tagConfig.fieldName,
						FieldValue:   fieldValue,
						FieldJsonTag: tagConfig,
						ExtendTag:    extendTag,
					})
				}
			} else {
				// 4.2 collision occurs, and we have created value into the `flattenFields`, removed them
				if createFromHere {
					index := flattenFieldsState[tagConfig.fieldName]
					flattenFields = append(flattenFields[:index], flattenFields[index+1:]...)
				}
			}
		}
		// 5. check if any of none json tag field would conflict with json tag fields, if not add them in
		for _, noneJsonfieldIndex := range noneJsonTagFields {
			field := item.Type().Field(noneJsonfieldIndex)
			if _, ok := flattenFieldsState[field.Name]; !ok {
				flattenFieldsState[field.Name] = len(flattenFields)
				FieldValue := item.Field(noneJsonfieldIndex)
				result := &JSONStructField{
					FieldName:    field.Name,
					FieldValue:   FieldValue,
					FieldJsonTag: nil,
				}
				jsonTag, ok := field.Tag.Lookup("json")
				if ok {
					// no json tag field name but have option
					tagConfig := GetFieldNameAndOptions(jsonTag)
					// check for omitempty
					if shouldDropField(FieldValue, tagConfig) {
						continue
					}
					result.FieldJsonTag = tagConfig
				}
				extendTag := getExtensionTags(field)
				result.ExtendTag = extendTag
				flattenFields = append(flattenFields, result)
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

func EncodePrimitiveValue(v interface{}) ([]byte, error) {
	if v == nil {
		return token.NullBytes, nil
	}
	switch data := v.(type) {
	case string:
		return EncodeToJsonString(data), nil
	case float32, float64:
		return []byte(fmt.Sprintf("%g", v)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return []byte(fmt.Sprintf("%d", v)), nil
	case bool:
		if data {
			return token.TrueBytes, nil
		}
		return token.FalseBytes, nil
	case nil:
		return token.NullBytes, nil
	case interface{}:
		return EncodePrimitiveValue(data)
	default:
		return nil, ErrorVariableDataKind
	}
}

// json input value is always float64, convert to different numeric value based on out element kind
func ConvertInterfaceNumberToFloat64(val interface{}) (float64, error) {

	value := reflect.ValueOf(val)
	if !value.IsValid() {
		return 0, ErrorInputNil
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// Continue to the conversion
	default:
		return 0, ErrorInputNotNumber
	}
	// Catch potential panics during the conversion
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		// Optionally log the panic message, e.g., using log package
	// 		// log.Println("panic recovered:", r)
	// 	}
	// }()
	f64Type := reflect.TypeOf(float64(0))
	f64Value := value.Convert(f64Type)
	return f64Value.Float(), nil
}

func ConvertInterfaceNumberToInt64(val interface{}) (int64, error) {
	value := reflect.ValueOf(val)
	if !value.IsValid() {
		return 0, ErrorInputNil
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// Continue to the conversion
	default:
		return 0, ErrorInputNotNumber
	}
	// Catch potential panics during the conversion
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		// Optionally log the panic message, e.g., using log package
	// 		// log.Println("panic recovered:", r)
	// 	}
	// }()
	i64Type := reflect.TypeOf(int64(0))
	i64Value := value.Convert(i64Type)
	return i64Value.Int(), nil
}
