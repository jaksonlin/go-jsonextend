package tokenizer

import (
	"bufio"
	"io"
	"strconv"
)

type BooleanState struct {
	PrimitiveValueTokenStateBase
}

var _ JzonePrimitiveTokenizer = &BooleanState{}

func (i *BooleanState) GetMode() StateMode {
	return BOOLEAN_MODE
}

func (i *BooleanState) ProcessData(dataSource *bufio.Reader) error {
	// in boolean state we will consume until it is the end of boolean.
	data, err := dataSource.Peek(1)
	if err != nil {
		return err
	}

	numberOfRead := 0
	if data[0] == 't' {
		numberOfRead = 4
	} else if data[0] == 'f' {
		numberOfRead = 5
	} else {
		return ErrorIncorrectCharacter
	}

	rs := make([]byte, numberOfRead)

	_, err = io.ReadFull(dataSource, rs)
	if err != nil {
		return err
	}

	rsBoolean, err := strconv.ParseBool(string(rs))
	if err != nil {
		return ErrorIncorrectValueForState
	}

	err = i.storeTokenValue(i.GetMode(), rsBoolean)
	if err != nil {
		return err
	}
	return i.switchState()

}
