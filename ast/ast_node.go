package ast

import "github.com/jaksonlin/go-jsonextend/util"

type AST_NODETYPE byte

const (
	AST_NODE_TYPE_BOUNDARY = 200
	AST_ARRAY              = 201
	AST_OBJECT             = 202
	AST_KVPAIR             = 203
	AST_VARIABLE           = 204
	AST_STRING_VARIABLE    = 205
	AST_STRING             = 206
	AST_NUMBER             = 207
	AST_BOOLEAN            = 208
	AST_NULL               = 209
	AST_NODE_UNDEFINED     = 210
)

func nodeFactory(t AST_NODETYPE, value interface{}) (JsonNode, error) {

	switch t {
	case AST_ARRAY:
		return &JsonArrayNode{
			Value: make([]JsonNode, 0),
		}, nil
	case AST_OBJECT:
		return &JsonObjectNode{
			Value: make([]*JsonKeyValuePairNode, 0),
		}, nil
	case AST_KVPAIR:
		node, ok := value.(JsonStringValueNode)
		if !ok {
			return nil, ErrorASTKeyValuePairNotStringAsKey
		}
		return &JsonKeyValuePairNode{
			Key: node,
		}, nil
	case AST_STRING:
		return &JsonStringNode{
			Value: value.([]byte),
		}, nil
	case AST_NUMBER:
		return &JsonNumberNode{
			Value: value.(float64),
		}, nil
	case AST_BOOLEAN:
		return &JsonBooleanNode{
			Value: value.(bool),
		}, nil
	case AST_NULL:
		return &JsonNullNode{
			Value: nil,
		}, nil
	case AST_VARIABLE:
		node := &JsonExtendedVariableNode{
			Value: value.([]byte),
		}
		node.extractVariable()
		return node, nil
	case AST_STRING_VARIABLE:
		node := &JsonExtendedStringWIthVariableNode{
			JsonStringNode: JsonStringNode{
				Value: value.([]byte),
			},
		}
		node.extractVariables()
		return node, nil
	default:
		return nil, ErrorASTIncorrectNodeType
	}
}

type JsonVisitor interface {
	VisitStringNode(node *JsonStringNode) error
	VisitNumberNode(node *JsonNumberNode) error
	VisitBooleanNode(node *JsonBooleanNode) error
	VisitNullNode(node *JsonNullNode) error
	VisitArrayNode(node *JsonArrayNode) error
	VisitKeyValuePairNode(node *JsonKeyValuePairNode) error
	VisitObjectNode(node *JsonObjectNode) error
	VisitVariableNode(node *JsonExtendedVariableNode) error
	VisitStringWithVariableNode(node *JsonExtendedStringWIthVariableNode) error
}

type JsonNode interface {
	GetNodeType() AST_NODETYPE
	Visit(visitor JsonVisitor) error
}

type JsonCollectionNode interface {
	JsonNode
	Length() int
}

type JsonStringValueNode interface {
	JsonNode
	GetValue() string
}

type JsonStringNode struct {
	Value []byte
}

var _ JsonNode = &JsonStringNode{}

func (node *JsonStringNode) GetNodeType() AST_NODETYPE {
	return AST_STRING
}

func (node *JsonStringNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitStringNode(node)
}

func (node *JsonStringNode) GetValue() string {
	if len(node.Value) == 2 {
		return "" // empty string with 2 double quotation marks only
	} else {
		return string(node.Value[1 : len(node.Value)-1])
	}
}

type JsonNumberNode struct {
	Value float64
}

var _ JsonNode = &JsonNumberNode{}

func (node *JsonNumberNode) GetNodeType() AST_NODETYPE {
	return AST_NUMBER
}

func (node *JsonNumberNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitNumberNode(node)
}

type JsonBooleanNode struct {
	Value bool
}

var _ JsonNode = &JsonBooleanNode{}

func (node *JsonBooleanNode) GetNodeType() AST_NODETYPE {
	return AST_BOOLEAN
}

func (node *JsonBooleanNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitBooleanNode(node)
}

type JsonNullNode struct {
	Value interface{}
}

var _ JsonNode = &JsonNullNode{}

func (node *JsonNullNode) GetNodeType() AST_NODETYPE {
	return AST_NULL
}

func (node *JsonNullNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitNullNode(node)
}

type JsonArrayNode struct {
	Value []JsonNode
}

var _ JsonNode = &JsonArrayNode{}

func (node *JsonArrayNode) GetNodeType() AST_NODETYPE {
	return AST_ARRAY
}

func (node *JsonArrayNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitArrayNode(node)
}

func (node *JsonArrayNode) Append(n JsonNode) {
	node.Value = append(node.Value, n)
}

func (node *JsonArrayNode) Length() int {
	return len(node.Value)
}

type JsonKeyValuePairNode struct {
	Key   JsonStringValueNode
	Value JsonNode
}

var _ JsonNode = &JsonKeyValuePairNode{}

func (node *JsonKeyValuePairNode) GetNodeType() AST_NODETYPE {
	return AST_KVPAIR
}

func (node *JsonKeyValuePairNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitKeyValuePairNode(node)
}

func (node *JsonKeyValuePairNode) IsFilled() bool {
	return node.Value != nil
}

type JsonObjectNode struct {
	Value []*JsonKeyValuePairNode
}

var _ JsonNode = &JsonObjectNode{}

func (node *JsonObjectNode) GetNodeType() AST_NODETYPE {
	return AST_OBJECT
}

func (node *JsonObjectNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitObjectNode(node)
}

func (node *JsonObjectNode) Append(kvNode *JsonKeyValuePairNode) {
	node.Value = append(node.Value, kvNode)
}

func (node *JsonObjectNode) Length() int {
	return len(node.Value)
}

type JsonExtendedVariableNode struct {
	Value    []byte
	Variable string
}

var _ JsonNode = &JsonExtendedVariableNode{}

func (node *JsonExtendedVariableNode) GetNodeType() AST_NODETYPE {
	return AST_VARIABLE
}

func (node *JsonExtendedVariableNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitVariableNode(node)
}

func (node *JsonExtendedVariableNode) extractVariable() {
	rs := util.RegStringWithVariable.FindSubmatch(node.Value)
	node.Variable = string(rs[1])
}

type JsonExtendedStringWIthVariableNode struct {
	JsonStringNode
	Variables map[string][]byte
}

var _ JsonNode = &JsonExtendedStringWIthVariableNode{}

func (node *JsonExtendedStringWIthVariableNode) GetNodeType() AST_NODETYPE {
	return AST_STRING_VARIABLE
}

func (node *JsonExtendedStringWIthVariableNode) Visit(visitor JsonVisitor) error {
	return visitor.VisitStringWithVariableNode(node)
}

func (node *JsonExtendedStringWIthVariableNode) extractVariables() {
	rs := util.RegStringWithVariable.FindAllSubmatch(node.Value, -1)
	if len(rs) > 0 {
		node.Variables = make(map[string][]byte)
	}
	for _, item := range rs {
		node.Variables[string(item[1])] = item[0]
	}
}

func (node *JsonExtendedStringWIthVariableNode) GetValue() string {
	return node.JsonStringNode.GetValue()
}
