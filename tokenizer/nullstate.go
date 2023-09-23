package tokenizer

import (
	"bufio"
	"io"
)

type NullState struct {
	PrimitiveValueTokenStateBase
}

var _ JzonePrimitiveTokenizer = &NullState{}

func (i *NullState) GetMode() StateMode {
	return NULL_MODE
}

func (i *NullState) ProcessData(dataSource *bufio.Reader) error {
	rs := make([]byte, 4)
	_, err := io.ReadFull(dataSource, rs)
	if err != nil {
		return err
	}

	if string(rs) != "null" {
		return ErrorIncorrectValueForState
	}

	err = i.storeTokenValue(i.GetMode(), "null")
	if err != nil {
		return err
	}

	return i.switchState()

}
