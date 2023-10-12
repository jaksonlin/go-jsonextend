package golang

import (
	"encoding/base64"

	"github.com/jaksonlin/go-jsonextend/ast"
	"github.com/jaksonlin/go-jsonextend/token"
	"github.com/jaksonlin/go-jsonextend/util"
)

// golang plugin should implement:
// string option, omit empty, byte slice to base64 and vice versa

// string option plugin
type stringOptionPlugin struct {
}

var _ ast.ASTNodePlugin = &stringOptionPlugin{}

func (plugin *stringOptionPlugin) PluginName() string {
	return "string_option_plugin"
}

func (plugin *stringOptionPlugin) PreVisitPlugin(visitor ast.JsonVisitor, node ast.JsonNode) error {
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
		tempNode.Value = val
	case *ast.JsonNumberNode:
		val, err := util.EncodePrimitiveValue(instance.Value)
		if err != nil {
			return err
		}
		tempNode.Value = val
	case *ast.JsonStringNode:
		strVal, err := instance.GetValue()
		if err != nil {
			return err
		}
		val, err := util.EncodePrimitiveValue(strVal)
		if err != nil {
			return err
		}
		tempNode.Value = val
	case *ast.JsonNullNode:
		val, err := util.EncodePrimitiveValue(token.NullBytes)
		if err != nil {
			return err
		}
		tempNode.Value = val
	default:
		return nil
	}
	node.SetVisited()
	return visitor.VisitStringNode(&tempNode)
}

func (plugin *stringOptionPlugin) PostVisitPlugin(visitor ast.JsonVisitor, node ast.JsonNode) error {
	return nil
}

// string option plugin
type byteSliceConversionPlugin struct {
}

func (plugin *byteSliceConversionPlugin) PluginName() string {
	return "byte_slice_conversion_plugin"
}

func (plugin *byteSliceConversionPlugin) PluginVisitor(visitor ast.JsonVisitor, node ast.JsonNode) error {
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
				byteSlices = append(byteSlices, byte(numNode.OriginValue.(byte)))
			default:
				return nil
			}
		}
		// hijack the node to visited
		newStringNode := &ast.JsonStringNode{
			Value: byteSlices,
		}
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

type omitEmptyPlugin struct {
}

func (plugin *omitEmptyPlugin) PluginVisitor(visitor ast.JsonVisitor, node ast.JsonNode) error {
	// this is a struct field, we need to omit it if the value is empty
	node, ok := node.(*ast.JsonKeyValuePairNode)
	if !ok {
		return nil
	}
	switch valueNode := node.(type) {
	case *ast.JsonStringNode:
		val, err := valueNode.GetValue()
		if err != nil {
			return err
		}
		if val == "" {
			node.SetVisited()
			return nil
		}
	case *ast.JsonNumberNode:
		if valueNode.Value == 0 {
			node.SetVisited()
			return nil
		}
	case *ast.JsonBooleanNode:
		if !valueNode.Value {
			node.SetVisited()
			return nil
		}
	case *ast.JsonNullNode:
		node.SetVisited()
		return nil
	case *ast.JsonArrayNode:
		if valueNode.Length() == 0 {
			node.SetVisited()
			return nil
		}
	case *ast.JsonObjectNode:
		isMapMeta := valueNode.GetMeta(OBJECT_FROM_MAP_META)
		// map with no key value pair in struct field when omitempty set, omit it
		if isMapMeta != nil && isMapMeta.(bool) {
			if valueNode.Length() == 0 {
				node.SetVisited()
				return nil
			}
		}
	default:
		return nil
	}
	return nil
}
