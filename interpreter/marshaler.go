package interpreter

import (
	"github.com/jaksonlin/go-jsonextend/tokenizer"
)

func marshal(v interface{}, depth int) ([]byte, error) {
	if depth > maxDepth {
		return nil, ErrorSelfCallTooDeep
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
		return nil, ErrorInvalidJson
	}
	ast := sm.GetAST()
	return InterpretAST(ast, nil, func(v interface{}) ([]byte, error) {
		return marshal(v, depth+1)
	})
}

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v, 1)
}
