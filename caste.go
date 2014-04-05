// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

func ToTimeE(i interface{}) (tim time.Time, ok bool) {
	switch s := i.(type) {
	case time.Time:
		return s, true
	case string:
		d, e := StringToDate(s)
		if e == nil {
			return d, true
		}

		jww.ERROR.Println("Could not parse Date/Time format:", e)
		return time.Time{}, false
	default:
		jww.ERROR.Printf("Unable to Cast %#v to Time", i)
		return time.Time{}, false
	}

	return time.Time{}, false
}

func ToBoolE(i interface{}) (bool, bool) {
	switch b := i.(type) {
	case bool:
		return b, true
	case nil:
		return false, true
	case int:
		if i.(int) > 0 {
			return true, true
		}
		return false, true
	default:
		return false, false
		jww.ERROR.Printf("Unable to Cast %#v to bool", i)
	}

	return false, false
}

func ToFloat64E(i interface{}) (float64, bool) {
	switch s := i.(type) {
	case float64:
		return s, true
	case float32:
		return float64(s), true
	case string:
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return float64(v), true
		} else {
			jww.ERROR.Printf("Unable to Cast %#v to float", i)
			jww.ERROR.Println(err)
		}

	default:
		jww.ERROR.Printf("Unable to Cast %#v to float", i)
	}

	return 0.0, false
}

func ToIntE(i interface{}) (int, bool) {
	switch s := i.(type) {
	case int:
		return s, true
	case int64:
		return int(s), true
	case int32:
		return int(s), true
	case int16:
		return int(s), true
	case int8:
		return int(s), true
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int(v), true
		} else {
			jww.ERROR.Printf("Unable to Cast %#v to int", i)
			jww.ERROR.Println(err)
		}
	case float64:
		return int(s), true
	case bool:
		if bool(s) {
			return 1, true
		} else {
			return 0, true
		}
	case nil:
		return 0, true
	default:
		jww.ERROR.Printf("Unable to Cast %#v to int", i)
	}

	return 0, false
}

func ToStringE(i interface{}) (string, bool) {
	switch s := i.(type) {
	case string:
		return s, true
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64), true
	case int:
		return strconv.FormatInt(int64(i.(int)), 10), true
	case []byte:
		return string(s), true
	case nil:
		return "", true
	default:
		jww.ERROR.Printf("Unable to Cast %#v to string", i)
	}

	return "", false
}

func ToStringMapStringE(i interface{}) (map[string]string, bool) {
	var m = map[string]string{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
	default:
		return m, false
	}

	return m, true
}

func ToStringMapE(i interface{}) (map[string]interface{}, bool) {
	var m = map[string]interface{}{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToString(k)] = val
		}
	default:
		return m, false
	}

	return m, true
}

func ToStringSliceE(i interface{}) ([]string, bool) {
	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, ToString(u))
		}
	default:
		return a, false
	}

	return a, true
}

func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05Z07:00",
		"02 Jan 06 15:04 MST",
		"2006-01-02",
		"02 Jan 2006",
	})
}

func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.Parse(dateType, s); e == nil {
			return
		}
	}
	return d, errors.New(fmt.Sprintf("Unable to parse date: %s", s))
}
