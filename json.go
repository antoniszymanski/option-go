// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package option

import (
	"reflect"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

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
