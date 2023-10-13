package golang

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
)

type astGolangConstructor struct {
	ast *ast.JsonextAST
}

var _ constructor.ASTStateManagement = (*astGolangConstructor)(nil)
var _ constructor.NodeConstructor = (*astGolangConstructor)(nil)

func newASTConstructor() *astGolangConstructor {
	return &astGolangConstructor{
		ast: ast.NewJsonextAST(),
	}
}

// when a json symbol is read, push it to syntax checker and construct the AST stack elements (as described in ast.go)
func (i *astGolangConstructor) RecordSyntaxSymbol(b token.TokenType) error {
	//routing base on symbol
	switch b {
	case token.TOKEN_LEFT_BRACE:
		_, err := i.ast.CreateNewASTNode(ast.AST_OBJECT, nil)
		return err
	case token.TOKEN_LEFT_BRACKET:
		_, err := i.ast.CreateNewASTNode(ast.AST_ARRAY, nil)
		return err
	case token.TOKEN_RIGHT_BRACKET:
		fallthrough
	case token.TOKEN_RIGHT_BRACE:
		_, err := i.ast.EncloseLatestElements()
		if err != nil {
			return err
		}
	default:
		return ErrorIncorrectSyntaxSymbolForConstructAST
	}
	return nil
}

func (i *astGolangConstructor) CreateNodeWithValue(valueType ast.AST_NODETYPE, nodeValue interface{}) (ast.JsonNode, error) {
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
