// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import "fmt"

// Error is a marker interface implemented by all domain errors returned by the library.
type Error interface {
	castError()
}

// InvalidCastError is returned when a value cannot be cast from its type to the target type.
type InvalidCastError struct {
	value any
	typ   string
}

func (e InvalidCastError) Error() string {
	// TODO: consider adding the value back
	// maybe with a length limit?
	return fmt.Sprintf("unable to cast from type %T to %s", e.value, e.typ)
}

// Implements the [Error] marker interface.
func (InvalidCastError) castError() {}

// RangeError is returned when a value cannot be cast from its type to the target type.
type RangeError struct {
	value any
	typ   string

	min any
	max any
}

func (e RangeError) Error() string {
	// TODO: consider adding the value back
	// maybe with a length limit?
	return fmt.Sprintf("value of type %T is out of range for %s [%v, %v]", e.value, e.typ, e.min, e.max)
}

// Implements the [Error] marker interface.
func (RangeError) castError() {}
