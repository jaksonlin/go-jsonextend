package jsonextend

import (
	"bytes"
	"errors"
	"io"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
	"github.com/jaksonlin/go-jsonextend/unmarshaler"
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
func unmarshal(reader io.Reader, variables map[string]interface{}, out interface{}, depth int) error {
	if depth > maxDepth {
		return errors.New("recursion depth exceeded")
	}
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)
	err := sm.ProcessData()
	if err != nil {
		return err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return errors.New("invalid json")
	}
	ast := sm.GetAST()
	return unmarshaler.UnmarshallAST(ast, variables, Marshal, func(v []byte, out interface{}) error {
		return unmarshal(bytes.NewReader(v), variables, out, depth+1)
	}, out)
}
func Unmarshal(reader io.Reader, variables map[string]interface{}, out interface{}) error {
	return unmarshal(reader, variables, out, 1)
}

const maxDepth = 4

func marshal(v interface{}, depth int) ([]byte, error) {
	if depth > maxDepth {
		return nil, errors.New("recursion depth exceeded")
	}
	sm, err := tokenizer.NewTokenizerStateMachineFromGoData(v)
	if err != nil {
		return nil, err
	}
	err = sm.ProcessData()
	if err != nil {
		return nil, err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return nil, errors.New("invalid object")
	}
	ast := sm.GetAST()
	return interpreter.InterpretAST(ast, nil, func(v interface{}) ([]byte, error) {
		return marshal(v, depth+1)
	})
}

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v, 1)
}
