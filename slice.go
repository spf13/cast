// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ToSliceE casts any value to a []any type.
func ToSliceE(i any) ([]any, error) {
	var s []any

	switch v := i.(type) {
	case []any:
		// TODO: use slices.Clone
		return append(s, v...), nil
	case []map[string]any:
		for _, u := range v {
			s = append(s, u)
		}

		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, s)
	}
}

func toSliceE[T any](i any, fn func(any) (T, error)) ([]T, error) {
	v, ok, err := toSliceEOk(i, fn)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
	}

	return v, nil
}

func toSliceEOk[T any](i any, fn func(any) (T, error)) ([]T, bool, error) {
	if i == nil {
		return nil, false, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
	}

	switch v := i.(type) {
	case []T:
		// TODO: clone slice
		return v, true, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]T, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := fn(s.Index(j).Interface())
			if err != nil {
				return nil, true, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
			}
			a[j] = val
		}
		return a, true, nil
	default:
		return nil, false, nil
	}
}

// ToBoolSliceE casts an interface to a []bool type.
func ToBoolSliceE(i interface{}) ([]bool, error) {
	return toSliceE(i, ToBoolE)
}

// ToStringSliceE casts any value to a []string type.
func ToStringSliceE(i any) ([]string, error) {
	if a, ok, err := toSliceEOk(i, ToStringE); ok {
		return a, err
	}

	var a []string

	switch v := i.(type) {
	case string:
		return strings.Fields(v), nil
	case any:
		str, err := ToStringE(v)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, a)
		}

		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, a)
	}
}

// ToIntSliceE casts an interface to a []int type.
func ToIntSliceE(i interface{}) ([]int, error) {
	return toSliceE(i, ToIntE)
}

// ToUintSliceE casts an interface to a []uint type.
func ToUintSliceE(i interface{}) ([]uint, error) {
	return toSliceE(i, ToUintE)
}

// ToFloat64SliceE casts an interface to a []float64 type.
func ToFloat64SliceE(i interface{}) ([]float64, error) {
	return toSliceE(i, ToFloat64E)
}

// ToInt64SliceE casts an interface to a []int64 type.
func ToInt64SliceE(i interface{}) ([]int64, error) {
	return toSliceE(i, ToInt64E)
}

// ToDurationSliceE casts an interface to a []time.Duration type.
func ToDurationSliceE(i interface{}) ([]time.Duration, error) {
	return toSliceE(i, ToDurationE)
}
