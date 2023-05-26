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

func Consume[T any](iter Iterator[T]) {
	for {
		if iter.Next().IsNone() {
			return
		}
	}
}

func Last[T any](iter Iterator[T]) option.Option[T] {
	var lastElement = option.None[T]()
	for {
		currentElement := iter.Next()
		if currentElement.IsNone() {
			return lastElement
		} else {
			lastElement = currentElement
		}
	}
}
