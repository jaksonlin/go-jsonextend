package tokenizer

import "github.com/jaksonlin/go-jsonextend/constructor"

type BooleanState struct {
	PrimitiveValueTokenStateBase
}

var _ PrimitiveValueTokenizer = &BooleanState{}

func (i *BooleanState) GetMode() StateMode {
	return BOOLEAN_MODE
}

func (i *BooleanState) ProcessData(provider constructor.TokenProvider) error {
	// in boolean state we will consume until it is the end of boolean.
	value, err := provider.ReadBool()
	if err != nil {
		return err
	}

	err = i.storeTokenValue(i.GetMode(), value)
	if err != nil {
		return err
	}
	return i.switchState()
}
