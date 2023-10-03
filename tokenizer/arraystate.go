package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
)

type ArrayState struct {
	TokenReader
}

var _ Tokenizer = &ArrayState{}

func (i *ArrayState) GetMode() StateMode {
	return ARRAY_MODE
}

func (i *ArrayState) ProcessData(provider constructor.TokenProvider) error {

	nextTokenType, err := provider.GetNextTokenType()
	if err != nil {
		return err
	}

	return i.switchState(nextTokenType)
}

func (i *ArrayState) switchState(nextTokenType token.TokenType) error {
	//route trigger token
	err := i.stateMachine.SwitchStateByToken(nextTokenType)
	if err != nil {
		switch nextTokenType {
		// valid but not trigger
		case token.TOKEN_SPACE:
		case token.TOKEN_COMMA:
		// enclose symbol
		case token.TOKEN_RIGHT_BRACKET:
			i.stateMachine.SwitchToLatestState()
		default:
			return NewErrorIncorrectToken(i.GetMode(), nextTokenType)
		}
	}
	return nil
}
