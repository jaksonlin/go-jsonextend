package golang

import (
	"errors"
)

var (
	ErrorUnknownData                          = errors.New("unknow data")
	ErrorUnsupportedDataKind                  = errors.New("unsupported data type for conversion")
	ErrorCyclicAccess                         = errors.New("cyclic access to the object")
	ErrorInvalidMapKey                        = errors.New("map key can only be string or int")
	ErrorInvalidTypeOnExportedField           = errors.New("invalid exported field type for marshaling")
	ErrNotNumericValueField                   = errors.New("field is not having value of numeric type")
	ErrorInvalidJsonTag                       = errors.New("invalid json tag")
	ErrorStringConfigTypeInvalid              = errors.New("json tag string config only support pritmive data type")
	ErrorIncorrectSyntaxSymbolForConstructAST = errors.New("incorrect character for construct ast")
)
