// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package option

import (
	"fmt"
	"reflect"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Option[T any] struct {
	value T
	valid bool
}

// Return a `Some` value containing the given value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, valid: true}
}

// Returns a `None` value of type T.
func None[T any]() Option[T] {
	return Option[T]{}
}

// Returns true if the option is a `Some` value.
func (o Option[T]) IsSome() bool {
	return o.valid
}

// Returns true if the option is a `None` value.
func (o Option[T]) IsNone() bool {
	return !o.valid
}

// Returns the contained `Some` value or panics with a custom panic message provided by msg.
func (o Option[T]) Expect(msg string) T {
	if o.valid {
		return o.value
	} else {
		panic(msg)
	}
}

// Returns the contained `Some` value or panics with a generic message.
func (o Option[T]) Unwrap() T {
	if o.valid {
		return o.value
	} else {
		panic("called `Option.Unwrap()` on a `None` value")
	}
}

// Returns the contained `Some` value or a provided fallback.
func (o Option[T]) UnwrapOr(fallback T) T {
	if o.valid {
		return o.value
	} else {
		return fallback
	}
}

// Returns the contained `Some` value or computes it from a closure.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.valid {
		return o.value
	} else {
		return fn()
	}
}

// Returns the contained `Some` value or the zero value for type T.
func (o Option[T]) UnwrapOrZero() T {
	return o.value
}

func (o Option[T]) String() string {
	if o.valid {
		return fmt.Sprintf("Some(%v)", o.value)
	} else {
		return "None"
	}
}

func (o Option[T]) GoString() string {
	if o.valid {
		return fmt.Sprintf("option.Some(%#v)", o.value)
	} else {
		return fmt.Sprintf("option.None[%T]()", o.value)
	}
}

var (
	_ json.MarshalerTo     = (*Option[int])(nil)
	_ json.UnmarshalerFrom = (*Option[int])(nil)
)

func (o *Option[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.valid {
		return json.MarshalEncode(enc, o.value)
	} else {
		return enc.WriteToken(jsontext.Null)
	}
}

func (o *Option[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() != 'n' {
		err := json.UnmarshalDecode(dec, &o.value)
		if err == nil {
			o.valid = true
		}
		return err
	} else {
		_, err := dec.ReadToken()
		if err == nil {
			*o = None[T]()
		}
		return err
	}
}

func (o Option[T]) IsZero() bool {
	if o.valid {
		if i, ok := any(o.value).(interface{ IsZero() bool }); ok {
			return i.IsZero()
		}
		return reflect.ValueOf(o.value).IsZero()
	} else {
		return true
	}
}
