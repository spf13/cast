// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var errNegativeNotAllowed = errors.New("unable to cast negative value")

type float64EProvider interface {
	Float64() (float64, error)
}

type float64Provider interface {
	Float64() float64
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type unsigned interface {
	uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func toNumber[T number](i any) (T, bool) {
	i = indirect(i)

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

func toNumberE[T number](i any, parseFn func(string) (T, error)) (T, error) {
	n, ok := toNumber[T](i)
	if ok {
		return n, nil
	}

	switch s := i.(type) {
	case string:
		v, err := parseFn(s)
		if err == nil {
			return v, nil
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
	case json.Number:
		v, err := parseFn(string(s))
		if err == nil {
			return v, nil
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
	case float64EProvider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
		}

		v, err := s.Float64()
		if err != nil {
			return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
		}

		return T(v), nil
	case float64Provider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
		}

		return T(s.Float64()), nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
	}
}

func toUnsignedNumber[T unsigned](i any) (T, bool, bool) {
	i = indirect(i)

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

func toUnsignedNumberE[T unsigned](i any, parseFn func(string) (T, error)) (T, error) {
	n, valid, ok := toUnsignedNumber[T](i)
	if ok {
		return n, nil
	}

	if !valid {
		return 0, errNegativeNotAllowed
	}

	switch s := i.(type) {
	case string:
		v, err := parseFn(s)
		if err == nil {
			return v, nil
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
	case json.Number:
		v, err := parseFn(string(s))
		if err == nil {
			return v, nil
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
	case float64EProvider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
		}

		v, err := s.Float64()
		if err != nil {
			return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
		}

		if v < 0 {
			return 0, errNegativeNotAllowed
		}

		return T(v), nil
	case float64Provider:
		if _, ok := any(n).(float64); !ok {
			return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
		}

		v := s.Float64()

		if v < 0 {
			return 0, errNegativeNotAllowed
		}

		return T(v), nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, n)
	}
}

func parseInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s string) (T, error) {
	v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
	if err != nil {
		return 0, err
	}

	return T(v), nil
}

func parseUint[T unsigned](s string) (T, error) {
	v, err := strconv.ParseUint(trimZeroDecimal(s), 0, 0)
	if err != nil {
		return 0, err
	}

	return T(v), nil
}

// ToFloat64E casts an interface to a float64 type.
func ToFloat64E(i any) (float64, error) {
	parseFn := func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	}

	return toNumberE[float64](i, parseFn)
}

// ToFloat32E casts an interface to a float32 type.
func ToFloat32E(i any) (float32, error) {
	parseFn := func(s string) (float32, error) {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, err
		}

		return float32(v), nil
	}

	return toNumberE[float32](i, parseFn)
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
