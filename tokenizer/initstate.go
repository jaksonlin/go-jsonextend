package tokenizer

import (
	"bufio"
)

type InitState struct {
	TokenReader
}

var _ JzoneTokenizer = &InitState{}

func (i *InitState) GetMode() StateMode {
	return INIT_MODE
}

func (i *InitState) ProcessData(dataSource *bufio.Reader) error {

	nextByteToken, err := i.PreprocessToken(dataSource)
	if err != nil {
		return err
	}

	return i.switchState(nextByteToken)
}

func (i *InitState) switchState(nextTokenType TokenType) error {
	err := i.stateMachine.SwitchStateByToken(nextTokenType)
	if err != nil {
		switch nextTokenType {
		case TOKEN_SPACE:
		default:
			return ErrorIncorrectCharacter
		}
	}
	return nil
}
