package result

import (
	"testing"
)
import "errors"

func Test1(t *testing.T) {
	var res = Map(Ok(32), func(val int) int { return val + 1 })
	res = FlatMap(res, func(val int) Result[int] { return Err[int](errors.New("")) })

	_, err := res.Get()
	if err != nil {
		t.Fail()
	}
}
