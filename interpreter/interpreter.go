package interpreter

import (
	"bytes"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/tokenizer"
	"github.com/jaksonlin/go-jsonextend/util"
)

type PrettyPrintVisitor struct {
	sb           *bytes.Buffer
	indentString string
	indent       int
	variables    map[string]interface{}
	stackNode    *util.Stack[ast.JsonNode]
	stackFormat  *util.Stack[byte]
	marshaler    ast.MarshalerFunc
}

var _ ast.NodeVisitor = &PrettyPrintVisitor{}

var colonFormat = []byte{' ', ':', ' '}

func NewPPInterpreter(variables map[string]interface{}, marshaler ast.MarshalerFunc) *PrettyPrintVisitor {

	return &PrettyPrintVisitor{
		sb:           bytes.NewBuffer(make([]byte, 0)),
		indentString: strings.Repeat(" ", 4),
		indent:       0,
		variables:    variables,
		stackNode:    util.NewStack[ast.JsonNode](),
		stackFormat:  util.NewStack[byte](),
		marshaler:    marshaler,
	}
}

func (s *PrettyPrintVisitor) getSymbolLength() int {
	return s.stackFormat.Length()
}

func (s *PrettyPrintVisitor) WriteSymbol() error {
	symbol, e := s.stackFormat.Pop()
	if e != nil {
		return e
	}
	// the caller is the last element in an object/array
	if symbol == ']' || symbol == '}' {
		//write return line
		s.sb.WriteByte('\n')
		// reduce indent
		s.indent--
		s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
		// put the closing symbol
		s.sb.WriteByte(symbol)
		// if we are in the middle of any collection,  break after writing the `comma`
		// else we will consume any number of closing symbols in between!
		for symbol != ',' {
			symbol, e = s.stackFormat.Pop()
			if e != nil {
				return e
			}
			if symbol != ',' {
				s.indent--
				s.sb.WriteByte('\n')
				s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
			}
			s.sb.WriteByte(symbol)
		}
		s.sb.WriteByte('\n')
		s.sb.WriteString(strings.Repeat(s.indentString, s.indent))

	} else {

		if symbol == ':' {
			s.sb.Write(colonFormat)
		} else {
			s.sb.WriteByte(symbol)
			s.sb.WriteByte('\n')
			s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
		}
	}
	return nil
}

func (s *PrettyPrintVisitor) VisitStringNode(node *ast.JsonStringNode) error {
	s.sb.Write(node.Value)
	return s.WriteSymbol()
}

func (s *PrettyPrintVisitor) VisitNumberNode(node *ast.JsonNumberNode) error {
	f64, err := util.ConvertInterfaceNumberToFloat64(node.Value)
	if err != nil {
		return err
	}
	s.sb.WriteString(strconv.FormatFloat(f64, 'f', -1, 64))
	return s.WriteSymbol()
}

func (s *PrettyPrintVisitor) VisitBooleanNode(node *ast.JsonBooleanNode) error {
	s.sb.WriteString(strconv.FormatBool(node.Value))
	return s.WriteSymbol()
}

func (s *PrettyPrintVisitor) VisitNullNode(node *ast.JsonNullNode) error {
	s.sb.WriteString("null")
	return s.WriteSymbol()
}

func (s *PrettyPrintVisitor) VisitStringWithVariableNode(node *ast.JsonExtendedStringWIthVariableNode) error {
	if s.marshaler == nil {
		s.sb.Write(node.Value)
		return s.WriteSymbol()
	}
	var result []byte = make([]byte, len(node.Value))
	copy(result, node.Value)
	for varName, varDollarName := range node.Variables {
		varVal, ok := s.variables[varName]
		if ok {
			content, err := s.marshalAndStripQuotes(varVal)
			if err != nil {
				return ErrorInterpretVariable
			}
			result = bytes.ReplaceAll(result, varDollarName, content)
		}
	}
	// the varaible value is of string type, remove the leading and trailing double quotation mark

	s.sb.Write(result)
	return s.WriteSymbol()
}

func (s *PrettyPrintVisitor) marshalAndStripQuotes(varVal interface{}) ([]byte, error) {
	content, err := s.marshalVariableValue(varVal)
	if err != nil {
		return nil, err
	}

	if content[0] == '"' {
		content = content[1 : len(content)-1]
	}
	return content, nil
}

func (s *PrettyPrintVisitor) marshalVariableValue(varVal interface{}) ([]byte, error) {
	var content []byte
	if util.IsPrimitiveType(reflect.ValueOf(varVal)) {
		c, err := util.EncodePrimitiveValue(varVal)
		if err != nil {
			return nil, err
		}
		content = c
	} else {
		c, err := s.marshaler(varVal)
		if err != nil {
			return nil, err
		}
		content = c
	}

	return content, nil
}

func (s *PrettyPrintVisitor) VisitVariableNode(node *ast.JsonExtendedVariableNode) error {
	if s.marshaler == nil {
		s.sb.Write(node.Value)
		return s.WriteSymbol()
	}
	varVal, ok := s.variables[node.Variable] // allow partial rendered
	if !ok {
		s.sb.Write(node.Value)
		return s.WriteSymbol()
	} else {
		content, err := s.marshalVariableValue(varVal)
		if err != nil {
			return ErrorInterpretVariable
		}
		s.sb.Write(content)
	}

	return s.WriteSymbol()
}

func (s *PrettyPrintVisitor) VisitArrayNode(node *ast.JsonArrayNode) error {
	s.sb.WriteString("[\n")
	s.indent++
	s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
	if len(node.Value) == 0 {
		s.stackFormat.Push(']')
		return s.WriteSymbol()
	}
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

func (s *PrettyPrintVisitor) VisitKeyValuePairNode(node *ast.JsonKeyValuePairNode) error {
	// stack, first in last out, value go first ;-)
	s.stackNode.Push(node.Value)
	s.stackNode.Push(node.Key)

	return nil
}

func (s *PrettyPrintVisitor) GetOutput() []byte {
	return s.sb.Bytes()
}

func (s *PrettyPrintVisitor) VisitObjectNode(node *ast.JsonObjectNode) error {
	s.sb.WriteString("{\n")
	s.indent++
	s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
	if len(node.Value) == 0 {
		s.stackFormat.Push('}')
		return s.WriteSymbol()
	}
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

func PrettyInterpret(node ast.JsonNode, variables map[string]interface{}, marshaler ast.MarshalerFunc) ([]byte, error) {
	// deep first traverse the AST

	visitor := NewPPInterpreter(variables, marshaler)
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

func ParseJsonExtendDocument(reader io.Reader, variables map[string]interface{}) ([]byte, error) {
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)
	err := sm.ProcessData()
	if err != nil {
		return nil, err
	}
	if sm.GetASTBuilder().HasOpenElements() {
		return nil, ErrorInvalidJson
	}
	ast := sm.GetAST()
	return PrettyInterpret(ast, variables, Marshal)
}
