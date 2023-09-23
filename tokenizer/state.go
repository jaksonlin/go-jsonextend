package tokenizer

import (
	"bufio"
)

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
type JzoneTokenizer interface {
	ProcessData(dataSource *bufio.Reader) error
	GetMode() StateMode
}

// for none-symbol token, they hold primitive value (not array/object), these value should be hold by array/object later,
// store their value to somewhere based on extracted state, able to switchState() without indicator, because there's no way from a primitive
// value to say where the state should go next
type JzonePrimitiveStateProcessor interface {
	switchState() error
	storeTokenValue(mode StateMode, value interface{}) error
}

// eventually as a state in the state machine, it also implement its way of ProcessData from the stream, but this would leave to the concrete state, case by case
type JzonePrimitiveTokenizer interface {
	JzoneTokenizer
	JzonePrimitiveStateProcessor
}

// use a statemachine to route the state and handle the storage of the token value
type PrimitiveValueTokenStateBase struct {
	stateMachine JzoneTokenizerStateMachine
}

var _ JzonePrimitiveStateProcessor = &PrimitiveValueTokenStateBase{}

func NewPrimitiveValueTokenStateBase(sm JzoneTokenizerStateMachine) PrimitiveValueTokenStateBase {
	return PrimitiveValueTokenStateBase{
		stateMachine: sm,
	}
}

// our state machine construct the AST on the fly, therefore we can use the AST for routing when after dealing with value
func (i *PrimitiveValueTokenStateBase) switchState() error {

	if err := i.stateMachine.SwitchToLatestState(); err != nil {
		return err
	}
	return nil
}

// it will also handle the store of the value, insert into the correct location in the AST with its corresponding StateMode(providing the meta information)
func (i *PrimitiveValueTokenStateBase) storeTokenValue(mode StateMode, value interface{}) error {
	return i.stateMachine.RecordSyntaxValue(mode, value)
}
