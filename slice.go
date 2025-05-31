// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"fmt"
	"reflect"
	"strings"
)

// ToSliceE casts any value to a []any type.
func ToSliceE(i any) ([]any, error) {
	i, _ = indirect(i)

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
	i, _ = indirect(i)
	if i == nil {
		return nil, true, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, []T{})
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

// ToStringSliceE casts any value to a []string type.
func ToStringSliceE(i any) ([]string, error) {
	if a, ok, err := toSliceEOk(i, ToStringE); ok {
		if err != nil {
			return nil, err
		}

		return a, nil
	}

	var a []string

	switch v := i.(type) {
	case string:
		return strings.Fields(v), nil
	case any:
		str, err := ToStringE(v)
		if err != nil {
			return nil, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, a)
		}

		return []string{str}, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to %T", i, i, a)
	}
}
