package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/util"
)

type VariableState struct {
	PrimitiveValueTokenStateBase
}

var _ PrimitiveValueTokenizer = &VariableState{}

func (i *VariableState) GetMode() StateMode {
	return VARIABLE_MODE
}

func (i *VariableState) ProcessData(provider constructor.TokenProvider) error {

	variable, err := provider.ReadVariable()
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
