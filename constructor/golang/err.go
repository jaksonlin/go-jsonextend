package golang

import (
	"errors"
)

var (
	ErrorUnknownData         = errors.New("unknow data")
	ErrorUnsupportedDataKind = errors.New("unsupported data type for conversion")
	ErrorCyclicAccess        = errors.New("cyclic access to the object")
	ErrorInvalidMapKey       = errors.New("Map key can only be string or int")
)
