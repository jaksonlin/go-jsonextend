package golang

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
)

type ASTGolangBaseBuilder struct {
	astConstructor *astGolangConstructor
	provider       *tokenProvider
}

func NewASTGolangBaseBuilder(obj interface{}) (constructor.ASTBuilder, error) {
	provider, err := newTokenProvider(obj)
	if err != nil {
		return nil, err
	}
	return &ASTGolangBaseBuilder{
		astConstructor: newASTConstructor(),
		provider:       provider,
	}, nil
}

var _ constructor.ASTBuilder = &ASTGolangBaseBuilder{}

// put the store to syntax symbol here, to decouple the relation of reader and writer
func (t *ASTGolangBaseBuilder) GetNextTokenType() (token.TokenType, error) {

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

func (t *ASTGolangBaseBuilder) ReadBool() (bool, error) {
	return t.provider.ReadBool()
}
func (t *ASTGolangBaseBuilder) ReadNull() error {
	return t.provider.ReadNull()
}

func (t *ASTGolangBaseBuilder) ReadNumber() (float64, error) {
	return t.provider.ReadNumber()
}

func (t *ASTGolangBaseBuilder) ReadString() ([]byte, error) {
	return t.provider.ReadString()
}

func (t *ASTGolangBaseBuilder) ReadVariable() ([]byte, error) {
	return t.provider.ReadVariable()
}

func (t *ASTGolangBaseBuilder) RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error {
	return t.astConstructor.RecordStateValue(valueType, nodeValue)
}

func (i *ASTGolangBaseBuilder) GetAST() ast.JsonNode {
	return i.astConstructor.GetAST()
}

func (i *ASTGolangBaseBuilder) HasComplete() bool {
	return i.astConstructor.HasComplete()
}

func (i *ASTGolangBaseBuilder) TopElementType() (ast.AST_NODETYPE, error) {
	return i.astConstructor.TopElementType()
}

func (i *ASTGolangBaseBuilder) HasOpenElements() bool {
	return i.astConstructor.HasOpenElements()
}
