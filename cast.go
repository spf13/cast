// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cast provides easy and safe casting in Go.
package cast

import (
	"fmt"
	"time"
)

const errorMsg = "unable to cast %#v of type %T to %T"
const errorMsgWith = "unable to cast %#v of type %T to %T: %w"

// Basic is a type parameter constraint for functions accepting basic types.
//
// It represents the supported basic types this package can cast to.
type Basic interface {
	string | bool | Number | time.Time | time.Duration
}

// ToE casts any value to a [Basic] type.
func ToE[T Basic](i any) (T, error) {
	var t T

	var v any
	var err error

	switch any(t).(type) {
	case string:
		v, err = ToStringE(i)
	case bool:
		v, err = ToBoolE(i)
	case int:
		v, err = toNumberE[int](i, parseInt[int])
	case int8:
		v, err = toNumberE[int8](i, parseInt[int8])
	case int16:
		v, err = toNumberE[int16](i, parseInt[int16])
	case int32:
		v, err = toNumberE[int32](i, parseInt[int32])
	case int64:
		v, err = toNumberE[int64](i, parseInt[int64])
	case uint:
		v, err = toUnsignedNumberE[uint](i, parseUint[uint])
	case uint8:
		v, err = toUnsignedNumberE[uint8](i, parseUint[uint8])
	case uint16:
		v, err = toUnsignedNumberE[uint16](i, parseUint[uint16])
	case uint32:
		v, err = toUnsignedNumberE[uint32](i, parseUint[uint32])
	case uint64:
		v, err = toUnsignedNumberE[uint64](i, parseUint[uint64])
	case float32:
		v, err = toNumberE[float32](i, parseFloat[float32])
	case float64:
		v, err = toNumberE[float64](i, parseFloat[float64])
	case time.Time:
		v, err = ToTimeE(i)
	case time.Duration:
		v, err = ToDurationE(i)
	default:
		return t, fmt.Errorf("unknown basic type: %T", t)
	}

	if err != nil {
		return t, err
	}

	return v.(T), nil
}

// To casts any value to a [Basic] type.
func To[T Basic](i any) T {
	v, _ := ToE[T](i)

	return v
}

// Must panics if there is an error, otherwise returns the value.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// ValueSetter is an interface for types that can provide custom conversion logic.
// When a conversion function encounters a type it cannot handle in its default case,
// it will check if the type implements ValueSetter and use it for custom conversion.
type ValueSetter interface {
	SetValue(any) error
}
