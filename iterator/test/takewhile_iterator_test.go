package test

import (
	"dot-parser/iterator"
	"testing"
)

func TestTakeWhile(t *testing.T) {
	iter := makeTestIterator("ab\nciao")
	iter = iterator.TakeWhile(iter, func(char rune) bool {
		return char != '\n'
	})

	testExpected(t, iter.Next(), 'a')
	testExpected(t, iter.Next(), 'b')
	testNotExpected(t, iter.Next())
	testNotExpected(t, iter.Next())
}

func TestTakeWhile2(t *testing.T) {
	iter := makeTestIterator("ab\nciao\n\nciao2")
	iter = iterator.TakeWhile(iter, func(char rune) bool {
		return char != '\n'
	})

	testExpected(t, iter.Next(), 'a')
	testExpected(t, iter.Next(), 'b')
	testNotExpected(t, iter.Next())
	testNotExpected(t, iter.Next())
}
