// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	var t T

	switch any(t).(type) {
	case int:
		return toNumberE[T](i, parseNumber[T])
	case int8:
		return toNumberE[T](i, parseNumber[T])
	case int16:
		return toNumberE[T](i, parseNumber[T])
	case int32:
		return toNumberE[T](i, parseNumber[T])
	case int64:
		return toNumberE[T](i, parseNumber[T])
	case uint:
		return toUnsignedNumberE[T](i, parseNumber[T])
	case uint8:
		return toUnsignedNumberE[T](i, parseNumber[T])
	case uint16:
		return toUnsignedNumberE[T](i, parseNumber[T])
	case uint32:
		return toUnsignedNumberE[T](i, parseNumber[T])
	case uint64:
		return toUnsignedNumberE[T](i, parseNumber[T])
	case float32:
		return toNumberE[T](i, parseNumber[T])
	case float64:
		return toNumberE[T](i, parseNumber[T])
	default:
		return 0, fmt.Errorf("unknown number type: %T", t)
	}
}

// ToNumber casts any value to a [Number] type.
func ToNumber[T Number](i any) T {
	v, _ := ToNumberE[T](i)

	return v
}

// toNumber's semantics differ from other "to" functions.
// It returns false as the second parameter if the conversion fails.
// This is to signal other callers that they should proceed with their own conversions.
func toNumber[T Number](i any) (T, bool) {
	i, _ = indirect(i)

	switch s := i.(type) {
	case T:
		return s, true
	case int:
		return narrowSigned[T](int64(s))
	case int8:
		return narrowSigned[T](int64(s))
	case int16:
		return narrowSigned[T](int64(s))
	case int32:
		return narrowSigned[T](int64(s))
	case int64:
		return narrowSigned[T](s)
	case uint:
		return narrowUnsigned[T](uint64(s))
	case uint8:
		return narrowUnsigned[T](uint64(s))
	case uint16:
		return narrowUnsigned[T](uint64(s))
	case uint32:
		return narrowUnsigned[T](uint64(s))
	case uint64:
		return narrowUnsigned[T](s)
	case float32:
		return narrowFloat[T](float64(s))
	case float64:
		return narrowFloat[T](s)
	case bool:
		if s {
			return 1, true
		}

		return 0, true
	case nil:
		return 0, true
	case time.Weekday:
		return narrowSigned[T](int64(s))
	case time.Month:
		return narrowSigned[T](int64(s))
	}

	return 0, false
}

func toNumberE[T Number](i any, parseFn func(string) (T, error)) (T, error) {
	n, ok := toNumber[T](i)
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
			return toNumberE(i, parseFn)
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
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case int8:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case int16:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case int32:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case int64:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case uint:
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case uint8:
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case uint16:
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case uint32:
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case uint64:
		n, ok := narrowUnsigned[T](s)
		return n, true, ok
	case float32:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowFloat[T](float64(s))
		return n, true, ok
	case float64:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowFloat[T](s)
		return n, true, ok
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
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
	case time.Month:
		if s < 0 {
			return 0, false, false
		}
		n, ok := narrowUnsigned[T](uint64(s))
		return n, true, ok
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

func parseNumber[T Number](s string) (T, error) {
	var t T

	switch any(t).(type) {
	case int:
		v, err := parseInt[int](s)

		return T(v), err
	case int8:
		v, err := parseInt[int8](s)

		return T(v), err
	case int16:
		v, err := parseInt[int16](s)

		return T(v), err
	case int32:
		v, err := parseInt[int32](s)

		return T(v), err
	case int64:
		v, err := parseInt[int64](s)

		return T(v), err
	case uint:
		v, err := parseUint[uint](s)

		return T(v), err
	case uint8:
		v, err := parseUint[uint8](s)

		return T(v), err
	case uint16:
		v, err := parseUint[uint16](s)

		return T(v), err
	case uint32:
		v, err := parseUint[uint32](s)

		return T(v), err
	case uint64:
		v, err := parseUint[uint64](s)

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

func parseInt[T integer](s string) (T, error) {
	v, err := strconv.ParseInt(trimDecimal(s), 0, intBitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(v), nil
}

func parseUint[T unsigned](s string) (T, error) {
	v, err := strconv.ParseUint(strings.TrimLeft(trimDecimal(s), "+"), 0, uintBitSize[T]())
	if err != nil {
		return 0, err
	}

	return T(v), nil
}

func intBitSize[T integer]() int {
	var zero T
	switch any(zero).(type) {
	case int:
		return strconv.IntSize
	case int8:
		return 8
	case int16:
		return 16
	case int32:
		return 32
	case int64:
		return 64
	default:
		return 0
	}
}

func uintBitSize[T unsigned]() int {
	var zero T
	switch any(zero).(type) {
	case uint:
		return strconv.IntSize
	case uint8:
		return 8
	case uint16:
		return 16
	case uint32:
		return 32
	case uint64:
		return 64
	default:
		return 0
	}
}

// narrowSigned converts v to T. An integer target that cannot hold v fails.
// Float targets keep the language conversion.
func narrowSigned[T Number](v int64) (T, bool) {
	var zero T
	switch any(zero).(type) {
	case float32:
		return any(float32(v)).(T), true
	case float64:
		return any(float64(v)).(T), true
	case int:
		if v < int64(math.MinInt) || v > int64(math.MaxInt) {
			return 0, false
		}
		return any(int(v)).(T), true
	case int8:
		if v < math.MinInt8 || v > math.MaxInt8 {
			return 0, false
		}
		return any(int8(v)).(T), true
	case int16:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return 0, false
		}
		return any(int16(v)).(T), true
	case int32:
		if v < math.MinInt32 || v > math.MaxInt32 {
			return 0, false
		}
		return any(int32(v)).(T), true
	case int64:
		return any(v).(T), true
	case uint:
		if v < 0 || uint64(v) > uint64(^uint(0)) {
			return 0, false
		}
		return any(uint(v)).(T), true
	case uint8:
		if v < 0 || v > math.MaxUint8 {
			return 0, false
		}
		return any(uint8(v)).(T), true
	case uint16:
		if v < 0 || v > math.MaxUint16 {
			return 0, false
		}
		return any(uint16(v)).(T), true
	case uint32:
		if v < 0 || v > math.MaxUint32 {
			return 0, false
		}
		return any(uint32(v)).(T), true
	case uint64:
		if v < 0 {
			return 0, false
		}
		return any(uint64(v)).(T), true
	default:
		return 0, false
	}
}

func narrowUnsigned[T Number](v uint64) (T, bool) {
	if v > math.MaxInt64 {
		var zero T
		switch any(zero).(type) {
		case uint64:
			return any(v).(T), true
		case uint:
			if v <= uint64(^uint(0)) {
				return any(uint(v)).(T), true
			}
			return 0, false
		default:
			return 0, false
		}
	}
	return narrowSigned[T](int64(v))
}

func narrowFloat[T Number](v float64) (T, bool) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}

	var zero T
	switch any(zero).(type) {
	case float32:
		return any(float32(v)).(T), true
	case float64:
		return any(v).(T), true
	case uint64:
		if v < 0 || v >= float64(math.MaxUint64) {
			return 0, false
		}
		return any(uint64(v)).(T), true
	case uint:
		if strconv.IntSize == 64 && v >= 0 && v < float64(math.MaxUint64) {
			return any(uint(v)).(T), true
		}
		if v < 0 || v >= float64(math.MaxInt64) {
			return 0, false
		}
		return narrowSigned[T](int64(v))
	default:
		if v < float64(math.MinInt64) || v >= float64(math.MaxInt64) {
			return 0, false
		}
		return narrowSigned[T](int64(v))
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

// ToFloat64E casts an interface to a float64 type.
func ToFloat64E(i any) (float64, error) {
	return toNumberE[float64](i, parseFloat[float64])
}

// ToFloat32E casts an interface to a float32 type.
func ToFloat32E(i any) (float32, error) {
	return toNumberE[float32](i, parseFloat[float32])
}

// ToInt64E casts an interface to an int64 type.
func ToInt64E(i any) (int64, error) {
	return toNumberE[int64](i, parseInt[int64])
}

// ToInt32E casts an interface to an int32 type.
func ToInt32E(i any) (int32, error) {
	return toNumberE[int32](i, parseInt[int32])
}

// ToInt16E casts an interface to an int16 type.
func ToInt16E(i any) (int16, error) {
	return toNumberE[int16](i, parseInt[int16])
}

// ToInt8E casts an interface to an int8 type.
func ToInt8E(i any) (int8, error) {
	return toNumberE[int8](i, parseInt[int8])
}

// ToIntE casts an interface to an int type.
func ToIntE(i any) (int, error) {
	return toNumberE[int](i, parseInt[int])
}

// ToUintE casts an interface to a uint type.
func ToUintE(i any) (uint, error) {
	return toUnsignedNumberE[uint](i, parseUint[uint])
}

// ToUint64E casts an interface to a uint64 type.
func ToUint64E(i any) (uint64, error) {
	return toUnsignedNumberE[uint64](i, parseUint[uint64])
}

// ToUint32E casts an interface to a uint32 type.
func ToUint32E(i any) (uint32, error) {
	return toUnsignedNumberE[uint32](i, parseUint[uint32])
}

// ToUint16E casts an interface to a uint16 type.
func ToUint16E(i any) (uint16, error) {
	return toUnsignedNumberE[uint16](i, parseUint[uint16])
}

// ToUint8E casts an interface to a uint type.
func ToUint8E(i any) (uint8, error) {
	return toUnsignedNumberE[uint8](i, parseUint[uint8])
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
