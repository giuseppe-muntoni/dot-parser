package parser

import (
	"dot-parser/iterator"
	. "dot-parser/lexer"
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
