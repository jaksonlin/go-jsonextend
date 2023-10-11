package tokenizer

import (
	"errors"
	"fmt"

	"github.com/jaksonlin/go-jsonextend/token"
)

const (
	incorrectToken = "unexpected token found in mode: %d, token type: %d"
)

var (
	ErrorIncorrectValueTypeForConstructAST = errors.New("incorrect value type for construct ast")
	ErrorInternalASTProcotolChanged        = errors.New("detect unexpect ast stack change, not kv, array, object, NodeType at top of stack")
	ErrorUnexpectedEOF                     = errors.New("unexpected EOF")
	ErrorTokenRouteNotConfigure            = errors.New("token route not configure")
	ErrorExtendedVariableFormatIncorrect   = errors.New("variable should be of ${variableName} format")
)

func NewErrorIncorrectToken(mode StateMode, token token.TokenType) error {
	return fmt.Errorf(incorrectToken, mode, token)
}
