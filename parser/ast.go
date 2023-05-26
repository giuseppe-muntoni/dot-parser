package parser

import "dot-parser/option"

type Graph struct {
	isStrict   bool
	isDirect   bool
	name       option.Option[string]
	statements []Statement
}

type Statement interface {
	isStatement() bool
}

type AttributeMap map[string]string

type Node struct {
	ID         NodeID
	attributes []AttributeMap
}

type NodeID struct {
	name string
	port option.Option[string]
}

func makeNodeID(name string, port option.Option[string]) NodeID {
	return NodeID{
		name: name,
		port: port,
	}
}

type Edge struct {
	lnode      NodeID
	rnode      NodeID
	attributes []AttributeMap
}

type AttributeLevel uint8

const (
	GRAPH_LEVEL AttributeLevel = iota
	NODE_LEVEL
	EDGE_LEVEL
)

type AttributeStmt struct {
	level      AttributeLevel
	attributes []AttributeMap
}

type SingleAttribute struct {
	key   string
	value string
}

func (n *Node) isStatement() bool            { return true }
func (e *Edge) isStatement() bool            { return true }
func (a *AttributeStmt) isStatement() bool   { return true }
func (a *SingleAttribute) isStatement() bool { return true }

func (attrs AttributeMap) String() string {
	var out_string string
	for key, value := range attrs {
		out_string += key + " : " + value + "; "
	}
	return "[ " + out_string + "]"
}

func (node NodeID) String() string {
	return node.name + ":" + node.port.OrElse("/")
}

func (node Node) String() string {
	out_string := node.ID.String() + " "
	for _, attributeMap := range node.attributes {
		out_string += attributeMap.String()
	}
	return out_string
}
