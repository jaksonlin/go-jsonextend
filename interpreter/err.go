package interpreter

import "errors"

var (
	ErrorInterpreSymbolFailure       = errors.New("have symbol not consumed")
	ErrorInterpretVariable           = errors.New("error when interpret variable values")
	ErrorInternalInterpreterOutdated = errors.New("interpreter not update as the ast grows")
)

var (
	ErrorUnmarshalNotSlice = errors.New("unmarshal to a non-slice variable")
)
