package result

type Result[T any] interface {
	Get() (T, error)
	OrElse(other T) T
	IsOk() bool
	Unwrap() T
}

func Make[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	} else {
		return Ok[T](value)
	}
}

// Ok Variant
type ok[T any] struct {
	val T
}

func (ok ok[T]) Get() (T, error) {
	return ok.val, nil
}

func (ok ok[T]) OrElse(other T) T {
	return ok.val
}

func (ok ok[T]) IsOk() bool {
	return true
}

func Ok[T any](value T) Result[T] {
	return ok[T]{val: value}
}

func (ok ok[T]) Unwrap() T {
	return ok.val
}

// Err Variant
type err[T any] struct {
	err error
}

func (err err[T]) Get() (T, error) {
	var result T
	return result, err.err
}

func (err err[T]) OrElse(other T) T {
	return other
}

func (err err[T]) IsOk() bool {
	return false
}

func Err[T any](error error) Result[T] {
	return err[T]{err: error}
}

func (err err[T]) Unwrap() T {
	panic(err.err)
}

// Map and FlatMap
func Map[U any, V any](result Result[U], fn func(U) V) Result[V] {
	switch res := result.(type) {
	case ok[U]:
		return Ok[V](fn(res.val))
	case err[U]:
		return Err[V](res.err)
	default:
		panic(nil)
	}
}

func FlatMap[U any, V any](result Result[U], fn func(U) Result[V]) Result[V] {
	switch res := result.(type) {
	case ok[U]:
		return fn(res.val)
	case err[U]:
		return Err[V](res.err)
	default:
		panic(nil)
	}
}
