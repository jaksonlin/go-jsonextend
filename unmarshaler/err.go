package unmarshaler

import (
	"errors"
	"fmt"
)

const (
	ExpectingStructFindOthers = "expecting struct but find %s"
	VariableNotFound          = "variable value for %s not found"
	FieldNotValid             = "field not exist %s"
	KVKindNotMatch            = "expect %s as key but value is not :%#v"
)

var (
	ErrOutNotPointer                                   = errors.New("out is not pointer")
	ErrorInternalNoneResolvable                        = errors.New("expecting dependendent element to resolve")
	ErrorInternalUnsupportedMapKeyKind                 = fmt.Errorf("unsupported map key to continue")
	ErrorStringVariableNotResolveOnKeyLocation         = errors.New("object key contain string variable that has no variable value")
	ErrorInternalDependentResolverHasOnResolveLocation = errors.New("dependent value has no idex or object key set to resolve")
	ErrorPrimitiveTypeCannotResolveDependency          = errors.New("pritimive type cannot resolve dependency")
	ErrorInternalExpectingArrayLikeObject              = errors.New("expecting array like object but find others")
	ErrorInvalidUnmarshalResult                        = errors.New("invalid unmarshal result")
	ErrorInvalidTag                                    = errors.New("invalid json tag")
)

type ErrorFieldNotExist struct {
	field string
}

func (e ErrorFieldNotExist) Error() string {
	return fmt.Sprintf(FieldNotValid, e.field)
}

func NewErrorFieldNotValid(field string) ErrorFieldNotExist {
	return ErrorFieldNotExist{field: field}
}
func NewErrorInternalExpectingStructButFindOthers(kind string) error {
	return fmt.Errorf(ExpectingStructFindOthers, kind)
}
func NewVariableNotFound(variable string) error {
	return fmt.Errorf(VariableNotFound, variable)
}
func NewErrorInternalMapKeyValueKindNotMatch(kind string, value interface{}) error {
	return fmt.Errorf(KVKindNotMatch, kind, value)
}
