package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/astbuilder"
	"github.com/jaksonlin/go-jsonextend/token"
)

type ObjectState struct {
	TokenReader
}

var _ Tokenizer = &ObjectState{}

func (i *ObjectState) GetMode() StateMode {
	return OBJECT_MODE
}

func (i *ObjectState) ProcessData(provider astbuilder.TokenProvider) error {

	nextTokenType, err := provider.GetNextTokenType()
	if err != nil {
		return err
	}

	return i.switchState(nextTokenType)
}

func (i *ObjectState) switchState(nextTokenType token.TokenType) error {
	//route trigger token
	err := i.stateMachine.SwitchStateByToken(nextTokenType)
	if err != nil {
		switch nextTokenType {
		// valid but not trigger symbol
		case token.TOKEN_SPACE:
		case token.TOKEN_COLON:
		case token.TOKEN_COMMA:
		// enclose symbol
		case token.TOKEN_RIGHT_BRACE:
			i.stateMachine.SwitchToLatestState()
		default:
			return NewErrorIncorrectToken(i.GetMode(), nextTokenType)
		}
	}
	return nil
}
