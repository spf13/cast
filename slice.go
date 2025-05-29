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

// ToSliceE casts an interface to a []interface{} type.
func ToSliceE(i interface{}) ([]interface{}, error) {
	var s []interface{}

	switch v := i.(type) {
	case []interface{}:
		return append(s, v...), nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
}

func toSliceE[T any](i any, fn func(any) (T, error)) ([]T, error) {
	if i == nil {
		return []T{}, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
	}

	switch v := i.(type) {
	case []T:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]T, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := fn(s.Index(j).Interface())
			if err != nil {
				return []T{}, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
			}
			a[j] = val
		}
		return a, nil
	default:
		return []T{}, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
	}
}

// ToBoolSliceE casts an interface to a []bool type.
func ToBoolSliceE(i interface{}) ([]bool, error) {
	return toSliceE(i, ToBoolE)
}

// ToStringSliceE casts an interface to a []string type.
func ToStringSliceE(i interface{}) ([]string, error) {
	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []string:
		return v, nil
	case []int8:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint8:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []uint64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case string:
		return strings.Fields(v), nil
	case []error:
		for _, err := range i.([]error) {
			a = append(a, err.Error())
		}
		return a, nil
	case interface{}:
		str, err := ToStringE(v)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		}
		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
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
