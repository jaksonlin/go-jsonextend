package ast

import (
	"github.com/jaksonlin/go-jsonextend/util"
)

type syntaxChecker struct {
	syntaxState *util.Stack[byte]
	length      int
}

func newSyntaxChecker() *syntaxChecker {
	return &syntaxChecker{
		syntaxState: &util.Stack[byte]{},
		length:      0,
	}
}

func (s *syntaxChecker) PushSymbol(b byte) {
	s.syntaxState.Push(b)
	s.length += 1
}

func (s *syntaxChecker) PushValue(val AST_NODETYPE) {
	s.syntaxState.Push(byte(val))
	s.length += 1
}

func (s *syntaxChecker) Length() int {
	return s.length
}

func (s *syntaxChecker) Enclose(b byte) error {

	t, err := s.syntaxState.Pop()
	if err == util.ErrorEndOfStack {
		return ErrorSyntaxEmptyStack
	}
	if t > AST_NODE_TYPE_BOUNDARY {
		return ErrorSyntaxEncloseIncorrectSymbol
	}
	if t != b {
		return ErrorSyntaxEncloseSymbolNotMatch
	}
	if t == ']' {
		return s.jsonArrayFormatCheck()
	} else if t == '}' {
		return s.jsonObjectCheck()
	} else {
		return ErrorSyntaxEncloseSymbolIncorrect
	}
}

func (s *syntaxChecker) jsonArrayFormatCheck() error {
	expectingValue := true
	lastIsValue := false
	hasEncounterValue := false
	for {
		t, err := s.syntaxState.Pop()
		if err == util.ErrorEndOfStack {
			return ErrorSyntaxEmptyStack
		}
		if t == '[' {
			if hasEncounterValue && !lastIsValue { // deal with [] | [,], the previous is ok hasNeverEncounterValue by pass to ok, later raise error
				return ErrorSyntaxCommaBehindLastItem
			}
			// mark that here is an array in the syntax checker
			s.syntaxState.Push(AST_ARRAY)
			return nil
		}
		if expectingValue { // ascii symbol
			if t < AST_NODE_TYPE_BOUNDARY {
				return ErrorSyntaxElementNotSeparatedByComma
			} else {
				lastIsValue = true
				hasEncounterValue = true
			}
		} else if !expectingValue {
			if t > AST_NODE_TYPE_BOUNDARY {
				return ErrorSyntaxElementNotSeparatedByComma
			}
			if t != 0x2C {
				return ErrorSyntaxUnexpectedSymbolInArray
			}
			lastIsValue = false
		}
		expectingValue = !expectingValue

	}
}

func (s *syntaxChecker) jsonObjectCheck() error {
	expectingValue := true // already pop the } | ]
	lastIsValue := false
	hasEncounterValue := false
	expectingSymbol := byte(':') // first symbol to expect is : then , then : then , ...
	for {
		t, err := s.syntaxState.Pop()
		if err == util.ErrorEndOfStack {
			return ErrorSyntaxEmptyStack
		}
		// check first otherwise drop into compare with allowed symbol
		if t == '{' {
			if hasEncounterValue && !lastIsValue {
				return ErrorSyntaxCommaBehindLastItem
			}
			// enclose the object as a value in the syntax checker, this will save our hands in handling }} or ]} in the syntax checker
			// this will collapse the checking of symbol into: always having symbol in between the value (braces and brakcets are collpased into value)
			// in our design, the array and object will be collapse into syntax_value, []{}
			s.syntaxState.Push(AST_OBJECT)
			return nil
		}
		if expectingValue {
			// expecting value but find symbol
			if t < AST_NODE_TYPE_BOUNDARY {
				return ErrorSyntaxElementNotSeparatedByComma
			} else {
				// the json key's previous symbol is either `{` or `,`, that means in the stack's next pop, if it is a ',' then this AST_STRING_vARIABLE is a json-key
				// which is invalid, because people may put arbitrary variable as key which may break the json format
				if t == AST_VARIABLE && expectingSymbol == ',' {
					return ErrorSyntaxExtendedSyntaxVariableAsKey
				}
				lastIsValue = true
				hasEncounterValue = true
			}
		} else if !expectingValue {
			if t > AST_NODE_TYPE_BOUNDARY {
				return ErrorSyntaxElementNotSeparatedByComma
			}
			if t != expectingSymbol {
				return ErrorSyntaxObjectSymbolNotMatch
			}
			// switch the symbol to expect, if expectingSymbol is : then next is , and vice versa
			if expectingSymbol == ':' {
				expectingSymbol = ','
			} else {
				expectingSymbol = ':'
			}
			lastIsValue = false
		}
		expectingValue = !expectingValue

	}
}
