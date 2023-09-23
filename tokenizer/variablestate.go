package tokenizer

import (
	"bufio"

	"github.com/jaksonlin/go-jsonextend/util"
)

type VariableState struct {
	PrimitiveValueTokenStateBase
}

var _ JzonePrimitiveTokenizer = &VariableState{}

func (i *VariableState) GetMode() StateMode {
	return VARIABLE_MODE
}

func (i *VariableState) ProcessData(dataSource *bufio.Reader) error {

	variable, err := dataSource.ReadBytes('}')
	if err != nil {
		return err
	}
	if !util.RegStringWithVariable.Match(variable) {
		return ErrorExtendedVariableFormatIncorrect
	}

	err = i.storeTokenValue(i.GetMode(), variable)
	if err != nil {
		return err
	}
	return i.switchState()
}
