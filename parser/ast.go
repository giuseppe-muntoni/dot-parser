package parser

import "dot-parser/option"

//Dot-Format Abstract Syntax Tree

type Graph struct {
	IsStrict   bool
	IsDirect   bool
	Name       option.Option[string]
	Statements []Statement
}

type Statement interface {
	isStatement() bool
}

type AttributeMap map[string]string

type Node struct {
	ID         NodeID
	Attributes []AttributeMap
}

type NodeID struct {
	Name string
	Port option.Option[string]
}

func makeNodeID(name string, port option.Option[string]) NodeID {
	return NodeID{
		Name: name,
		Port: port,
	}
}

type Edge struct {
	Lnode      NodeID
	Rnode      NodeID
	Attributes []AttributeMap
}

type AttributeLevel uint8

const (
	GRAPH_LEVEL AttributeLevel = iota
	NODE_LEVEL
	EDGE_LEVEL
)

type AttributeStmt struct {
	Level      AttributeLevel
	Attributes []AttributeMap
}

type SingleAttribute struct {
	Key   string
	Value string
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
	return node.Name + ":" + node.Port.OrElse("/")
}

func (node Node) String() string {
	out_string := node.ID.String() + " "
	for _, attributeMap := range node.Attributes {
		out_string += attributeMap.String()
	}
	return out_string
}
