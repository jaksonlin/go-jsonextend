package unmarshaler

import (
	"reflect"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/util"
)

func (resolver *unmarshallResolver) processKVKeyNode(node ast.JsonStringValueNode) (string, error) {
	var key string = node.GetValue()
	if node.GetNodeType() == ast.AST_STRING_VARIABLE {
		resultBytes, err := resolveStringVariable(node.(*ast.JsonExtendedStringWIthVariableNode), resolver.options)
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

func (resolver *unmarshallResolver) getFieldByTag(objKey string) (*util.JSONStructField, error) {
	resolver.collectAllFields()
	fieldInfo, ok := resolver.fields[objKey]
	if !ok {
		return nil, NewErrorFieldNotValid(objKey)
	}
	return fieldInfo, nil
}

// create resolver to resolving the things in kv's value
func (resolver *unmarshallResolver) processKVValueNode(key string, valueNode ast.JsonNode) (*unmarshallResolver, error) {
	// check if the root is a struct or map to hold our kv pair
	var kvParentElementType reflect.Type = resolver.ptrToActualValue.Elem().Type()

	var kvValueElementType reflect.Type = nil
	var tagOption string = ""
	if kvParentElementType.Kind() == reflect.Map {
		// when parent is a map, the child element type is the map's value type
		kvValueElementType = kvParentElementType.Elem()

	} else if kvParentElementType.Kind() == reflect.Struct {
		// when parent is a struct, the child element type is the struct's field type
		fieldInfo, err := resolver.getFieldByTag(key) // struct field
		if err != nil {
			return nil, err
		}
		kvValueElementType = fieldInfo.FieldValue.Type()
		tagOption = fieldInfo.Options

	} else {
		return nil, NewErrorInternalExpectingStructButFindOthers(kvParentElementType.Kind().String())
	}

	// 2. create the collection's reflection value representative
	newResolver, err := newUnmarshallResolver(valueNode, kvValueElementType, resolver.options, tagOption)
	if err != nil {
		return nil, err
	}

	// 3. create relation
	newResolver.bindObjectParent(key, resolver)

	return newResolver, nil
}

func (resolver *unmarshallResolver) createArrayElementResolver(index int, node ast.JsonNode) (*unmarshallResolver, error) {
	// root: slice, array, *array

	// 1. get the keys coresponding value type
	var childElementType reflect.Type = resolver.ptrToActualValue.Elem().Type()

	if resolver.ptrToActualValue.Elem().Kind() == reflect.Slice || resolver.ptrToActualValue.Elem().Kind() == reflect.Array {

		childElementType = resolver.ptrToActualValue.Elem().Type().Elem()

	}

	// 2. create the collection's reflection value representative
	newResolver, err := newUnmarshallResolver(node, childElementType, resolver.options, "")
	if err != nil {
		return nil, err
	}

	// 3. create relation
	newResolver.bindArrayLikeParent(index, resolver)

	return newResolver, nil
}

func (resolver *unmarshallResolver) resolve() error {
	// no parent, no need to enclose
	if resolver.parent == nil {
		return nil
	}

	// have unresolve child item, cannot enclose now
	if resolver.awaitingResolveCount > 0 {
		return nil
	}
	return resolver.parent.resolveDependency(resolver)
}

func (resolver *unmarshallResolver) process() error {
	node := resolver.astNode
	return node.Visit(resolver)
}

func UnmarshallAST(node ast.JsonNode, variables map[string]interface{}, marshaler ast.MarshalerFunc, unmarshaler ast.UnmarshalerFunc, out interface{}) error {
	// deep first traverse the AST
	valueItem := reflect.ValueOf(out)
	if valueItem.Kind() != reflect.Pointer || valueItem.IsNil() {
		return ErrOutNotPointer
	}

	options := NewUnMarshallOptions(variables, marshaler, unmarshaler)
	traverseStack := options.resolverStack
	resolver, err := newUnmarshallResolver(node, valueItem.Type(), options, "")
	if err != nil {
		return err
	}
	traverseStack.Push(resolver)

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
			if err := resolver.resolve(); err != nil {
				return err
			}
			traverseStack.Pop()
		}

	}
	actualValue := resolver.restoreValue().Elem()
	valueItem.Elem().Set(actualValue.Convert(valueItem.Elem().Type()))
	return nil
}
