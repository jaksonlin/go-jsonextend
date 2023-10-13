package ast

// the idea is to maintain a stack holding element of JsonArrayNode, JsonObjectNode and JsonKeyValuePairNode
// any element in json is either `element in array` or `element in object`
// 1. element in array, create a JsonArrayNode, push it to the stack;
// when the value of this array comes, if it is primitive value push it into the JsonArrayNodes's value;
// for none-primitive value, create a new JsonArrayNode | JsonObjectNode and push it to the stack
// 2. element in object, create a JsonObjectNode, push it to the stack;
// the json object's value is maintained in a JsonKeyValuePairNode, when the value comes,
// when a create node operation comes and find that the top element is JsonObjectNode, create a JsonKeyValuePairNode and push it to the stack
// then put the value of the requested create node operation into the JsonKeyValuePairNode's `key` field.
// when the key-value pairs' value comes, if it is primitive value push it into the JsonKeyValuePairNode's `value` field;
// otherwise, create a new JsonArrayNode | JsonObjectNode and push it to the stack.
// once the JsonKeyValuePairNode's `value` field is set, finalize the JsonKeyValuePairNode and append it to the JsonObjectNode's value field
// when to finalize the JsonObjectNode? when the tokenizer finds a } | ] symbol, it will call the EncloseLatestElements() method
// this means that the current element is finished, so we need to pop the top element from the stack and put it into the owner element

import (
	"github.com/jaksonlin/go-jsonextend/util"
)

type astState uint

const (
	AST_STATE_INIT astState = iota
	AST_STATE_INTERMIDIEATE
	AST_STATE_FINISHED
)

type JsonextAST struct {
	ast      JsonNode
	astTrace *util.Stack[JsonNode]
	state    astState
}

func NewJsonextAST() *JsonextAST {
	return &JsonextAST{
		ast:      nil,
		astTrace: util.NewStack[JsonNode](),
		state:    AST_STATE_INIT,
	}
}

func (i *JsonextAST) GetAST() JsonNode {
	return i.ast
}

func (i *JsonextAST) createRootNode(t AST_NODETYPE, value interface{}) (JsonNode, error) {
	n, err := NodeFactory(t, value)
	if err != nil {
		return nil, err
	}
	i.ast = n

	if t == AST_ARRAY || t == AST_OBJECT {
		i.astTrace.Push(n)
		i.state = AST_STATE_INTERMIDIEATE
	} else {
		i.state = AST_STATE_FINISHED
	}
	return n, nil
}

func (i *JsonextAST) CreateNewASTNode(t AST_NODETYPE, value interface{}) (JsonNode, error) {
	if i.state == AST_STATE_FINISHED {
		return nil, ErrorASTComplete
	}
	if i.ast == nil {
		return i.createRootNode(t, value)
	}
	latest, err := i.astTrace.Peek()
	if err != nil {
		return nil, err
	}

	switch realNode := latest.(type) {
	//stack have array at top, awaiting element
	case *JsonArrayNode:
		return i.createNewNodeForArrayObject(realNode, t, value)
		// stack have object at top, awaiting kv pair node
	case *JsonObjectNode:
		return i.createNewNodeForObject(realNode, t, value)
		// stack have kvpari at top, awaiting value node
	case *JsonKeyValuePairNode:
		return i.createValueNodeForKVPairs(realNode, t, value)
	default:
		return nil, ErrorASTUnexpectedElement
	}

}

func (i *JsonextAST) createNewNodeForArrayObject(owner *JsonArrayNode, t AST_NODETYPE, value interface{}) (JsonNode, error) {
	n, err := NodeFactory(t, value)
	if err != nil {
		return nil, err
	}
	if t == AST_ARRAY || t == AST_OBJECT {
		// none-primivtive value, create a new item to hold the furture values
		i.astTrace.Push(n)
	} else {
		// primitive value push to the owner on top of the stack
		owner.Append(n)
	}
	return n, nil
}

func (i *JsonextAST) createNewNodeForObject(owner *JsonObjectNode, t AST_NODETYPE, value interface{}) (JsonNode, error) {
	keyNode, err := NodeFactory(t, value)
	if err != nil {
		return nil, err
	}
	kvNode, err := NodeFactory(AST_KVPAIR, keyNode)
	if err != nil {
		return nil, err
	}
	i.astTrace.Push(kvNode)
	return kvNode, nil
}

func (i *JsonextAST) createValueNodeForKVPairs(owner *JsonKeyValuePairNode, t AST_NODETYPE, value interface{}) (JsonNode, error) {

	valueNode, err := NodeFactory(t, value)
	if err != nil {
		return nil, err
	}
	// the value of the kv is array|object, push it on top of the stack
	if t == AST_ARRAY || t == AST_OBJECT {
		i.astTrace.Push(valueNode)
	} else {
		// primivite value, finalize the k-v pair and append to the object node
		owner.Value = valueNode
		err = i.finlizeKVPair()
		if err != nil {
			return nil, err
		}
	}
	return owner, nil
}

// 2 reason to finalise, enclose of kv pair due to `,`, enclose of kv pair due to `}`
// {"1":2,"3":4}
func (i *JsonextAST) finlizeKVPair() error {
	kvElement, err := i.astTrace.Pop() // pop the kv, because it should be finalized to objet now.
	if err == util.ErrorEndOfStack {
		return ErrorASTStackEmpty
	}
	if kvElement.GetNodeType() != AST_KVPAIR {
		return ErrorASTUnexpectedElement
	}
	kvOwnerObj, err := i.astTrace.Peek()
	if err == util.ErrorEndOfStack {
		return ErrorASTStackEmpty
	}
	if kvOwnerObj.GetNodeType() != AST_OBJECT {
		return ErrorASTUnexpectedElement
	}
	el := kvOwnerObj.(*JsonObjectNode)
	el.Append(kvElement.(*JsonKeyValuePairNode))

	return nil
}

func (i *JsonextAST) EncloseLatestElements() (JsonNode, error) {

	itemToFinalize, err := i.astTrace.Pop()
	if err == util.ErrorEndOfStack {
		i.state = AST_STATE_FINISHED
		return nil, nil
	}
	return i.storeFinlizedItemToOwner(itemToFinalize)

}

func (i *JsonextAST) TopElementType() (AST_NODETYPE, error) {
	t, err := i.astTrace.Peek()
	if err != nil {
		return AST_NODE_UNDEFINED, err
	}
	return t.GetNodeType(), nil
}

func (i *JsonextAST) storeFinlizedItemToOwner(itemToFinalize JsonNode) (JsonNode, error) {
	nodeType := itemToFinalize.GetNodeType()
	switch nodeType {
	case AST_OBJECT: // item can only be value of kv or element of array
		fallthrough
	case AST_ARRAY:
		ownerElement, err := i.astTrace.Peek()
		if err == util.ErrorEndOfStack {
			i.state = AST_STATE_FINISHED
			return nil, nil // last element in the stack, no owner
		}
		switch ownerElement.GetNodeType() {
		case AST_ARRAY:
			el := ownerElement.(*JsonArrayNode)
			el.Append(itemToFinalize) // array case, put it in array
		case AST_KVPAIR: // kv case
			el := ownerElement.(*JsonKeyValuePairNode)
			el.Value = itemToFinalize
			err = i.finlizeKVPair()
			if err != nil {
				return nil, err
			}

		default:
			return nil, ErrorASTUnexpectedOwnerElement
		}
		return ownerElement, nil
	default:
		return nil, ErrorASTEncloseElementType
	}
}

func (i *JsonextAST) HasOpenElement() bool {
	return i.astTrace.Length() > 0
}

func (i *JsonextAST) HasComplete() bool {
	return i.state == AST_STATE_FINISHED
}
