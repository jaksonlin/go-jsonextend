package interpreter

import (
	"errors"
	"fmt"
)

const (
	fieldNotFoundConst = "field %s is not found"

	fieldNotFoundOrCannotSetConst      = "field %s is not found or cannot set"
	ExpectingStructFindOthers          = "expecting struct but find %s"
	ExpectingStructInPointerFindOthers = "expecting struct inside pointer but find %s"
	ExpectingArrayInPointerFindOthers  = "expecting array inside pointer but find %s"
	VariableNotFound                   = "variable value for %s not found"
	FieldTypeNotMatchAST               = "field type %s is not of ast collection node type"
)

var (
	ErrorInterpreSymbolFailure         = errors.New("have symbol not consumed")
	ErrorInterpretVariable             = errors.New("error when interpret variable values")
	ErrorInternalInterpreterOutdated   = errors.New("interpreter not update as the ast grows")
	ErrorKeyStringVariableNotResolve   = errors.New("find string variable as object key, but the variable value is missing")
	ErrorVariableValueNotJsonValueType = errors.New("variable value should be json value type")
	ErrorVariableValueNotFound         = errors.New("variable value is not found")
	ErrorInternalGetPrimitiveValue     = errors.New("try to get primitive value from none-primitive type")
	ErrorInternalShouldBeArrayOrSlice  = errors.New("should be array or slice")

	ErrorStringVariableNotResolveOnKeyLocation = errors.New("object key contain string variable that has no variable value")
	ErrorInternalExpectingPrimitive            = errors.New("expecting primitive values but find others")
	ErrorInternalPtrToArrayFindNotArray        = errors.New("expecting pointer to array, but underlying object not array")
	ErrorInternalExpectingArrayLikeObject      = errors.New("expecting array like object but find others")
	ErrSliceNotInit                            = errors.New("slice not init")
	ErrArrayNotInit                            = errors.New("array not init")
	ErrorUnmarshalStackNoKV                    = errors.New("there should not be kv pair in stack")
)

var (
	ErrorUnmarshalNotSlice = errors.New("unmarshal to a non-slice variable")
)

var (
	ErrorReflectNotObject     = errors.New("out element is not map nor struct")
	ErrorReflectInvalidMapKey = errors.New("map key type should be string")
	ErrOutNotPointer          = errors.New("out is not pointer")
	ErrOutNotNilPointer       = errors.New("out is nil pointer")
)

func NewErrorFiledNotFound(fieldName string) error {
	return fmt.Errorf(fieldNotFoundConst, fieldName)
}

func NewErrFieldCannotSetOrNotfound(fieldName string) error {
	return fmt.Errorf(fieldNotFoundOrCannotSetConst, fieldName)
}

func NewErrorInternalExpectingStructButFindOthers(kind string) error {
	return fmt.Errorf(ExpectingStructFindOthers, kind)
}
func NewErrorInternalExpectingStructInsidePointerButFindOthers(kind string) error {
	return fmt.Errorf(ExpectingStructInPointerFindOthers, kind)
}
func NewErrorInternalExpectingArrayInsidePointerButFindOthers(kind string) error {
	return fmt.Errorf(ExpectingArrayInPointerFindOthers, kind)
}
func NewVariableNotFound(variable string) error {
	return fmt.Errorf(VariableNotFound, variable)
}
func NewErrorInternalFieldTypeNotMatchAST(kind string) error {
	return fmt.Errorf(FieldTypeNotMatchAST, kind)
}
