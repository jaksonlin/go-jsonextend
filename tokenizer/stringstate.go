package tokenizer

import (
	"github.com/jaksonlin/go-jsonextend/constructor"
	"github.com/jaksonlin/go-jsonextend/util"
)

type StringState struct {
	PrimitiveValueTokenStateBase
}

var _ PrimitiveValueTokenizer = &StringState{}

func (i *StringState) GetMode() StateMode {
	return STRING_MODE
}

func (i *StringState) ProcessData(provider constructor.TokenProvider) error {

	rs, err := provider.ReadString()
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
