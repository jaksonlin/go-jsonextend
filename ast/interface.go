package ast

type MarshalerFunc func(v interface{}) ([]byte, error)
type UnmarshalerFunc func(v []byte, out interface{}) error

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
	SetVisited()
	IsVisited() bool
	UnsetVisited()
	String() string
	SetMeta(key string, value interface{})
	GetMeta(key string) interface{}
}

type ASTNodePlugin interface {
	//
	PreVisitPlugin(visitor JsonVisitor, pluginHolder JsonNode) error
	PostVisitPlugin(visitor JsonVisitor, pluginHolder JsonNode) error
	PluginName() string
}

type PluggableJsonNode interface {
	JsonNode
	AddPlugin(p ASTNodePlugin)
	RemovePlugin(name string)
}

type JsonCollectionNode interface {
	JsonNode
	Length() int
}

type JsonStringValueNode interface {
	JsonNode
	GetValue() (string, error)
}
