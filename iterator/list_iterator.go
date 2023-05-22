package iterator

import . "dot-parser/option"

type listIterator[T any] struct {
	list     []T
	last_pos int
}

func (iter *listIterator[T]) Next() Option[T] {
	next := None[T]()
	if iter.last_pos < len(iter.list) {
		next, iter.last_pos = Some(iter.list[iter.last_pos]), iter.last_pos+1
	}
	return next
}

func ListIterator[T any](list []T) Iterator[T] {
	return &listIterator[T]{
		list:     list,
		last_pos: 0,
	}
}
