package ast

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
