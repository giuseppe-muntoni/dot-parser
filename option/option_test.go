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

func TestFilter(t *testing.T) {
	some := Some(1)
	none := None[int]()

	if opt := some.Filter(func(val int) bool { return val != 1 }); opt.IsSome() {
		t.Errorf("Expected None Option, got %#v", opt)
	}

	if opt := some.Filter(func(val int) bool { return val == 1 }); opt.IsNone() {
		t.Errorf("Expected Some Option, got %#v", opt)
	}

	if opt := none.Filter(func(val int) bool { return val == 1 }); opt.IsSome() {
		t.Errorf("Expected None Option, got %#v", opt)
	}
}

func TestOrElse(t *testing.T) {
	some := Some(1)
	none := None[int]()

	if val := some.OrElse(50); val == 50 {
		t.Errorf("Expected 1, got %#v", val)
	}

	if val := none.OrElse(50); val != 50 {
		t.Errorf("Expected 50, got %#v", val)
	}
}

func TestMap(t *testing.T) {
	some := Some(1)
	none := None[int]()

	if opt := Map(some, func(val int) int { return val + 1 }); opt.IsNone() || opt.Unwrap() != 2 {
		t.Errorf("Expected Some(2), got %#v", opt)
	}

	if opt := Map(none, func(val int) int { return val + 1 }); opt.IsSome() {
		t.Errorf("Expected None, got %#v", opt)
	}
}

func TestFlatMap(t *testing.T) {
	some := Some(1)
	none := None[int]()

	if opt := FlatMap(some, func(val int) Option[int] { return Some(val + 1) }); opt.IsNone() || opt.Unwrap() != 2 {
		t.Errorf("Expected Some(2), got %#v", opt)
	}

	if opt := FlatMap(some, func(val int) Option[int] { return None[int]() }); opt.IsSome() {
		t.Errorf("Expected None, got %#v", opt)
	}

	if opt := FlatMap(none, func(val int) Option[int] { return Some(val + 1) }); opt.IsSome() {
		t.Errorf("Expected None, got %#v", opt)
	}
}
