package tokenizer

import "github.com/jaksonlin/go-jsonextend/util"

type TokenType uint

const (
	// The token type is a string
	TOKEN_STRING TokenType = iota
	// The token type is a number
	TOKEN_NUMBER
	// The token type is a boolean
	TOKEN_BOOLEAN
	// The token type is a null
	TOKEN_NULL
	// The token type is a left brace
	TOKEN_LEFT_BRACE
	// The token type is a right brace
	TOKEN_RIGHT_BRACE
	// The token type is a left bracket
	TOKEN_LEFT_BRACKET
	// The token type is a right bracket
	TOKEN_RIGHT_BRACKET
	// The token type is a colon
	TOKEN_COLON
	// The token type is a comma
	TOKEN_COMMA
	// customize token type
	TOKEN_VARIABLE
	// variable in string, a string can have multiple variable
	TOKEN_STRING_WITH_VARIABLE
	TOKEN_SPACE
	// deciaml token
	TOKEN_NUMBER_DECIMAL
	//
	TOKEN_DOUBLE_QUOTATION
	TOKEN_DUMMY = 98
	TOKEN_DROP  = 99
)

// in json k-v start bytes
func GetTokenTypeByStartCharacter(b byte) TokenType {
	switch {
	case util.IsNumberStartingCharacter(b):
		return TOKEN_NUMBER
	case b == '"':
		return TOKEN_STRING
	case b == 't' || b == 'f':
		return TOKEN_BOOLEAN
	case b == 'n':
		return TOKEN_NULL
	case b == '$':
		return TOKEN_VARIABLE
	case b == '{':
		return TOKEN_LEFT_BRACE
	case b == '[':
		return TOKEN_LEFT_BRACKET
	case b == '}':
		return TOKEN_RIGHT_BRACE
	case b == ']':
		return TOKEN_RIGHT_BRACKET
	case b == ':':
		return TOKEN_COLON
	case b == ',':
		return TOKEN_COMMA
	case util.IsSpaces(b):
		return TOKEN_SPACE
	default:
		return TOKEN_DROP
	}
}

// symbol token of json, these are the format token that need to use to construct the AST/ syntax checker
// double quotaion though is also symbol, it is value symbol, not json protocol symbol to hold the format
func IsSymbolToken(t TokenType) bool {
	return t == TOKEN_COMMA || t == TOKEN_COLON || t == TOKEN_LEFT_BRACE || t == TOKEN_LEFT_BRACKET || t == TOKEN_RIGHT_BRACE || t == TOKEN_RIGHT_BRACKET
}

// these symbols should be unread to buffer, they are read first to determine the state change,
// not using peek to collect them because there may be a long way to go till we see it.
func ShouldUnreadByte(t TokenType) bool {
	switch t {
	case TOKEN_BOOLEAN:
	case TOKEN_NUMBER:
	case TOKEN_STRING:
	case TOKEN_NULL:
	case TOKEN_VARIABLE:
	default: // unread value for value mode
		return false
	}
	return true
}
