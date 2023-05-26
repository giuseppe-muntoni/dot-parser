package parser

import (
	"strings"
	"testing"
)

func makeParser(input string) tokenIterator {
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
	if value.name != "NodeId" || !value.port.IsSome() || value.port.Unwrap() != "PortName" {
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
	if value.name != "NodeId" || value.port.IsSome() {
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
	if attribute.key != "first" || attribute.value != "second" {
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

	if nodeStmt.ID.name != "NodeId" || nodeStmt.ID.port.IsSome() {
		t.Fatalf("Expected NodeId with name NodeId and empty port, found %s", nodeStmt.ID)
	}
}
