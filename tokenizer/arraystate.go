package tokenizer

import (
	"bufio"
)

type ArrayState struct {
	TokenReader
}

var _ JzoneTokenizer = &ArrayState{}

func (i *ArrayState) GetMode() StateMode {
	return ARRAY_MODE
}

func (i *ArrayState) ProcessData(dataSource *bufio.Reader) error {

	nextByteToken, err := i.PreprocessToken(dataSource)
	if err != nil {
		return err
	}

	return i.switchState(nextByteToken)
}

func (i *ArrayState) switchState(nextTokenType TokenType) error {
	//route trigger token
	err := i.stateMachine.SwitchStateByToken(nextTokenType)
	if err != nil {
		switch nextTokenType {
		// valid but not trigger
		case TOKEN_SPACE:
		case TOKEN_COMMA:
		// enclose symbol
		case TOKEN_RIGHT_BRACKET:
			i.stateMachine.SwitchToLatestState()
		default:
			return ErrorIncorrectCharacter
		}
	}
	return nil
}
