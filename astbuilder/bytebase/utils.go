package bytebase

import (
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

// in json k-v start bytes
func GetTokenTypeByStartCharacter(b byte) token.TokenType {
	switch {
	case util.IsNumberStartingCharacter(b):
		return token.TOKEN_NUMBER
	case b == '"':
		return token.TOKEN_STRING
	case b == 't' || b == 'f':
		return token.TOKEN_BOOLEAN
	case b == 'n':
		return token.TOKEN_NULL
	case b == '$':
		return token.TOKEN_VARIABLE
	case b == '{':
		return token.TOKEN_LEFT_BRACE
	case b == '[':
		return token.TOKEN_LEFT_BRACKET
	case b == '}':
		return token.TOKEN_RIGHT_BRACE
	case b == ']':
		return token.TOKEN_RIGHT_BRACKET
	case b == ':':
		return token.TOKEN_COLON
	case b == ',':
		return token.TOKEN_COMMA
	case util.IsSpaces(b):
		return token.TOKEN_SPACE
	default:
		return token.TOKEN_DROP
	}
}

// these symbols should be unread to buffer, they are read first to determine the state change,
// not using peek to collect them because there may be a long way to go till we see it.
func ShouldUnreadByte(t token.TokenType) bool {
	switch t {
	case token.TOKEN_BOOLEAN:
	case token.TOKEN_NUMBER:
	case token.TOKEN_STRING:
	case token.TOKEN_NULL:
	case token.TOKEN_VARIABLE:
	default: // unread value for value mode
		return false
	}
	return true
}

func isJSONNumberByte(b byte) bool {
	return (b >= '0' && b <= '9') || b == '-' || b == '.' || b == 'e' || b == 'E' || b == '+'
}
