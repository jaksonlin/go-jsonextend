package bytebase

import (
	"github.com/jaksonlin/go-jsonextend/ast"
	astbuilder "github.com/jaksonlin/go-jsonextend/astbuilder"
	"github.com/jaksonlin/go-jsonextend/token"
)

type astByteBaseConstructor struct {
	ast           *ast.JsonextAST
	syntaxChecker *syntaxChecker
}

var _ astbuilder.ASTStateManagement = &astByteBaseConstructor{}
var _ astbuilder.NodeConstructor = &astByteBaseConstructor{}

func newASTConstructor() *astByteBaseConstructor {
	return &astByteBaseConstructor{
		ast:           ast.NewJsonextAST(),
		syntaxChecker: newSyntaxChecker(),
	}
}

// when a json symbol is read, push it to syntax checker and construct the AST stack elements (as described in ast.go)
func (i *astByteBaseConstructor) RecordSyntaxSymbol(b token.TokenType) error {
	//routing base on symbol
	switch b {
	case token.TOKEN_LEFT_BRACE:
		i.syntaxChecker.PushSymbol('{')
		_, err := i.ast.CreateNewASTNode(ast.AST_OBJECT, nil)
		return err
	case token.TOKEN_LEFT_BRACKET:
		i.syntaxChecker.PushSymbol('[')
		_, err := i.ast.CreateNewASTNode(ast.AST_ARRAY, nil)
		return err
	case token.TOKEN_RIGHT_BRACKET:
		i.syntaxChecker.PushSymbol(']')
		// check syntax before manipulate the AST
		err := i.syntaxChecker.Enclose(']')
		if err != nil {
			return err
		}
		_, err = i.ast.EncloseLatestElements()
		if err != nil {
			return err
		}
	case token.TOKEN_RIGHT_BRACE:
		i.syntaxChecker.PushSymbol('}')
		// check syntax before manipulate the AST
		err := i.syntaxChecker.Enclose('}')
		if err != nil {
			return err
		}
		_, err = i.ast.EncloseLatestElements()
		if err != nil {
			return err
		}
	case token.TOKEN_COLON:
		i.syntaxChecker.PushSymbol(':')
	case token.TOKEN_COMMA:
		i.syntaxChecker.PushSymbol(',')
	default:
		return ErrorIncorrectSyntaxSymbolForConstructAST
	}
	return nil
}

func (i *astByteBaseConstructor) CreateNodeWithValue(valueType ast.AST_NODETYPE, nodeValue interface{}) (ast.JsonNode, error) {
	i.syntaxChecker.PushValue(valueType)
	return i.ast.CreateNewASTNode(valueType, nodeValue)
}

func (i *astByteBaseConstructor) GetAST() ast.JsonNode {
	return i.ast.GetAST()
}

func (i *astByteBaseConstructor) HasComplete() bool {
	return i.ast.HasComplete()
}

func (i *astByteBaseConstructor) TopElementType() (ast.AST_NODETYPE, error) {
	return i.ast.TopElementType()
}

func (i *astByteBaseConstructor) HasOpenElements() bool {
	return i.ast.HasOpenElement()
}
