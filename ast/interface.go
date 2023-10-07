package ast

type MarshalerFunc func(v interface{}) ([]byte, error)

type NodeVisitor interface {
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
