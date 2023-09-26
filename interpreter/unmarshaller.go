package interpreter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/util"
)

type collectionResolver struct {
	astNode              ast.JsonNode
	out                  reflect.Value
	resolverStack        *util.Stack[*collectionResolver]
	variables            map[string]interface{}
	isArray              bool
	isPointer            bool
	awaitingResolveCount int
}

func newCollectionResolver(node ast.JsonNode, out reflect.Value, variables map[string]interface{}, resolverStack *util.Stack[*collectionResolver]) *collectionResolver {
	base := &collectionResolver{
		astNode:              node,
		out:                  out,
		resolverStack:        resolverStack,
		variables:            variables,
		isPointer:            out.Kind() == reflect.Pointer,
		awaitingResolveCount: 0,
	}
	if out.Kind() == reflect.Pointer {
		element := out.Elem()
		base.isArray = element.Kind() == reflect.Array
		base.isPointer = true
	} else {
		base.isArray = out.Kind() == reflect.Array
		base.isPointer = false
	}
	return base
}

func (resolver *collectionResolver) process() error {
	node := resolver.astNode
	switch node.GetNodeType() {
	case ast.AST_OBJECT:
		return resolver.processObject()
	}
	return nil
}

func removeBytes(b []byte, b2remove []byte) []byte {
	parts := bytes.Split(b, b2remove)
	if len(parts) == 1 {
		return b
	}
	var rebuilt []byte
	for _, part := range parts {
		rebuilt = append(rebuilt, part...)
	}
	return rebuilt
}

func (resolver *collectionResolver) resolveVariable(variableNode *ast.JsonExtendedVariableNode) (interface{}, error) {

	variableValue, ok := resolver.variables[variableNode.Variable]
	if !ok {
		return nil, NewVariableNotFound(variableNode.Variable)
	}
	return variableValue, nil
}
func (resolver *collectionResolver) resolveStringVariable(stringVariable *ast.JsonExtendedStringWIthVariableNode) ([]byte, error) {

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
				resultBytes = removeBytes(resultBytes, replacer)
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

func (resolver *collectionResolver) setKVPrimitiveValueBaseStruct(key string, value interface{}) error {
	field := resolver.out.FieldByName(key)
	if !field.IsValid() || !field.CanSet() {
		return NewErrFieldCannotSetOrNotfound(key)
	}
	field.Set(reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValuePointer(key string, value interface{}) error {
	if resolver.out.IsNil() { // *somestruct = nil
		elementType := resolver.out.Type().Elem() // somestruct
		if elementType.Kind() != reflect.Struct {
			return NewErrorInternalExpectingStructButFindOthers(elementType.String())
		}
		newObj := reflect.New(elementType) // *somestruct
		field := newObj.Elem().FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			return NewErrFieldCannotSetOrNotfound(key)
		}
		resolver.out.Set(newObj) // replace the new with above created pointer

	} else {
		if resolver.out.Elem().Kind() != reflect.Struct {
			return NewErrorInternalExpectingStructButFindOthers(resolver.out.Elem().String())
		}
		field := resolver.out.Elem().FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			return NewErrFieldCannotSetOrNotfound(key)
		}

	}
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValueMap(key string, value interface{}) error {
	if resolver.out.IsNil() {
		m := reflect.MakeMap(resolver.out.Type())
		resolver.out.Set(reflect.ValueOf(m))
	}
	resolver.out.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValue(key string, value interface{}) error {
	// for object: struct, *struct, map[string]interface{}, we set the primitive values here
	if resolver.out.Kind() == reflect.Struct {
		return resolver.setKVPrimitiveValueBaseStruct(key, value)
	} else if resolver.out.Kind() == reflect.Pointer { // for pointer it must not be Map, map is ref type in go, it can only be pointer to struct
		return resolver.setKVPrimitiveValuePointer(key, value)
	} else if resolver.out.Kind() == reflect.Map {
		return resolver.setKVPrimitiveValueMap(key, value)
	} else {
		return NewErrorInternalExpectingStructButFindOthers(resolver.out.String())
	}
}

func (resolver *collectionResolver) resolveKVPrimitive(node *ast.JsonKeyValuePairNode) error {

	var key string = node.Key.GetValue()
	if node.Key.GetNodeType() == ast.AST_STRING_VARIABLE {
		resultBytes, err := resolver.resolveStringVariable(node.Key.(*ast.JsonExtendedStringWIthVariableNode))
		if err != nil {
			return err
		}
		if util.RegStringWithVariable.Match(resultBytes) {
			return ErrorStringVariableNotResolveOnKeyLocation
		}
		key = string(resultBytes)
	}
	switch rType := node.Value.(type) {
	case *ast.JsonStringNode:
		err := resolver.setKVPrimitiveValue(key, string(rType.Value))
		if err != nil {
			return err
		}
	case *ast.JsonNumberNode:
		err := resolver.setKVPrimitiveValue(key, rType.Value)
		if err != nil {
			return err
		}
	case *ast.JsonBooleanNode:
		err := resolver.setKVPrimitiveValue(key, rType.Value)
		if err != nil {
			return err
		}
	case *ast.JsonNullNode:
		err := resolver.setKVPrimitiveValue(key, rType.Value)
		if err != nil {
			return err
		}
	case *ast.JsonExtendedStringWIthVariableNode:
		resultBytes, err := resolver.resolveStringVariable(node.Value.(*ast.JsonExtendedStringWIthVariableNode))
		if err != nil {
			return err
		}
		err = resolver.setKVPrimitiveValue(key, string(resultBytes))
		if err != nil {
			return err
		}
	case *ast.JsonExtendedVariableNode:
		variableValue, err := resolver.resolveVariable(node.Value.(*ast.JsonExtendedVariableNode))
		if err != nil {
			return err
		}
		err = resolver.setKVPrimitiveValue(key, variableValue)
		if err != nil {
			return err
		}
	default:
		return ErrorInternalExpectingPrimitive
	}
	return nil
}

func (resolver *collectionResolver) processObject() error {
	node := resolver.astNode.(*ast.JsonObjectNode)
	err := isValidObject(resolver.out)
	if err != nil {
		return err
	}
	for i := len(node.Value); i >= 0; i-- {
		kvpair := node.Value[i]
		switch kvpair.Value.GetNodeType() {
		case ast.AST_ARRAY:
			fallthrough
		case ast.AST_OBJECT:
			resolver.awaitingResolveCount += 1                                                                   // our resolve's awaiting ResolveCount increase, create a new node at the stack
			newResolver := newCollectionResolver(node, resolver.out, resolver.variables, resolver.resolverStack) // the key/value or field/value belongs to the current out, but they are defer to resolve
			resolver.resolverStack.Push(newResolver)
		default:
			err := resolver.resolveKVPrimitive(kvpair)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isValidObject(outValue reflect.Value) error {

	typeItem := reflect.TypeOf(outValue)

	valueItem := reflect.ValueOf(outValue)
	if typeItem.Kind() == reflect.Pointer && valueItem.IsNil() {
		return ErrOutNotNilPointer
	}
	// struct and map is ok, map needs to have string field
	if typeItem.Kind() != reflect.Struct {
		if typeItem.Kind() != reflect.Map {
			return ErrorReflectNotObject
		}
		if typeItem.Key().Kind() != reflect.String {
			return ErrorReflectInvalidMapKey
		}
	}
	return nil
}

func UnmarshallAST(node ast.JsonNode, variables map[string]interface{}, out interface{}) error {
	// deep first traverse the AST
	valueItem := reflect.ValueOf(out)
	if valueItem.Kind() != reflect.Pointer || valueItem.IsNil() {
		return ErrOutNotPointer
	}

	traverseStack := util.NewStack[*collectionResolver]()
	traverseStack.Push(newCollectionResolver(node, valueItem, variables, traverseStack))

	for {
		resolver, err := traverseStack.Peek()
		if err != nil {
			break
		}
		err = resolver.process()
		if err != nil {
			if err != util.ErrorEodOfStack {
				return err
			} else {
				break
			}
		}
		if resolver.awaitingResolveCount == 0 {
			s, _ := traverseStack.Pop() // no awaiting resolve items pop
			fmt.Printf("resolved: %#v \n", s)
		}

	}

	return nil
}
