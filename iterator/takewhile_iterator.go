package iterator

import "dot-parser/option"

type TakeWhileIterator[T any] struct {
	iterator  BufferedIterator[T]
	predicate func(T) bool
	finished  bool
}

func TakeWhile[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return &TakeWhileIterator[T]{
		iterator:  BufferedIterator[T]{iterator: iter},
		predicate: predicate,
		finished:  false,
	}
}

func (iter *TakeWhileIterator[T]) Next() option.Option[T] {
	if iter.finished {
		return option.None[T]()
	} else {
		next := iter.iterator.Next()
		valid := option.Map(next, iter.predicate).OrElse(false)
		if valid == false {
			iter.finished = true
			return option.None[T]()
		} else {
			return next
		}
	}
}
