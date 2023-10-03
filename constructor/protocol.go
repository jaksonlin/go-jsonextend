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

type ASTManager interface {
	RecordSyntaxSymbol(b token.TokenType) error
	RecordSyntaxValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error
	GetAST() ast.JsonNode
	HasComplete() bool
	TopElementType() (ast.AST_NODETYPE, error)
	HasOpenElements() bool
}

type ASTBuilder interface {
	ASTManager
	TokenProvider
}
