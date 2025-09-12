// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package option

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	jsonv1 "github.com/go-json-experiment/json/v1"
)

type Option[T any] struct {
	value T
	valid bool
}

var (
	_ fmt.Stringer         = (*Option[int])(nil)
	_ fmt.GoStringer       = (*Option[int])(nil)
	_ json.Marshaler       = (*Option[int])(nil)
	_ json.Unmarshaler     = (*Option[int])(nil)
	_ json.MarshalerTo     = (*Option[int])(nil)
	_ json.UnmarshalerFrom = (*Option[int])(nil)
)

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

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.valid {
		return jsonv1.Marshal(noEscape(&o.value)) // avoid boxing on the heap
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

func (o *Option[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) (err error) {
	if kind := dec.PeekKind(); isKindValid(kind) && kind != 'n' {
		if err = json.UnmarshalDecode(dec, &o.value); err == nil {
			o.valid = true
		} else {
			*o = Option[T]{}
		}
	} else {
		_, err = dec.ReadToken()
		*o = Option[T]{}
	}
	return
}

func isKindValid(k jsontext.Kind) bool {
	return k == 'n' || k == 'f' || k == 't' || k == '"' || k == '0' || k == '{' || k == '}' || k == '[' || k == ']'
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
