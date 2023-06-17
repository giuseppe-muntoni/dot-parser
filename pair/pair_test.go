package pair_test

import (
	. "dot-parser/pair"
	. "dot-parser/result"
	"errors"
	"testing"
)

func TestMakePair(t *testing.T) {
	pair := NewPair(10, 20)
	left, right := pair.Get()

	if left != 10 {
		t.Errorf("Expected left to be 10")
	}

	if right != 20 {
		t.Errorf("Expected right to be 20")
	}
}

func TestZipResultPair(t *testing.T) {
	ok := Ok(5)
	err := Err[string](errors.New(""))

	ok_pair := Zip(ok, ok)
	if ok_pair.IsErr() {
		t.Errorf("Expected ok_pair to be Ok Result")
	}

	err_pair := Zip(ok, err)
	if err_pair.IsOk() {
		t.Errorf("Expected err_pair to be Err Result, pair: {%#v}", err_pair)
	}

	err_pair2 := Zip(err, ok)
	if err_pair2.IsOk() {
		t.Errorf("Expected err_pair2 to be Err Result, pair: {%#v}", err_pair2)
	}
}
