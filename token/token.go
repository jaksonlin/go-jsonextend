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

// stop evil pointer and interface
func removePointersAndInterfaces(v reflect.Value, hasInterface *bool) (reflect.Value, bool) {
	rv := v
	if !rv.IsValid() {
		return rv, true
	}
	// incase interface hoding pointer
	if rv.Kind() == reflect.Interface {
		if !(*hasInterface) {
			*hasInterface = true
		}
		if rv.IsNil() {
			return rv, true
		}

		rv = rv.Elem()

	}
	// removes all the pointers
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return rv, true
		}
		rv = rv.Elem()
	}
	// incase it is interface at last
	if rv.Kind() == reflect.Interface {
		if !(*hasInterface) {
			*hasInterface = true
		}
		if rv.IsNil() {
			return rv, true
		} else {
			elem := rv.Elem()
			return removePointersAndInterfaces(elem, hasInterface)
		}
	}
	return rv, false
}

// about interface{}, when doing marshaling, we won't care how many pointers or interfaces a value is wrapped,
// all we need to do is to get the actual value, and then we can return the token type.
// when unmarshaling, there's 2 situation about interface{},
// 1. when a field is interface{}, we just put what ever value into it;
// 2. a field is *interface{}, then we just need to put the value into the interface, and restore the level of pointers.(which we have already do)
// otherthings is about whether a field is captured by interface{} or not, if it does, the json string tag option won't work.
// this is typically not easy to detect when the interface is wrapped deep by pointer, therefore we return this indicator from this function.
// as for howmany level of a pointer is wrapped
// we will not return the value from here, if we does, we will bypass the cyclic access check. (e.g. when a pointer points to itself, if we do return the value from here
// we will not be able to detect the cyclic access)

// get the acutal token type underneath the reflect.Value, and return the indicator whether the value is wrapped by interface{}
func GetTokenTypeByReflection(v reflect.Value) (TokenType, bool) {

	hasInterface := false
	val, isNil := removePointersAndInterfaces(v, &hasInterface)

	if isNil {
		return TOKEN_NULL, hasInterface
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return TOKEN_NUMBER, hasInterface
	case reflect.String:
		return TOKEN_STRING, hasInterface
	case reflect.Bool:
		return TOKEN_BOOLEAN, hasInterface
	case reflect.Slice, reflect.Array:
		return TOKEN_LEFT_BRACKET, hasInterface
	case reflect.Struct, reflect.Map:
		return TOKEN_LEFT_BRACE, hasInterface
	default:
		return TOKEN_UNKNOWN, hasInterface
	}
}

// func GetTokenTypeByReflection(v reflect.Value) TokenType {
// 	rv := v
// 	if !rv.IsValid() {
// 		return TOKEN_NULL
// 	}
// 	for rv.Kind() == reflect.Pointer {
// 		if rv.IsNil() {
// 			return TOKEN_NULL
// 		}
// 		rv = rv.Elem()
// 	}
// 	switch rv.Kind() {
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
// 		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
// 		reflect.Float32, reflect.Float64:
// 		return TOKEN_NUMBER
// 	case reflect.String:
// 		return TOKEN_STRING
// 	case reflect.Bool:
// 		return TOKEN_BOOLEAN
// 	case reflect.Slice, reflect.Array:
// 		return TOKEN_LEFT_BRACKET
// 	case reflect.Struct, reflect.Map:
// 		return TOKEN_LEFT_BRACE
// 	case reflect.Interface:
// 		if rv.IsNil() {
// 			return TOKEN_NULL
// 		} else {
// 			elem := rv.Elem()
// 			return GetTokenTypeByReflection(elem)
// 		}
// 	default:
// 		return TOKEN_UNKNOWN
// 	}
// }

var (
	NullBytes  []byte = []byte("null")
	FalseBytes []byte = []byte("false")
	TrueBytes  []byte = []byte("true")
)
