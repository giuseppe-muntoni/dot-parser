package option

import "errors"

// The Option type can either be a valid value of any given type or empty.
type Option[T any] struct {
	value   T
	isValid bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{value: value, isValid: true}
}

func None[T any]() Option[T] {
	return Option[T]{isValid: false}
}

func (opt Option[T]) IsSome() bool {
	return opt.isValid
}

func (opt Option[T]) IsNone() bool {
	return !opt.IsSome()
}

func (opt Option[T]) OrElse(other T) T {
	if opt.IsSome() {
		return opt.value
	} else {
		return other
	}
}

func (opt Option[T]) Unwrap() T {
	if opt.IsSome() {
		return opt.value
	} else {
		panic(errors.New("tried to unwrap a None Option"))
	}
}

func (opt Option[T]) Filter(predicate func(T) bool) Option[T] {
	if opt.IsSome() && predicate(opt.value) {
		return opt
	} else {
		return None[T]()
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

// Map and FlatMap
func Map[U any, V any](opt Option[U], fn func(U) V) Option[V] {
	if opt.IsSome() {
		return Some(fn(opt.value))
	} else {
		return None[V]()
	}
}

func FlatMap[U any, V any](opt Option[U], fn func(U) Option[V]) Option[V] {
	if opt.IsSome() {
		return fn(opt.value)
	} else {
		return None[V]()
	}
}
