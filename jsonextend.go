package jsonextend

import (
	"errors"
	"io"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

func Parse(reader io.Reader, variables map[string]interface{}) ([]byte, error) {
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)
	err := sm.ProcessData()
	if err != nil {
		return nil, err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return nil, errors.New("invalid json")
	}
	ast := sm.GetAST()
	return interpreter.PrettyInterpret(ast, variables, Marshal)
}

func Unmarshal(reader io.Reader, variables map[string]interface{}, out interface{}) error {
	return interpreter.Unmarshal(reader, variables, out)
}

func Marshal(v interface{}) ([]byte, error) {
	return interpreter.Marshal(v)
}
