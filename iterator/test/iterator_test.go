package test_test

import (
	"dot-parser/iterator"
	"dot-parser/option"
	"testing"
)

type testIterator struct {
	data []rune
	pos  int
}

func iter(data string) iterator.Iterator[rune] {
	return &testIterator{data: []rune(data), pos: 0}
}

func (iter *testIterator) HasNext() bool {
	return iter.pos < len(iter.data)
}

func (iter *testIterator) GetNext() option.Option[rune] {
	if iter.HasNext() {
		char := iter.data[iter.pos]
		iter.pos += 1
		return option.Some(char)
	} else {
		return option.None[rune]()
	}
}

func TestTakeWhile(t *testing.T) {
	iter := iter("ab\nciao")
	iter = iterator.TakeWhile(iter, func(char rune) bool {
		return char != '\n'
	})

	if next := iter.GetNext(); !next.IsSome() || next.Unwrap() != 'a' {
		t.Errorf("Expected a, got %v", next)
	}

	if next := iter.GetNext(); !next.IsSome() || next.Unwrap() != 'b' {
		t.Errorf("Expected b, got %v", next)
	}

	if next := iter.GetNext(); next.IsSome() {
		t.Errorf("Expected None, got %v", next)
	}
}

func TestTakeWhile2(t *testing.T) {
	iter := iter("ab\nciao")
	iter = iterator.TakeWhile(iter, func(char rune) bool {
		return char != '\n'
	})

	if next := iter.GetNext(); !next.IsSome() || next.Unwrap() != 'a' {
		t.Errorf("Expected a, got %v", next)
	}

	if next := iter.GetNext(); !next.IsSome() || next.Unwrap() != 'b' {
		t.Errorf("Expected b, got %v", next)
	}

	if next := iter.GetNext(); next.IsSome() {
		t.Errorf("Expected None, got %v", next)
	}
}
