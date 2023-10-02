package ast

import (
	"errors"
)

var (
	ErrorASTComplete                          = errors.New("cannot create new ast node when the ast has finished")
	ErrorASTStackInvalidElement               = errors.New("ast trace stack will only contain: kv, array, object, no primitive values")
	ErrorASTUnexpectedElement                 = errors.New("unexecpted stack element type")
	ErrorASTUnexpectedOwnerElement            = errors.New("unexecpted stack owner element type")
	ErrorASTStackEmpty                        = errors.New("unexecpted stack empty")
	ErrorASTEncloseElementType                = errors.New("enclose element type must be array or object")
	ErrorASTIncorrectNodeType                 = errors.New("incorrect node type")
	ErrorIncorrectSyntaxSymbolForConstructAST = errors.New("incorrect character for construct ast")
	ErrorASTKeyValuePairNotStringAsKey        = errors.New("object key should be string")
)

var (
	ErrorSyntaxEmptyStack                  = errors.New("empty syntax stack")
	ErrorSyntaxEncloseIncorrectSymbol      = errors.New("invalid operation, not ]|} to enclose")
	ErrorSyntaxEncloseSymbolNotMatch       = errors.New("enclose symbol not match")
	ErrorSyntaxEncloseSymbolIncorrect      = errors.New("enclose symbol incorrect")
	ErrorSyntaxCommaBehindLastItem         = errors.New("find `,]` or `,}` in the syntax checker")
	ErrorSyntaxElementNotSeparatedByComma  = errors.New("syntax element not separated by comma")
	ErrorSyntaxUnexpectedSymbolInArray     = errors.New("unexpected symbol in array")
	ErrorSyntaxExtendedSyntaxVariableAsKey = errors.New("extended syntax variable as key")
	ErrorSyntaxObjectSymbolNotMatch        = errors.New("object symbol not match")
)
