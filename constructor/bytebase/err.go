package bytebase

import "errors"

var (
	ErrorSyntaxEmptyStack                     = errors.New("empty syntax stack")
	ErrorSyntaxEncloseIncorrectSymbol         = errors.New("invalid operation, not ]|} to enclose")
	ErrorSyntaxEncloseSymbolNotMatch          = errors.New("enclose symbol not match")
	ErrorSyntaxEncloseSymbolIncorrect         = errors.New("enclose symbol incorrect")
	ErrorSyntaxCommaBehindLastItem            = errors.New("find `,]` or `,}` in the syntax checker")
	ErrorSyntaxElementNotSeparatedByComma     = errors.New("syntax element not separated by comma")
	ErrorSyntaxUnexpectedSymbolInArray        = errors.New("unexpected symbol in array")
	ErrorSyntaxExtendedSyntaxVariableAsKey    = errors.New("extended syntax variable as key")
	ErrorSyntaxObjectSymbolNotMatch           = errors.New("object symbol not match")
	ErrorIncorrectSyntaxSymbolForConstructAST = errors.New("incorrect character for construct ast")

	ErrorIncorrectCharacter     = errors.New("incorrect character")
	ErrorIncorrectValueForState = errors.New("extracted value not match state")
)
