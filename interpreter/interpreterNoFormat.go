package interpreter

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

type standardVisitor struct {
	sb          *bytes.Buffer
	variables   map[string]interface{}
	stackNode   *util.Stack[ast.JsonNode]
	stackFormat *util.Stack[byte]
}

var _ JsonVisitor = &standardVisitor{}

func NewInterpreter(variables map[string]interface{}) *standardVisitor {

	return &standardVisitor{
		sb:          bytes.NewBuffer(make([]byte, 0)),
		variables:   variables,
		stackNode:   util.NewStack[ast.JsonNode](),
		stackFormat: util.NewStack[byte](),
	}
}

func (s *standardVisitor) getSymbolLength() int {
	return s.stackFormat.Length()
}

func (s *standardVisitor) WriteSymbol() error {
	symbol, e := s.stackFormat.Pop()
	if e != nil {
		return e
	}
	s.sb.WriteByte(symbol)
	// the caller is the last element in an object/array
	if symbol == ']' || symbol == '}' {
		// if we are in the middle of any collection, write one more `comma` after the closing symbol
		// stack top [`]` , `}`, ... `,` , `]`], the first `]` has been poped above, we need to write one more `,`
		for symbol != ',' {
			symbol, e = s.stackFormat.Pop()
			if e != nil {
				return e
			}
			s.sb.WriteByte(symbol)
		}

	}
	return nil
}

func (s *standardVisitor) VisitStringNode(node *ast.JsonStringNode) error {
	s.sb.Write(node.Value)
	return s.WriteSymbol()
}

func (s *standardVisitor) VisitNumberNode(node *ast.JsonNumberNode) error {
	s.sb.WriteString(strconv.FormatFloat(node.Value, 'f', -1, 64))
	return s.WriteSymbol()
}

func (s *standardVisitor) VisitBooleanNode(node *ast.JsonBooleanNode) error {
	if node.Value {
		s.sb.Write(token.TrueBytes)
	} else {
		s.sb.Write(token.FalseBytes)
	}
	return s.WriteSymbol()
}

func (s *standardVisitor) VisitNullNode(node *ast.JsonNullNode) error {
	s.sb.Write(token.NullBytes)
	return s.WriteSymbol()
}

func (s *standardVisitor) VisitStringWithVariableNode(node *ast.JsonExtendedStringWIthVariableNode) error {
	var result []byte = make([]byte, len(node.Value))
	copy(result, node.Value)
	for varName, varDollarName := range node.Variables {
		varVal, ok := s.variables[varName]
		if ok {
			content, err := json.Marshal(varVal)
			if err != nil {
				return ErrorInterpretVariable
			}
			// remove the json string's leading and trailing double quotation mark, otherwise you will get something ""value"", which is invalid string
			if content[0] == '"' {
				content = content[1 : len(content)-1]
			}
			result = bytes.ReplaceAll(result, varDollarName, content)
		}
	}
	// the varaible value is of string type, remove the leading and trailing double quotation mark

	s.sb.Write(result)
	return s.WriteSymbol()
}

func (s *standardVisitor) VisitVariableNode(node *ast.JsonExtendedVariableNode) error {

	varVal, ok := s.variables[node.Variable] // allow partial rendered
	if !ok {
		s.sb.Write(node.Value)
		return s.WriteSymbol()
	} else {
		content, err := json.Marshal(varVal)
		if err != nil {
			return ErrorInterpretVariable
		}
		if content[0] == '"' {
			content = content[1 : len(content)-1]
		}
		s.sb.Write(content)
	}

	return s.WriteSymbol()
}

func (s *standardVisitor) VisitArrayNode(node *ast.JsonArrayNode) error {
	s.sb.WriteByte('[')

	for i := len(node.Value) - 1; i >= 0; i-- {
		s.stackNode.Push(node.Value[i])
		if i == len(node.Value)-1 {
			s.stackFormat.Push(']')
		} else {
			s.stackFormat.Push(',')
		}
	}
	return nil
}

func (s *standardVisitor) VisitKeyValuePairNode(node *ast.JsonKeyValuePairNode) error {
	// stack, first in last out, value go first ;-)
	s.stackNode.Push(node.Value)
	s.stackNode.Push(node.Key)

	return nil
}

func (s *standardVisitor) GetOutput() []byte {
	return s.sb.Bytes()
}

func (s *standardVisitor) VisitObjectNode(node *ast.JsonObjectNode) error {
	s.sb.WriteByte('{')
	for i := len(node.Value) - 1; i >= 0; i-- {
		s.stackNode.Push(node.Value[i])
		if i == len(node.Value)-1 { // stack, first in last out
			s.stackFormat.Push('}')
			s.stackFormat.Push(':')
		} else {
			s.stackFormat.Push(',')
			s.stackFormat.Push(':')
		}
	}
	return nil
}

func InterpretAST(node ast.JsonNode, variables map[string]interface{}) ([]byte, error) {
	// deep first traverse the AST

	visitor := NewInterpreter(variables)
	visitor.stackNode.Push(node)

	for {
		node, err := visitor.stackNode.Pop()
		if err != nil {
			break
		}
		err = node.Visit(visitor)
		if err != nil {
			if err != util.ErrorEndOfStack {
				return nil, err
			} else {
				break
			}
		}

	}

	if visitor.getSymbolLength() > 0 {
		return nil, ErrorInterpreSymbolFailure
	}
	rs := visitor.GetOutput()
	return rs, nil
}
