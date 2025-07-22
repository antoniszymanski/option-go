// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package option

type Option[T any] struct {
	value T
	valid bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{value: value, valid: true}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) IsSome() bool {
	return o.valid
}

func (o Option[T]) IsNone() bool {
	return !o.valid
}

func (o Option[T]) Unwrap() T {
	if o.valid {
		return o.value
	} else {
		panic("called `Option.Unwrap()` on a `None` value")
	}
}

func (o Option[T]) UnwrapOr(fallback T) T {
	if o.valid {
		return o.value
	} else {
		return fallback
	}
}

func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.valid {
		return o.value
	} else {
		return fn()
	}
}

func (o Option[T]) UnwrapOrZero() T {
	return o.value
}

func (o Option[T]) UnwrapUnchecked() T {
	return o.value
}
