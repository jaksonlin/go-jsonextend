package ast

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/jaksonlin/go-jsonextend/util"
)

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
			OriginValue: value,
			Value:       value.(float64),
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

type astNodeBase struct {
	visited     bool
	nodePlugins nodePlugins
	meta        map[string]interface{}
}

func (node *astNodeBase) SetMeta(key string, value interface{}) {
	if node.meta == nil {
		node.meta = make(map[string]interface{})
	}
	node.meta[key] = value
}

func (node *astNodeBase) GetMeta(key string) interface{} {
	if node.meta == nil {
		return nil
	}
	val, ok := node.meta[key]
	if !ok {
		return nil
	}
	return val
}

func (node *astNodeBase) SetVisited() {
	node.visited = true
}

func (node *astNodeBase) IsVisited() bool {
	return node.visited
}

func (node *astNodeBase) UnsetVisited() {
	node.visited = false
}

type JsonStringNode struct {
	astNodeBase
	Value       []byte
	stringValue string
}

var _ JsonNode = &JsonStringNode{}

func (node *JsonStringNode) GetNodeType() AST_NODETYPE {
	return AST_STRING
}

func (node *JsonStringNode) Visit(visitor JsonVisitor) error {
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitStringNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
}

func (node *JsonStringNode) String() string {
	return fmt.Sprintf("string node, value: %s\n", node.Value)
}

func (node *JsonStringNode) GetValue() (string, error) {
	if node.stringValue != "" {
		return node.stringValue, nil
	}
	val, err := strconv.Unquote(string(node.Value))
	if err != nil {
		return "", err
	}
	node.stringValue = val
	return val, nil
}

func (node *JsonStringNode) ToArrayNode() (*JsonArrayNode, error) {
	bytesData, err := node.GetValue()
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(bytesData)
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
	astNodeBase
	OriginValue interface{}
	Value       float64
}

var _ JsonNode = &JsonNumberNode{}

func (node *JsonNumberNode) GetNodeType() AST_NODETYPE {
	return AST_NUMBER
}

func (node *JsonNumberNode) Visit(visitor JsonVisitor) error {
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitNumberNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
}

func (node *JsonNumberNode) String() string {
	return fmt.Sprintf("number node, value: %f\n", node.Value)
}

type JsonBooleanNode struct {
	astNodeBase
	Value bool
}

var _ JsonNode = &JsonBooleanNode{}

func (node *JsonBooleanNode) GetNodeType() AST_NODETYPE {
	return AST_BOOLEAN
}

func (node *JsonBooleanNode) Visit(visitor JsonVisitor) error {
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitBooleanNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
}

func (node *JsonBooleanNode) String() string {
	return fmt.Sprintf("boolean node, value: %t\n", node.Value)
}

type JsonNullNode struct {
	astNodeBase
	Value interface{}
}

var _ JsonNode = &JsonNullNode{}

func (node *JsonNullNode) GetNodeType() AST_NODETYPE {
	return AST_NULL
}

func (node *JsonNullNode) Visit(visitor JsonVisitor) error {
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitNullNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
}

func (node *JsonNullNode) String() string {
	return fmt.Sprintf("null node, value: %v\n", node.Value)
}

type JsonArrayNode struct {
	astNodeBase
	Value []JsonNode
}

var _ JsonNode = &JsonArrayNode{}

func (node *JsonArrayNode) GetNodeType() AST_NODETYPE {
	return AST_ARRAY
}

func (node *JsonArrayNode) Visit(visitor JsonVisitor) error {
	// allow user to shutdown the visit of the node
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	// it is possible that the plugin set the node as visited, so we need to check again
	if !node.IsVisited() {
		err = visitor.VisitArrayNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
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
	astNodeBase
	Key   JsonStringValueNode
	Value JsonNode
}

var _ JsonNode = &JsonKeyValuePairNode{}

func (node *JsonKeyValuePairNode) GetNodeType() AST_NODETYPE {
	return AST_KVPAIR
}

func (node *JsonKeyValuePairNode) UnsetVisited() {
	node.visited = false
	node.Key.UnsetVisited()
	node.Value.UnsetVisited()
}

func (node *JsonKeyValuePairNode) Visit(visitor JsonVisitor) error {
	// allow user to shutdown the visit of the node
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitKeyValuePairNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
}

func (node *JsonKeyValuePairNode) IsFilled() bool {
	return node.Value != nil
}

func (node *JsonKeyValuePairNode) String() string {
	return fmt.Sprintf("key value pair node, key: [%s], value: [%s]\n", node.Key.String(), node.Value.String())
}

type JsonObjectNode struct {
	astNodeBase
	Value []*JsonKeyValuePairNode
}

var _ JsonNode = &JsonObjectNode{}

func (node *JsonObjectNode) GetNodeType() AST_NODETYPE {
	return AST_OBJECT
}

func (node *JsonObjectNode) Visit(visitor JsonVisitor) error {
	// allow user to shutdown the visit of the node
	if node.IsVisited() {
		return nil
	}

	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitObjectNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
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
	astNodeBase
	Value    []byte
	Variable string
}

var _ JsonNode = &JsonExtendedVariableNode{}

func (node *JsonExtendedVariableNode) GetNodeType() AST_NODETYPE {
	return AST_VARIABLE
}

func (node *JsonExtendedVariableNode) Visit(visitor JsonVisitor) error {
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitVariableNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
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
	if node.IsVisited() {
		return nil
	}
	err := node.nodePlugins.PreVisitPlugin(visitor, node)
	if err != nil {
		return err
	}
	if !node.IsVisited() {
		err = visitor.VisitStringWithVariableNode(node)
		if err != nil {
			return err
		}
		node.SetVisited()
	}
	return node.nodePlugins.PostVisitPlugin(visitor, node)
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

func (node *JsonExtendedStringWIthVariableNode) GetValue() (string, error) {
	return node.JsonStringNode.GetValue()
}

func (node *JsonExtendedStringWIthVariableNode) String() string {
	return fmt.Sprintf("string with variable node, value: %s\n", node.Value)
}
