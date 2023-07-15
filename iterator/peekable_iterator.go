package iterator

import "dot-parser/option"

//A peekable iterator allows to query for the next element without advancing the iterator.
type PeekableIterator[T any] interface {
	Iterator[T]
	Peek() option.Option[T]
}

//A multi-peekable iterator, similarly to a peekable iterator, allows to query for the next
//n-ths elements without advancing the iterator.
type MultiPeekableIterator[T any] interface {
	PeekableIterator[T]
	PeekNth(n int32) option.Option[T]
}

//Exactly like the Fold operation, but the iteration can be arbitrarly stopped on a condition.
//Other than the accumulated value, an iterator over the skipped elements is returned.
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

//Eagerly advances the iterator while a certan predicate on the elements holds true.
func SkipWhile[T any, ITER interface{ PeekableIterator[T] }](iter ITER, predicate func(T) bool) ITER {
	_, newIter := FoldWhile(nil, iter, func(accum any, elem T) (bool, any) {
		return predicate(elem), nil
	})

	return newIter
}
