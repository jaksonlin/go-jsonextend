package tokenizer

import (
	"errors"
)

var (
	ErrorIncorrectCharacter                = errors.New("incorrect character")
	ErrorIncorrectValueForState            = errors.New("extracted value not match state")
	ErrorIncorrectValueTypeForConstructAST = errors.New("incorrect value type for construct ast")
	ErrorInternalASTProcotolChanged        = errors.New("detect unexpect ast stack change, not kv, array, object, NodeType at top of stack")
	ErrorUnexpectedEOF                     = errors.New("unexpected EOF")
	ErrorTokenRouteNotConfigure            = errors.New("token route not configure")
	ErrorExtendedVariableFormatIncorrect   = errors.New("variable should be of ${variableName} format")
)
