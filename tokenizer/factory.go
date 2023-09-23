package tokenizer

func NewTokenizerStateMachine() JzoneTokenizerStateMachine {

	sm := TokenizerStateMachine{}
	sm.initState = &InitState{NewTokenReader(&sm)}
	sm.arrayState = &ArrayState{NewTokenReader(&sm)}
	sm.objectState = &ObjectState{NewTokenReader(&sm)}
	sm.stringState = &StringState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.numberState = &NumberState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.booleanState = &BooleanState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.nullState = &NullState{NewPrimitiveValueTokenStateBase(&sm)}
	sm.variableState = &VariableState{NewPrimitiveValueTokenStateBase(&sm)}
	//construct a route table instead of using switch every where.
	sm.defaultRoute = map[TokenType]stateChangeFunc{

		TOKEN_STRING: func() error {
			sm.currentState = sm.stringState
			return nil
		},
		TOKEN_NUMBER: func() error {
			sm.currentState = sm.numberState
			return nil
		},
		TOKEN_BOOLEAN: func() error {
			sm.currentState = sm.booleanState
			return nil
		},
		TOKEN_NULL: func() error {
			sm.currentState = sm.nullState
			return nil
		},
		TOKEN_LEFT_BRACKET: func() error {
			sm.currentState = sm.arrayState
			return nil
		},
		TOKEN_LEFT_BRACE: func() error {
			sm.currentState = sm.objectState
			return nil
		},
		TOKEN_VARIABLE: func() error {
			sm.currentState = sm.variableState
			return nil
		},
	}
	sm.Reset()
	return &sm
}
