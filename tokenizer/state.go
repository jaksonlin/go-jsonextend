package tokenizer

import "github.com/jaksonlin/go-jsonextend/astbuilder"

type StateMode uint

const (
	INIT_MODE StateMode = iota
	OBJECT_MODE
	ARRAY_MODE
	STRING_MODE
	NUMBER_MODE
	BOOLEAN_MODE
	NULL_MODE
	VARIABLE_MODE
	STRING_VARIABLE_MODE
)

// common tokenizer implementation
type Tokenizer interface {
	ProcessData(provider astbuilder.TokenProvider) error
	GetMode() StateMode
}

// for primitive value token, they hold primitive value (not array/object)
// store their value to somewhere based on extracted value,
// and the switch of state is based on AST current state, therefore no need for parameterized state change
type PrimitiveStateProcessor interface {
	switchState() error
	storeTokenValue(mode StateMode, value interface{}) error
}

type PrimitiveValueTokenizer interface {
	Tokenizer
	PrimitiveStateProcessor
}

// for the base, only provides what `PrimitiveStateProcessor` needs
type PrimitiveValueTokenStateBase struct {
	stateMachine *TokenizerStateMachine
}

var _ PrimitiveStateProcessor = &PrimitiveValueTokenStateBase{}

func NewPrimitiveValueTokenStateBase(sm *TokenizerStateMachine) PrimitiveValueTokenStateBase {
	return PrimitiveValueTokenStateBase{
		stateMachine: sm,
	}
}

// state machine construct the AST on the fly, use the AST for state change after producing a primivite value
func (i *PrimitiveValueTokenStateBase) switchState() error {

	if err := i.stateMachine.SwitchToLatestState(); err != nil {
		return err
	}
	return nil
}

// it will also handle the store of the value to update the AST
func (i *PrimitiveValueTokenStateBase) storeTokenValue(mode StateMode, value interface{}) error {
	return i.stateMachine.RecordStateValue(mode, value)
}
