// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package option

import (
	"fmt"
)

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
