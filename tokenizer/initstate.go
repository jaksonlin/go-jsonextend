package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/astbuilder"
	"github.com/jaksonlin/go-jsonextend/token"
)

type InitState struct {
	TokenReader
}

var _ Tokenizer = &InitState{}

func (i *InitState) GetMode() StateMode {
	return INIT_MODE
}

func (i *InitState) ProcessData(provider astbuilder.TokenProvider) error {

	nextTokenType, err := provider.GetNextTokenType()
	if err != nil {
		return err
	}

	return i.switchState(nextTokenType)
}

func (i *InitState) switchState(nextTokenType token.TokenType) error {
	err := i.stateMachine.SwitchStateByToken(nextTokenType)
	if err != nil {
		switch nextTokenType {
		case token.TOKEN_SPACE:
		default:
			return NewErrorIncorrectToken(i.GetMode(), nextTokenType)
		}
	}
	return nil
}
