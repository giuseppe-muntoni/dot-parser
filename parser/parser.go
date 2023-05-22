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

func parseAttrList(iter tokenIterator) Result[parserData[AttributeMap]] {
	var attributes []AttributeMap
	newIter := parse(iter,
		skip(matchToken(OPEN_SQUARE_BRACKET)),
		keep(&attributes, list(parseAList, []Token{ID})),
		skip(matchToken(CLOSE_SQUARE_BRACKET)),
	)

	var finalAttributes AttributeMap = make(AttributeMap, len(attributes))
	for _, attribute := range attributes {
		for k, v := range attribute {
			finalAttributes[k] = v
		}
	}

	return makeParserDataRes(newIter, finalAttributes)
}

func parseAList(iter tokenIterator) Result[parserData[AttributeMap]] {
	var firstId string
	var secondId string
	newIter := parse(iter,
		keep(&firstId, matchToken(ID)),
		skip(matchToken(EQUAL)),
		keep(&secondId, matchToken(ID)),
		skip(optional(matchToken(SEMICOLON, COMMA), []Token{SEMICOLON, COMMA})),
	)

	return makeParserDataRes[AttributeMap](newIter, map[string]string{firstId: secondId})
}

func parseNodeID(iter tokenIterator) Result[parserData[NodeID]] {
	var nodeName string
	var port option.Option[string]
	newIter := parse(iter,
		keep(&nodeName, matchToken(ID)),
		keep(&port, optional(parsePort, []Token{COLON})),
	)

	return makeParserDataRes(newIter, makeNodeID(nodeName, port))
}

func parsePort(iter tokenIterator) Result[parserData[string]] {
	var port string
	newIter := parse(iter,
		skip(matchToken(COLON)),
		keep(&port, matchToken(ID)),
	)

	return makeParserDataRes(newIter, port)
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

func matchToken(expectedTokens ...Token) func(tokenIterator) Result[parserData[string]] {
	return func(iter tokenIterator) Result[parserData[string]] {
		token := iter.Next().Unwrap()
		return FlatMap(token, func(token TokenData) Result[parserData[string]] {
			for _, expectedToken := range expectedTokens {
				if token.Token() == expectedToken {
					return makeParserData(iter, string(token.Lexeme()))
				}
			}
			return makeParserError[string](token, expectedTokens[0])
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
