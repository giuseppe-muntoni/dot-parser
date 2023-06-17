package test

import (
	. "dot-parser/iterator"
	"testing"
)

func TestFold(t *testing.T) {
	in_str := "test_string"
	iter := makeTestIterator(in_str)
	out_str := Fold("", iter, func(accum string, char rune) string {
		return accum + string(char)
	})

	if out_str != in_str {
		t.Errorf("Expected input and output string to match.")
	}
}
