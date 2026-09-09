//go:build go1.27

// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"strconv"
	"strings"
	"time"
)

// Caster carries conversion configuration for the generic [Caster.To] and
// [Caster.ToE] methods.
//
// The zero value is usable and behaves identically to the package-level cast
// functions (for example [To], [ToE], [ToIntE], [ToTimeE]).
//
// Configuration is applied through value-semantics builder methods such as
// [Caster.WithBase] and [Caster.WithLocation], each of which returns a modified
// copy and leaves the receiver untouched.
type Caster struct {
	base     int
	location *time.Location
}

// Default is a zero-value [Caster].
var Default Caster

// WithBase returns a copy of c that parses integer strings in the given base.
//
// A base of 0 means auto-detection from the string prefix (the default):
// "0x" or "0X" selects base 16, "0o"/"0O" or a leading "0" selects base 8,
// "0b"/"0B" selects base 2, otherwise base 10. See [strconv.ParseInt].
//
// The base is ignored for float and non-integer targets.
//
// WithBase does not validate the base. An invalid base (for example 1, 37 or a
// negative number) is only observed when a string is actually parsed: in that
// case [Caster.ToE] returns the [strconv] error wrapped in the usual cast
// error. Inputs that never reach the string parser (numeric values, booleans,
// nil, ...) convert successfully regardless of the configured base.
func (c Caster) WithBase(base int) Caster {
	c.base = base

	return c
}

// WithLocation returns a copy of c that interprets timezone-less time strings
// in the given location.
//
// Passing nil resets to [time.UTC] (the default); pass [time.Local] explicitly
// for local time.
//
// The location is ignored for non-[time.Time] targets.
func (c Caster) WithLocation(location *time.Location) Caster {
	c.location = location

	return c
}

// locationOrUTC normalizes the zero value (nil) to time.UTC at the point of use.
//
// It must not be stored back: a nil location is what makes Caster{} equivalent
// to ToTimeE, and ToTimeInDefaultLocationE treats nil as local time, not UTC.
func (c Caster) locationOrUTC() *time.Location {
	if c.location == nil {
		return time.UTC
	}

	return c.location
}

type numberKind uint8

const (
	kindSigned numberKind = iota
	kindUnsigned
	kindFloat
)

type numberInfo struct {
	kind    numberKind
	bitSize int
}

// numberInfoOf is the single source of truth for signedness and bit size in
// the Caster code path.
func numberInfoOf[T Number]() numberInfo {
	var t T

	switch any(t).(type) {
	case int:
		return numberInfo{kind: kindSigned, bitSize: strconv.IntSize}
	case int8:
		return numberInfo{kind: kindSigned, bitSize: 8}
	case int16:
		return numberInfo{kind: kindSigned, bitSize: 16}
	case int32:
		return numberInfo{kind: kindSigned, bitSize: 32}
	case int64:
		return numberInfo{kind: kindSigned, bitSize: 64}
	case uint:
		return numberInfo{kind: kindUnsigned, bitSize: strconv.IntSize}
	case uint8:
		return numberInfo{kind: kindUnsigned, bitSize: 8}
	case uint16:
		return numberInfo{kind: kindUnsigned, bitSize: 16}
	case uint32:
		return numberInfo{kind: kindUnsigned, bitSize: 32}
	case uint64:
		return numberInfo{kind: kindUnsigned, bitSize: 64}
	case float32:
		return numberInfo{kind: kindFloat, bitSize: 32}
	case float64:
		return numberInfo{kind: kindFloat, bitSize: 64}
	default:
		// Unreachable: Number has no tilde, so T is exactly one of the cases above.
		return numberInfo{kind: kindSigned, bitSize: strconv.IntSize}
	}
}

// parseNumber parses a string into T using the configured base.
//
// It mirrors the package-level parseInt, parseUint and parseFloat helpers:
// signed and unsigned integers go through trimDecimal (unsigned additionally
// drops a leading "+"), floats are parsed as-is and ignore the base.
func (c Caster) parseNumber[T Number](s string) (T, error) {
	info := numberInfoOf[T]()

	switch info.kind {
	case kindUnsigned:
		v, err := strconv.ParseUint(strings.TrimLeft(trimDecimal(s), "+"), c.base, info.bitSize)
		if err != nil {
			return 0, err
		}

		return T(v), nil
	case kindFloat:
		v, err := strconv.ParseFloat(s, info.bitSize)
		if err != nil {
			return 0, err
		}

		return T(v), nil
	default: // kindSigned
		v, err := strconv.ParseInt(trimDecimal(s), c.base, info.bitSize)
		if err != nil {
			return 0, err
		}

		return T(v), nil
	}
}

// toNumberE routes to the existing number helpers, supplying the configured
// string parser. Everything except string parsing (indirection, aliases,
// json.Number, float64 providers, bool, nil, negative checks) is handled there.
//
// Note: the package-level toNumberE (called for signed integers and floats) is
// a different function from this method; Go keeps method and function
// namespaces separate.
func (c Caster) toNumberE[T Number](i any) (T, error) {
	if numberInfoOf[T]().kind == kindUnsigned {
		return toUnsignedNumberE[T](i, c.parseNumber[T])
	}

	return toNumberE[T](i, c.parseNumber[T])
}

// ToE casts any value to a [Basic] type using the configuration carried by c.
//
// With a zero-value receiver it is equivalent to [ToE].
func (c Caster) ToE[T Basic](i any) (T, error) {
	var t T

	var v any
	var err error

	switch any(t).(type) {
	case string:
		v, err = ToStringE(i)
	case bool:
		v, err = ToBoolE(i)
	case int:
		v, err = c.toNumberE[int](i)
	case int8:
		v, err = c.toNumberE[int8](i)
	case int16:
		v, err = c.toNumberE[int16](i)
	case int32:
		v, err = c.toNumberE[int32](i)
	case int64:
		v, err = c.toNumberE[int64](i)
	case uint:
		v, err = c.toNumberE[uint](i)
	case uint8:
		v, err = c.toNumberE[uint8](i)
	case uint16:
		v, err = c.toNumberE[uint16](i)
	case uint32:
		v, err = c.toNumberE[uint32](i)
	case uint64:
		v, err = c.toNumberE[uint64](i)
	case float32:
		v, err = c.toNumberE[float32](i)
	case float64:
		v, err = c.toNumberE[float64](i)
	case time.Time:
		v, err = ToTimeInDefaultLocationE(i, c.locationOrUTC())
	case time.Duration:
		v, err = ToDurationE(i)
	}

	if err != nil {
		return t, err
	}

	return v.(T), nil
}

// To casts any value to a [Basic] type using the configuration carried by c.
//
// With a zero-value receiver it is equivalent to [To].
func (c Caster) To[T Basic](i any) T {
	v, _ := c.ToE[T](i)

	return v
}
