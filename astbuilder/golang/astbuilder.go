package golang

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/astbuilder"
	"github.com/jaksonlin/go-jsonextend/token"
)

type ASTGolangBaseBuilder struct {
	astConstructor *astGolangConstructor
	provider       *tokenProvider
}

func NewASTGolangBaseBuilder(obj interface{}, options []astbuilder.TokenProviderOptions) (astbuilder.ASTBuilder, error) {
	provider, err := newRootTokenProvider(obj, options)
	if err != nil {
		return nil, err
	}
	return &ASTGolangBaseBuilder{
		astConstructor: newASTConstructor(),
		provider:       provider,
	}, nil
}

var _ astbuilder.ASTBuilder = &ASTGolangBaseBuilder{}

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

func (t *ASTGolangBaseBuilder) ReadNumber() (interface{}, error) {
	return t.provider.ReadNumber()
}

func (t *ASTGolangBaseBuilder) ReadString() ([]byte, error) {
	return t.provider.ReadString()
}

func (t *ASTGolangBaseBuilder) ReadVariable() ([]byte, error) {
	return t.provider.ReadVariable()
}

func (t *ASTGolangBaseBuilder) RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error {
	node, err := t.astConstructor.CreateNodeWithValue(valueType, nodeValue)
	if err != nil {
		return err
	}
	// we do peek for tokenizer's ReadValue, now we pop it out and add any meta/plugin if needed
	valueWorkItem, err := t.provider.workingStack.Pop()
	if err != nil {
		return err
	}
	valueWorkItem.SetMetaAndPlugins(node)
	return nil

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
