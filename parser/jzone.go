package parser

import (
	"errors"
	"io"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

func ParseJzoneString(reader io.Reader, variables map[string]interface{}) (string, error) {
	sm := tokenizer.NewTokenizerStateMachine()
	err := sm.ProcessData(reader)
	if err != nil {
		return "", err
	}
	if sm.GetASTConstructor().HasOpenElements() {
		return "", errors.New("invalid json")
	}
	ast := sm.GetASTConstructor().GetAST()
	return interpreter.Interpret(ast, variables)
}
