// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var errNegativeNotAllowed = errors.New("unable to cast negative value")

type float64EProvider interface {
	Float64() (float64, error)
}

type float64Provider interface {
	Float64() float64
}

// Number is a type parameter constraint for functions accepting number types.
//
// It represents the supported number types this package can cast to.
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type integer interface {
	int | int8 | int16 | int32 | int64
}

type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type float interface {
	float32 | float64
}

// ToNumberE casts any value to a [Number] type.
func ToNumberE[T Number](i any) (T, error) {
	return toNumberE[T](i, 0)
}

// ToNumber casts any value to a [Number] type.
func ToNumber[T Number](i any) T {
	v, _ := ToNumberE[T](i)

	return v
}

// ToNumberBaseE casts any value to a [Number] type.
//
// For integer types the second parameter is used as a base instead of auto detection.
func ToNumberBaseE[T Number](i any, base int) (T, error) {
	return toNumberE[T](i, base)
}

// ToNumberBase casts any value to a [Number] type.
//
// For integer types the second parameter is used as a base instead of auto detection.
func ToNumberBase[T Number](i any, base int) T {
	v, _ := ToNumberBaseE[T](i, base)

	return v
}

// toNumberE casts any value to a [Number] type.
//
// For integer types the second parameter is used as a base.
func toNumberE[T Number](i any, base int) (T, error) {
	var t T

	switch any(t).(type) {
	case int:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	case int8:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	case int16:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	case int32:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	case int64:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	case uint:
		return toUnsignedNumberE[T](i, parseNumberBase[T](base))
	case uint8:
		return toUnsignedNumberE[T](i, parseNumberBase[T](base))
	case uint16:
		return toUnsignedNumberE[T](i, parseNumberBase[T](base))
	case uint32:
		return toUnsignedNumberE[T](i, parseNumberBase[T](base))
	case uint64:
		return toUnsignedNumberE[T](i, parseNumberBase[T](base))
	case float32:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	case float64:
		return toSignedNumberE[T](i, parseNumberBase[T](base))
	default:
		return 0, fmt.Errorf("unknown number type: %T", t)
	}
}

// toSignedNumber's semantics differ from other "to" functions.
// It returns false as the second parameter if the conversion fails.
// This is to signal other callers that they should proceed with their own conversions.
func toSignedNumber[T Number](i any) (T, bool) {
	i, _ = indirect(i)

	switch s := i.(type) {
	case T:
		return s, true
	case int:
		return T(s), true
	case int8:
		return T(s), true
	case int16:
		return T(s), true
	case int32:
		return T(s), true
	case int64:
		return T(s), true
	case uint:
		return T(s), true
	case uint8:
		return T(s), true
	case uint16:
		return T(s), true
	case uint32:
		return T(s), true
	case uint64:
		return T(s), true
	case float32:
		return T(s), true
	case float64:
		return T(s), true
	case bool:
		if s {
			return 1, true
		}

		return 0, true
	case nil:
		return 0, true
	case time.Weekday:
		return T(s), true
	case time.Month:
		return T(s), true
	}

	return 0, false
}

func toSignedNumberE[T Number](i any, parseFn func(string) (T, error)) (T, error) {
	n, ok := toSignedNumber[T](i)
	if ok {
		return n, nil
	}

	i, _ = indirect(i)

	switch s := i.(type) {
	case string:
		if s == "" {
			return 0, nil
		}

		v, err := parseFn(s)
		if err != nil {
			return 0, fmt.Errorf(errorMsgWith, i, i, n, err)
		}

		return v, nil
	case json.Number:
		if s == "" {
			return 0, nil
		}

		v, err := parseFn(string(s))
		if err != nil {
			return 0, fmt.Errorf(errorMsgWith, i, i, n, err)
		}

		return v, nil
	case float64EProvider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf(errorMsg, i, i, n)
		}

		v, err := s.Float64()
		if err != nil {
			return 0, fmt.Errorf(errorMsg, i, i, n)
		}

		return T(v), nil
	case float64Provider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf(errorMsg, i, i, n)
		}

		return T(s.Float64()), nil
	default:
		if i, ok := resolveAlias(i); ok {
			return toSignedNumberE(i, parseFn)
		}

		return 0, fmt.Errorf(errorMsg, i, i, n)
	}
}

func toUnsignedNumber[T Number](i any) (T, bool, bool) {
	i, _ = indirect(i)

	switch s := i.(type) {
	case T:
		return s, true, true
	case int:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case int8:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case int16:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case int32:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case int64:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case uint:
		return T(s), true, true
	case uint8:
		return T(s), true, true
	case uint16:
		return T(s), true, true
	case uint32:
		return T(s), true, true
	case uint64:
		return T(s), true, true
	case float32:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case float64:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case bool:
		if s {
			return 1, true, true
		}

		return 0, true, true
	case nil:
		return 0, true, true
	case time.Weekday:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	case time.Month:
		if s < 0 {
			return 0, false, false
		}

		return T(s), true, true
	}

	return 0, true, false
}

func toUnsignedNumberE[T Number](i any, parseFn func(string) (T, error)) (T, error) {
	n, valid, ok := toUnsignedNumber[T](i)
	if ok {
		return n, nil
	}

	i, _ = indirect(i)

	if !valid {
		return 0, errNegativeNotAllowed
	}

	switch s := i.(type) {
	case string:
		if s == "" {
			return 0, nil
		}

		v, err := parseFn(s)
		if err != nil {
			return 0, fmt.Errorf(errorMsgWith, i, i, n, err)
		}

		return v, nil
	case json.Number:
		if s == "" {
			return 0, nil
		}

		v, err := parseFn(string(s))
		if err != nil {
			return 0, fmt.Errorf(errorMsgWith, i, i, n, err)
		}

		return v, nil
	case float64EProvider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf(errorMsg, i, i, n)
		}

		v, err := s.Float64()
		if err != nil {
			return 0, fmt.Errorf(errorMsg, i, i, n)
		}

		if v < 0 {
			return 0, errNegativeNotAllowed
		}

		return T(v), nil
	case float64Provider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf(errorMsg, i, i, n)
		}

		v := s.Float64()

		if v < 0 {
			return 0, errNegativeNotAllowed
		}

		return T(v), nil
	default:
		if i, ok := resolveAlias(i); ok {
			return toUnsignedNumberE(i, parseFn)
		}

		return 0, fmt.Errorf(errorMsg, i, i, n)
	}
}

func parseNumberBase[T Number](base int) func(s string) (T, error) {
	return func(s string) (T, error) {
		var t T

		switch any(t).(type) {
		case int:
			v, err := parseInt[int](base, strconv.IntSize)(s)

			return T(v), err
		case int8:
			v, err := parseInt[int8](base, 8)(s)

			return T(v), err
		case int16:
			v, err := parseInt[int16](base, 16)(s)

			return T(v), err
		case int32:
			v, err := parseInt[int32](base, 32)(s)

			return T(v), err
		case int64:
			v, err := parseInt[int64](base, 64)(s)

			return T(v), err
		case uint:
			v, err := parseUint[uint](base, strconv.IntSize)(s)

			return T(v), err
		case uint8:
			v, err := parseUint[uint8](base, 8)(s)

			return T(v), err
		case uint16:
			v, err := parseUint[uint16](base, 16)(s)

			return T(v), err
		case uint32:
			v, err := parseUint[uint32](base, 32)(s)

			return T(v), err
		case uint64:
			v, err := parseUint[uint64](base, 64)(s)

			return T(v), err
		case float32:
			v, err := strconv.ParseFloat(s, 32)

			return T(v), err
		case float64:
			v, err := strconv.ParseFloat(s, 64)

			return T(v), err

		default:
			return 0, fmt.Errorf("unknown number type: %T", t)
		}
	}
}

func parseInt[T integer](base int, bitSize int) func(s string) (T, error) {
	return func(s string) (T, error) {
		v, err := strconv.ParseInt(trimDecimal(s), base, bitSize)
		if err != nil {
			return 0, err
		}

		return T(v), nil
	}
}

func parseUint[T unsigned](base int, bitSize int) func(s string) (T, error) {
	return func(s string) (T, error) {
		v, err := strconv.ParseUint(strings.TrimLeft(trimDecimal(s), "+"), base, bitSize)
		if err != nil {
			return 0, err
		}

		return T(v), nil
	}
}

func parseFloat[T float](s string) (T, error) {
	var t T

	var v any
	var err error

	switch any(t).(type) {
	case float32:
		n, e := strconv.ParseFloat(s, 32)

		v = float32(n)
		err = e
	case float64:
		n, e := strconv.ParseFloat(s, 64)

		v = float64(n)
		err = e
	}

	return v.(T), err
}

// ToIntE casts an interface to an int type.
func ToIntE(i any) (int, error) {
	return toSignedNumberE[int](i, parseInt[int](0, strconv.IntSize))
}

// ToInt8E casts an interface to an int8 type.
func ToInt8E(i any) (int8, error) {
	return toSignedNumberE[int8](i, parseInt[int8](0, 8))
}

// ToInt16E casts an interface to an int16 type.
func ToInt16E(i any) (int16, error) {
	return toSignedNumberE[int16](i, parseInt[int16](0, 16))
}

// ToInt32E casts an interface to an int32 type.
func ToInt32E(i any) (int32, error) {
	return toSignedNumberE[int32](i, parseInt[int32](0, 32))
}

// ToInt64E casts an interface to an int64 type.
func ToInt64E(i any) (int64, error) {
	return toSignedNumberE[int64](i, parseInt[int64](0, 64))
}

// ToUintE casts an interface to a uint type.
func ToUintE(i any) (uint, error) {
	return toUnsignedNumberE[uint](i, parseUint[uint](0, strconv.IntSize))
}

// ToUint8E casts an interface to a uint type.
func ToUint8E(i any) (uint8, error) {
	return toUnsignedNumberE[uint8](i, parseUint[uint8](0, 8))
}

// ToUint16E casts an interface to a uint16 type.
func ToUint16E(i any) (uint16, error) {
	return toUnsignedNumberE[uint16](i, parseUint[uint16](0, 16))
}

// ToUint32E casts an interface to a uint32 type.
func ToUint32E(i any) (uint32, error) {
	return toUnsignedNumberE[uint32](i, parseUint[uint32](0, 32))
}

// ToUint64E casts an interface to a uint64 type.
func ToUint64E(i any) (uint64, error) {
	return toUnsignedNumberE[uint64](i, parseUint[uint64](0, 64))
}

// ToFloat32E casts an interface to a float32 type.
func ToFloat32E(i any) (float32, error) {
	return toSignedNumberE[float32](i, parseFloat[float32])
}

// ToFloat64E casts an interface to a float64 type.
func ToFloat64E(i any) (float64, error) {
	return toSignedNumberE[float64](i, parseFloat[float64])
}

func trimZeroDecimal(s string) string {
	var foundZero bool
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				return s[:i-1]
			}
		case '0':
			foundZero = true
		default:
			return s
		}
	}
	return s
}

var stringNumberRe = regexp.MustCompile(`^([-+]?\d*)(\.\d*)?$`)

// see [BenchmarkDecimal] for details about the implementation
func trimDecimal(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}

	matches := stringNumberRe.FindStringSubmatch(s)
	if matches != nil {
		// matches[1] is the captured integer part with sign
		s = matches[1]

		// handle special cases
		switch s {
		case "-", "+":
			s += "0"
		case "":
			s = "0"
		}

		return s
	}

	return s
}
