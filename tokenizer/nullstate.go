package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/constructor"
)

type NullState struct {
	PrimitiveValueTokenStateBase
}

var _ PrimitiveValueTokenizer = &NullState{}

func (i *NullState) GetMode() StateMode {
	return NULL_MODE
}

func (i *NullState) ProcessData(provider constructor.TokenProvider) error {
	err := provider.ReadNull()
	if err != nil {
		return err
	}
	err = i.storeTokenValue(i.GetMode(), nil)
	if err != nil {
		return err
	}

	return i.switchState()

}
