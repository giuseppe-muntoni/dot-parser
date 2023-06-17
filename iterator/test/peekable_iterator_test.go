package test

import (
	. "dot-parser/iterator"
	"testing"
)

func TestFoldWhile(t *testing.T) {
	in_str := "test_string"
	iter := Buffered(makeTestIterator(in_str + "-" + in_str))
	out_str, new_iter := FoldWhile("", iter, func(accum string, char rune) (bool, string) {
		if char != '-' {
			return true, accum + string(char)
		} else {
			return false, accum
		}
	})

	testExpected(t, new_iter.Next(), '-')

	if out_str != in_str {
		t.Errorf("Expected input and output string to match.")
	}
}

func TestSkipWhile(t *testing.T) {
	in_str := "test_string"
	iter := Buffered(makeTestIterator(in_str + "-" + in_str))
	new_iter := SkipWhile(iter, func(char rune) bool {
		return char != '-'
	})

	testExpected(t, new_iter.Next(), '-')
}
