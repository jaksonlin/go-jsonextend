package tokenizer

import (
	"bufio"

	"github.com/jaksonlin/go-jsonextend/ast"

	"io"
)

type JzoneTokenizerStateMachine interface {
	ProcessData(dataSource io.Reader) error
	Reset()
	SwitchToLatestState() error
	SwitchStateByToken(tokenType TokenType) error
	GetCurrentMode() StateMode
	GetASTConstructor() *ast.ASTConstructor
	RecordSyntaxValue(t StateMode, value interface{}) error
	RecordSyntaxSymbol(b byte) error
}

type stateChangeFunc func() error

type TokenizerStateMachine struct {
	stringState   JzoneTokenizer
	numberState   JzoneTokenizer
	booleanState  JzoneTokenizer
	nullState     JzoneTokenizer
	arrayState    JzoneTokenizer
	objectState   JzoneTokenizer
	initState     JzoneTokenizer
	variableState JzoneTokenizer
	currentState  JzoneTokenizer
	// use a route table to route the default state other than a large switch case
	defaultRoute map[TokenType]stateChangeFunc // replace of ToxxxState
	// consturct the AST and check syntax when processing the token in the fly
	astConstructor *ast.ASTConstructor
}

func (i *TokenizerStateMachine) Reset() {
	i.astConstructor = ast.NewASTConstructor()
	i.currentState = i.initState
}

func (i *TokenizerStateMachine) GetASTConstructor() *ast.ASTConstructor {
	return i.astConstructor
}

var _ JzoneTokenizerStateMachine = &TokenizerStateMachine{}

func (i *TokenizerStateMachine) SwitchStateByToken(tokenType TokenType) error {
	proxy, ok := i.defaultRoute[tokenType]
	if !ok {
		return ErrorTokenRouteNotConfigure
	}
	proxy()
	return nil
}

func (i *TokenizerStateMachine) RecordSyntaxSymbol(b byte) error {
	return i.astConstructor.RecordSyntaxSymbol(b)
}

func (i *TokenizerStateMachine) RecordSyntaxValue(valueType StateMode, nodeValue interface{}) error {
	// keeps a matching between the state mode and the ast node type, may change in the future
	switch valueType {
	case STRING_MODE:
		return i.astConstructor.RecordSyntaxValue(ast.AST_STRING, nodeValue)
	case BOOLEAN_MODE:
		return i.astConstructor.RecordSyntaxValue(ast.AST_BOOLEAN, nodeValue)
	case NUMBER_MODE:
		return i.astConstructor.RecordSyntaxValue(ast.AST_NUMBER, nodeValue)
	case NULL_MODE:
		return i.astConstructor.RecordSyntaxValue(ast.AST_NULL, nodeValue)
	case VARIABLE_MODE:
		return i.astConstructor.RecordSyntaxValue(ast.AST_VARIABLE, nodeValue)
	case STRING_VARIABLE_MODE:
		return i.astConstructor.RecordSyntaxValue(ast.AST_STRING_VARIABLE, nodeValue)
	default:
		return ErrorIncorrectValueTypeForConstructAST
	}
}

// use AST to switch the state of the machine when primitive values end of their processing
func (i *TokenizerStateMachine) SwitchToLatestState() error {
	if i.astConstructor.HasComplete() {
		// cannot and no need to route, the ast has parsed an json object
		return nil
	}
	n, err := i.astConstructor.TopElementType()
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

func (i *TokenizerStateMachine) ProcessData(dataSource io.Reader) error {
	bufferReader := bufio.NewReader(dataSource)
	for {
		// 1. the ast complete parsing json, end and not read the rest of bytes
		if i.astConstructor.HasComplete() {
			return nil
		}
		err := i.currentState.ProcessData(bufferReader)
		// 2. the stream ends, and ast is still expecting content, fail.

		if err != nil {
			if err == io.EOF {
				if !i.astConstructor.HasComplete() {
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
	return i.astConstructor.GetAST()
}
