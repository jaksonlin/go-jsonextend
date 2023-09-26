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
	arrayIndex           int
	objectKey            string
	awaitingResolveCount int
	hasInit              bool
	isPointerValue       bool // when parent is not nil, the `out` is always pointer to something, to distinguish whether it is really a pointer.
	parent               *collectionResolver
}

func newCollectionResolver(node ast.JsonNode, out reflect.Value, variables map[string]interface{}, resolverStack *util.Stack[*collectionResolver]) *collectionResolver {
	base := &collectionResolver{
		astNode:              node,
		out:                  out,
		variables:            variables,
		resolverStack:        resolverStack,
		awaitingResolveCount: 0,
		hasInit:              false,
		parent:               nil,
		arrayIndex:           -1,
	}
	return base
}

func (c *collectionResolver) createReflectionValueForResolver(kvValueElementType reflect.Type, node ast.JsonCollectionNode) (*collectionResolver, error) {
	// we only deal with the kv value being collection, primitives are already set
	var resultValue reflect.Value
	isPointerValue := false
	if kvValueElementType.Kind() == reflect.Slice {
		elementType := kvValueElementType.Elem()
		sliceType := reflect.SliceOf(elementType)
		resultValue = reflect.MakeSlice(sliceType, node.Length(), node.Length())
	} else if kvValueElementType.Kind() == reflect.Array {
		elementType := kvValueElementType.Elem()
		arrayType := reflect.ArrayOf(node.Length(), elementType)
		resultValue = reflect.New(arrayType)
	} else if kvValueElementType.Kind() == reflect.Map {
		resultValue = reflect.MakeMap(kvValueElementType)
	} else if kvValueElementType.Kind() == reflect.Struct {
		resultValue = reflect.New(kvValueElementType)
	} else if kvValueElementType.Kind() == reflect.Pointer {
		if kvValueElementType.Elem().Kind() != reflect.Struct && kvValueElementType.Elem().Kind() != reflect.Array {
			return nil, NewErrorInternalFieldTypeNotMatchAST(kvValueElementType.Elem().Name())
		}
		isPointerValue = true
		resultValue = reflect.New(kvValueElementType)
	} else {
		return nil, NewErrorInternalFieldTypeNotMatchAST(kvValueElementType.Name())
	}
	newResolver := newCollectionResolver(node, resultValue, c.variables, c.resolverStack)
	newResolver.isPointerValue = isPointerValue
	return newResolver, nil
}

// create resolver to resolving the things in kv's value
func (c *collectionResolver) createChildResolverWithObjectKey(key string, node ast.JsonCollectionNode) (*collectionResolver, error) {
	// root: struct, *struct, map

	// 1. get the keys coresponding value type
	var childElementType reflect.Type
	if c.out.Kind() == reflect.Map {
		childElementType = c.out.Type().Elem() // map[string]something, something
	} else if c.out.Kind() == reflect.Struct {
		fieldInfo := c.out.FieldByName(key) // struct field
		childElementType = fieldInfo.Type()
	} else if c.out.Kind() == reflect.Pointer {
		if c.out.Elem().Kind() != reflect.Struct { // *struct, deference field
			return nil, NewErrorInternalExpectingStructInsidePointerButFindOthers(c.out.Elem().Kind().String())
		}
		fieldInfo := c.out.Elem().FieldByName(key)
		childElementType = fieldInfo.Type()
	} else {
		// current resolver's out is not object type capable
		return nil, NewErrorInternalExpectingStructButFindOthers(childElementType.Name())
	}
	// 2. create the collection's reflection value representative
	newResolver, err := c.createReflectionValueForResolver(childElementType, node)
	if err != nil {
		return nil, err
	}
	// 3. bind parent
	newResolver.objectKey = key
	newResolver.parent = c
	c.awaitingResolveCount += 1
	return newResolver, nil
}

func (c *collectionResolver) createChildResolverWithArrayIndex(index int, node ast.JsonCollectionNode) (*collectionResolver, error) {
	// root: slice, array, *array

	// 1. get the keys coresponding value type
	var childElementType reflect.Type
	if c.out.Kind() == reflect.Slice || c.out.Kind() == reflect.Array {
		childElementType = c.out.Type().Elem()
	} else if c.out.Kind() == reflect.Pointer {
		if c.out.Type().Elem().Kind() != reflect.Array {
			return nil, NewErrorInternalExpectingArrayInsidePointerButFindOthers(c.out.Type().Elem().Kind().String())
		}
		childElementType = c.out.Type().Elem().Elem() // we do not use c.out.Elem() to avoid nil pointer deferenc, use c.out.Type().Elem() to get the Array type and use one more Elem() to get the array element type
	}

	// 2. create the collection's reflection value representative
	newResolver, err := c.createReflectionValueForResolver(childElementType, node)
	if err != nil {
		return nil, err
	}
	newResolver.arrayIndex = index
	newResolver.parent = c
	c.awaitingResolveCount += 1
	return newResolver, nil
}

func (c *collectionResolver) SetInit() {
	c.hasInit = true
}

func (c *collectionResolver) EncloseArray() error {
	if c.parent.out.Kind() == reflect.Array {
		c.parent.out.Index(c.arrayIndex).Set(c.out.Elem()) // deference the pointer created
	} else if c.parent.out.Kind() == reflect.Slice {
		c.parent.out = reflect.Append(c.parent.out, c.out.Elem())
	} else if c.parent.out.Kind() == reflect.Pointer {
		if c.parent.out.Elem().Kind() == reflect.Array {
			c.parent.out.Elem().Index(c.arrayIndex).Set(c.out.Elem())
		} else {
			return NewErrorInternalExpectingArrayInsidePointerButFindOthers(c.parent.out.Elem().String())
		}
	} else {
		return ErrorInternalExpectingArrayLikeObject
	}
	return nil
}

func (c *collectionResolver) EncloseObject() error {
	if c.parent.out.Kind() == reflect.Struct {
		field := c.parent.out.FieldByName(c.objectKey)
		field.Set(c.out.Elem())
	} else if c.parent.out.Kind() == reflect.Map {
		c.parent.out.SetMapIndex(reflect.ValueOf(c.objectKey), c.out.Elem())
	} else if c.parent.out.Kind() == reflect.Pointer {
		if c.parent.out.Elem().Kind() == reflect.Struct {
			field := c.parent.out.Elem().FieldByName(c.objectKey)
			field.Set(c.out.Elem())
		} else {
			return NewErrorInternalExpectingStructInsidePointerButFindOthers(c.parent.out.Elem().String())
		}
	} else {
		return NewErrorInternalExpectingStructButFindOthers(c.parent.out.Kind().String())
	}
	return nil
}

func (c *collectionResolver) Enclose() error {
	// no parent, no need to enclose
	if c.parent == nil {
		return nil
	}

	// have unresolve child item, cannot enclose now
	if c.awaitingResolveCount > 0 {
		return nil
	}

	// awaitngResolveCount == 0 && c.parent !=nil
	if c.arrayIndex >= 0 {
		err := c.EncloseArray()
		if err != nil {
			return err
		}
		c.parent.awaitingResolveCount -= 1
	} else if c.objectKey != "" {
		err := c.EncloseObject()
		if err != nil {
			return err
		}
		c.parent.awaitingResolveCount -= 1
	}
	return nil
}

func (resolver *collectionResolver) process() error {
	node := resolver.astNode
	switch node.GetNodeType() {
	case ast.AST_OBJECT:
		return resolver.processObject()
	case ast.AST_ARRAY:
		return resolver.processArray()
	default:
		return ErrorUnmarshalStackNoKV
	}
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
			if field.Kind() == reflect.Int {
				field.SetInt(int64(value.(float64)))
			} else {
				field.Set(reflect.ValueOf(value))
			}
		} else {
			return NewErrFieldCannotSetOrNotfound(key)
		}

	}
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValueMap(key string, value interface{}) error {
	resolver.out.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValue(key string, value interface{}) error {
	// for object: struct, *struct, map[string]interface{}, we set the primitive values here
	if resolver.out.Kind() == reflect.Struct { // root is struct
		return resolver.setKVPrimitiveValueBaseStruct(key, value)
	} else if resolver.out.Kind() == reflect.Pointer { // for pointer it must not be Map, map is ref type in go, it can only be pointer to struct
		return resolver.setKVPrimitiveValuePointer(key, value)
	} else if resolver.out.Kind() == reflect.Map {
		return resolver.setKVPrimitiveValueMap(key, value)
	} else {
		return NewErrorInternalExpectingStructButFindOthers(resolver.out.String())
	}
}
func (resolver *collectionResolver) getKeyValueKeyFromKvPair(node *ast.JsonKeyValuePairNode) (string, error) {
	var key string = node.Key.GetValue()
	if node.Key.GetNodeType() == ast.AST_STRING_VARIABLE {
		resultBytes, err := resolver.resolveStringVariable(node.Key.(*ast.JsonExtendedStringWIthVariableNode))
		if err != nil {
			return "", err
		}
		if util.RegStringWithVariable.Match(resultBytes) {
			return "", ErrorStringVariableNotResolveOnKeyLocation
		}
		key = string(resultBytes)
	}
	return key, nil
}

func (resolver *collectionResolver) resolveKVPrimitiveValue(key string, node ast.JsonNode) error {

	switch rType := node.(type) {
	case *ast.JsonStringNode:
		err := resolver.setKVPrimitiveValue(key, rType.GetValue())
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
		resultBytes, err := resolver.resolveStringVariable(node.(*ast.JsonExtendedStringWIthVariableNode))
		if err != nil {
			return err
		}
		err = resolver.setKVPrimitiveValue(key, string(resultBytes))
		if err != nil {
			return err
		}
	case *ast.JsonExtendedVariableNode:
		variableValue, err := resolver.resolveVariable(node.(*ast.JsonExtendedVariableNode))
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

func (resolver *collectionResolver) setPrimitiveToArray(index int, value interface{}) error {
	resolver.out.Index(index).Set(reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setPrimitiveValueToPtrArray(index int, value interface{}) error {
	resolver.out.Elem().Index(index).Set(reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setPrimitiveValueToSlice(value interface{}) error {
	resolver.out = reflect.AppendSlice(resolver.out, reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setArrayPrimitiveValue(index int, value interface{}) error {
	// for array: array, *array, slice
	if resolver.out.Kind() == reflect.Array { // root is struct
		return resolver.setPrimitiveToArray(index, value)
	} else if resolver.out.Kind() == reflect.Pointer { // for pointer it must not be Map, map is ref type in go, it can only be pointer to struct
		return resolver.setPrimitiveValueToPtrArray(index, value)
	} else if resolver.out.Kind() == reflect.Slice {
		return resolver.setPrimitiveValueToSlice(value)
	} else {
		return ErrorInternalExpectingArrayLikeObject
	}
}

func (resolver *collectionResolver) resolveArrayElementPrimitive(index int, node ast.JsonNode) error {

	switch rType := node.(type) {
	case *ast.JsonStringNode:
		err := resolver.setArrayPrimitiveValue(index, string(rType.Value))
		if err != nil {
			return err
		}
	case *ast.JsonNumberNode:
		err := resolver.setArrayPrimitiveValue(index, rType.Value)
		if err != nil {
			return err
		}
	case *ast.JsonBooleanNode:
		err := resolver.setArrayPrimitiveValue(index, rType.Value)
		if err != nil {
			return err
		}
	case *ast.JsonNullNode:
		err := resolver.setArrayPrimitiveValue(index, rType.Value)
		if err != nil {
			return err
		}
	case *ast.JsonExtendedStringWIthVariableNode:
		resultBytes, err := resolver.resolveStringVariable(node.(*ast.JsonExtendedStringWIthVariableNode))
		if err != nil {
			return err
		}
		err = resolver.setArrayPrimitiveValue(index, string(resultBytes))
		if err != nil {
			return err
		}
	case *ast.JsonExtendedVariableNode:
		variableValue, err := resolver.resolveVariable(node.(*ast.JsonExtendedVariableNode))
		if err != nil {
			return err
		}
		err = resolver.setArrayPrimitiveValue(index, variableValue)
		if err != nil {
			return err
		}
	default:
		return ErrorInternalExpectingPrimitive
	}
	return nil
}

func (resolver *collectionResolver) initObject() error {
	if resolver.out.Kind() == reflect.Struct {
		return nil // already struct, seems nothing to do
	} else if resolver.out.Kind() == reflect.Map {
		if resolver.out.IsNil() {
			m := reflect.MakeMap(resolver.out.Type())
			resolver.out.Set(reflect.ValueOf(m))
		}
	} else if resolver.out.Kind() == reflect.Pointer {
		if resolver.out.IsNil() {
			elementType := resolver.out.Type().Elem() // somestruct
			if elementType.Kind() != reflect.Struct {
				return NewErrorInternalExpectingStructButFindOthers(elementType.String())
			}
			newObj := reflect.New(elementType) // *somestruct
			resolver.out.Set(newObj)           // replace the new with above created pointer
		}
	} else {
		NewErrorInternalExpectingStructButFindOthers(resolver.out.String())
	}

	resolver.SetInit()
	return nil
}

func (resolver *collectionResolver) processObject() error {
	node := resolver.astNode.(*ast.JsonObjectNode)

	err := resolver.initObject()
	if err != nil {
		return err
	}
	err = isValidObject(resolver.out)
	if err != nil {
		return err
	}

	for i := len(node.Value) - 1; i >= 0; i-- {
		kvpair := node.Value[i]
		key, err := resolver.getKeyValueKeyFromKvPair(kvpair)
		if err != nil {
			return err
		}
		switch kvpair.Value.GetNodeType() {
		case ast.AST_ARRAY:
			fallthrough
		case ast.AST_OBJECT:

			newResolver, err := resolver.createChildResolverWithObjectKey(key, kvpair.Value.(ast.JsonCollectionNode))
			if err != nil {
				return err
			}
			resolver.resolverStack.Push(newResolver)
		default:
			err := resolver.resolveKVPrimitiveValue(key, kvpair.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (resolver *collectionResolver) initArrayLikeResolver(cap int) error {
	if resolver.out.Kind() == reflect.Slice {
		if resolver.out.IsNil() {
			sliceType := resolver.out.Type().Elem()
			sliceValue := reflect.MakeSlice(sliceType, cap, cap)
			resolver.out.Set(sliceValue)
		}
	} else if resolver.out.Kind() == reflect.Array {
		arrElementType := resolver.out.Type().Elem()
		arrayType := reflect.ArrayOf(cap, arrElementType)
		arrayValue := reflect.New(arrayType).Elem()
		resolver.out.Set(arrayValue)
	} else if resolver.out.Kind() == reflect.Pointer {
		if resolver.out.Elem().Kind() != reflect.Array {
			return ErrorInternalPtrToArrayFindNotArray
		}
		arrElementType := resolver.out.Elem().Type().Elem() // get the underlying type in *ptrArray([]int). get the int
		arrayType := reflect.ArrayOf(cap, arrElementType)   // []int
		arrValue := reflect.New(arrayType)                  // ptr to some type, just use it as resolver.out is Pointer
		resolver.out.Set(arrValue)
	} else {
		return ErrorInternalExpectingArrayLikeObject
	}
	resolver.SetInit()
	return nil
}

func (resolver *collectionResolver) processArray() error {
	node := resolver.astNode.(*ast.JsonArrayNode)

	err := resolver.initArrayLikeResolver(len(node.Value))
	if err != nil {
		return err
	}
	err = isValidArray(resolver.out)
	if err != nil {
		return err
	}
	for i := len(node.Value) - 1; i >= 0; i-- {
		element := node.Value[i]
		switch element.GetNodeType() {
		case ast.AST_ARRAY:
			fallthrough
		case ast.AST_OBJECT:
			resolver, err := resolver.createChildResolverWithArrayIndex(i, node)
			if err != nil {
				return err
			}
			resolver.resolverStack.Push(resolver)
		default:
			err := resolver.resolveArrayElementPrimitive(i, element)
			if err != nil {
				return err
			}
		}
	}

	return resolver.Enclose()
}

func isValidObject(outValue reflect.Value) error {

	typeItem := reflect.TypeOf(outValue)

	valueItem := reflect.ValueOf(outValue)
	if typeItem.Kind() == reflect.Pointer {
		typeItem = typeItem.Elem()
		if valueItem.IsNil() {
			return ErrOutNotNilPointer
		}
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

func isValidArray(outValue reflect.Value) error {

	typeItem := reflect.TypeOf(outValue)

	valueItem := reflect.ValueOf(outValue)

	if typeItem.Kind() == reflect.Slice {
		if valueItem.IsNil() {
			return ErrSliceNotInit
		}
		return nil
	}

	if typeItem.Kind() == reflect.Pointer {
		typeItem = typeItem.Elem()
		if typeItem.Kind() != reflect.Array {
			return ErrorInternalExpectingArrayLikeObject
		}
		if valueItem.IsNil() {
			return ErrArrayNotInit
		}
		return nil
	}

	// struct and map is ok, map needs to have string field
	if typeItem.Kind() != reflect.Array {
		return ErrorInternalExpectingArrayLikeObject
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
