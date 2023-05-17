package iterator

import (
	"dot-parser/option"
)

// Iterator interface
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

// TakeWhile iterator
type TakeWhileIterator[T any] struct {
	iterator  Iterator[T]
	predicate func(T) bool
	buffer    option.Option[T]
}

func TakeWhile[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return &TakeWhileIterator[T]{
		iterator:  iter,
		predicate: predicate,
		buffer:    option.None[T](),
	}
}

func (iter *TakeWhileIterator[T]) fillBuffer() {
	if !iter.buffer.IsSome() {
		iter.buffer = iter.iterator.GetNext()
	}
}

func (iter *TakeWhileIterator[T]) HasNext() bool {
	iter.fillBuffer()

	return option.Map(iter.buffer, iter.predicate).OrElse(false)
}

func (iter *TakeWhileIterator[T]) GetNext() option.Option[T] {
	iter.fillBuffer()

	return option.Filter(option.Take(&iter.buffer), iter.predicate)
}
