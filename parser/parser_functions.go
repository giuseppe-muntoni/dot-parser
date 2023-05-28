package parser

import (
	"dot-parser/iterator"
	. "dot-parser/lexer"
	"dot-parser/option"
	. "dot-parser/result"
)

func parse(iter TokenIterator, fns ...func(TokenIterator) Result[TokenIterator]) Result[TokenIterator] {
	functions := iterator.ListIterator(fns)
	return iterator.Fold(Ok(iter), functions, FlatMap[TokenIterator, TokenIterator])
}

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

func skip[T any](fn func(TokenIterator) Result[parserData[T]]) func(TokenIterator) Result[TokenIterator] {
	return func(iter TokenIterator) Result[TokenIterator] {
		return Map(fn(iter),
			func(data parserData[T]) TokenIterator {
				return data.iter
			},
		)
	}
}

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
