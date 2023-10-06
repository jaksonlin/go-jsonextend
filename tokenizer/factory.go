package tokenizer

import (
	"io"

	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/constructor/bytebase"
	"github.com/jaksonlin/go-jsonextend/constructor/golang"
	"github.com/jaksonlin/go-jsonextend/token"
)

func NewTokenizerStateMachineFromIOReader(reader io.Reader) *TokenizerStateMachine {
	astMan := bytebase.NewASTByteBaseBuilder(reader)
	return newTokenizerStateMachine(astMan)
}

func NewTokenizerStateMachineFromGoData(obj interface{}) (*TokenizerStateMachine, error) {
	astMan, err := golang.NewASTGolangBaseBuilder(obj)
	if err != nil {
		return nil, err
	}
	return newTokenizerStateMachine(astMan), nil
}

func newTokenizerStateMachine(builder constructor.ASTBuilder) *TokenizerStateMachine {

	sm := TokenizerStateMachine{}
	sm.astBuilder = builder
	sm.initState = &InitState{NewTokenReader(&sm)}
	sm.arrayState = &ArrayState{NewTokenReader(&sm)}
	sm.objectState = &ObjectState{NewTokenReader(&sm)}
	sm.stringState = &StringState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.numberState = &NumberState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.booleanState = &BooleanState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.nullState = &NullState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.variableState = &VariableState{NewPrimitiveValueTokenStateBase(&sm)}
	//construct a route table instead of using switch every where.
	sm.defaultRoute = map[token.TokenType]stateChangeFunc{

		token.TOKEN_STRING: func() error {
			sm.currentState = sm.stringState
			return nil
		},
		token.TOKEN_NUMBER: func() error {
			sm.currentState = sm.numberState
			return nil
		},
		token.TOKEN_BOOLEAN: func() error {
			sm.currentState = sm.booleanState
			return nil
		},
		token.TOKEN_NULL: func() error {
			sm.currentState = sm.nullState
			return nil
		},
		token.TOKEN_LEFT_BRACKET: func() error {
			sm.currentState = sm.arrayState
			return nil
		},
		token.TOKEN_LEFT_BRACE: func() error {
			sm.currentState = sm.objectState
			return nil
		},
		token.TOKEN_VARIABLE: func() error {
			sm.currentState = sm.variableState
			return nil
		},
	}
	sm.currentState = sm.initState
	return &sm
}
