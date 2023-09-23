package tokenizer

import "bufio"

type TokenPreProcessor interface {
	PreprocessToken(dataSource *bufio.Reader) (TokenType, error)
}

type TokenReader struct {
	stateMachine JzoneTokenizerStateMachine
}

var _ TokenPreProcessor = &TokenReader{}

func NewTokenReader(sm JzoneTokenizerStateMachine) TokenReader {
	return TokenReader{
		stateMachine: sm,
	}
}

func (i *TokenReader) PreprocessToken(dataSource *bufio.Reader) (TokenType, error) {

	nextByte, err := dataSource.ReadByte()
	if err != nil {
		return TOKEN_DUMMY, err
	}

	nextTokenType := GetTokenTypeByStartCharacter(nextByte)

	if ShouldUnreadByte(nextTokenType) {
		err = dataSource.UnreadByte()
		if err != nil {
			return TOKEN_DUMMY, err
		}
	}

	if IsSymbolToken(nextTokenType) { // note symbol token will be parse in the corresponding primitive value state
		err = i.stateMachine.RecordSyntaxSymbol(nextByte)
		if err != nil {
			return TOKEN_DUMMY, err
		}
	}

	return nextTokenType, nil
}
