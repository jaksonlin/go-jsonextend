package token

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

// symbol token of json, these are the format token that need to use to construct the AST/ syntax checker
// double quotaion though is also symbol, it is value symbol, not json protocol symbol to hold the format
func IsSymbolToken(t TokenType) bool {
	return t == TOKEN_COMMA || t == TOKEN_COLON || t == TOKEN_LEFT_BRACE || t == TOKEN_LEFT_BRACKET || t == TOKEN_RIGHT_BRACE || t == TOKEN_RIGHT_BRACKET
}
