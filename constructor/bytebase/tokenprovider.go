package bytebase

import (
	"bufio"
	"io"
	"strconv"

	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/token"
)

type tokenProvider struct {
	dataSource     *bufio.Reader
	CurrentOffset  int
	LastReadLength int // this can give us the correct startoffset of current element
}

func newTokenProvider(reader io.Reader) *tokenProvider {
	return &tokenProvider{
		dataSource: bufio.NewReader(reader),
	}
}

var _ constructor.TokenProvider = &tokenProvider{}

func (t *tokenProvider) ReadBool() (bool, error) {
	// in boolean state we will consume until it is the end of boolean.
	data, err := t.dataSource.Peek(1)
	if err != nil {
		return false, err
	}

	numberOfRead := 0
	if data[0] == 't' {
		numberOfRead = 4
	} else if data[0] == 'f' {
		numberOfRead = 5
	} else {
		return false, ErrorIncorrectCharacter
	}

	rs := make([]byte, numberOfRead)

	_, err = io.ReadFull(t.dataSource, rs)
	if err != nil {
		return false, err
	}

	t.LastReadLength = numberOfRead
	t.CurrentOffset += t.LastReadLength

	rsBoolean, err := strconv.ParseBool(string(rs))
	if err != nil {
		return false, ErrorIncorrectValueForState
	}
	return rsBoolean, nil
}

func (t *tokenProvider) GetNextTokenType() (token.TokenType, error) {

	nextByte, err := t.dataSource.ReadByte()
	if err != nil {
		return token.TOKEN_DUMMY, err
	}

	nextTokenType := GetTokenTypeByStartCharacter(nextByte)

	if ShouldUnreadByte(nextTokenType) {
		err = t.dataSource.UnreadByte()
		if err != nil {
			return token.TOKEN_DUMMY, err
		}
	} else {
		t.LastReadLength = 1
		t.CurrentOffset += t.LastReadLength
	}

	return nextTokenType, nil
}

func (t *tokenProvider) ReadNull() error {
	rs := make([]byte, 4)
	_, err := io.ReadFull(t.dataSource, rs)
	if err != nil {
		return err
	}
	t.LastReadLength = 4
	t.CurrentOffset += t.LastReadLength
	if string(rs) != "null" {
		return ErrorIncorrectValueForState
	}
	return nil

}

func (t *tokenProvider) ReadNumber() (float64, error) {
	lengthOfNumber := 1
	for {
		nextByte, err := t.dataSource.Peek(lengthOfNumber)
		if err != nil {
			if err != io.EOF {
				return 0, err
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
	_, err := io.ReadFull(t.dataSource, result)
	if err != nil {
		return 0, err
	}
	t.LastReadLength = lengthOfNumber
	t.CurrentOffset += t.LastReadLength
	f64, err := strconv.ParseFloat(string(result), 64)
	if err != nil {
		return 0, ErrorIncorrectValueForState
	}
	return f64, nil
}

func (t *tokenProvider) ReadString() ([]byte, error) {

	// in order to deal with the multiple slash/escape sequence, we need a flag to check the string state
	isSlashEnclosed := true
	stringLength := 1
	validQuotationCount := 0
	for {
		nextByte, err := t.dataSource.Peek(stringLength)
		if err != nil {
			return nil, err
		}

		// slashes are closed
		if isSlashEnclosed {
			// is double quotation mark, end of string
			if nextByte[stringLength-1] == '"' {
				validQuotationCount += 1
				if validQuotationCount == 2 {
					rs := make([]byte, stringLength)
					_, err := io.ReadFull(t.dataSource, rs)
					if err != nil {
						return nil, err
					}
					t.LastReadLength = stringLength
					t.CurrentOffset += t.LastReadLength
					return rs, nil
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

func (t *tokenProvider) ReadVariable() ([]byte, error) {
	variable, err := t.dataSource.ReadBytes('}')
	if err != nil {
		return nil, err
	}
	t.LastReadLength = len(variable)
	t.CurrentOffset += t.LastReadLength
	return variable, nil
}
