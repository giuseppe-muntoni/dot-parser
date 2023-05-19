package test

import (
	"dot-parser/iterator"
	"testing"
)

func TestPeekNth(t *testing.T) {
	iter := iterator.Buffered(makeTestIterator("abcdefgh"))

	testExpected(t, iter.Peek(), 'a')
	testExpected(t, iter.PeekNth(1), 'a')
	testExpected(t, iter.PeekNth(3), 'c')
	testExpected(t, iter.PeekNth(3), 'c')
	testExpected(t, iter.PeekNth(2), 'b')
}

func TestNextAfterPeek(t *testing.T) {
	iter := iterator.Buffered(makeTestIterator("abcdefgh"))

	testExpected(t, iter.Peek(), 'a')
	testExpected(t, iter.PeekNth(3), 'c')
	testExpected(t, iter.Next(), 'a')
	testExpected(t, iter.Peek(), 'b')
	testExpected(t, iter.PeekNth(2), 'c')
}

func TestPeekOnNone(t *testing.T) {
	iter := iterator.Buffered(makeTestIterator("abc"))
	testExpected(t, iter.Peek(), 'a')
	testExpected(t, iter.PeekNth(3), 'c')
	testNotExpected(t, iter.PeekNth(4))

	testExpected(t, iter.Next(), 'a')
	testExpected(t, iter.PeekNth(2), 'c')
	testNotExpected(t, iter.PeekNth(3))
	testNotExpected(t, iter.PeekNth(4))
}
