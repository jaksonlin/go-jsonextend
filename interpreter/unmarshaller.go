package interpreter

import (
	"bytes"
	"encoding/json"
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
	awaitingResolve      bool
	isPointerValue       bool // when parent is not nil, the `out` is always pointer to something, to distinguish whether it is really a pointer.
	parent               *collectionResolver
}

func (c *collectionResolver) getResolveLocation() reflect.Value {
	// when we use reflect.New to create the field, it creates a pointer to the field's type
	if c.isPointerValue {
		return c.out.Elem()
	}
	return c.out
}

func newCollectionResolver(node ast.JsonNode, out reflect.Value, variables map[string]interface{}, resolverStack *util.Stack[*collectionResolver]) *collectionResolver {
	base := &collectionResolver{
		astNode:              node,
		out:                  out,
		variables:            variables,
		resolverStack:        resolverStack,
		awaitingResolveCount: 0,
		awaitingResolve:      false,
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
	c.awaitingResolve = true
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
	c.awaitingResolve = true
	return newResolver, nil
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
	field := resolver.getResolveLocation().FieldByName(key)
	if !field.IsValid() || !field.CanSet() {
		return NewErrFieldCannotSetOrNotfound(key)
	}
	field.Set(reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValuePointer(key string, value interface{}) error {
	resolveLocation := resolver.getResolveLocation()
	if resolveLocation.Elem().Kind() != reflect.Struct {
		return NewErrorInternalExpectingStructButFindOthers(resolveLocation.Elem().String())
	}
	field := resolveLocation.Elem().FieldByName(key)
	if field.IsValid() && field.CanSet() {
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				ptrValue := reflect.New(field.Type().Elem())
				field.Set(ptrValue)
			}
			field = field.Elem()
		}

		if value == nil {
			field.Set(reflect.Zero(field.Type()))
		} else {
			realValue := convertNumberBaseOnKind(field.Kind(), value)
			field.Set(reflect.ValueOf(realValue))
		}

	} else {
		return NewErrFieldCannotSetOrNotfound(key)
	}

	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValueMap(key string, value interface{}) error {
	resolver.getResolveLocation().SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setKVPrimitiveValue(key string, value interface{}) error {
	// for object: struct, *struct, map[string]interface{}, we set the primitive values here
	if resolver.getResolveLocation().Kind() == reflect.Struct { // root is struct
		return resolver.setKVPrimitiveValueBaseStruct(key, value)
	} else if resolver.getResolveLocation().Kind() == reflect.Pointer { // for pointer it must not be Map, map is ref type in go, it can only be pointer to struct
		return resolver.setKVPrimitiveValuePointer(key, value)
	} else if resolver.getResolveLocation().Kind() == reflect.Map {
		return resolver.setKVPrimitiveValueMap(key, value)
	} else {
		return NewErrorInternalExpectingStructButFindOthers(resolver.getResolveLocation().String())
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
	resolver.getResolveLocation().Index(index).Set(reflect.ValueOf(value))
	return nil
}

func (resolver *collectionResolver) setPrimitiveValueToPtrArray(index int, value interface{}) error {
	resolveLocation := resolver.getResolveLocation()
	//.Type().Elem() get the array and further Elem get the array element type this prevent deference on reflect.Value which may be nil
	realValue := convertNumberBaseOnKind(resolveLocation.Type().Elem().Elem().Kind(), value)

	resolver.getResolveLocation().Elem().Index(index).Set(reflect.ValueOf(realValue))

	return nil
}
func convertNumberBaseOnKind(k reflect.Kind, value interface{}) interface{} {
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
	default:
		return value
	}
}
func (resolver *collectionResolver) setPrimitiveValueToPtrSlice(index int, value interface{}) error {
	resolveLocation := resolver.getResolveLocation() // get the correct location to resolve, when the ptrToSlice is in a struct, then we will have **[]someType, where the first * comes form reflect.New, and pointing to the ptrToSlice
	realSlice := resolveLocation.Elem()              // dereference, get the real slice convert *[]int to []int
	if value != nil {
		elemKind := realSlice.Type().Elem().Kind()
		realValue := convertNumberBaseOnKind(elemKind, value)
		realSlice.Index(index).Set(reflect.ValueOf(realValue))
	} else {
		nilValue := reflect.Zero(realSlice.Type().Elem())
		realSlice.Index(index).Set(nilValue)
	}
	return nil
}

func (resolver *collectionResolver) setPrimitiveValueToSlice(index int, value interface{}) error {
	realSlice := resolver.getResolveLocation() //1. bare slice; 2. *[]int, * comes from the reflect.New for some struct
	if value != nil {
		realSlice.Index(index).Set(reflect.ValueOf(value))
	} else {
		nilValue := reflect.Zero(realSlice.Type().Elem())
		realSlice.Index(index).Set(nilValue)
	}
	return nil
}

func (resolver *collectionResolver) setArrayPrimitiveValue(index int, value interface{}) error {
	// for array: array, *array, slice
	resolveLocation := resolver.getResolveLocation()
	if resolveLocation.Kind() == reflect.Array { // root is struct
		return resolver.setPrimitiveToArray(index, value)
	} else if resolveLocation.Kind() == reflect.Pointer { // for pointer it must not be Map, map is ref type in go, it can only be pointer to struct
		if resolveLocation.Elem().Kind() == reflect.Array {
			return resolver.setPrimitiveValueToPtrArray(index, value)
		} else if resolveLocation.Elem().Kind() == reflect.Slice {
			return resolver.setPrimitiveValueToPtrSlice(index, value)
		} else {
			return ErrorInternalExpectingArrayLikeObject
		}
	} else if resolveLocation.Kind() == reflect.Slice {
		return resolver.setPrimitiveValueToSlice(index, value)
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
	if resolver.getResolveLocation().Kind() == reflect.Struct {
		return nil // already struct, seems nothing to do
	} else if resolver.getResolveLocation().Kind() == reflect.Map {
		if resolver.getResolveLocation().IsNil() {
			m := reflect.MakeMap(resolver.getResolveLocation().Type())
			resolver.getResolveLocation().Set(reflect.ValueOf(m))
		}
	} else if resolver.getResolveLocation().Kind() == reflect.Pointer {
		if resolver.getResolveLocation().IsNil() {
			elementType := resolver.getResolveLocation().Type().Elem() // somestruct
			if elementType.Kind() != reflect.Struct {
				return NewErrorInternalExpectingStructButFindOthers(elementType.String())
			}
			newObj := reflect.New(elementType)        // *somestruct
			resolver.getResolveLocation().Set(newObj) // replace the new with above created pointer
		}
	} else {
		NewErrorInternalExpectingStructButFindOthers(resolver.getResolveLocation().String())
	}

	return nil
}

func (resolver *collectionResolver) processObject() error {
	node := resolver.astNode.(*ast.JsonObjectNode)

	err := resolver.initObject()
	if err != nil {
		return err
	}
	err = isValidObject(resolver.getResolveLocation())
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
	return resolver.Enclose()
}

func (resolver *collectionResolver) initArrayLikeResolver(cap int) error {
	resolverOutElement := resolver.getResolveLocation()
	if resolverOutElement.Kind() == reflect.Pointer {
		resolverOutElement = resolverOutElement.Elem()
	}
	if resolverOutElement.Kind() == reflect.Slice {
		if resolverOutElement.IsNil() {
			sliceType := reflect.SliceOf(resolverOutElement.Type().Elem())
			sliceValue := reflect.MakeSlice(sliceType, cap, cap)
			resolverOutElement.Set(sliceValue)
		}
	} else if resolver.getResolveLocation().Kind() == reflect.Array {
		arrElementType := resolver.getResolveLocation().Type().Elem()
		arrayType := reflect.ArrayOf(cap, arrElementType)
		arrayValue := reflect.New(arrayType).Elem()
		resolver.getResolveLocation().Set(arrayValue)
	} else {
		return ErrorInternalExpectingArrayLikeObject
	}
	return nil
}

func (resolver *collectionResolver) processArray() error {
	node := resolver.astNode.(*ast.JsonArrayNode)

	err := resolver.initArrayLikeResolver(len(node.Value))
	if err != nil {
		return err
	}
	err = isValidArray(resolver.getResolveLocation())
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

	checkerValue := outValue
	if checkerValue.Kind() == reflect.Pointer {
		if checkerValue.IsNil() {
			return ErrSliceOrArrayNotInit
		}
		checkerValue = checkerValue.Elem()
	}

	if checkerValue.Kind() != reflect.Slice && checkerValue.Kind() != reflect.Array {
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
		// no dependency waiting
		if !resolver.awaitingResolve {
			// process elements
			err = resolver.process()
			if err != nil {
				if err != util.ErrorEndOfStack {
					return err
				} else {
					break
				}
			}
			// if there's no awaiting dependency, pop
			if !resolver.awaitingResolve {
				traverseStack.Pop() // no awaiting resolve items pop
			}
		} else {
			// when the dependecies are resolved, enclose and pop
			if resolver.awaitingResolveCount != 0 {
				return ErrorInternalNoneResolvable
			}
			if err := resolver.Enclose(); err != nil {
				return err
			}
			traverseStack.Pop()
		}

	}

	return nil
}
