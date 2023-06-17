package iterator

import (
	"dot-parser/option"
)

type Iterator[T any] interface {
	Next() option.Option[T]
}

func Fold[A any, T any](accumulator A, iter Iterator[T], fn func(A, T) A) A {
	for {
		if next := iter.Next(); next.IsSome() {
			accumulator = fn(accumulator, next.Unwrap())
		} else {
			return accumulator
		}
	}
}
