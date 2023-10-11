package token

import "reflect"

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

	TOKEN_UNKNOWN TokenType = 97
	TOKEN_DUMMY   TokenType = 98
	TOKEN_DROP    TokenType = 99
)

// symbol token of json, these are the format token that need to use to construct the AST/ syntax checker
// double quotaion though is also symbol, it is value symbol, not json protocol symbol to hold the format
func IsSymbolToken(t TokenType) bool {
	return t == TOKEN_COMMA || t == TOKEN_COLON || t == TOKEN_LEFT_BRACE || t == TOKEN_LEFT_BRACKET || t == TOKEN_RIGHT_BRACE || t == TOKEN_RIGHT_BRACKET
}

type ValueState struct {
	HasPointer   bool
	HasInterface bool
	Field        reflect.Value
	TokenType    TokenType
}

// stop evil pointer and interface
func removePointersAndInterfaces(state *ValueState) {
	rv := state.Field
	if !rv.IsValid() {
		state.TokenType = TOKEN_NULL
		return
	}
	// incase interface hoding pointer
	if rv.Kind() == reflect.Interface {
		state.HasInterface = true
		if rv.IsNil() {
			state.TokenType = TOKEN_NULL
			return
		}
		rv = rv.Elem()
		state.Field = rv
	}
	// removes all the pointers
	for rv.Kind() == reflect.Pointer {
		state.HasPointer = true
		if rv.IsNil() {
			state.TokenType = TOKEN_NULL
			return
		}
		rv = rv.Elem()
		state.Field = rv
	}
	// incase it is interface at last
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			state.TokenType = TOKEN_NULL
			return
		} else {
			elem := rv.Elem()
			state.Field = elem
			removePointersAndInterfaces(state)
		}
	}
}

func GetValueState(v reflect.Value) *ValueState {
	rs := &ValueState{Field: v}
	removePointersAndInterfaces(rs)

	switch rs.Field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		rs.TokenType = TOKEN_NUMBER
	case reflect.String:
		rs.TokenType = TOKEN_STRING
	case reflect.Bool:
		rs.TokenType = TOKEN_BOOLEAN
	case reflect.Slice, reflect.Array:
		rs.TokenType = TOKEN_LEFT_BRACKET
	case reflect.Struct, reflect.Map:
		rs.TokenType = TOKEN_LEFT_BRACE
	default:
		rs.TokenType = TOKEN_UNKNOWN
	}
	return rs
}

var (
	NullBytes  []byte = []byte("null")
	FalseBytes []byte = []byte("false")
	TrueBytes  []byte = []byte("true")
)
