package option_test

import (
	. "dot-parser/option"
	"testing"
)

func TestIsSome(t *testing.T) {
	some := Some(1)
	if !some.IsSome() {
		t.Errorf("some is expected to be some after construction")
	}

	none := None[int]()
	if none.IsSome() {
		t.Errorf("none is expected to be none after construction")
	}
}

func TestIsNone(t *testing.T) {
	some := Some(1)
	if some.IsNone() {
		t.Errorf("some is expected to be some after construction")
	}

	none := None[int]()
	if !none.IsNone() {
		t.Errorf("none is expected to be none after construction")
	}
}

func TestTake(t *testing.T) {
	opt := Some(1)
	val := Take(&opt)
	if !val.IsSome() {
		t.Errorf("new owner was expected to be some")
	}
	if opt.IsSome() {
		t.Errorf("old owner was expected to be none")
	}
}
