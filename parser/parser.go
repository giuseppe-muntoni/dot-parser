package parser

import (
	"dot-parser/iterator"
	. "dot-parser/lexer"
	"dot-parser/option"
	. "dot-parser/result"
	"fmt"
	"io"
)

type tokenIterator iterator.MultiPeekableIterator[Result[TokenData]]

func makeTokenIterator(reader io.Reader) tokenIterator {
	lex := MakeLexer(reader)
	iter := iterator.TakeWhile(lex, func(token Result[TokenData]) bool {
		eofFound := false
		return Map(token, func(token TokenData) bool {
			if token.Token() == EOF {
				if eofFound {
					return false
				} else {
					eofFound = true
					return true
				}
			} else {
				return true
			}
		}).OrElse(false)
	})
	return iterator.Buffered(iter)
}

type ParserError struct {
	token         TokenData
	expectedToken Token
}

func (err *ParserError) Error() string {
	return fmt.Sprintf(
		"Parsing error at line %d column %d: Got token %s with lexeme \"%s\", but %s was expected",
		err.token.Position().Line(),
		err.token.Position().Column(),
		err.token.Token(),
		err.token.Lexeme(),
		err.expectedToken)
}

func makeParserError[T any](token TokenData, expectedToken Token) Result[parserData[T]] {
	return Err[parserData[T]](
		&ParserError{
			token:         token,
			expectedToken: expectedToken,
		},
	)
}

type parserData[T any] struct {
	value T
	iter  tokenIterator
}

func makeParserData[T any](iter tokenIterator, value T) Result[parserData[T]] {
	return Ok(parserData[T]{
		iter:  iter,
		value: value,
	})
}

func makeParserDataRes[T any](iter Result[tokenIterator], value T) Result[parserData[T]] {
	return FlatMap(iter, func(iter tokenIterator) Result[parserData[T]] {
		return makeParserData(iter, value)
	})
}

func parseGraph(iter tokenIterator) Result[parserData[Graph]] {
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
			keep(&stmts, list(parseStmtInList, []Token{ID, GRAPH, NODE, EDGE})),
			skip(matchToken(CLOSE_BRACE)),
			skip(matchToken(EOF)),
		)
	})

	var graphStmts []Statement
	for _, stmts := range stmts {
		for _, stmt := range stmts {
			graphStmts = append(graphStmts, stmt)
		}
	}

	return makeParserDataRes(newIter, Graph{
		isStrict:   strict,
		isDirect:   isDirect,
		name:       option.Map(name, func(token TokenData) string { return string(token.Lexeme()) }),
		statements: graphStmts,
	})
}

func parseStmtInList(iter tokenIterator, isDirect bool) Result[parserData[[]Statement]] {
	var stmt []Statement

	newIter := parse(iter,
		keep(&stmt, partialApply(isDirect, parseStmt)),
		skip(optional(matchToken(SEMICOLON), []Token{SEMICOLON})),
	)

	return makeParserDataRes(newIter, stmt)
}

func parseStmt(iter tokenIterator, isDirect bool) Result[parserData[[]Statement]] {
	var stmt []Statement

	var newIter Result[tokenIterator]
	if peekToken(1, ID)(iter) {
		if peekToken(2, EQUAL)(iter) {
			var attrib SingleAttribute
			newIter = parse(iter, keep(&attrib, parseSingleAttribStmt))
			stmt = []Statement{&attrib}
		} else if peekToken(2, ARC, DIRECTED_ARC)(iter) || peekToken(4, ARC, DIRECTED_ARC)(iter) {
			var edgeStmt []Edge
			parseEdgeStmt := partialApply(isDirect, parseEdgeStmt)
			newIter = parse(iter, keep(&edgeStmt, parseEdgeStmt))

			stmt = []Statement{}
			for _, edge := range edgeStmt {
				stmt = append(stmt, &edge)
			}
		} else {
			var nodeStmt Node
			newIter = parse(iter, keep(&nodeStmt, parseNodeStmt))
			stmt = []Statement{&nodeStmt}
		}
	} else {
		var attrStmt AttributeStmt
		newIter = parse(iter, keep(&attrStmt, parseAttrStmt))
		stmt = []Statement{&attrStmt}
	}

	return makeParserDataRes(newIter, stmt)
}

func parseSingleAttribStmt(iter tokenIterator) Result[parserData[SingleAttribute]] {
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

func parseAttrStmt(iter tokenIterator) Result[parserData[AttributeStmt]] {
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

	return makeParserDataRes(newIter, AttributeStmt{
		level:      level,
		attributes: attrList,
	})
}

func parseEdgeStmt(iter tokenIterator, isDirect bool) Result[parserData[[]Edge]] {
	var firstLhs NodeID
	var firstRhs NodeID
	var nodes []NodeID
	var edges []Edge
	var attributes []AttributeMap

	parseEdgeRhs := partialApply(isDirect, parseEdgeRhs)

	newIter := parse(iter,
		keep(&firstLhs, parseNodeID),
		keep(&firstRhs, parseEdgeRhs),
		keep(&nodes, list(parseEdgeRhs, []Token{ARC, DIRECTED_ARC})),
		keep(&attributes, list(parseAttrList, []Token{OPEN_SQUARE_BRACKET})),
	)

	edges = append(edges, Edge{
		lnode:      firstLhs,
		rnode:      firstRhs,
		attributes: attributes,
	})

	for _, node := range nodes {
		edges = append(edges, Edge{
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

func parseNodeStmt(iter tokenIterator) Result[parserData[Node]] {
	var nodeID NodeID
	var attrList []AttributeMap

	newIter := parse(iter,
		keep(&nodeID, parseNodeID),
		keep(&attrList, list(parseAttrList, []Token{OPEN_SQUARE_BRACKET})),
	)

	return makeParserDataRes(newIter, Node{
		ID:         nodeID,
		attributes: attrList,
	})
}

func parseAttrList(iter tokenIterator) Result[parserData[AttributeMap]] {
	var attributes []SingleAttribute
	newIter := parse(iter,
		skip(matchToken(OPEN_SQUARE_BRACKET)),
		keep(&attributes, list(parseAttribute, []Token{ID})),
		skip(matchToken(CLOSE_SQUARE_BRACKET)),
	)

	var finalAttributes = make(AttributeMap, len(attributes))
	for _, attribute := range attributes {
		finalAttributes[attribute.key] = attribute.value
	}

	return makeParserDataRes(newIter, finalAttributes)
}

func parseAttribute(iter tokenIterator) Result[parserData[SingleAttribute]] {
	var firstId TokenData
	var secondId TokenData
	newIter := parse(iter,
		keep(&firstId, matchToken(ID)),
		skip(matchToken(EQUAL)),
		keep(&secondId, matchToken(ID)),
		skip(optional(matchToken(SEMICOLON, COMMA), []Token{SEMICOLON, COMMA})),
	)

	return makeParserDataRes(newIter, SingleAttribute{
		key:   string(firstId.Lexeme()),
		value: string(secondId.Lexeme()),
	})
}

func parseNodeID(iter tokenIterator) Result[parserData[NodeID]] {
	var nodeName TokenData
	var port option.Option[string]
	newIter := parse(iter,
		keep(&nodeName, matchToken(ID)),
		keep(&port, optional(parsePort, []Token{COLON})),
	)

	return makeParserDataRes(newIter, makeNodeID(string(nodeName.Lexeme()), port))
}

func parsePort(iter tokenIterator) Result[parserData[string]] {
	var port TokenData
	newIter := parse(iter,
		skip(matchToken(COLON)),
		keep(&port, matchToken(ID)),
	)

	return makeParserDataRes(newIter, string(port.Lexeme()))
}

func parse(iter tokenIterator, fns ...func(tokenIterator) Result[tokenIterator]) Result[tokenIterator] {
	functions := iterator.ListIterator(fns)
	return iterator.Fold(Ok(iter), functions, FlatMap[tokenIterator, tokenIterator])
}

func keep[T any](pointer *T, fn func(tokenIterator) Result[parserData[T]]) func(tokenIterator) Result[tokenIterator] {
	return func(iter tokenIterator) Result[tokenIterator] {
		return Map(fn(iter),
			func(data parserData[T]) tokenIterator {
				*pointer = data.value
				return data.iter
			},
		)
	}
}

func skip[T any](fn func(tokenIterator) Result[parserData[T]]) func(tokenIterator) Result[tokenIterator] {
	return func(iter tokenIterator) Result[tokenIterator] {
		return Map(fn(iter),
			func(data parserData[T]) tokenIterator {
				return data.iter
			},
		)
	}
}

func optional[T any](fn func(tokenIterator) Result[parserData[T]], expectedTokens ...[]Token) func(tokenIterator) Result[parserData[option.Option[T]]] {
	return func(iter tokenIterator) Result[parserData[option.Option[T]]] {
		var depth int32 = 1
		expectedTokensList := iterator.ListIterator(expectedTokens)
		isPresent := iterator.Fold(true, expectedTokensList, func(accum bool, expectedTokens []Token) bool {
			if accum {
				accum = peekToken(depth, expectedTokens...)(iter)
				depth += 1
			}
			return accum
		})

		if isPresent {
			return FlatMap(fn(iter),
				func(data parserData[T]) Result[parserData[option.Option[T]]] {
					return makeParserData(data.iter, option.Some(data.value))
				},
			)
		} else {
			return makeParserData(iter, option.None[T]())
		}
	}
}

func list[T any](fn func(tokenIterator) Result[parserData[T]], expectedTokens ...[]Token) func(tokenIterator) Result[parserData[[]T]] {
	return func(iter tokenIterator) Result[parserData[[]T]] {
		var out_list []T
		for {
			if res := optional(fn, expectedTokens...)(iter); res.IsOk() {
				iter = res.Unwrap().iter
				value := res.Unwrap().value
				if value.IsSome() {
					out_list = append(out_list, value.Unwrap())
				} else {
					return makeParserData(iter, out_list)
				}
			} else {
				return Err[parserData[[]T]](res.UnwrapErr())
			}
		}
	}
}

func matchToken(expectedTokens ...Token) func(tokenIterator) Result[parserData[TokenData]] {
	return func(iter tokenIterator) Result[parserData[TokenData]] {
		token := iter.Next().Unwrap()
		return FlatMap(token, func(token TokenData) Result[parserData[TokenData]] {
			for _, expectedToken := range expectedTokens {
				if token.Token() == expectedToken {
					return makeParserData(iter, token)
				}
			}
			return makeParserError[TokenData](token, expectedTokens[0])
		})
	}
}

func peekToken(depth int32, expectedTokens ...Token) func(tokenIterator) bool {
	return func(iter tokenIterator) bool {
		token := iter.PeekNth(depth)
		if !token.IsSome() {
			return false
		}

		return Map(token.Unwrap(),
			func(token TokenData) bool {
				for _, expectedToken := range expectedTokens {
					if token.Token() == expectedToken {
						return true
					}
				}
				return false
			},
		).OrElse(false)
	}
}
