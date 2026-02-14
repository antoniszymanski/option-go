// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package option

import (
	"fmt"
	"reflect"
	"structs"
	"unsafe"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	jsonv1 "github.com/go-json-experiment/json/v1"
)

type Option[T any] struct {
	_     structs.HostLayout
	valid bool
	value T
}

func Some[T any](value T) Option[T] {
	return Option[T]{valid: true, value: value}
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

func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	return o.valid && f(o.value)
}

func (o Option[T]) IsNoneOr(f func(T) bool) bool {
	return !o.valid || f(o.value)
}

func (o Option[T]) AsSlice() []T {
	if o.valid {
		return unsafe.Slice(&o.value, 1)
	} else {
		return nil
	}
}

func (o Option[T]) Expect(msg string) T {
	if o.valid {
		return o.value
	} else {
		panic(msg)
	}
}

func (o Option[T]) Unwrap() T {
	if o.valid {
		return o.value
	} else {
		panic("called Unwrap on a None value")
	}
}

func (o Option[T]) UnwrapOr(fallback T) T {
	if o.valid {
		return o.value
	} else {
		return fallback
	}
}

func (o Option[T]) UnwrapOrZero() T {
	return o.value
}

func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.valid {
		return o.value
	} else {
		return f()
	}
}

func (o Option[T]) Filter(predicate func(*T) bool) Option[T] {
	if o.valid && predicate(&o.value) {
		return o
	} else {
		return Option[T]{}
	}
}

func (o Option[T]) Inspect(f func(*T)) Option[T] {
	if o.valid {
		f(&o.value)
	}
	return o
}

func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.valid {
		return Option[T]{valid: true, value: f(o.value)}
	} else {
		return Option[T]{}
	}
}

func (o Option[T]) MapOr(fallback T, f func(T) T) T {
	if o.valid {
		return f(o.value)
	} else {
		return fallback
	}
}

func (o Option[T]) MapOrElse(fallback func() T, f func(T) T) T {
	if o.valid {
		return f(o.value)
	} else {
		return fallback()
	}
}

func (o Option[T]) And(other Option[T]) Option[T] {
	if o.valid {
		return other
	} else {
		return Option[T]{}
	}
}

func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.valid {
		return o
	} else {
		return other
	}
}

func (o Option[T]) Xor(other Option[T]) Option[T] {
	switch {
	case o.valid && !other.valid:
		return o
	case !o.valid && other.valid:
		return other
	default:
		return Option[T]{}
	}
}

func (o Option[T]) AndThen(f func(T) Option[T]) Option[T] {
	if o.valid {
		return f(o.value)
	} else {
		return Option[T]{}
	}
}

func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.valid {
		return o
	} else {
		return f()
	}
}

var (
	_ fmt.Stringer   = Option[int]{}
	_ fmt.GoStringer = Option[int]{}
)

func (o Option[T]) String() string {
	if o.valid {
		return fmt.Sprintf("Some(%v)", elem(&o.value))
	} else {
		return "None"
	}
}

func (o Option[T]) GoString() string {
	if o.valid {
		return fmt.Sprintf("option.Some(%#v)", elem(&o.value))
	} else {
		return fmt.Sprintf("option.None[%T]()", elem(&o.value))
	}
}

var (
	_ json.Marshaler       = Option[int]{}
	_ json.Unmarshaler     = &Option[int]{}
	_ json.MarshalerTo     = &Option[int]{}
	_ json.UnmarshalerFrom = &Option[int]{}
)

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.valid {
		return jsonv1.Marshal(&o.value) // avoid boxing on the heap
	} else {
		return []byte("null"), nil
	}
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = Option[T]{}
		return nil
	}
	if err := jsonv1.Unmarshal(data, &o.value); err != nil {
		*o = Option[T]{}
		return err
	}
	o.valid = true
	return nil
}

func (o *Option[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.valid {
		return json.MarshalEncode(enc, &o.value) // avoid boxing on the heap
	} else {
		return enc.WriteToken(jsontext.Null)
	}
}

func (o *Option[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindInvalid:
		*o = Option[T]{}
		_, err := dec.ReadToken()
		return err
	case jsontext.KindNull:
		*o = Option[T]{}
		return nil
	default:
		if err := json.UnmarshalDecode(dec, &o.value); err != nil {
			*o = Option[T]{}
			return err
		}
		o.valid = true
		return nil
	}
}

func (o Option[T]) IsZero() bool {
	if !o.valid {
		return true
	}
	if i, ok := elem(&o.value).(interface{ IsZero() bool }); ok {
		return i.IsZero()
	}
	return false
}

//go:nosplit
func elem[P ~*E, E any](p P) any {
	if p == nil {
		return nil
	}
	typ := reflect.TypeFor[E]()
	return *(*any)(unsafe.Pointer(&iface{
		Type: (*iface)(unsafe.Pointer(&typ)).Data,
		Data: unsafe.Pointer(noEscape(p)),
	}))
}

type iface struct {
	Type, Data unsafe.Pointer
}

//go:nosplit
func noEscape[P ~*E, E any](p P) P {
	x := uintptr(unsafe.Pointer(p))
	return P(unsafe.Pointer(x ^ 0)) //nolint:all
}
