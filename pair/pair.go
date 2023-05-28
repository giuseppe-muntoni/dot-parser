package pair

import "dot-parser/result"

type Pair[L any, R any] struct {
	First  L
	Second R
}

func NewPair[L any, R any](first L, second R) Pair[L, R] {
	return Pair[L, R]{First: first, Second: second}
}

func (pair Pair[L, R]) Get() (L, R) {
	return pair.First, pair.Second
}

func Zip[L any, R any](first result.Result[L], second result.Result[R]) result.Result[Pair[L, R]] {
	if first.IsErr() {
		return result.Err[Pair[L, R]](first.UnwrapErr())
	} else if second.IsErr() {
		return result.Err[Pair[L, R]](first.UnwrapErr())
	} else {
		return result.Ok(NewPair(first.Unwrap(), second.Unwrap()))
	}
}
