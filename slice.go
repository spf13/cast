// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// ToSliceE casts any value to a []any type.
func ToSliceE(i any) ([]any, error) {
	i, _ = indirect(i)

	var s []any

	if i == nil {
		return s, fmt.Errorf(errorMsg, i, i, s)
	}

	switch v := i.(type) {
	case []any:
		return slices.Clone(v), nil
	case []map[string]any:
		s = make([]any, len(v))
		for idx, u := range v {
			s[idx] = u
		}

		return s, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		val := reflect.ValueOf(i)
		s = make([]any, val.Len())
		for j := 0; j < val.Len(); j++ {
			s[j] = val.Index(j).Interface()
		}

		return s, nil
	default:
		return s, fmt.Errorf(errorMsg, i, i, s)
	}
}

func toSliceE[T Basic](i any) ([]T, error) {
	v, ok, err := toSliceEOk[T](i)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf(errorMsg, i, i, []T{})
	}

	return v, nil
}

func toSliceEOk[T Basic](i any) ([]T, bool, error) {
	i, _ = indirect(i)
	if i == nil {
		return nil, true, fmt.Errorf(errorMsg, i, i, []T{})
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
			val, err := ToE[T](s.Index(j).Interface())
			if err != nil {
				return nil, true, fmt.Errorf(errorMsg, i, i, []T{})
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
	if a, ok, err := toSliceEOk[string](i); ok {
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
			return nil, fmt.Errorf(errorMsg, i, i, a)
		}

		return []string{str}, nil
	default:
		return nil, fmt.Errorf(errorMsg, i, i, a)
	}
}
