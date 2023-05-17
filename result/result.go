package result

import (
	"dot-parser/option"
	"errors"
)

type Result[T any] interface {
	Get() (T, error)
	OrElse(other T) T
	IsOk() bool
	Unwrap() T
	UnwrapErr() error
}

func Make[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	} else {
		return Ok(value)
	}
}

// Ok Variant
type ok[T any] struct {
	val T
}

func Ok[T any](value T) Result[T] {
	return ok[T]{val: value}
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

func (ok ok[T]) Unwrap() T {
	return ok.val
}

func (ok ok[T]) UnwrapErr() error {
	panic(errors.New("tried to unwrap an error on an Ok Result"))
}

// Err Variant
type err[T any] struct {
	err error
}

func Err[T any](error error) Result[T] {
	return err[T]{err: error}
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

func (err err[T]) Unwrap() T {
	panic(err.err)
}

func (err err[T]) UnwrapErr() error {
	return err.err
}

// Map and FlatMap
func Map[U any, V any](result Result[U], fn func(U) V) Result[V] {
	switch res := result.(type) {
	case ok[U]:
		return Ok(fn(res.val))
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

func Contains[T comparable](result Result[T], value T) bool {
	switch res := result.(type) {
	case ok[T]:
		return res.val == value
	case err[T]:
		return false
	default:
		panic(nil)
	}
}

func ToOption[T any](result Result[T]) option.Option[T] {
	if result.IsOk() {
		return option.Some(result.Unwrap())
	} else {
		return option.None[T]()
	}
}
