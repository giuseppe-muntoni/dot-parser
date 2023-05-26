package parser

import (
	. "dot-parser/lexer"
	"dot-parser/option"
	. "dot-parser/result"
)

// Graph:
// | STRICT? GRAPH ID? '{' StatementInList(false)+ '}' EOF
// | STRICT? DIGRAPH ID? '{' StatementInList(true)+ '}' EOF
func ParseGraph(iter tokenIterator) Result[parserData[Graph]] {
	var strictT option.Option[TokenData]
	var isDirectT TokenData
	var name option.Option[TokenData]
	var stmts [][]Statement

	newIter := parse(iter,
		keep(&strictT, optional(matchToken(STRICT), []Token{STRICT})),
		keep(&isDirectT, matchToken(GRAPH, DIGRAPH)),
	)

	var strict = strictT.IsSome()
	var isDirect = isDirectT.Token() == DIGRAPH

	newIter = FlatMap(newIter, func(iter tokenIterator) Result[tokenIterator] {
		parseStmtInList := partialApply(isDirect, parseStmtInList)
		return parse(iter,
			keep(&name, optional(matchToken(ID), []Token{ID})),
			skip(matchToken(OPEN_BRACE)),
			keep(&stmts, nonEmptyList(parseStmtInList, []Token{ID, GRAPH, NODE, EDGE})),
			skip(matchToken(CLOSE_BRACE)),
			skip(matchToken(EOF)),
		)
	})

	var graphStmts []Statement
	for _, stmts := range stmts {
		graphStmts = append(graphStmts, stmts...)
	}

	return makeParserDataRes(newIter, Graph{
		isStrict:   strict,
		isDirect:   isDirect,
		name:       option.Map(name, func(token TokenData) string { return string(token.Lexeme()) }),
		statements: graphStmts,
	})
}

// StatementInList: Statement ';'?
func parseStmtInList(iter tokenIterator, isDirect bool) Result[parserData[[]Statement]] {
	var stmt []Statement

	newIter := parse(iter,
		keep(&stmt, partialApply(isDirect, parseStmt)),
		skip(optional(matchToken(SEMICOLON), []Token{SEMICOLON})),
	)

	return makeParserDataRes(newIter, stmt)
}

// Statement(isDirect bool): NodeStatement | EdgeStatement(isDirect) | AttributeStatement | SingleAttributeStatement
func parseStmt(iter tokenIterator, isDirect bool) Result[parserData[[]Statement]] {
	var stmt []Statement

	var newIter Result[tokenIterator]
	if peekToken(1, ID)(iter) {
		if peekToken(2, EQUAL)(iter) {
			var attrib SingleAttribute
			newIter = parse(iter, keep(&attrib, parseAttribute))
			stmt = []Statement{&attrib}
		} else if peekToken(2, ARC, DIRECTED_ARC)(iter) || peekToken(4, ARC, DIRECTED_ARC)(iter) {
			newIter = parse(iter, keep(&stmt, partialApply(isDirect, parseEdgeStmt)))
		} else {
			newIter = parse(iter, keep(&stmt, parseNodeStmt))
		}
	} else {
		newIter = parse(iter, keep(&stmt, parseAttrStmt))
	}

	return makeParserDataRes(newIter, stmt)
}

// AttributeStatement: (GRAPH | NODE | EDGE) AttributeList*
func parseAttrStmt(iter tokenIterator) Result[parserData[[]Statement]] {
	var attrType TokenData
	var attrList []AttributeMap

	newIter := parse(iter,
		keep(&attrType, matchToken(GRAPH, NODE, EDGE)),
		keep(&attrList, list(parseAttrList, []Token{OPEN_SQUARE_BRACKET})),
	)

	var level AttributeLevel
	switch attrType.Token() {
	case NODE:
		level = NODE_LEVEL
	case EDGE:
		level = EDGE_LEVEL
	case GRAPH:
		level = GRAPH_LEVEL
	}

	attribute := AttributeStmt{level: level, attributes: attrList}
	return makeParserDataRes(newIter, []Statement{&attribute})
}

// EdgeStatement(isDirect bool): NodeId EdgeRHS(isDirect)+ AttributeList*
func parseEdgeStmt(iter tokenIterator, isDirect bool) Result[parserData[[]Statement]] {
	var firstLhs NodeID
	var nodes []NodeID
	var attributes []AttributeMap

	parseEdgeRhs := partialApply(isDirect, parseEdgeRhs)

	newIter := parse(iter,
		keep(&firstLhs, parseNodeID),
		keep(&nodes, nonEmptyList(parseEdgeRhs, []Token{ARC, DIRECTED_ARC})),
		keep(&attributes, list(parseAttrList, []Token{OPEN_SQUARE_BRACKET})),
	)

	var edges []Statement
	firstRhs, nodes := nodes[0], nodes[1:]
	edges = append(edges, &Edge{
		lnode:      firstLhs,
		rnode:      firstRhs,
		attributes: attributes,
	})

	for _, node := range nodes {
		edges = append(edges, &Edge{
			lnode:      firstRhs,
			rnode:      node,
			attributes: attributes,
		})
		firstRhs = node
	}
	return makeParserDataRes(newIter, edges)
}

func partialApply[T any](isDirect bool, fn func(tokenIterator, bool) Result[parserData[T]]) func(tokenIterator) Result[parserData[T]] {
	return func(iter tokenIterator) Result[parserData[T]] {
		return fn(iter, isDirect)
	}
}

// EdgeRHS(isDirect bool):
// if isDirect: DIRECTED_ARC NodeId
// else: ARC NodeId
func parseEdgeRhs(iter tokenIterator, isDirect bool) Result[parserData[NodeID]] {
	var nodeID NodeID
	var newIter Result[tokenIterator]

	var matchArc func(tokenIterator) Result[parserData[TokenData]]
	if isDirect {
		matchArc = matchToken(DIRECTED_ARC)
	} else {
		matchArc = matchToken(ARC)
	}

	newIter = parse(iter,
		skip(matchArc),
		keep(&nodeID, parseNodeID),
	)

	return makeParserDataRes(newIter, nodeID)
}

// NodeStatement: NodeId AttributeList*
func parseNodeStmt(iter tokenIterator) Result[parserData[[]Statement]] {
	var nodeID NodeID
	var attrList []AttributeMap

	newIter := parse(iter,
		keep(&nodeID, parseNodeID),
		keep(&attrList, list(parseAttrList, []Token{OPEN_SQUARE_BRACKET})),
	)

	node := Node{ID: nodeID, attributes: attrList}
	return makeParserDataRes(newIter, []Statement{&node})
}

// AttributeList: '[' SingleAttribute* ']'
func parseAttrList(iter tokenIterator) Result[parserData[AttributeMap]] {
	var attributes []SingleAttribute
	newIter := parse(iter,
		skip(matchToken(OPEN_SQUARE_BRACKET)),
		keep(&attributes, list(parseAttributeInList, []Token{ID})),
		skip(matchToken(CLOSE_SQUARE_BRACKET)),
	)

	var finalAttributes = make(AttributeMap, len(attributes))
	for _, attribute := range attributes {
		finalAttributes[attribute.key] = attribute.value
	}

	return makeParserDataRes(newIter, finalAttributes)
}

// SingleAttribute: ID '=' ID (';' | ',')?
func parseAttributeInList(iter tokenIterator) Result[parserData[SingleAttribute]] {
	var attrib SingleAttribute
	newIter := parse(iter,
		keep(&attrib, parseAttribute),
		skip(optional(matchToken(SEMICOLON, COMMA), []Token{SEMICOLON, COMMA})),
	)

	return makeParserDataRes(newIter, attrib)
}

// SingleAttribute: ID '=' ID
func parseAttribute(iter tokenIterator) Result[parserData[SingleAttribute]] {
	var firstId TokenData
	var secondId TokenData
	newIter := parse(iter,
		keep(&firstId, matchToken(ID)),
		skip(matchToken(EQUAL)),
		keep(&secondId, matchToken(ID)),
	)

	return makeParserDataRes(newIter, SingleAttribute{
		key:   string(firstId.Lexeme()),
		value: string(secondId.Lexeme()),
	})
}

// NodeId: ID Port?
func parseNodeID(iter tokenIterator) Result[parserData[NodeID]] {
	var nodeName TokenData
	var port option.Option[string]
	newIter := parse(iter,
		keep(&nodeName, matchToken(ID)),
		keep(&port, optional(parsePort, []Token{COLON})),
	)

	return makeParserDataRes(newIter, makeNodeID(string(nodeName.Lexeme()), port))
}

// Port: ':' ID
func parsePort(iter tokenIterator) Result[parserData[string]] {
	var port TokenData
	newIter := parse(iter,
		skip(matchToken(COLON)),
		keep(&port, matchToken(ID)),
	)

	return makeParserDataRes(newIter, string(port.Lexeme()))
}
