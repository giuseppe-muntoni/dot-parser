package parser

import (
	"strings"
	"testing"
)

func makeParser(input string) TokenIterator {
	return makeTokenIterator(strings.NewReader(input))
}

func TestParsePort(t *testing.T) {
	iter := makeParser(": PortName")

	res := parsePort(iter)
	if res.IsErr() {
		t.Fatalf("Expected port, failed with %s", res.UnwrapErr())
	}

	value := res.Unwrap().value
	if value != "PortName" {
		t.Fatalf("Expected port with name 'PortName', found %s", value)
	}
}

func TestParseNodeIdWithPort(t *testing.T) {
	iter := makeParser("NodeId : PortName")

	res := parseNodeID(iter)
	if res.IsErr() {
		t.Fatalf("Expected NodeId, failed with %s", res.UnwrapErr())
	}

	value := res.Unwrap().value
	if value.Name != "NodeId" || !value.Port.IsSome() || value.Port.Unwrap() != "PortName" {
		t.Fatalf("Expected NodeId with name 'NodeId' and port 'PortName', found %v", value)
	}
}

func TestParseNodeIdWithoutPort(t *testing.T) {
	iter := makeParser("NodeId")

	res := parseNodeID(iter)
	if res.IsErr() {
		t.Fatalf("Expected NodeId, failed with %s", res.UnwrapErr())
	}

	value := res.Unwrap().value
	if value.Name != "NodeId" || value.Port.IsSome() {
		t.Fatalf("Expected NodeId with name 'NodeId' without port, found %v", value)
	}
}

func TestParseAttribute(t *testing.T) {
	iter := makeParser("first = second")

	res := parseAttribute(iter)
	if res.IsErr() {
		t.Fatalf("Expected Attribute, failed with %s", res.UnwrapErr())
	}

	attribute := res.Unwrap().value
	if attribute.Key != "first" || attribute.Value != "second" {
		t.Fatalf("Expected Attribute with key 'first' and value 'second', found %s", attribute)
	}
}

func TestParseAttributeList(t *testing.T) {
	iter := makeParser("[ a0 = a0 a1= a1, a2 = a2; a3 = a3; ]")

	res := parseAttrList(iter)
	if res.IsErr() {
		t.Fatalf("Expected Attribute List, failed with %s", res.UnwrapErr())
	}

	attributeMap := res.Unwrap().value
	if value, contains := attributeMap["a0"]; !contains || value != "a0" {
		t.Fatalf("Expected Attribute with key 'a0' and value 'a0', found %s", attributeMap)
	}
	if value, contains := attributeMap["a1"]; !contains || value != "a1" {
		t.Fatalf("Expected Attribute with key 'a1' and value 'a1', found %s", attributeMap)
	}
	if value, contains := attributeMap["a2"]; !contains || value != "a2" {
		t.Fatalf("Expected Attribute with key 'a2' and value 'a2', found %s", attributeMap)
	}
	if value, contains := attributeMap["a3"]; !contains || value != "a3" {
		t.Fatalf("Expected Attribute with key 'a3' and value 'a3', found %s", attributeMap)
	}
}

func TestParseNodeStatementWithAttributes(t *testing.T) {
	iter := makeParser("NodeId [ a0 = a0 ]")

	res := parseNodeStmt(iter)
	if res.IsErr() {
		t.Fatalf("Expected Node Statement, failed with %s", res.UnwrapErr())
	}

	var nodeStmt Node
	switch val := res.Unwrap().value[0].(type) {
	case *Node:
		nodeStmt = *val
	default:
		t.Fatalf("Expected Node Statement, but got %v", val)
	}

	if nodeStmt.ID.Name != "NodeId" || nodeStmt.ID.Port.IsSome() {
		t.Fatalf("Expected NodeId with name NodeId and empty port, found %s", nodeStmt.ID)
	}

	if len(nodeStmt.Attributes) != 1 {
		t.Fatalf("Expected Node Statement with one attribute, but got %v", nodeStmt)
	}

	if attr := nodeStmt.Attributes[0]; attr["a0"] != "a0" {
		t.Fatalf("Expected Node Statement with attribute 'a0':'a0', but got %v", attr)
	}
}

func TestParseNodeStatementWithoutAttributes(t *testing.T) {
	iter := makeParser("NodeId ")

	res := parseNodeStmt(iter)
	if res.IsErr() {
		t.Fatalf("Expected Node Statement, failed with %s", res.UnwrapErr())
	}

	var nodeStmt Node
	switch val := res.Unwrap().value[0].(type) {
	case *Node:
		nodeStmt = *val
	default:
		t.Fatalf("Expected Node Statement, but got %v", val)
	}

	if nodeStmt.ID.Name != "NodeId" || nodeStmt.ID.Port.IsSome() {
		t.Fatalf("Expected NodeId with name NodeId and empty port, found %s", nodeStmt.ID)
	}
}

func TestParseIndirectEdgeRHS(t *testing.T) {
	iter := makeParser("-- NodeId : PortName")

	res := parseEdgeRhs(iter, false)
	if res.IsErr() {
		t.Fatalf("Expected Indirect Edge RHS, failed with %s", res.UnwrapErr())
	}

	value := res.Unwrap().value
	if value.Name != "NodeId" || !value.Port.IsSome() || value.Port.Unwrap() != "PortName" {
		t.Fatalf("Expected NodeId with name 'NodeId' and port 'PortName', found %v", value)
	}
}

func TestParseDirectEdgeRHS(t *testing.T) {
	iter := makeParser("-> NodeId : PortName")

	res := parseEdgeRhs(iter, true)
	if res.IsErr() {
		t.Fatalf("Expected Direct Edge RHS, failed with %s", res.UnwrapErr())
	}

	value := res.Unwrap().value
	if value.Name != "NodeId" || !value.Port.IsSome() || value.Port.Unwrap() != "PortName" {
		t.Fatalf("Expected NodeId with name 'NodeId' and port 'PortName', found %v", value)
	}
}

func TestParseEdgeStmt(t *testing.T) {
	iter := makeParser("NodeId0 -- NodeId1 -- NodeId2 [ class = test ]")

	res := parseEdgeStmt(iter, false)
	if res.IsErr() {
		t.Fatalf("Expected Edge Statement, failed with %s", res.UnwrapErr())
	}

	var edgeStmts []Edge
	for _, stmt := range res.Unwrap().value {
		switch val := stmt.(type) {
		case *Edge:
			edgeStmts = append(edgeStmts, *val)
		default:
			t.Fatalf("Expected Edge Statement, but got %v", val)
		}
	}

	if len(edgeStmts) != 2 {
		t.Fatalf("Expected 2 Edge Statements, but got %v", edgeStmts)
	}

	if edge := edgeStmts[0]; edge.Lnode.Name != "NodeId0" || edge.Rnode.Name != "NodeId1" {
		t.Fatalf("Expected First Edge to be 'NodeId0 -- NodeId1', but got %v", edge)
	}

	if edge := edgeStmts[1]; edge.Lnode.Name != "NodeId1" || edge.Rnode.Name != "NodeId2" {
		t.Fatalf("Expected First Edge to be 'NodeId1 -- NodeId2', but got %v", edge)
	}

	if len(edgeStmts[0].Attributes) != 1 || len(edgeStmts[1].Attributes) != 1 {
		t.Fatalf("Expected Edge Statements to have one attribute map each, but got %v", edgeStmts)
	}

	if attr := edgeStmts[0].Attributes[0]; attr["class"] != "test" {
		t.Fatalf("Expected attribute 'class':'test', but got %v", attr)
	}

	if attr := edgeStmts[1].Attributes[0]; attr["class"] != "test" {
		t.Fatalf("Expected attribute 'class':'test', but got %v", attr)
	}
}

func TestParseGraphAttrStmt(t *testing.T) {
	iter := makeParser("graph []")

	res := parseAttrStmt(iter)
	if res.IsErr() {
		t.Fatalf("Expected Attribute Statement, failed with %s", res.UnwrapErr())
	}

	var attrStmt AttributeStmt
	switch val := res.Unwrap().value[0].(type) {
	case *AttributeStmt:
		attrStmt = *val
	default:
		t.Fatalf("Expected Attribute Statement, but got %v", val)
	}

	if attrStmt.Level != GRAPH_LEVEL {
		t.Fatalf("Expected Attribute to be at Graph level, but got %v", attrStmt)
	}
}

func TestParseNodeAttrStmt(t *testing.T) {
	iter := makeParser("node [ class = test ]")

	res := parseAttrStmt(iter)
	if res.IsErr() {
		t.Fatalf("Expected Attribute Statement, failed with %s", res.UnwrapErr())
	}

	var attrStmt AttributeStmt
	switch val := res.Unwrap().value[0].(type) {
	case *AttributeStmt:
		attrStmt = *val
	default:
		t.Fatalf("Expected Attribute Statement, but got %v", val)
	}

	if attrStmt.Level != NODE_LEVEL {
		t.Fatalf("Expected Attribute to be at Graph level, but got %v", attrStmt)
	}
}

func TestParseEdgeAttrStmt(t *testing.T) {
	iter := makeParser("edge [ class = test ][][][ test = test ]")

	res := parseAttrStmt(iter)
	if res.IsErr() {
		t.Fatalf("Expected Attribute Statement, failed with %s", res.UnwrapErr())
	}

	var attrStmt AttributeStmt
	switch val := res.Unwrap().value[0].(type) {
	case *AttributeStmt:
		attrStmt = *val
	default:
		t.Fatalf("Expected Attribute Statement, but got %v", val)
	}

	if attrStmt.Level != EDGE_LEVEL {
		t.Fatalf("Expected Attribute to be at Graph level, but got %v", attrStmt)
	}
}

func TestParseSingleAttributeStmt(t *testing.T) {
	iter := makeParser("class = test")

	res := parseStmt(iter, false)
	if res.IsErr() {
		t.Fatalf("Expected Statement, failed with %s", res.UnwrapErr())
	}

	var attrStmt SingleAttribute
	switch val := res.Unwrap().value[0].(type) {
	case *SingleAttribute:
		attrStmt = *val
	default:
		t.Fatalf("Expected Single Attribute Statement, but got %v", val)
	}

	if attrStmt.Key != "class" && attrStmt.Value != "test" {
		t.Fatalf("Expected attribute 'class':'test', but got %v", attrStmt)
	}
}

func TestParseStmtNode(t *testing.T) {
	iter := makeParser("NodeId []")

	res := parseStmt(iter, false)
	if res.IsErr() {
		t.Fatalf("Expected Statement, failed with %s", res.UnwrapErr())
	}

	switch val := res.Unwrap().value[0].(type) {
	case *Node:
		break
	default:
		t.Fatalf("Expected Node Statement, but got %v", val)
	}
}

func TestParseStmtEdge(t *testing.T) {
	iter := makeParser("NodeId -> NodeId -> NodeId")

	res := parseStmt(iter, true)
	if res.IsErr() {
		t.Fatalf("Expected Statement, failed with %s", res.UnwrapErr())
	}

	for _, val := range res.Unwrap().value {
		switch val := val.(type) {
		case *Edge:
			break
		default:
			t.Fatalf("Expected Node Statement, but got %v", val)
		}
	}
}

func TestParseStmtAttributeStmt(t *testing.T) {
	iter := makeParser("graph []")

	res := parseStmt(iter, true)
	if res.IsErr() {
		t.Fatalf("Expected Statement, failed with %s", res.UnwrapErr())
	}

	switch val := res.Unwrap().value[0].(type) {
	case *AttributeStmt:
		break
	default:
		t.Fatalf("Expected Node Statement, but got %v", val)
	}
}

func TestParseDirectGraph(t *testing.T) {
	iter := makeParser("digraph GraphName { NodeId []; NodeId -> NodeId -> NodeId; edge [][]; NodeId -> NodeId }")

	res := parseGraph(iter)
	if res.IsErr() {
		t.Fatalf("Expected Graph, failed with %s", res.UnwrapErr())
	}

	graph := res.Unwrap().value

	if !graph.IsDirect {
		t.Fatalf("Expected Direct Graph, got %#v", graph)
	}

	if graph.IsStrict {
		t.Fatalf("Expected Non Strict Graph, got %#v", graph)
	}

	if graph.Name.IsNone() || graph.Name.Unwrap() != "GraphName" {
		t.Fatalf("Expected Graph with name 'GraphName', got %#v", graph)
	}

	if len(graph.Statements) != 5 {
		t.Fatalf("Expected Graph with 5 statements, got %#v", graph)
	}
}

func TestParseIndirectGraph(t *testing.T) {
	iter := makeParser("strict graph { }")

	res := parseGraph(iter)
	if res.IsErr() {
		t.Fatalf("Expected Graph, failed with %s", res.UnwrapErr())
	}

	graph := res.Unwrap().value

	if graph.IsDirect {
		t.Fatalf("Expected Non Direct Graph, got %#v", graph)
	}

	if !graph.IsStrict {
		t.Fatalf("Expected Strict Graph, got %#v", graph)
	}

	if graph.Name.IsSome() {
		t.Fatalf("Expected Graph without name, got %#v", graph)
	}

	if len(graph.Statements) != 0 {
		t.Fatalf("Expected Graph with 0 statements, got %#v", graph)
	}
}
