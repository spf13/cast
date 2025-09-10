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
	}

	if err != nil {
		return t, err
	}

	return v.(T), nil
}

// Must is a helper that wraps a call to a cast function and panics if the error is non-nil.
func Must[T any](i any, err error) T {
	if err != nil {
		panic(err)
	}

	return i.(T)
}

// To casts any value to a [Basic] type.
func To[T Basic](i any) T {
	v, _ := ToE[T](i)

	return v
}

// ValueSetter is an interface that types can implement to provide custom conversion logic.
// When ToValue is called with a pointer to a type that implements ValueSetter,
// the SetValue method will be called to perform the conversion.
type ValueSetter interface {
	SetValue(any) error
}

// ToValue attempts to convert a value and assign it to the target.
// If target implements ValueSetter, it will use the custom conversion logic.
// Otherwise, it falls back to standard cast conversion based on the target type.
func ToValue(target any, value any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	// Check if target implements ValueSetter
	if setter, ok := target.(ValueSetter); ok {
		return setter.SetValue(value)
	}

	// Use reflection-like approach based on target type
	switch t := target.(type) {
	case *string:
		v, err := ToStringE(value)
		if err != nil {
			return err
		}
		*t = v
	case *bool:
		v, err := ToBoolE(value)
		if err != nil {
			return err
		}
		*t = v
	case *int:
		v, err := toNumberE[int](value, parseInt[int])
		if err != nil {
			return err
		}
		*t = v
	case *int8:
		v, err := toNumberE[int8](value, parseInt[int8])
		if err != nil {
			return err
		}
		*t = v
	case *int16:
		v, err := toNumberE[int16](value, parseInt[int16])
		if err != nil {
			return err
		}
		*t = v
	case *int32:
		v, err := toNumberE[int32](value, parseInt[int32])
		if err != nil {
			return err
		}
		*t = v
	case *int64:
		v, err := toNumberE[int64](value, parseInt[int64])
		if err != nil {
			return err
		}
		*t = v
	case *uint:
		v, err := toUnsignedNumberE[uint](value, parseUint[uint])
		if err != nil {
			return err
		}
		*t = v
	case *uint8:
		v, err := toUnsignedNumberE[uint8](value, parseUint[uint8])
		if err != nil {
			return err
		}
		*t = v
	case *uint16:
		v, err := toUnsignedNumberE[uint16](value, parseUint[uint16])
		if err != nil {
			return err
		}
		*t = v
	case *uint32:
		v, err := toUnsignedNumberE[uint32](value, parseUint[uint32])
		if err != nil {
			return err
		}
		*t = v
	case *uint64:
		v, err := toUnsignedNumberE[uint64](value, parseUint[uint64])
		if err != nil {
			return err
		}
		*t = v
	case *float32:
		v, err := toNumberE[float32](value, parseFloat[float32])
		if err != nil {
			return err
		}
		*t = v
	case *float64:
		v, err := toNumberE[float64](value, parseFloat[float64])
		if err != nil {
			return err
		}
		*t = v
	case *time.Time:
		v, err := ToTimeE(value)
		if err != nil {
			return err
		}
		*t = v
	case *time.Duration:
		v, err := ToDurationE(value)
		if err != nil {
			return err
		}
		*t = v
	default:
		return fmt.Errorf("unsupported target type %T", target)
	}

	return nil
}
