package interpreter

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/util"
)

type unmarshallOptions struct {
	ensureInt     bool
	resolverStack *util.Stack[*unmarshallResolver]
	variables     map[string]interface{}
}

func NewUnMarshallOptions(variables map[string]interface{}) *unmarshallOptions {
	options := &unmarshallOptions{
		ensureInt:     true,
		variables:     variables,
		resolverStack: util.NewStack[*unmarshallResolver](),
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

	var resultBytes []byte
	copy(resultBytes, stringVariable.Value)
	for variableName, replacer := range stringVariable.Variables {
		variableValue, ok := resolver.variables[variableName]
		if !ok {
			continue
		}
		variableValueBytes, err := json.Marshal(variableValue)
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
	return resultBytes, nil

}

func convertNumberBaseOnKind(k reflect.Kind, value interface{}, resolver *unmarshallOptions) interface{} {
	switch k {
	case reflect.Int:
		return int(value.(float64))
	case reflect.Int16:
		return int16(value.(float64))
	case reflect.Int32:
		return int32(value.(float64))
	case reflect.Int64:
		return int64(value.(float64))
	case reflect.Int8:
		return int8(value.(float64))
	case reflect.Float32:
		return float32(value.(float64))
	case reflect.Float64:
		return value
	case reflect.Uint:
		return uint(value.(float64))
	case reflect.Uint8:
		return uint8(value.(float64))
	case reflect.Uint16:
		return uint16(value.(float64))
	case reflect.Uint32:
		return uint32(value.(float64))
	case reflect.Uint64:
		return uint64(value.(float64))
	case reflect.Interface:
		floatVal, ok := value.(float64)
		if ok {
			return int(floatVal)
		} else {
			return value
		}
	default:
		return value
	}
}
