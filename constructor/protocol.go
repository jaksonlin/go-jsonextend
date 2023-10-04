package constructor

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/token"
)

type TokenProvider interface {
	ReadBool() (bool, error)
	ReadNull() error
	ReadNumber() (float64, error)
	ReadString() ([]byte, error)
	ReadVariable() ([]byte, error)
	GetNextTokenType() (token.TokenType, error)
}

// user interface hide the implementation of offset recording from tokenizer
type ASTBuilderFacade interface {
	RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error
}

type ASTStateRecorder interface {
	RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}, currentPosition int, lastReadLength int) error
}

type ASTManagerBase interface {
	GetAST() ast.JsonNode
	HasComplete() bool
	TopElementType() (ast.AST_NODETYPE, error)
	HasOpenElements() bool
}
type ASTByteReaderManager interface {
	ASTStateRecorder
	ASTManagerBase
}

type ASTBuilder interface {
	ASTBuilderFacade
	ASTManagerBase
	TokenProvider
}
