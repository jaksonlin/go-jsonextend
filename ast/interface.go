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
	AddPlugin(p ASTNodePlugin)
	RemovePlugin(name string)
	PrependPlugin(p ASTNodePlugin)
}

type ASTNodePluginFunc func(visitor JsonVisitor, pluginHolder JsonNode) error

type ASTNodePlugin interface {
	PreVisitPlugin(visitor JsonVisitor, node JsonNode) error
	PostVisitPlugin(visitor JsonVisitor, node JsonNode) error
	PluginName() string
}

type astNodePluginImpl struct {
	preVisitFunc  ASTNodePluginFunc
	postVisitFunc ASTNodePluginFunc
	pluginName    string
}

var _ ASTNodePlugin = (*astNodePluginImpl)(nil)

func NewASTNodePlugin(name string, preVisitFunc ASTNodePluginFunc, postVisitFunc ASTNodePluginFunc) ASTNodePlugin {
	return &astNodePluginImpl{
		preVisitFunc:  preVisitFunc,
		postVisitFunc: postVisitFunc,
		pluginName:    name,
	}
}

func (p *astNodePluginImpl) PreVisitPlugin(visitor JsonVisitor, node JsonNode) error {
	if p.preVisitFunc == nil {
		return nil
	}
	return p.preVisitFunc(visitor, node)
}

func (p *astNodePluginImpl) PostVisitPlugin(visitor JsonVisitor, node JsonNode) error {
	if p.postVisitFunc == nil {
		return nil
	}
	return p.postVisitFunc(visitor, node)
}

func (p *astNodePluginImpl) PluginName() string {
	return p.pluginName
}

type JsonCollectionNode interface {
	JsonNode
	Length() int
	SetChildVisited()
	ResetVisited()
	IsChildVisited() bool
}

type JsonStringValueNode interface {
	JsonNode
	GetValue() (string, error)
}
