package bytebase

import (
	"io"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/astbuilder"
	"github.com/jaksonlin/go-jsonextend/token"
)

type ASTByteBaseBuilder struct {
	astConstructor *astByteBaseConstructor
	provider       *tokenProvider
}

func NewASTByteBaseBuilder(reader io.Reader) *ASTByteBaseBuilder {
	return &ASTByteBaseBuilder{
		astConstructor: newASTConstructor(),
		provider:       newTokenProvider(reader),
	}
}

var _ astbuilder.ASTBuilder = &ASTByteBaseBuilder{}

// put the store to syntax symbol here, to decouple the relation of reader and writer
func (t *ASTByteBaseBuilder) GetNextTokenType() (token.TokenType, error) {

	nextTokenType, err := t.provider.GetNextTokenType()
	if err != nil {
		return token.TOKEN_DUMMY, err
	}

	if token.IsSymbolToken(nextTokenType) { // note symbol token will be parse in the corresponding primitive value state
		err = t.astConstructor.RecordSyntaxSymbol(nextTokenType)
		if err != nil {
			return token.TOKEN_DUMMY, err
		}
	}

	return nextTokenType, nil
}

func (t *ASTByteBaseBuilder) ReadBool() (bool, error) {
	return t.provider.ReadBool()
}
func (t *ASTByteBaseBuilder) ReadNull() error {
	return t.provider.ReadNull()
}

func (t *ASTByteBaseBuilder) ReadNumber() (interface{}, error) {
	return t.provider.ReadNumber()
}

func (t *ASTByteBaseBuilder) ReadString() ([]byte, error) {
	return t.provider.ReadString()
}

func (t *ASTByteBaseBuilder) ReadVariable() ([]byte, error) {
	return t.provider.ReadVariable()
}

func (t *ASTByteBaseBuilder) RecordSyntaxSymbol(b token.TokenType) error {
	return t.astConstructor.RecordSyntaxSymbol(b)
}

func (t *ASTByteBaseBuilder) RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error {
	_, err := t.astConstructor.CreateNodeWithValue(valueType, nodeValue)
	return err
}

func (i *ASTByteBaseBuilder) GetAST() ast.JsonNode {
	return i.astConstructor.GetAST()
}

func (i *ASTByteBaseBuilder) HasComplete() bool {
	return i.astConstructor.HasComplete()
}

func (i *ASTByteBaseBuilder) TopElementType() (ast.AST_NODETYPE, error) {
	return i.astConstructor.TopElementType()
}

func (i *ASTByteBaseBuilder) HasOpenElements() bool {
	return i.astConstructor.HasOpenElements()
}
