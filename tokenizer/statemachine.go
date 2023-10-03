package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"

	"io"
)

type stateChangeFunc func() error

type TokenizerStateMachine struct {
	stringState   Tokenizer
	numberState   Tokenizer
	booleanState  Tokenizer
	nullState     Tokenizer
	arrayState    Tokenizer
	objectState   Tokenizer
	initState     Tokenizer
	variableState Tokenizer
	currentState  Tokenizer
	// use a route table to route the default state other than a large switch case
	defaultRoute map[token.TokenType]stateChangeFunc // replace of ToxxxState
	// consturct the AST and check syntax when processing the token in the fly
	astBuilder constructor.ASTBuilder
}

func (i *TokenizerStateMachine) SwitchStateByToken(tokenType token.TokenType) error {
	proxy, ok := i.defaultRoute[tokenType]
	if !ok {
		return ErrorTokenRouteNotConfigure
	}
	proxy()
	return nil
}

func (i *TokenizerStateMachine) RecordSyntaxValue(valueType StateMode, nodeValue interface{}) error {
	// keeps a matching between the state mode and the ast node type, may change in the future
	switch valueType {
	case STRING_MODE:
		return i.astBuilder.RecordSyntaxValue(ast.AST_STRING, nodeValue)
	case BOOLEAN_MODE:
		return i.astBuilder.RecordSyntaxValue(ast.AST_BOOLEAN, nodeValue)
	case NUMBER_MODE:
		return i.astBuilder.RecordSyntaxValue(ast.AST_NUMBER, nodeValue)
	case NULL_MODE:
		return i.astBuilder.RecordSyntaxValue(ast.AST_NULL, nodeValue)
	case VARIABLE_MODE:
		return i.astBuilder.RecordSyntaxValue(ast.AST_VARIABLE, nodeValue)
	case STRING_VARIABLE_MODE:
		return i.astBuilder.RecordSyntaxValue(ast.AST_STRING_VARIABLE, nodeValue)
	default:
		return ErrorIncorrectValueTypeForConstructAST
	}
}

// use AST to switch the state of the machine when primitive values end of their processing
func (i *TokenizerStateMachine) SwitchToLatestState() error {
	if i.astBuilder.HasComplete() {
		// cannot and no need to route, the ast has parsed an json object
		return nil
	}
	n, err := i.astBuilder.TopElementType()
	if err != nil {
		return err
	}
	switch n {
	case ast.AST_ARRAY:
		i.currentState = i.arrayState
		return nil
	case ast.AST_OBJECT:
		fallthrough
	case ast.AST_KVPAIR:
		i.currentState = i.objectState
		return nil
	default:
		return ErrorInternalASTProcotolChanged
	}
}

func (i *TokenizerStateMachine) ProcessData() error {
	for {
		// 1. the ast complete parsing json, end and not read the rest of bytes
		if i.astBuilder.HasComplete() {
			return nil
		}
		err := i.currentState.ProcessData(i.astBuilder)
		// 2. the stream ends, and ast is still expecting content, fail.

		if err != nil {
			if err == io.EOF {
				if !i.astBuilder.HasComplete() {
					return ErrorUnexpectedEOF
				}
				return nil
			} else {
				return err
			}
		}
	}
}

func (i *TokenizerStateMachine) GetCurrentMode() StateMode {
	return i.currentState.GetMode()
}

func (i *TokenizerStateMachine) GetAST() ast.JsonNode {
	return i.astBuilder.GetAST()
}

func (i *TokenizerStateMachine) GetASTBuilder() constructor.ASTManager {
	return i.astBuilder
}
