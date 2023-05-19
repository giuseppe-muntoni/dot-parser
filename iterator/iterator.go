package iterator

import (
	"dot-parser/option"
)

// todo: remove hasNext and leave only next()
type Iterator[T any] interface {
	HasNext() bool
	GetNext() option.Option[T]
}

func Fold[A any, T any](accumulator A, iter Iterator[T], fn func(A, T) A) A {
	for {
		if next := iter.GetNext(); next.IsSome() {
			accumulator = fn(accumulator, next.Unwrap())
		} else {
			return accumulator
		}
	}
}

func Consume[T any](iter Iterator[T]) {
	for {
		if !iter.GetNext().IsSome() {
			return
		}
	}
}

func Last[T any](iter Iterator[T]) option.Option[T] {
	var lastElement = option.None[T]()
	for {
		currentElement := iter.GetNext()
		if !currentElement.IsSome() {
			return lastElement
		} else {
			lastElement = currentElement
		}
	}
}
