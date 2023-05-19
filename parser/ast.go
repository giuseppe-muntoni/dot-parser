package parser

import "dot-parser/option"

type Graph struct {
	isStrict   bool
	isDirect   bool
	name       string
	statements []Statement
}

type Statement interface {
	isStatement() bool
}

type AttributeMap map[string]string

type Node struct {
	name       string
	port       option.Option[string]
	attributes []AttributeMap
}

type Edge struct {
	lnode      string
	rnode      string
	attributes []AttributeMap
}

type AttributeLevel uint8

const (
	GRAPH AttributeLevel = iota
	NODE
	EDGE
)

type AttributeStmt struct {
	level AttributeLevel
}

func (n *Node) isStatement() bool          { return true }
func (e *Edge) isStatement() bool          { return true }
func (a *AttributeStmt) isStatement() bool { return true }
