package iterator

import "dot-parser/option"

type TakeWhileIterator[T any] struct {
	iterator  BufferedIterator[T]
	predicate func(T) bool
}

func TakeWhile[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return &TakeWhileIterator[T]{
		iterator:  BufferedIterator[T]{iterator: iter, buffer: option.None[T]()},
		predicate: predicate,
	}
}

func (iter *TakeWhileIterator[T]) HasNext() bool {
	return option.Map(iter.iterator.Peek(), iter.predicate).OrElse(false)
}

// todo: fix the invalidation of the iterator after the predicate is false
func (iter *TakeWhileIterator[T]) GetNext() option.Option[T] {
	return option.Filter(iter.iterator.GetNext(), iter.predicate)
}
