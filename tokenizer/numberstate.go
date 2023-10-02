package tokenizer

import (
	"bufio"
	"io"
	"strconv"
)

type NumberState struct {
	PrimitiveValueTokenStateBase
}

var _ JzonePrimitiveTokenizer = &NumberState{}

func (i *NumberState) GetMode() StateMode {
	return NUMBER_MODE
}

func (i *NumberState) ProcessData(dataSource *bufio.Reader) error {
	lengthOfNumber := 1
	for {
		nextByte, err := dataSource.Peek(lengthOfNumber)
		if err != nil {
			if err != io.EOF {
				return err
			}
			lengthOfNumber -= 1
			//bare number
			break
		}
		if !isJSONNumberByte(nextByte[lengthOfNumber-1]) {
			lengthOfNumber -= 1 // remove the invalid location
			break
		}
		lengthOfNumber += 1
	}

	result := make([]byte, lengthOfNumber)
	_, err := io.ReadFull(dataSource, result)
	if err != nil {
		return err
	}

	f64, err := strconv.ParseFloat(string(result), 64)
	if err != nil {
		return ErrorIncorrectValueForState
	}
	err = i.storeTokenValue(i.GetMode(), f64)

	if err != nil {
		return err
	}

	return i.switchState()

}

func isJSONNumberByte(b byte) bool {
	return (b >= '0' && b <= '9') || b == '-' || b == '.' || b == 'e' || b == 'E' || b == '+'
}
