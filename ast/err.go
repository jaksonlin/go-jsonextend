package ast

import (
	"errors"
)

var (
	ErrorASTComplete                   = errors.New("cannot create new ast node when the ast has finished")
	ErrorASTStackInvalidElement        = errors.New("ast trace stack will only contain: kv, array, object, no primitive values")
	ErrorASTUnexpectedElement          = errors.New("unexecpted stack element type")
	ErrorASTUnexpectedOwnerElement     = errors.New("unexecpted stack owner element type")
	ErrorASTStackEmpty                 = errors.New("unexecpted stack empty")
	ErrorASTEncloseElementType         = errors.New("enclose element type must be array or object")
	ErrorASTIncorrectNodeType          = errors.New("incorrect node type")
	ErrorASTKeyValuePairNotStringAsKey = errors.New("object key should be string")
)
