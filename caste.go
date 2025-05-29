// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"reflect"
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

// ToBoolE casts an interface to a bool type.
func ToBoolE(i interface{}) (bool, error) {
	i = indirect(i)

	switch b := i.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int:
		return b != 0, nil
	case int64:
		return b != 0, nil
	case int32:
		return b != 0, nil
	case int16:
		return b != 0, nil
	case int8:
		return b != 0, nil
	case uint:
		return b != 0, nil
	case uint64:
		return b != 0, nil
	case uint32:
		return b != 0, nil
	case uint16:
		return b != 0, nil
	case uint8:
		return b != 0, nil
	case float64:
		return b != 0, nil
	case float32:
		return b != 0, nil
	case time.Duration:
		return b != 0, nil
	case string:
		return strconv.ParseBool(i.(string))
	case json.Number:
		v, err := ToInt64E(b)
		if err == nil {
			return v != 0, nil
		}
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type unsigned interface {
	uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func toNumber[T number](i any) (T, bool) {
	i = indirect(i)

	intv, ok := toInt(i)
	if ok {
		return T(intv), true
	}

	switch s := i.(type) {
	case T:
		return s, true
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

	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, false, false
		}

		return T(intv), true, true
	}

	switch s := i.(type) {
	case T:
		return s, true, true
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

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// ToStringE casts an interface to a string type.
func ToStringE(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		v := reflect.ValueOf(s)
		for v.Kind() == reflect.Ptr && !v.IsNil() {
			return ToStringE(v.Elem().Interface())
		}

		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

func toMapE[K comparable, V any](i any, keyFn func(any) K, valFn func(any) V) (map[K]V, error) {
	m := map[K]V{}

	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, m)
	}

	switch v := i.(type) {
	case map[K]V:
		return v, nil

	case map[K]any:
		for k, val := range v {
			m[k] = valFn(val)
		}

		return m, nil

	case map[any]V:
		for k, val := range v {
			m[keyFn(k)] = val
		}

		return m, nil

	case map[any]any:
		for k, val := range v {
			m[keyFn(k)] = valFn(val)
		}

		return m, nil

	case string:
		err := jsonStringToObject(v, &m)

		return m, err

	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, m)
	}
}

func toStringMapE[T any](i any, fn func(any) T) (map[string]T, error) {
	return toMapE(i, ToString, fn)
}

// ToStringMapStringE casts an interface to a map[string]string type.
func ToStringMapStringE(i any) (map[string]string, error) {
	return toStringMapE(i, ToString)
}

// ToStringMapStringSliceE casts an interface to a map[string][]string type.
func ToStringMapStringSliceE(i interface{}) (map[string][]string, error) {
	m := map[string][]string{}

	switch v := i.(type) {
	case map[string][]string:
		return v, nil
	case map[string][]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[string]string:
		for k, val := range v {
			m[ToString(k)] = []string{val}
		}
	case map[string]interface{}:
		for k, val := range v {
			switch vt := val.(type) {
			case []interface{}:
				m[ToString(k)] = ToStringSlice(vt)
			case []string:
				m[ToString(k)] = vt
			default:
				m[ToString(k)] = []string{ToString(val)}
			}
		}
		return m, nil
	case map[interface{}][]string:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}][]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			key, err := ToStringE(k)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			value, err := ToStringSliceE(val)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			m[key] = value
		}
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
	}
	return m, nil
}

// ToStringMapBoolE casts an interface to a map[string]bool type.
func ToStringMapBoolE(i interface{}) (map[string]bool, error) {
	return toStringMapE(i, ToBool)
}

// ToStringMapE casts an interface to a map[string]interface{} type.
func ToStringMapE(i interface{}) (map[string]interface{}, error) {
	fn := func(i any) any { return i }

	return toStringMapE(i, fn)
}

func toStringMapIntE[T int | int64](i any, fn func(any) T, fnE func(any) (T, error)) (map[string]T, error) {
	m := map[string]T{}

	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, m)
	}

	switch v := i.(type) {
	case map[string]T:
		return v, nil

	case map[string]any:
		for k, val := range v {
			m[k] = fn(val)
		}

		return m, nil

	case map[any]T:
		for k, val := range v {
			m[ToString(k)] = val
		}

		return m, nil

	case map[any]any:
		for k, val := range v {
			m[ToString(k)] = fn(val)
		}

		return m, nil

	case string:
		err := jsonStringToObject(v, &m)

		return m, err
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, m)
	}

	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)

	for _, keyVal := range v.MapKeys() {
		val, err := fnE(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, m)
		}

		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}

	return m, nil
}

// ToStringMapIntE casts an interface to a map[string]int{} type.
func ToStringMapIntE(i any) (map[string]int, error) {
	return toStringMapIntE(i, ToInt, ToIntE)
}

// ToStringMapInt64E casts an interface to a map[string]int64{} type.
func ToStringMapInt64E(i interface{}) (map[string]int64, error) {
	return toStringMapIntE(i, ToInt64, ToInt64E)
}

// jsonStringToObject attempts to unmarshall a string as JSON into
// the object passed as pointer.
func jsonStringToObject(s string, v interface{}) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}

// toInt returns the int value of v if v or v's underlying type
// is an int.
// Note that this will return false for int64 etc. types.
func toInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
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
