package iterator

import "dot-parser/option"

type BufferedIterator[T any] struct {
	iterator Iterator[T]
	buffer   []T
}

//A Buffered Iterator implements the multi-peekable and peekable iterator interfaces from any basic iterator.
//Note that while the peek operations do not advance the iterator position, they advance the internal iterator.
func Buffered[T any](iter Iterator[T]) MultiPeekableIterator[T] {
	return &BufferedIterator[T]{iterator: iter, buffer: make([]T, 0)}
}

func (iter *BufferedIterator[T]) Next() option.Option[T] {
	if len(iter.buffer) > 0 {
		var res option.Option[T]
		iter.buffer, res = iter.buffer[1:], option.Some(iter.buffer[0])
		return res
	} else {
		return iter.iterator.Next()
	}
}

func (iter *BufferedIterator[T]) Peek() option.Option[T] {
	return iter.PeekNth(1)
}

func (iter *BufferedIterator[T]) PeekNth(n int32) option.Option[T] {
	for int32(len(iter.buffer)) < n {
		next := iter.iterator.Next()
		if next.IsSome() {
			iter.buffer = append(iter.buffer, next.Unwrap())
		} else {
			return option.None[T]()
		}
	}

	return option.Some(iter.buffer[n-1])
}
