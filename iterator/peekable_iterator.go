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

func FoldWhile[A any, T any, ITER interface{ PeekableIterator[T] }](accumulator A, iter ITER, fn func(A, T) (bool, A)) (A, ITER) {
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

func SkipWhile[T any, ITER interface{ PeekableIterator[T] }](iter ITER, predicate func(T) bool) ITER {
	_, newIter := FoldWhile(nil, iter, func(accum any, elem T) (bool, any) {
		return predicate(elem), nil
	})

	return newIter
}
