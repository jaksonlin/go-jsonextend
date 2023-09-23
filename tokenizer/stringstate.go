package tokenizer

import (
	"bufio"
	"io"

	"github.com/jaksonlin/go-jsonextend/util"
)

type StringState struct {
	PrimitiveValueTokenStateBase
}

var _ JzonePrimitiveTokenizer = &StringState{}

func (i *StringState) GetMode() StateMode {
	return STRING_MODE
}

func (i *StringState) ProcessData(dataSource *bufio.Reader) error {

	// in order to deal with the multiple slash/escape sequence, we need a flag to check the string state
	isSlashEnclosed := true
	stringLength := 1
	validQuotationCount := 0
	for {
		nextByte, err := dataSource.Peek(stringLength)
		if err != nil {
			return err
		}

		// slashes are closed
		if isSlashEnclosed {
			// is double quotation mark, end of string
			if nextByte[stringLength-1] == '"' {
				validQuotationCount += 1
				if validQuotationCount == 2 {
					return i.handleStringType(dataSource, stringLength)
				}
			} else if nextByte[stringLength-1] == 0x5c { // is slash
				isSlashEnclosed = false
			}
		} else {
			isSlashEnclosed = true // flip the slash when there's nextByte, in escape mode skip quotation count check
		}
		stringLength += 1

	}
}

func (i *StringState) handleStringType(dataSource *bufio.Reader, stringLen int) error {

	rs := make([]byte, stringLen)
	_, err := io.ReadFull(dataSource, rs)
	if err != nil {
		return err
	}

	if util.RegStringWithVariable.Match(rs) {
		err = i.storeTokenValue(STRING_VARIABLE_MODE, rs)
	} else {
		err = i.storeTokenValue(i.GetMode(), rs)
	}

	if err != nil {
		return err
	}
	return i.switchState()
}
