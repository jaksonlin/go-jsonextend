package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/astbuilder"
)

type NumberState struct {
	PrimitiveValueTokenStateBase
}

var _ PrimitiveValueTokenizer = &NumberState{}

func (i *NumberState) GetMode() StateMode {
	return NUMBER_MODE
}

func (i *NumberState) ProcessData(provider astbuilder.TokenProvider) error {
	f64, err := provider.ReadNumber()
	if err != nil {
		return err
	}

	err = i.storeTokenValue(i.GetMode(), f64)

	if err != nil {
		return err
	}

	return i.switchState()

}
