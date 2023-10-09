package interpreter

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

var dummyMap map[string]interface{}

// we hold a pointer to these field rather than keeping the field's reflect.Value directly
// number of pointer of the original field, not including the one we create to hold the filed
type unmarshallResolver struct {
	options              *unmarshallOptions
	astNode              ast.JsonNode
	outElementKind       reflect.Kind
	arrayIndex           int
	awaitingResolveCount int
	numberOfPointer      int
	awaitingResolve      bool
	isPointerValue       bool
	IsNil                bool // when our ptrToActualValue holds a struct, there's no way to tell if we are accepting a nil for it.
	objectKey            string
	parent               *unmarshallResolver
	ptrToActualValue     reflect.Value // single ptr to no matter what actual value is (for *****int, keeps only *int to the actual value)
	fields               map[string]*util.JSONStructField
	hasUnmarshaller      bool
	tagOption            string
}

func (resolver *unmarshallResolver) collectAllFields() error {
	if len(resolver.fields) > 0 {
		return nil
	}

	result := util.FlattenJsonStructForUnmarshal(resolver.ptrToActualValue.Elem())
	if result == nil {
		return NewErrorInternalExpectingStructButFindOthers(resolver.ptrToActualValue.Elem().Kind().String())
	}
	resolver.fields = result
	return nil

}

// a story for align to the go's json unmarshall is that, when the field is a pointer, and it points to a nil value, the unmarshall will resolve to a `nil pointer` not `pointer to nil value`
// therefore below is a if check to set the value to `zero` rather than `nilValue of some type`
func (resolver *unmarshallResolver) resolveSliceDependency(dependentResolver *unmarshallResolver) error {
	dependentValue := dependentResolver.restoreValue()
	if dependentResolver.IsNil || dependentResolver.isPointerValue && (dependentValue.Elem().Kind() == reflect.Slice ||
		dependentValue.Elem().Kind() == reflect.Interface ||
		dependentValue.Elem().Kind() == reflect.Map) && dependentValue.Elem().IsNil() {
		resolver.ptrToActualValue.Elem().Index(dependentResolver.arrayIndex).SetZero()
	} else {
		resolver.ptrToActualValue.Elem().Index(dependentResolver.arrayIndex).Set(dependentValue.Convert(resolver.ptrToActualValue.Elem().Type().Elem()))
	}
	return nil
}
func (resolver *unmarshallResolver) resolveStructDependency(dependentResolver *unmarshallResolver) error {
	field, err := resolver.getFieldByTag(dependentResolver.objectKey)
	if err != nil {
		return err
	}

	dependentValue := dependentResolver.restoreValue()
	if dependentResolver.IsNil || dependentResolver.isPointerValue && (dependentValue.Elem().Kind() == reflect.Slice ||
		dependentValue.Elem().Kind() == reflect.Interface ||
		dependentValue.Elem().Kind() == reflect.Map) && dependentValue.Elem().IsNil() {
		field.FieldValue.SetZero()
	} else {
		field.FieldValue.Set(dependentValue.Convert(field.FieldValue.Type()))
	}
	return nil
}

func (resolver *unmarshallResolver) resolveMapDependency(dependentResolver *unmarshallResolver) error {
	key, err := resolver.createMapKeyValueByMapKeyKind(dependentResolver.objectKey)
	if err != nil {
		return err
	}

	dependentValue := dependentResolver.restoreValue()
	if dependentResolver.isPointerValue && (dependentValue.Elem().Kind() == reflect.Slice ||
		dependentValue.Elem().Kind() == reflect.Interface ||
		dependentValue.Elem().Kind() == reflect.Map) && dependentValue.Elem().IsNil() {
		mapElementType := resolver.ptrToActualValue.Elem().Type().Elem()
		mapElementZero := reflect.Zero(mapElementType)
		resolver.ptrToActualValue.Elem().SetMapIndex(key, mapElementZero)
	} else {
		resolver.ptrToActualValue.Elem().SetMapIndex(key, dependentValue.Convert(resolver.ptrToActualValue.Elem().Type().Elem()))
	}
	return nil
}
func (resolver *unmarshallResolver) resolveInterfaceDependency(dependentResolver *unmarshallResolver) error {

	if dependentResolver.arrayIndex != -1 { //interface holding slice
		return resolver.resolveSliceDependency(dependentResolver)
	} else {
		if len(dependentResolver.objectKey) > 0 { // interface holding map
			return resolver.resolveMapDependency(dependentResolver)
		} else { // we need array index or key to resolve the location of the dependent Value
			return ErrorInternalDependentResolverHasOnResolveLocation
		}
	}

}
func (resolver *unmarshallResolver) resolveDependency(dependentResolver *unmarshallResolver) error {
	resolver.awaitingResolveCount -= 1
	if resolver.outElementKind == reflect.Array || resolver.outElementKind == reflect.Slice {

		return resolver.resolveSliceDependency(dependentResolver)

	} else if resolver.outElementKind == reflect.Struct {

		return resolver.resolveStructDependency(dependentResolver)

	} else if resolver.outElementKind == reflect.Map {

		return resolver.resolveMapDependency(dependentResolver)

	} else if resolver.outElementKind == reflect.Interface {
		return resolver.resolveInterfaceDependency(dependentResolver)
	} else {
		return ErrorPrimitiveTypeCannotResolveDependency
	}

}

func (resolver *unmarshallResolver) setValue(value interface{}) {
	if value == nil {
		nilValue := reflect.Zero(resolver.ptrToActualValue.Elem().Type())
		resolver.ptrToActualValue.Elem().Set(nilValue)
	} else {
		resolver.ptrToActualValue.Elem().Set(reflect.ValueOf(value).Convert(resolver.ptrToActualValue.Elem().Type()))
	}
}

// return the actual reflect.Value in the resolver, the resolver is desinged to hold a pointer to anything it keeps
// when the actual field is pointer type, you need a pointer to the actual field to set its Elem to the retrun value from this func
func (resolver *unmarshallResolver) restoreValue() reflect.Value {
	if !resolver.isPointerValue {
		return resolver.ptrToActualValue.Elem() // remove the pointer we add (newUnmarshallResolver)
	} else {
		// the field is *interface{}
		if resolver.outElementKind == reflect.Interface {
			var value interface{} = resolver.ptrToActualValue.Elem().Interface()
			resolver.ptrToActualValue = reflect.ValueOf(&value)
		}
		// actual value
		if resolver.numberOfPointer == 1 {
			return resolver.ptrToActualValue // just use our value holder
		}
		var resultPtr = resolver.ptrToActualValue // what we hold *someStruct, actual ***int (isPointerValue=true, numberOfPointer=3), then we start from *someStruct
		var tmpPtr reflect.Value
		for i := resolver.numberOfPointer; i > 1; i-- {
			tmpPtr = reflect.New(resultPtr.Type())
			tmpPtr.Elem().Set(resultPtr)
			resultPtr = tmpPtr
		}
		return resultPtr // when the original field is ***int, and we create the ***int, you need to take the address of the origianl field to set this value in.(use a ****int to set its element to this func's return value)
	}
}

func (resolver *unmarshallResolver) bindObjectParent(key string, parent *unmarshallResolver) {
	resolver.objectKey = key
	resolver.parent = parent
	parent.awaitingResolveCount += 1
	parent.awaitingResolve = true
}
func (resolver *unmarshallResolver) bindArrayLikeParent(index int, parent *unmarshallResolver) {
	resolver.arrayIndex = index
	resolver.parent = parent
	parent.awaitingResolveCount += 1
	parent.awaitingResolve = true
}

func createPtrToSliceValue(nodeToWork ast.JsonNode, someOutType reflect.Type) (reflect.Value, ast.JsonNode, error) {
	var isNil = nodeToWork.GetNodeType() == ast.AST_NULL
	numberOfElement := 0
	var convertedNode ast.JsonNode = nil
	if !isNil {
		// in golang json processing for slice of Uint8 it will convert to base64
		if someOutType.Elem().Kind() == reflect.Uint8 && nodeToWork.GetNodeType() == ast.AST_STRING {
			n, err := nodeToWork.(*ast.JsonStringNode).ToArrayNode()
			if err != nil {
				return reflect.Value{}, nil, err
			}
			numberOfElement = n.Length()
			convertedNode = n
		} else {
			n, ok := nodeToWork.(*ast.JsonArrayNode)
			if !ok {
				return reflect.Value{}, nil, ErrorInternalExpectingArrayLikeObject
			}
			numberOfElement = n.Length()
			convertedNode = n
		}
	}

	sliceType := reflect.SliceOf(someOutType.Elem())
	sliceValue := reflect.MakeSlice(sliceType, numberOfElement, numberOfElement) // use index to manipulate the slice
	ptrToActualValue := reflect.New(sliceValue.Type())
	ptrToActualValue.Elem().Set(sliceValue)
	return ptrToActualValue, convertedNode, nil
}

func createPtrToArrayValue(nodeToWork ast.JsonNode, someOutType reflect.Type) (reflect.Value, error) {
	n, ok := nodeToWork.(*ast.JsonArrayNode)
	if !ok {
		return reflect.Value{}, ErrorInternalExpectingArrayLikeObject
	}
	numberOfElement := n.Length()
	arrayType := reflect.ArrayOf(numberOfElement, someOutType.Elem())
	ptrToActualValue := reflect.New(arrayType)
	return ptrToActualValue, nil
}

func createPtrToInterfaceValue(nodeToWork ast.JsonNode, someOutType reflect.Type) (reflect.Value, error) {
	var ptrToActualValue reflect.Value
	// someField: interface{}
	if nodeToWork.GetNodeType() == ast.AST_ARRAY {
		numberOfElement := nodeToWork.(*ast.JsonArrayNode).Length()
		sliceType := reflect.SliceOf(reflect.TypeOf((*interface{})(nil)).Elem())
		sliceValue := reflect.MakeSlice(sliceType, numberOfElement, numberOfElement) // use index to manipulate the slice
		ptrToActualValue = reflect.New(sliceValue.Type())
		ptrToActualValue.Elem().Set(sliceValue)
	} else if nodeToWork.GetNodeType() == ast.AST_OBJECT {
		newMap := reflect.MakeMap(reflect.TypeOf(dummyMap))
		ptrToActualValue = reflect.New(newMap.Type())
		ptrToActualValue.Elem().Set(newMap)
	} else {
		ptrToActualValue = reflect.New(someOutType)
		ptrToActualValue.Elem().Set(reflect.Zero(someOutType))
	}
	return ptrToActualValue, nil
}

func newUnmarshallResolver(
	node ast.JsonNode,
	outType reflect.Type,
	options *unmarshallOptions,
	tagOption string) (*unmarshallResolver, error) {
	var nodeToWork ast.JsonNode = node
	someOutType := outType
	numberOfPointer := 0
	var elementKind reflect.Kind

	isPointerValue := someOutType.Kind() == reflect.Pointer
	for someOutType.Kind() == reflect.Pointer {
		someOutType = someOutType.Elem()
		numberOfPointer += 1
	}
	var ptrToActualValue reflect.Value
	// use a pointer to hold no matter what it is inside
	switch someOutType.Kind() {
	case reflect.Slice:
		ptr, convertedNode, err := createPtrToSliceValue(nodeToWork, someOutType)
		if err != nil {
			return nil, err
		}
		if convertedNode != nil {
			nodeToWork = convertedNode
		}
		ptrToActualValue = ptr
		elementKind = reflect.Slice
	case reflect.Array:
		ptr, err := createPtrToArrayValue(nodeToWork, someOutType)
		if err != nil {
			return nil, err
		}
		ptrToActualValue = ptr
		elementKind = reflect.Array
	case reflect.Map:
		newMap := reflect.MakeMap(someOutType)
		ptrToActualValue = reflect.New(newMap.Type())
		ptrToActualValue.Elem().Set(newMap)
		elementKind = reflect.Map
	case reflect.Struct:
		ptrToActualValue = reflect.New(someOutType) //*Struct
		elementKind = reflect.Struct
	case reflect.Interface:
		// someField: interface{}
		ptr, err := createPtrToInterfaceValue(nodeToWork, someOutType)
		if err != nil {
			return nil, err
		}
		ptrToActualValue = ptr
		elementKind = reflect.Interface
	default: // primitives
		ptrToActualValue = reflect.New(someOutType)
		ptrToActualValue.Elem().Set(reflect.Zero(someOutType))
		elementKind = someOutType.Kind()
	}
	// we only support pointer receiver unmarshaler, therefore pass in the Pointer not the pointer to element
	hasUnmarshaller := implementsUnmarshaler(ptrToActualValue.Type())

	base := &unmarshallResolver{
		options:              options,
		astNode:              nodeToWork,
		ptrToActualValue:     ptrToActualValue,
		awaitingResolveCount: 0,
		awaitingResolve:      false,
		parent:               nil,
		arrayIndex:           -1,
		numberOfPointer:      numberOfPointer,
		isPointerValue:       isPointerValue,
		outElementKind:       elementKind,
		IsNil:                nodeToWork.GetNodeType() == ast.AST_NULL,
		hasUnmarshaller:      hasUnmarshaller,
		tagOption:            tagOption,
	}
	return base, nil
}
func implementsUnmarshaler(t reflect.Type) bool {
	// Check for pointer type if the provided type isn't a pointer.
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}

	method, ok := t.MethodByName("UnmarshalJSON")
	if !ok {
		return false
	}

	// Check method signature
	if method.Type.NumIn() != 2 || method.Type.In(1) != reflect.TypeOf([]byte{}) || method.Type.NumOut() != 1 || method.Type.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return false
	}

	return true
}

var _ ast.JsonVisitor = &unmarshallResolver{}

func (resolver *unmarshallResolver) resolveByCustomizeObjectUnmarshal(node ast.JsonNode) error {

	unmarshalMethod := resolver.ptrToActualValue.MethodByName("UnmarshalJSON")

	payload, err := InterpretAST(node, resolver.options.variables, resolver.options.marshaler)
	if err != nil {
		return err
	}
	result := unmarshalMethod.Call([]reflect.Value{reflect.ValueOf(payload)})
	unmarshalError := result[0].Interface()
	if unmarshalError == nil {
		return resolver.resolve()
	} else {
		if err, ok := unmarshalError.(error); ok {
			return err
		} else {
			return ErrorInvalidUnmarshalResult
		}
	}
}

func (resolver *unmarshallResolver) VisitArrayNode(node *ast.JsonArrayNode) error {
	// fill the values in the reflection.Value
	if resolver.hasUnmarshaller {
		return resolver.resolveByCustomizeObjectUnmarshal(node)
	}
	for i := len(node.Value) - 1; i >= 0; i-- {
		resolver, err := resolver.createArrayElementResolver(i, node.Value[i])
		if err != nil {
			return err
		}
		resolver.options.resolverStack.Push(resolver)
	}
	return resolver.resolve()

}

// this is only visit from VisitObjectNode, it does not resolve any value, no need to check unmarshaler call
func (resolver *unmarshallResolver) VisitKeyValuePairNode(node *ast.JsonKeyValuePairNode) error {
	// here the resolver.ptrToActualValue is a pointer to the object that holds this kv
	// translate the key to the field name
	key, err := resolver.processKVKeyNode(node.Key)
	if err != nil {
		return err
	}

	newResolver, err := resolver.processKVValueNode(key, node.Value)
	if err != nil {
		var notFind ErrorFieldNotExist
		if errors.As(err, &notFind) {
			// pass this field
			return nil
		}
		return err
	}

	resolver.options.resolverStack.Push(newResolver)

	return nil
}

func (resolver *unmarshallResolver) VisitObjectNode(node *ast.JsonObjectNode) error {
	if resolver.hasUnmarshaller {
		return resolver.resolveByCustomizeObjectUnmarshal(node)
	}
	// fill the values in the reflection.Value
	for i := len(node.Value) - 1; i >= 0; i-- {
		kvNode := node.Value[i]
		if err := kvNode.Visit(resolver); err != nil {
			return err
		}
	}
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitBooleanNode(node *ast.JsonBooleanNode) error {
	if resolver.hasUnmarshaller {
		if node.Value {
			return resolver.resolveByCustomizePrimitiveUnmarshal(token.TrueBytes)
		} else {
			return resolver.resolveByCustomizePrimitiveUnmarshal(token.FalseBytes)
		}
	}
	resolver.setValue(node.Value)
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitNullNode(node *ast.JsonNullNode) error {
	if resolver.hasUnmarshaller {
		// fast unmarshal instead of using interpreter for primitive values
		return resolver.resolveByCustomizePrimitiveUnmarshal(token.NullBytes)
	}
	resolver.setValue(node.Value)
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitNumberNode(node *ast.JsonNumberNode) error {
	if resolver.hasUnmarshaller {
		// fast unmarshal instead of using interpreter for primitive values
		numStr := strconv.FormatFloat(node.Value, 'f', -1, 64)
		return resolver.resolveByCustomizePrimitiveUnmarshal([]byte(numStr))
	}
	realValue := resolver.convertNumberBaseOnKind(node.Value)
	resolver.setValue(realValue)
	return resolver.resolve()
}

// this is design to call the customize unmarshaler and `resolve` the resolver
func (resolver *unmarshallResolver) resolveByCustomizePrimitiveUnmarshal(payload []byte) error {
	// fast unmarshal instead of using interpreter for primitive values
	unmarshalMethod := resolver.ptrToActualValue.MethodByName("UnmarshalJSON")
	if unmarshalMethod.IsValid() {
		result := unmarshalMethod.Call([]reflect.Value{reflect.ValueOf(payload)})
		if unmarshalErr, ok := result[0].Interface().(error); ok {
			if unmarshalErr != nil {
				return unmarshalErr
			}
		}
		return resolver.resolve()
	}
	return nil
}

func (resolver *unmarshallResolver) VisitStringNode(node *ast.JsonStringNode) error {
	if resolver.hasUnmarshaller {
		return resolver.resolveByCustomizePrimitiveUnmarshal(node.Value)
	}
	valueToUnmarshal := util.RepairUTF8(string(node.GetValue()))
	if resolver.tagOption != "string" {

		resolver.setValue(valueToUnmarshal)
	} else {
		// this will only happen at AST string node when the tag is `string`
		switch resolver.outElementKind {
		case reflect.Bool:
			if valueToUnmarshal == "true" {
				resolver.setValue(true)
			} else {
				resolver.setValue(false)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(valueToUnmarshal, 10, 64)
			if err != nil {
				return err
			}
			resolver.setValue(intValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(valueToUnmarshal, 64)
			if err != nil {
				return err
			}
			resolver.setValue(floatValue)
		case reflect.String:
			var result string
			err := resolver.options.unmarshaler(node.Value, &result)
			if err != nil {
				return err
			}
			decodedString, err := strconv.Unquote(result)
			if err != nil {
				return err
			}
			resolver.setValue(decodedString)
		default:
			return ErrorUnsupportedDataKind
		}
	}
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitStringWithVariableNode(node *ast.JsonExtendedStringWIthVariableNode) error {
	if resolver.hasUnmarshaller {
		valueToUnmarshal := util.RepairUTF8(string(node.GetValue()))
		return resolver.resolveByCustomizePrimitiveUnmarshal([]byte(valueToUnmarshal))
	}
	result, err := resolveStringVariable(node, resolver.options)
	if err != nil {
		return err
	}
	valueToSet := util.RepairUTF8(string(result))
	resolver.setValue(valueToSet)
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitVariableNode(node *ast.JsonExtendedVariableNode) error {
	if resolver.hasUnmarshaller {
		return resolver.resolveByCustomizePrimitiveUnmarshal(node.Value)
	}
	result, err := resolveVariable(node, resolver.options)
	if err != nil {
		return err
	}
	if result != nil && reflect.TypeOf(result).Kind() == reflect.String {
		valueToSet := util.RepairUTF8(result.(string))
		resolver.setValue(valueToSet)
	} else {
		resolver.setValue(result)
	}
	return resolver.resolve()
}
