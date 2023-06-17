package result

import (
	"dot-parser/option"
	"errors"
	"testing"
)

func TestMakeOk(t *testing.T) {
	res := Ok(50)

	if res.IsErr() {
		t.Fatalf("The result was expected to be Ok, got instead %#v", res)
	}

	if res.Unwrap() != 50 {
		t.Fatalf("Result expected to contain 50, got instead %#v", res)
	}
}

func TestMakeErr(t *testing.T) {
	res := Err[int](errors.New("err"))

	if res.IsOk() {
		t.Fatalf("The result was expected to be Err, got instead %#v", res)
	}

	if res.UnwrapErr().Error() != "err" {
		t.Fatalf("Result error expected to be 'err', got instead %#v", res)
	}
}

func TestMakeNilErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("The code did not panic")
		}
	}()

	Err[int](nil)
}

func TestMakeResult(t *testing.T) {
	res := Make(50, errors.New("err"))

	if res.IsOk() {
		t.Errorf("The result was expected to be Err, got instead %#v", res)
	}

	res = Make(50, nil)

	if res.IsErr() {
		t.Errorf("The result was expected to be Ok, got instead %#v", res)
	}
}

func TestGetMethod(t *testing.T) {
	res := Make(50, errors.New("err"))

	if _, err := res.Get(); err == nil {
		t.Errorf("The result was expected to be Err, got instead %#v", res)
	}

	res = Make(50, nil)

	if value, err := res.Get(); err != nil || value != 50 {
		t.Errorf("The result was expected to be Ok, got instead %#v", res)
	}
}

func TestOrElseMethod(t *testing.T) {
	res := Make(50, errors.New("err"))

	if res.OrElse(51) != 51 {
		t.Errorf("Expected 51, result was %#v", res)
	}

	res = Make(50, nil)

	if res.OrElse(51) != 50 {
		t.Errorf("Expected 50, result was %#v", res)
	}
}

func TestMap(t *testing.T) {
	res := Map(Ok(32), func(val int) int { return val + 1 })

	if res.Unwrap() != 33 {
		t.Errorf("Expected 33, result was %#v", res)
	}

	res = Map(Err[int](errors.New("err")), func(val int) int { return val + 1 })

	if res.IsOk() {
		t.Errorf("Expected Err Result, result was %#v", res)
	}
}

func TestFlatMap(t *testing.T) {
	res := FlatMap(Ok(32), func(val int) Result[int] { return Err[int](errors.New("")) })

	if res.IsOk() {
		t.Errorf("Expected Err Result, result was %#v", res)
	}

	res = FlatMap(Ok(32), func(val int) Result[int] { return Ok(val + 1) })

	if res.Unwrap() != 33 {
		t.Errorf("Expected 33, result was %#v", res)
	}
}

func TestContains(t *testing.T) {
	res := Make(50, errors.New("err"))

	if Contains(res, 50) {
		t.Errorf("Expected not contained, result was %#v", res)
	}

	res = Make(50, nil)

	if !Contains(res, 50) {
		t.Errorf("Expected contained, result was %#v", res)
	}
}

func TestConvertToOption(t *testing.T) {
	res := Make(50, errors.New("err"))

	if ToOption(res).IsSome() {
		t.Errorf("Expected None, result was %#v", res)
	}

	res = Make(50, nil)

	if ToOption(res).IsNone() {
		t.Fatalf("Expected Some, result was %#v", res)
	}

	if ToOption(res).Unwrap() != 50 {
		t.Errorf("Expected 50, result was %#v", res)
	}
}

func TestConvertFromOption(t *testing.T) {
	err := errors.New("err")

	if res := FromOption(option.Some(50), err); res.IsErr() {
		t.Fatalf("Expected Ok Result, result was %#v", res)
	}

	if res := FromOption(option.Some(50), err); res.Unwrap() != 50 {
		t.Errorf("Expected 50, result was %#v", res)
	}

	if res := FromOption(option.None[int](), err); res.IsOk() {
		t.Fatalf("Expected Err Result, result was %#v", res)
	}

	if res := FromOption(option.None[int](), err); res.UnwrapErr().Error() != "err" {
		t.Errorf("Expected error to be 'err', result was %#v", res)
	}
}
