package option_test

import (
	"dot-parser/option"
	"testing"
)

func TestTake(t *testing.T) {
	opt := option.Some(1)
	val := option.Take(&opt)
	if !val.IsSome() {
		t.Errorf("new owner was expected to be some")
	}
	if opt.IsSome() {
		t.Errorf("old owner was expected to be none")
	}
}
