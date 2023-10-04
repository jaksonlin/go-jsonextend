package golang

import (
	"errors"
)

var (
	ErrorUnknownData         = errors.New("unknow data")
	ErrorUnsupportedDataKind = errors.New("unsupported data type for conversion")
	ErrorCyclicAccess        = errors.New("cyclic access to the object")
)
