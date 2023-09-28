package interpreter

import (
	"reflect"

	"github.com/jaksonlin/go-jsonextend/ast"
)

// we hold a pointer to these field rather than keeping the field's reflect.Value directly
// number of pointer of the original field, not including the one we create to hold the filed
type unmarshallResolver struct {
	options              *unmarshallOptions
	astNode              ast.JsonNode
	collectionDataType   reflect.Type
	outElementKind       reflect.Kind
	arrayIndex           int
	awaitingResolveCount int
	numberOfPointer      int
	awaitingResolve      bool
	isPointerValue       bool
	objectKey            string
	parent               *unmarshallResolver
	ptrToActualValue     reflect.Value // single ptr to no matter what actual value is (for *****int, keeps only *int to the actual value)
}

func (resolver *unmarshallResolver) resolveDependency(dependentResolver *unmarshallResolver) error {
	dependentValue := dependentResolver.restoreValue()
	if resolver.outElementKind == reflect.Array || resolver.outElementKind == reflect.Slice {
		resolver.ptrToActualValue.Elem().Index(dependentResolver.arrayIndex).Set(dependentValue)
	} else if resolver.outElementKind == reflect.Struct {
		field := resolver.ptrToActualValue.Elem().FieldByName(dependentResolver.objectKey)
		if !field.IsValid() || !field.CanSet() {
			return NewErrFieldCannotSetOrNotfound(dependentResolver.objectKey)
		}
		field.Set(dependentValue)
	} else if resolver.outElementKind == reflect.Map {
		key, err := resolver.createMapKeyValueByMapKeyKind(dependentResolver.objectKey)
		if err != nil {
			return err
		}
		resolver.ptrToActualValue.Elem().SetMapIndex(key, dependentValue)
	} else {
		return ErrorPrimitiveTypeCannotResolveDependency
	}
	resolver.awaitingResolveCount -= 1
	return nil
}

func (resolver *unmarshallResolver) setValue(value interface{}) {
	if value == nil {
		nilValue := reflect.Zero(resolver.ptrToActualValue.Elem().Type())
		resolver.ptrToActualValue.Elem().Set(nilValue)
	} else {
		resolver.ptrToActualValue.Elem().Set(reflect.ValueOf(value))
	}
}

// return the actual reflect.Value in the resolver, the resolver is desinged to hold a pointer to anything it keeps
// when the actual field is pointer type, you need a pointer to the actual field to set its Elem to the retrun value from this func
func (resolver *unmarshallResolver) restoreValue() reflect.Value {
	if !resolver.isPointerValue {
		return resolver.ptrToActualValue.Elem() // remove the pointer we add (newunmarshallResolver)
	} else {
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
func newunmarshallResolver(
	node ast.JsonNode,
	outType reflect.Type,
	options *unmarshallOptions) *unmarshallResolver {

	someOutType := outType
	numberOfPointer := 0
	var elementKind reflect.Kind
	var collectionDataType reflect.Type = nil
	isPointerValue := someOutType.Kind() == reflect.Pointer
	for someOutType.Kind() == reflect.Pointer {
		someOutType = someOutType.Elem()
		numberOfPointer += 1
	}
	var ptrToActualValue reflect.Value
	// use a pointer to hold no matter what it is inside
	switch someOutType.Kind() {
	case reflect.Slice:
		numberOfElement := node.(*ast.JsonArrayNode).Length()
		sliceType := reflect.SliceOf(someOutType.Elem())
		sliceValue := reflect.MakeSlice(sliceType, numberOfElement, numberOfElement) // use index to manipulate the slice
		ptrToActualValue = reflect.New(sliceValue.Type())
		ptrToActualValue.Elem().Set(sliceValue)
		elementKind = reflect.Slice
		collectionDataType = sliceValue.Type().Elem()
	case reflect.Array:
		numberOfElement := node.(*ast.JsonArrayNode).Length()
		arrayType := reflect.ArrayOf(numberOfElement, someOutType.Elem())
		ptrToActualValue = reflect.New(arrayType)
		elementKind = reflect.Array
		collectionDataType = ptrToActualValue.Type().Elem()
	case reflect.Map:
		newMap := reflect.MakeMap(someOutType)
		ptrToActualValue = reflect.New(newMap.Type())
		ptrToActualValue.Elem().Set(newMap)
		elementKind = reflect.Map
		collectionDataType = newMap.Type().Elem()
	case reflect.Struct:
		ptrToActualValue = reflect.New(someOutType) //*Struct
		elementKind = reflect.Struct
	default: // primitives
		ptrToActualValue = reflect.New(someOutType)
		ptrToActualValue.Elem().Set(reflect.Zero(someOutType))
		elementKind = someOutType.Kind()
	}

	base := &unmarshallResolver{
		options:              options,
		astNode:              node,
		ptrToActualValue:     ptrToActualValue,
		awaitingResolveCount: 0,
		awaitingResolve:      false,
		parent:               nil,
		arrayIndex:           -1,
		numberOfPointer:      numberOfPointer,
		isPointerValue:       isPointerValue,
		outElementKind:       elementKind,
		collectionDataType:   collectionDataType,
	}
	return base
}

var _ ast.JsonVisitor = &unmarshallResolver{}

func (resolver *unmarshallResolver) VisitArrayNode(node *ast.JsonArrayNode) error {
	// fill the values in the reflection.Value

	for i := len(node.Value) - 1; i >= 0; i-- {
		resolver, err := resolver.createArrayElementResolver(i, node.Value[i])
		if err != nil {
			return err
		}
		resolver.options.resolverStack.Push(resolver)
	}
	return resolver.resolve()

}

func (resolver *unmarshallResolver) VisitKeyValuePairNode(node *ast.JsonKeyValuePairNode) error {

	key, err := resolver.processKVKeyNode(node.Key)
	if err != nil {
		return err
	}

	newResolver, err := resolver.processKVValueNode(key, node.Value)
	if err != nil {
		return err
	}

	resolver.options.resolverStack.Push(newResolver)

	return nil
}

func (resolver *unmarshallResolver) VisitObjectNode(node *ast.JsonObjectNode) error {
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
	resolver.setValue(node.Value)
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitNullNode(node *ast.JsonNullNode) error {
	resolver.setValue(node.Value)
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitNumberNode(node *ast.JsonNumberNode) error {
	realValue := resolver.convertNumberBaseOnKind(node.Value)
	resolver.setValue(realValue)
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitStringNode(node *ast.JsonStringNode) error {
	resolver.setValue(node.GetValue())
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitStringWithVariableNode(node *ast.JsonExtendedStringWIthVariableNode) error {
	result, err := resolveStringVariable(node, resolver.options)
	if err != nil {
		return err
	}
	resolver.setValue(string(result))
	return resolver.resolve()
}

func (resolver *unmarshallResolver) VisitVariableNode(node *ast.JsonExtendedVariableNode) error {
	result, err := resolveVariable(node, resolver.options)
	if err != nil {
		return err
	}
	resolver.setValue(result)
	return resolver.resolve()
}
