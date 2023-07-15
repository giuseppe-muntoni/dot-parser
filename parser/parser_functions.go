package parser

import (
	"dot-parser/iterator"
	. "dot-parser/lexer"
	"dot-parser/option"
	. "dot-parser/result"
)

//Higher-Order parser functions

// The function takes an iterator and other functions to describe the structure of the parser.
// The input functions are executed on the iterator in order, and the final iterator is returned.
func parse(iter TokenIterator, fns ...func(TokenIterator) Result[TokenIterator]) Result[TokenIterator] {
	functions := iterator.ListIterator(fns)
	return iterator.Fold(Ok(iter), functions, FlatMap[TokenIterator, TokenIterator])
}

// Wraps a function which recognises data of a certain type from an iterator and saves this data in the location
// provided by the input pointer.
func keep[T any](pointer *T, fn func(TokenIterator) Result[parserData[T]]) func(TokenIterator) Result[TokenIterator] {
	return func(iter TokenIterator) Result[TokenIterator] {
		return Map(fn(iter),
			func(data parserData[T]) TokenIterator {
				*pointer = data.value
				return data.iter
			},
		)
	}
}

// Wraps a function which recognises data of a certain type from an iterator without saving it.
// This function is used to advance the iterator without producing any useful value.
func skip[T any](fn func(TokenIterator) Result[parserData[T]]) func(TokenIterator) Result[TokenIterator] {
	return func(iter TokenIterator) Result[TokenIterator] {
		return Map(fn(iter),
			func(data parserData[T]) TokenIterator {
				return data.iter
			},
		)
	}
}

// Wraps a function which recognises data of a certain type from an iterator.
// The function is executed if the expected tokens are available, otherwise the iterator unchanged is returned.
// The expected tokens parameter is a list of list of tokens. The function is exectued if the first peeked character is contained in the first list of tokens,
// the second peeked character is contained in the second list of tokens and so on.
func optional[T any](fn func(TokenIterator) Result[parserData[T]], expectedTokens ...[]Token) func(TokenIterator) Result[parserData[option.Option[T]]] {
	return func(iter TokenIterator) Result[parserData[option.Option[T]]] {
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

// Similarly to the optional function, here a list of results is returned.
// The provided parser function is called iteratively while the expected tokens are available.
func list[T any](fn func(TokenIterator) Result[parserData[T]], expectedTokens ...[]Token) func(TokenIterator) Result[parserData[[]T]] {
	return func(iter TokenIterator) Result[parserData[[]T]] {
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

// Similarly to the list function, here it is required that the first element is always present for the parsing to be successful.
func nonEmptyList[T any](fn func(TokenIterator) Result[parserData[T]], expectedTokens ...[]Token) func(TokenIterator) Result[parserData[[]T]] {
	return func(iter TokenIterator) Result[parserData[[]T]] {
		var out_list []T

		if res := fn(iter); res.IsOk() {
			iter = res.Unwrap().iter
			value := res.Unwrap().value
			out_list = append(out_list, value)
		} else {
			return Err[parserData[[]T]](res.UnwrapErr())
		}

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

// Provides a higher-order function that matches one of the tokens provided as input.
func matchToken(expectedTokens ...Token) func(TokenIterator) Result[parserData[TokenData]] {
	return func(iter TokenIterator) Result[parserData[TokenData]] {
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

// Provides a higher-order function that matches one of the tokens provided as input when peeking at the given depth.
func peekToken(depth int32, expectedTokens ...Token) func(TokenIterator) bool {
	return func(iter TokenIterator) bool {
		token := iter.PeekNth(depth)
		if token.IsNone() {
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
