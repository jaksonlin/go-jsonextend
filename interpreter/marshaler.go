package interpreter

import (
	"github.com/jaksonlin/go-jsonextend/astbuilder"
	"github.com/jaksonlin/go-jsonextend/astbuilder/golang"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

func marshal(v interface{}, depth int, variables map[string]interface{}, options []astbuilder.TokenProviderOptions) ([]byte, error) {
	if depth > maxDepth {
		return nil, ErrorSelfCallTooDeep
	}
	sm, err := tokenizer.NewTokenizerStateMachineFromGoData(v, options)
	if err != nil {
		return nil, err
	}
	err = sm.ProcessData()
	if err != nil {
		return nil, err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return nil, ErrorInvalidJson
	}
	ast := sm.GetAST()
	return InterpretAST(ast, variables, func(v interface{}) ([]byte, error) {
		return marshal(v, depth+1, variables, options)
	})
}

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v, 1, nil, nil)
}

func MarshalWithVariables(v interface{}, variables map[string]interface{}) ([]byte, error) {
	return marshal(v, 1, variables, []astbuilder.TokenProviderOptions{golang.EnableJsonExtTag})
}

func MarshalIntoTemplate(v interface{}) ([]byte, error) {
	return marshal(v, 1, nil, []astbuilder.TokenProviderOptions{golang.EnableJsonExtTag})
}
