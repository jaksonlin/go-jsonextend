package astbuilder

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/token"
)

type TokenProvider interface {
	ReadBool() (bool, error)
	ReadNull() error
	ReadNumber() (interface{}, error)
	ReadString() ([]byte, error)
	ReadVariable() ([]byte, error)
	GetNextTokenType() (token.TokenType, error)
}

type NodeConstructor interface {
	CreateNodeWithValue(valueType ast.AST_NODETYPE, nodeValue interface{}) (ast.JsonNode, error)
}

type ASTStateManagement interface {
	GetAST() ast.JsonNode
	HasComplete() bool
	TopElementType() (ast.AST_NODETYPE, error)
	HasOpenElements() bool
}

type ASTBuilder interface {
	ASTStateManagement
	TokenProvider
	RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error
}
