package golang

import (
	"encoding/base64"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

// golang plugin should implement:
// string option, omit empty, byte slice to base64 and vice versa

var stringOptPlugin ast.ASTNodePlugin = ast.NewASTNodePlugin(PLUGIN_STRING_ENCODE, stringOptionConversion, nil)

func stringOptionConversion(visitor ast.JsonVisitor, node ast.JsonNode) error {
	// if node is visited, skip
	if node.IsVisited() {
		return nil
	}
	// create a temp node for visitor to visit
	var tempNode ast.JsonStringNode
	switch instance := node.(type) {
	case *ast.JsonBooleanNode:
		val, err := util.EncodePrimitiveValue(instance.Value)
		if err != nil {
			return err
		}
		tempNode.Value = util.EncodeToJsonString(string(val))
	case *ast.JsonNumberNode:
		val, err := util.EncodePrimitiveValue(instance.Value)
		if err != nil {
			return err
		}
		tempNode.Value = util.EncodeToJsonString(string(val))
	case *ast.JsonStringNode:
		val, err := util.EncodePrimitiveValue(string(instance.Value))
		if err != nil {
			return err
		}
		tempNode.Value = val
	case *ast.JsonNullNode:
		val, err := util.EncodePrimitiveValue(token.NullBytes)
		if err != nil {
			return err
		}
		tempNode.Value = util.EncodeToJsonString(string(val))
	default:
		return nil
	}
	node.SetVisited()
	return visitor.VisitStringNode(&tempNode)
}

var sliceConversionPlugin ast.ASTNodePlugin = ast.NewASTNodePlugin(PLUGIN_SLICE_CONVERSION, sliceByteConversion, nil)

func sliceByteConversion(visitor ast.JsonVisitor, node ast.JsonNode) error {
	// if node is visited, skip
	if node.IsVisited() {
		return nil
	}

	switch instance := node.(type) {
	case *ast.JsonStringNode:
		nodeValue, err := instance.GetValue()
		if err != nil {
			return err
		}
		arrayNode, err := arrayNodeFromStringNode(nodeValue)
		if err != nil {
			return err
		}
		// hijack the node to visited
		node.SetVisited()
		return visitor.VisitArrayNode(arrayNode)
	case *ast.JsonArrayNode:
		byteSlices := make([]byte, 0, instance.Length())
		for _, item := range instance.Value {
			switch numNode := item.(type) {
			case *ast.JsonNumberNode:
				byteSlices = append(byteSlices, byte(numNode.Value.(byte)))
			default:
				return nil
			}
		}

		newStringNode := &ast.JsonStringNode{
			Value: byteSlices,
		}
		// filp the flag to visited, and let the visitor receive value from a string node
		node.SetVisited()
		return visitor.VisitStringNode(newStringNode)
	}
	return nil

}

func arrayNodeFromStringNode(nodeValue string) (*ast.JsonArrayNode, error) {

	data, err := base64.StdEncoding.DecodeString(nodeValue)
	if err != nil {
		return nil, err
	}
	rs := &ast.JsonArrayNode{
		Value: make([]ast.JsonNode, 0, len(data)),
	}
	for _, n := range data {
		v := uint8(n)
		rs.Value = append(rs.Value, &ast.JsonNumberNode{
			Value: float64(v),
		})
	}
	return rs, nil

}
