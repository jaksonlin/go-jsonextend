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

type jzoneAST struct {
	ast      JsonNode
	astTrace *util.Stack[JsonNode]
	state    astState
}

func newJzoneAST() *jzoneAST {
	return &jzoneAST{
		ast:      nil,
		astTrace: util.NewStack[JsonNode](),
		state:    AST_STATE_INIT,
	}
}

func (i *jzoneAST) GetAST() JsonNode {
	return i.ast
}

func (i *jzoneAST) createRootNode(t AST_NODETYPE, value interface{}) error {
	n, err := nodeFactory(t, value)
	if err != nil {
		return err
	}
	i.ast = n

	if t == AST_ARRAY || t == AST_OBJECT {
		i.astTrace.Push(n)
		i.state = AST_STATE_INTERMIDIEATE
	} else {
		i.state = AST_STATE_FINISHED
	}
	return nil
}

func (i *jzoneAST) CreateNewASTNode(t AST_NODETYPE, value interface{}) error {
	if i.state == AST_STATE_FINISHED {
		return ErrorASTComplete
	}
	if i.ast == nil {
		return i.createRootNode(t, value)
	}
	latest, err := i.astTrace.Peek()
	if err != nil {
		return err
	}

	switch realNode := latest.(type) {
	case *JsonArrayNode:
		return i.createNewNodeForArrayObject(realNode, t, value)
	case *JsonObjectNode:
		return i.createNewNodeForObject(realNode, t, value)
	case *JsonKeyValuePairNode:
		return i.createValueNodeForKVPairs(realNode, t, value)
	default:
		return ErrorASTUnexpectedElement
	}

}

func (i *jzoneAST) createNewNodeForArrayObject(owner *JsonArrayNode, t AST_NODETYPE, value interface{}) error {
	n, err := nodeFactory(t, value)
	if err != nil {
		return err
	}
	if t == AST_ARRAY || t == AST_OBJECT {
		// none-primivtive value, create a new item to hold the furture values
		i.astTrace.Push(n)
	} else {
		// primitive value push to the owner on top of the stack
		owner.Append(n)
	}
	return nil
}

func (i *jzoneAST) createNewNodeForObject(owner *JsonObjectNode, t AST_NODETYPE, value interface{}) error {
	keyNode, err := nodeFactory(t, value)
	if err != nil {
		return err
	}
	n, err := nodeFactory(AST_KVPAIR, keyNode)
	if err != nil {
		return err
	}
	i.astTrace.Push(n)
	return nil
}

func (i *jzoneAST) createValueNodeForKVPairs(owner *JsonKeyValuePairNode, t AST_NODETYPE, value interface{}) error {

	n, err := nodeFactory(t, value)
	if err != nil {
		return err
	}
	// the value of the kv is array|object, push it on top of the stack
	if t == AST_ARRAY || t == AST_OBJECT {
		i.astTrace.Push(n)
	} else {
		// primivite value, finalize the k-v pair and append to the object node
		owner.Value = n
		err = i.finlizeKVPair()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *jzoneAST) finlizeKVPair() error {
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
	kvOwnerObj.(*JsonObjectNode).Append(kvElement.(*JsonKeyValuePairNode))
	return nil
}

func (i *jzoneAST) EncloseLatestElements() error {

	itemToFinalize, err := i.astTrace.Pop()
	if err == util.ErrorEndOfStack {
		i.state = AST_STATE_FINISHED
		return nil
	}
	err = i.storeFinlizedItemToOwner(itemToFinalize)
	if err != nil {
		return err
	}
	return nil

}

func (i *jzoneAST) TopElementType() (AST_NODETYPE, error) {
	t, err := i.astTrace.Peek()
	if err != nil {
		return AST_NODE_UNDEFINED, err
	}
	return t.GetNodeType(), nil
}

func (i *jzoneAST) storeFinlizedItemToOwner(itemToFinalize JsonNode) error {
	nodeType := itemToFinalize.GetNodeType()
	switch nodeType {
	case AST_OBJECT: // item can only be value of kv or element of array
		fallthrough
	case AST_ARRAY:
		ownerElement, err := i.astTrace.Peek()
		if err == util.ErrorEndOfStack {
			i.state = AST_STATE_FINISHED
			return nil // last element in the stack, no owner
		}
		switch ownerElement.GetNodeType() {
		case AST_ARRAY:
			ownerElement.(*JsonArrayNode).Append(itemToFinalize) // array case, put it in array
		case AST_KVPAIR: // kv case
			ownerElement.(*JsonKeyValuePairNode).Value = itemToFinalize
			err = i.finlizeKVPair()
			if err != nil {
				return err
			}
		default:
			return ErrorASTUnexpectedOwnerElement
		}
	default:
		return ErrorASTEncloseElementType
	}
	return nil
}

func (i *jzoneAST) HasOpenElement() bool {
	return i.astTrace.Length() > 0
}

func (i *jzoneAST) HasComplete() bool {
	return i.state == AST_STATE_FINISHED
}
