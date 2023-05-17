package option

import (
	"errors"
)

type Option[T any] interface {
	IsSome() bool
	OrElse(other T) T
	Unwrap() T
}

// Some variant
type some[T any] struct {
	val T
}

func Some[T any](value T) Option[T] {
	return &some[T]{val: value}
}

func (opt *some[T]) IsSome() bool {
	return true
}

func (opt *some[T]) OrElse(other T) T {
	return opt.val
}

func (opt *some[T]) Unwrap() T {
	return opt.val
}

// None variant
type none[T any] struct{}

func None[T any]() Option[T] {
	return &none[T]{}
}

func (opt *none[T]) IsSome() bool {
	return false
}

func (opt *none[T]) OrElse(other T) T {
	return other
}

func (opt *none[T]) Unwrap() T {
	panic(errors.New("tried to unwrap a None Option"))
}

// Map and FlatMap
func Filter[T any](opt Option[T], predicate func(T) bool) Option[T] {
	return FlatMap(opt, func(val T) Option[T] {
		if predicate(val) {
			return Some(val)
		} else {
			return None[T]()
		}
	})
}

func Map[U any, V any](opt Option[U], fn func(U) V) Option[V] {
	switch opt := opt.(type) {
	case *some[U]:
		return Some(fn(opt.val))
	case *none[U]:
		return None[V]()
	default:
		panic(nil)
	}
}

func FlatMap[U any, V any](opt Option[U], fn func(U) Option[V]) Option[V] {
	switch opt := opt.(type) {
	case *some[U]:
		return fn(opt.val)
	case *none[U]:
		return None[V]()
	default:
		panic(nil)
	}
}

func Take[T any](opt *Option[T]) Option[T] {
	if (*opt).IsSome() {
		copy := *opt
		*opt = None[T]()
		return copy
	} else {
		return None[T]()
	}
}
