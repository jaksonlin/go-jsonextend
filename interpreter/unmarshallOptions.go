package interpreter

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/config"
	"github.com/jaksonlin/go-jsonextend/util"
)

type unmarshallOptions struct {
	ensureInt     bool
	resolverStack *util.Stack[*unmarshallResolver]
	variables     map[string]interface{}
	marshaler     ast.MarshalerFunc
	unmarshaler   ast.UnmarshalerFunc
}

func NewUnMarshallOptions(variables map[string]interface{}, marshaler ast.MarshalerFunc, unmarshaler ast.UnmarshalerFunc) *unmarshallOptions {
	options := &unmarshallOptions{
		ensureInt:     config.EnsureInt,
		variables:     variables,
		resolverStack: util.NewStack[*unmarshallResolver](),
		marshaler:     marshaler,
		unmarshaler:   unmarshaler,
	}
	return options
}

func resolveVariable(variableNode *ast.JsonExtendedVariableNode, resolver *unmarshallOptions) (interface{}, error) {

	variableValue, ok := resolver.variables[variableNode.Variable]
	if !ok {
		return nil, NewVariableNotFound(variableNode.Variable)
	}
	return variableValue, nil
}

func resolveStringVariable(stringVariable *ast.JsonExtendedStringWIthVariableNode, resolver *unmarshallOptions) ([]byte, error) {

	var resultBytes []byte = make([]byte, len(stringVariable.Value))
	copy(resultBytes, stringVariable.Value)
	for variableName, replacer := range stringVariable.Variables {
		variableValue, ok := resolver.variables[variableName]
		if !ok {
			continue
		}
		variableValueBytes, err := resolver.marshaler(variableValue)
		if err != nil {
			return nil, err
		}
		if variableValueBytes[0] == '"' {
			if len(variableValueBytes) == 2 {
				// empty string
				resultBytes = util.RemoveBytes(resultBytes, replacer)
			} else {
				// remove leading tailing double quotation mark to prevent invalid string
				variableValueBytes = variableValueBytes[1 : len(variableValueBytes)-1]
				resultBytes = bytes.ReplaceAll(resultBytes, replacer, variableValueBytes)
			}
		} else {
			resultBytes = bytes.ReplaceAll(resultBytes, replacer, variableValueBytes)
		}
	}
	if resultBytes[0] == '"' {
		if len(resultBytes) == 2 {
			resultBytes = []byte("")
		} else {
			resultBytes = resultBytes[1 : len(resultBytes)-1]
		}
	}
	return resultBytes, nil

}

func (resolver *unmarshallResolver) createMapKeyValueByMapKeyKind(value string) (reflect.Value, error) {
	mapKeyType := resolver.ptrToActualValue.Elem().Type().Key()
	mapKeyKind := mapKeyType.Kind()
	// Helper function to convert a string to a numeric type
	convertToNumeric := func(value string) (reflect.Value, error) {
		switch mapKeyKind {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			val, err := strconv.ParseInt(value, 10, 64)
			return reflect.ValueOf(val).Convert(mapKeyType), err
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(value, 10, 64)
			return reflect.ValueOf(val).Convert(mapKeyType), err
		default:
			return reflect.Value{}, ErrorInternalUnsupportedMapKeyKind
		}
	}

	// Convert string to the appropriate type based on mapKeyKind
	switch mapKeyKind {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		if numericValue, err := convertToNumeric(value); err == nil {
			return numericValue, nil
		}
		return reflect.Value{}, NewErrorInternalMapKeyValueKindNotMatch(mapKeyKind.String(), value)

	default:
		return reflect.Value{}, ErrorInternalUnsupportedMapKeyKind
	}
}
