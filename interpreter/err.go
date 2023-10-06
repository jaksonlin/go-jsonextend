package interpreter

import (
	"errors"
	"fmt"
)

const (
	fieldNotFoundConst = "field %s is not found"

	fieldNotFoundOrCannotSetConst      = "field %s is not found or cannot set"
	ExpectingStructInPointerFindOthers = "expecting struct inside pointer but find %s"
	ExpectingArrayInPointerFindOthers  = "expecting array inside pointer but find %s"
	FieldTypeNotMatchAST               = "field type %s is not of ast collection node type"
	FieldNotFound                      = "field with tag %s not found"
)

var (
	ErrorInterpreSymbolFailure           = errors.New("have symbol not consumed")
	ErrorInterpretVariable               = errors.New("error when interpret variable values")
	ErrorInternalInterpreterOutdated     = errors.New("interpreter not update as the ast grows")
	ErrorKeyStringVariableNotResolve     = errors.New("find string variable as object key, but the variable value is missing")
	ErrorVariableValueNotJsonValueType   = errors.New("variable value should be json value type")
	ErrorVariableValueNotFound           = errors.New("variable value is not found")
	ErrorInternalGetPrimitiveValue       = errors.New("try to get primitive value from none-primitive type")
	ErrorInternalShouldBeArrayOrSlice    = errors.New("should be array or slice")
	ErrorInternalMapKeyTypeNotMatchValue = fmt.Errorf("expect bool as key but value is not bool")

	ErrorInternalExpectingPrimitive     = errors.New("expecting primitive values but find others")
	ErrorInternalPtrToArrayFindNotArray = errors.New("expecting pointer to array, but underlying object not array")
	ErrSliceOrArrayNotInit              = errors.New("slice/array not init")
	ErrorUnmarshalStackNoKV             = errors.New("there should not be kv pair in stack")
	ErrorNotSupportedASTNode            = errors.New("not supported ast node")
	ErrorInternalSymbolStackIncorrect   = errors.New("interpret symbol stack invalid")
)

var (
	ErrorUnmarshalNotSlice = errors.New("unmarshal to a non-slice variable")
)

var (
	ErrorReflectNotObject     = errors.New("out element is not map nor struct")
	ErrorReflectInvalidMapKey = errors.New("map key type should be string")
	ErrOutNotNilPointer       = errors.New("out is nil pointer")
)

func NewErrorFiledNotFound(fieldName string) error {
	return fmt.Errorf(fieldNotFoundConst, fieldName)
}

func NewErrFieldCannotSetOrNotfound(fieldName string) error {
	return fmt.Errorf(fieldNotFoundOrCannotSetConst, fieldName)
}

func NewErrorInternalExpectingStructInsidePointerButFindOthers(kind string) error {
	return fmt.Errorf(ExpectingStructInPointerFindOthers, kind)
}
func NewErrorInternalExpectingArrayInsidePointerButFindOthers(kind string) error {
	return fmt.Errorf(ExpectingArrayInPointerFindOthers, kind)
}

func NewErrorInternalFieldTypeNotMatchAST(kind string) error {
	return fmt.Errorf(FieldTypeNotMatchAST, kind)
}

func NewErrorFieldNotFound(field string) error {
	return fmt.Errorf(FieldNotFound, field)
}
