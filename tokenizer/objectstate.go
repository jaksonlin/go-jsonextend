package tokenizer

import (
	"bufio"
)

type ObjectState struct {
	TokenReader
}

var _ JzoneTokenizer = &ObjectState{}

func (i *ObjectState) GetMode() StateMode {
	return OBJECT_MODE
}

func (i *ObjectState) ProcessData(dataSource *bufio.Reader) error {

	nextByteToken, err := i.PreprocessToken(dataSource)
	if err != nil {
		return err
	}

	return i.switchState(nextByteToken)
}

func (i *ObjectState) switchState(nextTokenType TokenType) error {
	//route trigger token
	err := i.stateMachine.SwitchStateByToken(nextTokenType)
	if err != nil {
		switch nextTokenType {
		// valid but not trigger symbol
		case TOKEN_SPACE:
		case TOKEN_COLON:
		case TOKEN_COMMA:
		// enclose symbol
		case TOKEN_RIGHT_BRACE:
			i.stateMachine.SwitchToLatestState()
		default:
			return ErrorIncorrectCharacter
		}
	}
	return nil
}
