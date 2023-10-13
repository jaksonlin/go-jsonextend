package util

import "errors"

var (
	ErrorInputNil                         = errors.New("input is nil")
	ErrorInputNotNumber                   = errors.New("input is not a number")
	ErrorVariableDataKind                 = errors.New("unsupported variable data kind")
	ErrorUnsupportedDataKindConvertNumber = errors.New("unsupported data kind for number conversion")
	ErrorEndOfStack                       = errors.New("end of stack")
)
