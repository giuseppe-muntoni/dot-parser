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

func makeParserError(token TokenData, expectedToken Token) Result[parserData[TokenData]] {
	return Err[parserData[TokenData]](
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

func parseAttrList[T any](iter parserData[T]) Result[parserData[AttributeMap]] {
	matchOpenSquareBracket := matchToken[T](OPEN_SQUARE_BRACKET)(iter)
	attributeMap := map[string]string{}
	for {
		isIDPresent := FlatMap(matchOpenSquareBracket, peekToken[TokenData](1, ID))
		attribute := FlatMap(isIDPresent, optionalParse(parseAList[bool]))
		if attribute.IsOk() {
			if attribute.Unwrap().value.IsSome() {
				for k, v := range attribute.Unwrap().value.Unwrap() {
					attributeMap[k] = v
				}
			} else {
				matchClosedSquareBracket := FlatMap(attribute, matchToken[option.Option[AttributeMap]](CLOSE_SQUARE_BRACKET))
				return FlatMap(matchClosedSquareBracket, func(res parserData[TokenData]) Result[parserData[AttributeMap]] {
					return makeParserData[AttributeMap](attribute.Unwrap().iter, attributeMap)
				})
			}
		} else {
			return Err[parserData[AttributeMap]](attribute.UnwrapErr())
		}
	}
}

func parseAList[T any](iter parserData[T]) Result[parserData[AttributeMap]] {
	matchFirstId := matchToken[T](ID)(iter)
	firstId := Map(matchFirstId, getLexeme).OrElse("")
	matchEqual := FlatMap(matchFirstId, matchToken[TokenData](EQUAL))
	matchSecondId := FlatMap(matchEqual, matchToken[TokenData](ID))
	secondId := Map(matchSecondId, getLexeme).OrElse("")

	isSemicolonOrColonPresent := FlatMap(matchSecondId, peekToken[TokenData](1, SEMICOLON, COMMA))
	matchSemicolonOrColon := FlatMap(isSemicolonOrColonPresent, optionalParse(matchToken[bool](SEMICOLON, COMMA)))

	return FlatMap(matchSemicolonOrColon, func(res parserData[option.Option[TokenData]]) Result[parserData[AttributeMap]] {
		attributeMap := map[string]string{firstId: secondId}
		return makeParserData[AttributeMap](res.iter, attributeMap)
	})
}

func parseNodeID[T any](iter parserData[T]) Result[parserData[NodeID]] {
	matchID := matchToken[T](ID)(iter)
	nodeName := Map(matchID, getLexeme).OrElse("")
	isPortPresent := FlatMap(matchID, peekToken[TokenData](1, COLON))
	port := FlatMap(isPortPresent, optionalParse(parsePort[bool]))
	return FlatMap(port, func(port parserData[option.Option[string]]) Result[parserData[NodeID]] {
		return makeParserData(port.iter, makeNodeID(nodeName, port.value))
	})
}

func parseNodeID[T any](iter tokenIterator) Result[parserData[NodeID]] {
	nodeName := Map(matchToken(iter, ID), getLexeme).OrElse("")

	isPortPresent := (peekToken1, COLON))
	port := FlatMap(isPortPresent, optionalParse(parsePort[bool]))
	return FlatMap(port, func(port parserData[option.Option[string]]) Result[parserData[NodeID]] {
		return makeParserData(port.iter, makeNodeID(nodeName, port.value))
	})
}

func parseNodeID(iter tokenIterator) Result[parserData[NodeID]] {
	var result Result[parserData[NodeID]] = nil
	defer func(result *Result[parserData[NodeID]]) Result[parserData[NodeID]] {
		if err, _ := recover().(error); err != nil  {
			*result = Err[parserData[NodeID]](err)
		}
		return *result
	}(&result)

	nodeName := matchToken(iter, ID)
	isPortPresent = peekToken(iter, 1, COLON)

	port := option.None[string]()
	if isPortPresent {
		port = option.Some(parsePort(iter).Unwrap())
	}

	result = makeParserData(iter, makeNodeID(nodeName, port))
}

func parsePort[T any](iter parserData[T]) Result[parserData[string]] {
	matchColon := matchToken[T](COLON)(iter)
	matchPort := FlatMap(matchColon, matchToken[TokenData](ID))
	return FlatMap(matchPort,
		func(res parserData[TokenData]) Result[parserData[string]] {
			return makeParserData[string](res.iter, string(res.value.Lexeme()))
		})
}

func getLexeme(token parserData[TokenData]) string {
	return string(token.value.Lexeme())
}

func optionalParse[T any](fn func(parserData[bool]) Result[parserData[T]]) func(iter parserData[bool]) Result[parserData[option.Option[T]]] {
	return func(iter parserData[bool]) Result[parserData[option.Option[T]]] {
		if iter.value {
			parsed := fn(iter)
			return FlatMap(parsed, func(res parserData[T]) Result[parserData[option.Option[T]]] {
				return makeParserData(res.iter, option.Some(res.value))
			})
		} else {
			return makeParserData(iter.iter, option.None[T]())
		}
	}
}

func peekToken[T any](depth int32, expectedTokens ...Token) func(parserData[T]) Result[parserData[bool]] {
	return func(res parserData[T]) Result[parserData[bool]] {
		token := res.iter.PeekNth(depth)
		if token.IsSome() {
			return FlatMap(token.Unwrap(), func(token TokenData) Result[parserData[bool]] {
				for _, expectedToken := range expectedTokens {
					if token.Token() == expectedToken {
						return makeParserData(res.iter, true)
					}
				}
				return makeParserData(res.iter, false)
			})
		} else {
			return makeParserData(res.iter, false)
		}
	}
}

func matchToken[T any](iter tokenIterator, expectedTokens ...Token) func(parserData[T]) Result[parserData[TokenData]] {
	return func(res parserData[T]) Result[parserData[TokenData]] {
		token := res.iter.Next().Unwrap()
		return FlatMap(token, func(token TokenData) Result[parserData[TokenData]] {
			for _, expectedToken := range expectedTokens {
				if token.Token() == expectedToken {
					return makeParserData(res.iter, token)
				}
			}
			return makeParserError(token, expectedTokens[0])
		})
	}
}

func matchToken[T any](iter tokenIterator, expectedTokens ...Token) Result[TokenData] {
	token := iter.Next().Unwrap()
	return FlatMap(token, func(token TokenData) Result[TokenData] {
		for _, expectedToken := range expectedTokens {
			if token.Token() == expectedToken {
				return Ok(token)
			}
		}
		return Err[TokenData](&ParserError{token: token, expectedToken: expectedTokens[0]})
	})
}
