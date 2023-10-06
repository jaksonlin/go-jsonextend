package golang

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/constructor"
)

type astGolangConstructor struct {
	ast *ast.JsonextAST
}

var _ constructor.ASTManager = &astGolangConstructor{}

func newASTConstructor() *astGolangConstructor {
	return &astGolangConstructor{
		ast: ast.NewJsonextAST(),
	}
}

func (i *astGolangConstructor) RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error {
	return i.ast.CreateNewASTNode(valueType, nodeValue)
}

func (i *astGolangConstructor) GetAST() ast.JsonNode {
	return i.ast.GetAST()
}

func (i *astGolangConstructor) HasComplete() bool {
	return i.ast.HasComplete()
}

func (i *astGolangConstructor) TopElementType() (ast.AST_NODETYPE, error) {
	return i.ast.TopElementType()
}

func (i *astGolangConstructor) HasOpenElements() bool {
	return i.ast.HasOpenElement()
}
