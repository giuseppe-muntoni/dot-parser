package iterator

import "dot-parser/option"

type BufferedIterator[T any] struct {
	iterator Iterator[T]
	buffer   option.Option[T]
}

func Buffered[T any](iter Iterator[T], value option.Option[T]) PeekableIterator[T] {
	return &BufferedIterator[T]{iterator: iter, buffer: value}
}

func (iter *BufferedIterator[T]) fillBuffer() {
	if !iter.buffer.IsSome() {
		iter.buffer = iter.iterator.GetNext()
	}
}

func (iter *BufferedIterator[T]) HasNext() bool {
	iter.fillBuffer()
	return iter.buffer.IsSome()
}

func (iter *BufferedIterator[T]) GetNext() option.Option[T] {
	iter.fillBuffer()
	return option.Take(&iter.buffer)
}

func (iter *BufferedIterator[T]) Peek() option.Option[T] {
	iter.fillBuffer()
	return iter.buffer
}
