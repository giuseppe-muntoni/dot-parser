package iterator

import "dot-parser/option"

type PeekableIterator[T any] interface {
	Iterator[T]
	Peek() option.Option[T]
}

type MultiPeekableIterator[T any] interface {
	PeekableIterator[T]
	PeekNth(n int32) option.Option[T]
}

func FoldWhile[A any, T any](accumulator A, iter PeekableIterator[T], fn func(A, T) (bool, A)) (A, PeekableIterator[T]) {
	keepIterating := true
	for {
		if next := iter.Peek(); next.IsSome() {
			keepIterating, accumulator = fn(accumulator, next.Unwrap())
			if !keepIterating {
				return accumulator, iter
			} else {
				iter.Next()
			}
		} else {
			return accumulator, iter
		}
	}
}

func SkipWhile[T any](iter PeekableIterator[T], predicate func(T) bool) PeekableIterator[T] {
	_, newIter := FoldWhile(nil, iter, func(accum any, elem T) (bool, any) {
		return predicate(elem), nil
	})

	return newIter
}
