package test

import (
	"dot-parser/iterator"
	"dot-parser/option"
	"testing"
)

type testIterator struct {
	data []rune
	pos  int
}

func makeTestIterator(data string) iterator.Iterator[rune] {
	return &testIterator{data: []rune(data), pos: 0}
}

func (iter *testIterator) Next() option.Option[rune] {
	if iter.pos < len(iter.data) {
		char := iter.data[iter.pos]
		iter.pos += 1
		return option.Some(char)
	} else {
		return option.None[rune]()
	}
}

func testExpected[T comparable](t *testing.T, next option.Option[T], expected T) {
	if next.IsNone() || next.Unwrap() != expected {
		t.Errorf("Expected %v, got %v", expected, next)
	}
}

func testNotExpected[T any](t *testing.T, next option.Option[T]) {
	if next.IsSome() {
		t.Errorf("Expected None, got %v", next)
	}
}
