package ast

type ASTConstructor struct {
	ast           *jzoneAST
	syntaxChecker *syntaxChecker
}

func NewASTConstructor() *ASTConstructor {
	return &ASTConstructor{
		ast:           newJzoneAST(),
		syntaxChecker: newSyntaxChecker(),
	}
}

// when a json symbol is read, push it to syntax checker and construct the AST stack elements (as described in ast.go)
func (i *ASTConstructor) RecordSyntaxSymbol(b byte) error {
	i.syntaxChecker.PushSymbol(b)
	//routing base on symbol
	switch b {
	case '{':
		return i.ast.CreateNewASTNode(AST_OBJECT, nil)
	case '[':
		return i.ast.CreateNewASTNode(AST_ARRAY, nil)
	case ']':
		fallthrough
	case '}':
		// check syntax before manipulate the AST
		err := i.syntaxChecker.Enclose(b)
		if err != nil {
			return err
		}
		err = i.ast.EncloseLatestElements()
		if err != nil {
			return err
		}
	case ':':
	case ',':
	default:
		return ErrorIncorrectSyntaxSymbolForConstructAST
	}
	return nil
}

func (i *ASTConstructor) RecordSyntaxValue(valueType AST_NODETYPE, nodeValue interface{}) error {
	i.syntaxChecker.PushValue(valueType)
	return i.ast.CreateNewASTNode(valueType, nodeValue)
}

func (i *ASTConstructor) GetAST() JsonNode {
	return i.ast.GetAST()
}

func (i *ASTConstructor) HasComplete() bool {
	return i.ast.HasComplete()
}

func (i *ASTConstructor) TopElementType() (AST_NODETYPE, error) {
	return i.ast.TopElementType()
}

func (i *ASTConstructor) HasOpenElements() bool {
	return i.ast.HasOpenElement()
}
