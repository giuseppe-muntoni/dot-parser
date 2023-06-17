package test

import (
	. "dot-parser/iterator"
	"testing"
)

func TestEmptyListIterator(t *testing.T) {
	iter := ListIterator([]int{})

	testNotExpected(t, iter.Next())
}

func TestListIterator(t *testing.T) {
	iter := ListIterator([]int{10, 25, 30})

	testExpected(t, iter.Next(), 10)
	testExpected(t, iter.Next(), 25)
	testExpected(t, iter.Next(), 30)
	testNotExpected(t, iter.Next())
}
