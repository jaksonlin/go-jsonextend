package tokenizer

type TokenReader struct {
	stateMachine *TokenizerStateMachine
}

func NewTokenReader(stateMachine *TokenizerStateMachine) TokenReader {
	return TokenReader{
		stateMachine: stateMachine,
	}
}
