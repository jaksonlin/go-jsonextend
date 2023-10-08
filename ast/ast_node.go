package ast

import (
	"encoding/base64"
	"fmt"

	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

type AST_NODETYPE byte

func (a AST_NODETYPE) Byte() byte {
	return byte(a)
}

const (
	AST_NODE_TYPE_BOUNDARY AST_NODETYPE = 200
	AST_ARRAY              AST_NODETYPE = 201
	AST_OBJECT             AST_NODETYPE = 202
	AST_KVPAIR             AST_NODETYPE = 203
	AST_VARIABLE           AST_NODETYPE = 204
	AST_STRING_VARIABLE    AST_NODETYPE = 205
	AST_STRING             AST_NODETYPE = 206
	AST_NUMBER             AST_NODETYPE = 207
	AST_BOOLEAN            AST_NODETYPE = 208
	AST_NULL               AST_NODETYPE = 209
	AST_NODE_UNDEFINED     AST_NODETYPE = 210
)

func ConvertTokenTypeToNodeType(t token.TokenType) AST_NODETYPE {
	switch t {
	case token.TOKEN_BOOLEAN:
		return AST_NUMBER
	case token.TOKEN_STRING:
		return AST_STRING
	case token.TOKEN_NUMBER:
		return AST_NUMBER
	case token.TOKEN_NULL:
		return AST_NULL
	case token.TOKEN_VARIABLE:
		return AST_VARIABLE
	case token.TOKEN_STRING_WITH_VARIABLE:
		return AST_STRING_VARIABLE
	default:
		return AST_NODE_UNDEFINED
	}
}

func NodeFactory(t AST_NODETYPE, value interface{}) (JsonNode, error) {

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
	String() string
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

func (node *JsonStringNode) String() string {
	return fmt.Sprintf("string node, value: %s\n", node.Value)
}

func (node *JsonStringNode) GetValue() string {
	if len(node.Value) == 2 {
		return "" // empty string with 2 double quotation marks only
	} else {
		return string(node.Value[1 : len(node.Value)-1])
	}
}

func (node *JsonStringNode) ToArrayNode() (*JsonArrayNode, error) {

	data, err := base64.StdEncoding.DecodeString(node.GetValue())
	if err != nil {
		return nil, err
	}
	rs := &JsonArrayNode{
		Value: make([]JsonNode, 0, len(data)),
	}
	for _, n := range data {
		v := uint8(n)
		rs.Value = append(rs.Value, &JsonNumberNode{
			Value: float64(v),
		})
	}
	return rs, nil

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

func (node *JsonNumberNode) String() string {
	return fmt.Sprintf("number node, value: %f\n", node.Value)
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

func (node *JsonBooleanNode) String() string {
	return fmt.Sprintf("boolean node, value: %t\n", node.Value)
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

func (node *JsonNullNode) String() string {
	return fmt.Sprintf("null node, value: %v\n", node.Value)
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

func (node *JsonArrayNode) String() string {
	return fmt.Sprintf("array node, length: %d\n", len(node.Value))
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

func (node *JsonKeyValuePairNode) String() string {
	return fmt.Sprintf("key value pair node, key: [%s], value: [%s]\n", node.Key.String(), node.Value.String())
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

func (node *JsonObjectNode) String() string {
	return fmt.Sprintf("object node, length: %d\n", len(node.Value))
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

func (node *JsonExtendedVariableNode) String() string {
	return fmt.Sprintf("variable node, value: %s\n", node.Value)
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

func (node *JsonExtendedStringWIthVariableNode) String() string {
	return fmt.Sprintf("string with variable node, value: %s\n", node.Value)
}
