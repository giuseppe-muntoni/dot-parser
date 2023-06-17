package result

import (
	"dot-parser/option"
	"errors"
)

type Result[T any] struct {
	value T
	err   error
}

// Constructors
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

func Err[T any](err error) Result[T] {
	if err == nil {
		panic("Please provide the error when creating an Err Result")
	}

	return Result[T]{err: err}
}

func Make[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	} else {
		return Ok(value)
	}
}

// Methods
func (res Result[T]) Get() (T, error) {
	return res.value, res.err
}

func (res Result[T]) IsOk() bool {
	return res.err == nil
}

func (res Result[T]) IsErr() bool {
	return !res.IsOk()
}

func (res Result[T]) OrElse(other T) T {
	if res.IsOk() {
		return res.value
	} else {
		return other
	}
}

func (res Result[T]) Unwrap() T {
	if res.IsOk() {
		return res.value
	} else {
		panic(res.err)
	}
}

func (res Result[T]) UnwrapErr() error {
	if res.IsOk() {
		panic(errors.New("tried to unwrap the error on an Ok Result"))
	} else {
		return res.err
	}
}

// Map and FlatMap
func Map[U any, V any](result Result[U], fn func(U) V) Result[V] {
	if result.IsOk() {
		return Ok(fn(result.value))
	} else {
		return Err[V](result.err)
	}
}

func FlatMap[U any, V any](result Result[U], fn func(U) Result[V]) Result[V] {
	if result.IsOk() {
		return fn(result.value)
	} else {
		return Err[V](result.err)
	}
}

func Contains[T comparable](result Result[T], value T) bool {
	if result.IsOk() {
		return result.value == value
	} else {
		return false
	}
}

func ToOption[T any](result Result[T]) option.Option[T] {
	if result.IsOk() {
		return option.Some(result.Unwrap())
	} else {
		return option.None[T]()
	}
}

func FromOption[T any](option option.Option[T], err error) Result[T] {
	if option.IsSome() {
		return Ok(option.Unwrap())
	} else {
		return Err[T](err)
	}
}
