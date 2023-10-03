package jsonextend

import (
	"errors"
	"io"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

func Parse(reader io.Reader, variables map[string]interface{}) (string, error) {
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)
	err := sm.ProcessData()
	if err != nil {
		return "", err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return "", errors.New("invalid json")
	}
	ast := sm.GetAST()
	return interpreter.Interpret(ast, variables)
}

func Unmarshal(reader io.Reader, variables map[string]interface{}, out interface{}) error {
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)
	err := sm.ProcessData()
	if err != nil {
		return err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return errors.New("invalid json")
	}
	ast := sm.GetAST()
	return interpreter.UnmarshallAST(ast, variables, out)
}
