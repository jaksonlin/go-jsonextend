package interpreter

// import (
// 	"bytes"
// 	"encoding/json"
// 	"strconv"
// 	"strings"

// 	"github.com/jaksonlin/go-jsonextend/ast"
// 	"github.com/jaksonlin/go-jsonextend/util"
// )

// type goStructdVisitor struct {
// 	out         interface{}
// 	variables   map[string]interface{}
// 	stackNode   *util.Stack[ast.JsonNode]
// 	stackFormat *util.Stack[byte]
// 	stackOut    *util.Stack[interface{}]
// }

// var _ JsonVisitor = &goStructdVisitor{}

// func NewGoStructInterpreter(variables map[string]interface{}, out interface{}) *goStructdVisitor {
// 	return &goStructdVisitor{
// 		out:         out,
// 		variables:   variables,
// 		stackNode:   util.NewStack[ast.JsonNode](),
// 		stackFormat: util.NewStack[byte](),
// 	}
// }

// func (s *goStructdVisitor) getSymbolLength() int {
// 	return s.stackFormat.Length()
// }

// func (s *goStructdVisitor) WriteSymbol() error {
// 	symbol, e := s.stackFormat.Pop()
// 	if e != nil {
// 		return e
// 	}
// 	// the caller is the last element in an object/array
// 	if symbol == ']' || symbol == '}' {
// 		//write return line
// 		s.sb.WriteByte('\n')
// 		// reduce indent
// 		s.indent--
// 		s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
// 		// put the closing symbol
// 		s.sb.WriteByte(symbol)
// 		// if we are in the middle of any collection,  break after writing the `comma`
// 		// else we will consume the closing symbols
// 		for symbol != ',' {
// 			symbol, e = s.stackFormat.Pop()
// 			if e != nil {
// 				return e
// 			}
// 			if symbol != ',' {
// 				s.indent--
// 				s.sb.WriteByte('\n')
// 				s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
// 			}
// 			s.sb.WriteByte(symbol)
// 		}
// 		s.sb.WriteByte('\n')
// 		s.sb.WriteString(strings.Repeat(s.indentString, s.indent))

// 	} else {

// 		if symbol == ':' {
// 			s.sb.Write(colonFormat)
// 		} else {
// 			s.sb.WriteByte(symbol)
// 			s.sb.WriteByte('\n')
// 			s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
// 		}
// 	}
// 	return nil
// }

// func (s *goStructdVisitor) VisitStringNode(node *ast.JsonStringNode) error {
// 	s.sb.Write(node.Value)
// 	return s.WriteSymbol()
// }

// func (s *goStructdVisitor) VisitNumberNode(node *ast.JsonNumberNode) error {
// 	s.sb.WriteString(strconv.FormatFloat(node.Value, 'f', -1, 64))
// 	return s.WriteSymbol()
// }

// func (s *goStructdVisitor) VisitBooleanNode(node *ast.JsonBooleanNode) error {
// 	s.sb.WriteString(strconv.FormatBool(node.Value))
// 	return s.WriteSymbol()
// }

// func (s *goStructdVisitor) VisitNullNode(node *ast.JsonNullNode) error {
// 	s.sb.WriteString("null")
// 	return s.WriteSymbol()
// }

// func (s *goStructdVisitor) VisitStringWithVariableNode(node *ast.JsonExtendedStringWIthVariableNode) error {
// 	var result []byte = make([]byte, len(node.Value))
// 	copy(result, node.Value)
// 	for varName, varDollarName := range node.Variables {
// 		varVal, ok := s.variables[varName]
// 		if ok {
// 			content, err := json.Marshal(varVal)
// 			if err != nil {
// 				return ErrorInterpretVariable
// 			}
// 			// remove the json string's leading and trailing double quotation mark, otherwise you will get something ""value"", which is invalid string
// 			if content[0] == '"' {
// 				content = content[1 : len(content)-1]
// 			}
// 			result = bytes.ReplaceAll(result, varDollarName, content)
// 		}
// 	}
// 	// the varaible value is of string type, remove the leading and trailing double quotation mark

// 	s.sb.Write(result)
// 	return s.WriteSymbol()
// }

// func (s *goStructdVisitor) VisitVariableNode(node *ast.JsonExtendedVariableNode) error {

// 	varVal, ok := s.variables[node.Variable] // allow partial rendered
// 	if !ok {
// 		s.sb.Write(node.Value)
// 		return s.WriteSymbol()
// 	} else {
// 		content, err := json.Marshal(varVal)
// 		if err != nil {
// 			return ErrorInterpretVariable
// 		}
// 		if content[0] == '"' {
// 			content = content[1 : len(content)-1]
// 		}
// 		s.sb.Write(content)
// 	}

// 	return s.WriteSymbol()
// }

// func (s *goStructdVisitor) Visit(node ast.JsonNode) error {

// 	switch node.GetNodeType() {
// 	case ast.AST_ARRAY:
// 		return s.VisitArrayNode(node.(*ast.JsonArrayNode))
// 	case ast.AST_OBJECT:
// 		return s.VisitObjectNode(node.(*ast.JsonObjectNode))
// 	case ast.AST_STRING:
// 		return s.VisitStringNode(node.(*ast.JsonStringNode))
// 	case ast.AST_NUMBER:
// 		return s.VisitNumberNode(node.(*ast.JsonNumberNode))
// 	case ast.AST_BOOLEAN:
// 		return s.VisitBooleanNode(node.(*ast.JsonBooleanNode))
// 	case ast.AST_NULL:
// 		return s.VisitNullNode(node.(*ast.JsonNullNode))
// 	case ast.AST_STRING_VARIABLE:
// 		return s.VisitStringWithVariableNode(node.(*ast.JsonExtendedStringWIthVariableNode))
// 	case ast.AST_VARIABLE:
// 		return s.VisitVariableNode(node.(*ast.JsonExtendedVariableNode))
// 	case ast.AST_KVPAIR:
// 		return s.VisitKeyValuePairNode(node.(*ast.JsonKeyValuePairNode))
// 	default:
// 		return ErrorInternalInterpreterOutdated
// 	}
// }

// func (s *goStructdVisitor) VisitArrayNode(node *ast.JsonArrayNode) error {
// 	s.sb.WriteString("[\n")
// 	s.indent++
// 	s.sb.WriteString(strings.Repeat(s.indentString, s.indent))
// 	for i := len(node.Value) - 1; i >= 0; i-- {
// 		s.stackNode.Push(node.Value[i])
// 		if i == len(node.Value)-1 {
// 			s.stackFormat.Push(']')
// 		} else {
// 			s.stackFormat.Push(',')
// 		}
// 	}
// 	return nil
// }

// func (s *goStructdVisitor) VisitKeyValuePairNode(node *ast.JsonKeyValuePairNode) error {
// 	// stack, first in last out, value go first ;-)
// 	s.stackNode.Push(node.Value)
// 	s.stackNode.Push(node.Key)

// 	return nil
// }

// func (s *goStructdVisitor) GetOutput() string {
// 	return s.sb.String()
// }

// func (s *goStructdVisitor) VisitObjectNode(node *ast.JsonObjectNode) error {
// 	s.stackOut.Push(node)
// 	for i := len(node.Value) - 1; i >= 0; i-- {
// 		s.stackNode.Push(node.Value[i])
// 		if i == len(node.Value)-1 { // stack, first in last out
// 			s.stackFormat.Push('}')
// 			s.stackFormat.Push(':')
// 		} else {
// 			s.stackFormat.Push(',')
// 			s.stackFormat.Push(':')
// 		}
// 	}
// 	return nil
// }

// func GoStructInterpret(node ast.JsonNode, variables map[string]interface{}, out interface{}) (string, error) {
// 	// deep first traverse the AST

// 	visitor := NewGoStructInterpreter(variables, out)
// 	visitor.stackNode.Push(node)

// 	for {
// 		node, err := visitor.stackNode.Pop()
// 		if err != nil {
// 			break
// 		}
// 		err = visitor.Visit(node)
// 		if err != nil {
// 			if err != util.ErrorEodOfStack {
// 				return "", err
// 			} else {
// 				break
// 			}
// 		}

// 	}

// 	if visitor.getSymbolLength() > 0 {
// 		return "", ErrorInterpreSymbolFailure
// 	}
// 	rs := visitor.GetOutput()
// 	return rs, nil
// }
