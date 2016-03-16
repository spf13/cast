// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"fmt"
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	var eight interface{} = 8
	assert.Equal(t, ToInt(8), 8)
	assert.Equal(t, ToInt(8.31), 8)
	assert.Equal(t, ToInt("8"), 8)
	assert.Equal(t, ToInt(true), 1)
	assert.Equal(t, ToInt(false), 0)
	assert.Equal(t, ToInt(eight), 8)
}

func TestToInt64(t *testing.T) {
	var eight interface{} = 8
	assert.Equal(t, ToInt64(int64(8)), int64(8))
	assert.Equal(t, ToInt64(8), int64(8))
	assert.Equal(t, ToInt64(8.31), int64(8))
	assert.Equal(t, ToInt64("8"), int64(8))
	assert.Equal(t, ToInt64(true), int64(1))
	assert.Equal(t, ToInt64(false), int64(0))
	assert.Equal(t, ToInt64(eight), int64(8))
}

func TestToFloat64(t *testing.T) {
	var eight interface{} = 8
	assert.Equal(t, ToFloat64(8), 8.00)
	assert.Equal(t, ToFloat64(8.31), 8.31)
	assert.Equal(t, ToFloat64("8.31"), 8.31)
	assert.Equal(t, ToFloat64(eight), 8.0)
}

func TestToString(t *testing.T) {
	var foo interface{} = "one more time"
	assert.Equal(t, ToString(8), "8")
	assert.Equal(t, ToString(int64(16)), "16")
	assert.Equal(t, ToString(8.12), "8.12")
	assert.Equal(t, ToString([]byte("one time")), "one time")
	assert.Equal(t, ToString(template.HTML("one time")), "one time")
	assert.Equal(t, ToString(template.URL("http://somehost.foo")), "http://somehost.foo")
	assert.Equal(t, ToString(template.JS("(1+2)")), "(1+2)")
	assert.Equal(t, ToString(template.CSS("a")), "a")
	assert.Equal(t, ToString(template.HTMLAttr("a")), "a")
	assert.Equal(t, ToString(foo), "one more time")
	assert.Equal(t, ToString(nil), "")
	assert.Equal(t, ToString(true), "true")
	assert.Equal(t, ToString(false), "false")
}

type foo struct {
	val string
}

func (x foo) String() string {
	return x.val
}

func TestStringerToString(t *testing.T) {

	var x foo
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
}

type fu struct {
	val string
}

func (x fu) Error() string {
	return x.val
}

func TestErrorToString(t *testing.T) {
	var x fu
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
}

func TestMaps(t *testing.T) {
	var taxonomies = map[interface{}]interface{}{"tag": "tags", "group": "groups"}
	var stringMapBool = map[interface{}]interface{}{"v1": true, "v2": false}

	// ToStringMapString inputs/outputs
	var stringMapString = map[string]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var stringMapInterface = map[string]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapString = map[interface{}]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapInterface = map[interface{}]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}

	// ToStringMapStringSlice inputs/outputs
	var stringMapStringSlice = map[string][]string{"key 1": []string{"value 1", "value 2", "value 3"}, "key 2": []string{"value 1", "value 2", "value 3"}, "key 3": []string{"value 1", "value 2", "value 3"}}
	var stringMapInterfaceSlice = map[string][]interface{}{"key 1": []interface{}{"value 1", "value 2", "value 3"}, "key 2": []interface{}{"value 1", "value 2", "value 3"}, "key 3": []interface{}{"value 1", "value 2", "value 3"}}
	var stringMapStringSingleSliceFieldsResult = map[string][]string{"key 1": []string{"value", "1"}, "key 2": []string{"value", "2"}, "key 3": []string{"value", "3"}}
	var interfaceMapStringSlice = map[interface{}][]string{"key 1": []string{"value 1", "value 2", "value 3"}, "key 2": []string{"value 1", "value 2", "value 3"}, "key 3": []string{"value 1", "value 2", "value 3"}}
	var interfaceMapInterfaceSlice = map[interface{}][]interface{}{"key 1": []interface{}{"value 1", "value 2", "value 3"}, "key 2": []interface{}{"value 1", "value 2", "value 3"}, "key 3": []interface{}{"value 1", "value 2", "value 3"}}

	var stringMapStringSliceMultiple = map[string][]string{"key 1": []string{"value 1", "value 2", "value 3"}, "key 2": []string{"value 1", "value 2", "value 3"}, "key 3": []string{"value 1", "value 2", "value 3"}}
	var stringMapStringSliceSingle = map[string][]string{"key 1": []string{"value 1"}, "key 2": []string{"value 2"}, "key 3": []string{"value 3"}}

	assert.Equal(t, ToStringMap(taxonomies), map[string]interface{}{"tag": "tags", "group": "groups"})
	assert.Equal(t, ToStringMapBool(stringMapBool), map[string]bool{"v1": true, "v2": false})

	// ToStringMapString tests
	assert.Equal(t, ToStringMapString(stringMapString), stringMapString)
	assert.Equal(t, ToStringMapString(stringMapInterface), stringMapString)
	assert.Equal(t, ToStringMapString(interfaceMapString), stringMapString)
	assert.Equal(t, ToStringMapString(interfaceMapInterface), stringMapString)

	// ToStringMapStringSlice tests
	assert.Equal(t, ToStringMapStringSlice(stringMapStringSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapInterfaceSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapStringSliceMultiple), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapStringSliceMultiple), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapString), stringMapStringSliceSingle)
	assert.Equal(t, ToStringMapStringSlice(stringMapInterface), stringMapStringSliceSingle)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapStringSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapInterfaceSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapString), stringMapStringSingleSliceFieldsResult)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapInterface), stringMapStringSingleSliceFieldsResult)
}

func TestSlices(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, ToStringSlice([]string{"a", "b"}))
	assert.Equal(t, []string{"1", "3"}, ToStringSlice([]interface{}{1, 3}))
	assert.Equal(t, []int{1, 3}, ToIntSlice([]int{1, 3}))
	assert.Equal(t, []int{1, 3}, ToIntSlice([]interface{}{1.2, 3.2}))
	assert.Equal(t, []int{2, 3}, ToIntSlice([]string{"2", "3"}))
	assert.Equal(t, []int{2, 3}, ToIntSlice([2]string{"2", "3"}))
}

func TestToBool(t *testing.T) {
	assert.Equal(t, ToBool(0), false)
	assert.Equal(t, ToBool(nil), false)
	assert.Equal(t, ToBool("false"), false)
	assert.Equal(t, ToBool("FALSE"), false)
	assert.Equal(t, ToBool("False"), false)
	assert.Equal(t, ToBool("f"), false)
	assert.Equal(t, ToBool("F"), false)
	assert.Equal(t, ToBool(false), false)
	assert.Equal(t, ToBool("foo"), false)

	assert.Equal(t, ToBool("true"), true)
	assert.Equal(t, ToBool("TRUE"), true)
	assert.Equal(t, ToBool("True"), true)
	assert.Equal(t, ToBool("t"), true)
	assert.Equal(t, ToBool("T"), true)
	assert.Equal(t, ToBool(1), true)
	assert.Equal(t, ToBool(true), true)
	assert.Equal(t, ToBool(-1), true)
}

func TestIndirectPointers(t *testing.T) {
	x := 13
	y := &x
	z := &y

	assert.Equal(t, ToInt(y), 13)
	assert.Equal(t, ToInt(z), 13)
}

func TestToDuration(t *testing.T) {
	a := time.Second * 5
	ai := int64(a)
	b := time.Second * 5
	bf := float64(b)
	assert.Equal(t, ToDuration(ai), a)
	assert.Equal(t, ToDuration(bf), b)
}

func TestToTime(t *testing.T) {
	est, err := time.LoadLocation("EST")
	if !assert.NoError(t, err) {
		return
	}

	irn, err := time.LoadLocation("Iran")
	if !assert.NoError(t, err) {
		return
	}

	swd, err := time.LoadLocation("Europe/Stockholm")
	if !assert.NoError(t, err) {
		return
	}

	// time.Parse*() fns handle the target & local timezones being the same
	// differently, so make sure we use one of the timezones as local by
	// temporarily change it.
	if !locationEqual(time.Local, swd) {
		var originalLocation *time.Location
		originalLocation, time.Local = time.Local, swd
		defer func() {
			time.Local = originalLocation
		}()
	}

	// Test same local time in different timezones
	utc2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	est2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, est)
	irn2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, irn)
	swd2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, swd)

	for _, format := range timeFormats {
		t.Logf("Checking time format '%s', has timezone: %v", format.format, format.hasTimezone)

		est2016str := est2016.Format(format.format)
		if !assert.NotEmpty(t, est2016str) {
			continue
		}

		swd2016str := swd2016.Format(format.format)
		if !assert.NotEmpty(t, swd2016str) {
			continue
		}

		// Test conversion without a default location
		converted, err := ToTimeE(est2016str)
		if assert.NoError(t, err) {
			if format.hasTimezone {
				// Converting inputs with a timezone should preserve it
				assertTimeEqual(t, est2016, converted)
				assertLocationEqual(t, est, converted.Location())
			} else {
				// Converting inputs without a timezone should be interpreted
				// as a local time in UTC.
				assertTimeEqual(t, utc2016, converted)
				assertLocationEqual(t, time.UTC, converted.Location())
			}
		}

		// Test conversion of a time in the local timezone without a default
		// location
		converted, err = ToTimeE(swd2016str)
		if assert.NoError(t, err) {
			if format.hasTimezone {
				// Converting inputs with a timezone should preserve it
				assertTimeEqual(t, swd2016, converted)
				assertLocationEqual(t, swd, converted.Location())
			} else {
				// Converting inputs without a timezone should be interpreted
				// as a local time in UTC.
				assertTimeEqual(t, utc2016, converted)
				assertLocationEqual(t, time.UTC, converted.Location())
			}
		}

		// Conversion with a nil default location sould have same behavior
		converted, err = ToTimeInDefaultLocationE(est2016str, nil)
		if assert.NoError(t, err) {
			if format.hasTimezone {
				// Converting inputs with a timezone should preserve it
				assertTimeEqual(t, est2016, converted)
				assertLocationEqual(t, est, converted.Location())
			} else {
				// Converting inputs without a timezone should be interpreted
				// as a local time in the local timezone.
				assertTimeEqual(t, swd2016, converted)
				assertLocationEqual(t, swd, converted.Location())
			}
		}

		// Test conversion with a default location that isn't UTC
		converted, err = ToTimeInDefaultLocationE(est2016str, irn)
		if assert.NoError(t, err) {
			if format.hasTimezone {
				// Converting inputs with a timezone should preserve it
				assertTimeEqual(t, est2016, converted)
				assertLocationEqual(t, est, converted.Location())
			} else {
				// Converting inputs without a timezone should be interpreted
				// as a local time in the given location.
				assertTimeEqual(t, irn2016, converted)
				assertLocationEqual(t, irn, converted.Location())
			}
		}

		// Test conversion of a time in the local timezone with a default
		// location that isn't UTC
		converted, err = ToTimeInDefaultLocationE(swd2016str, irn)
		if assert.NoError(t, err) {
			if format.hasTimezone {
				// Converting inputs with a timezone should preserve it
				assertTimeEqual(t, swd2016, converted)
				assertLocationEqual(t, swd, converted.Location())
			} else {
				// Converting inputs without a timezone should be interpreted
				// as a local time in the given location.
				assertTimeEqual(t, irn2016, converted)
				assertLocationEqual(t, irn, converted.Location())
			}
		}
	}
}

func assertTimeEqual(t *testing.T, expected, actual time.Time, msgAndArgs ...interface{}) bool {
	if !expected.Equal(actual) {
		return assert.Fail(t, fmt.Sprintf("Expected time '%s', got '%s'", expected, actual), msgAndArgs...)
	}
	return true
}

func assertLocationEqual(t *testing.T, expected, actual *time.Location, msgAndArgs ...interface{}) bool {
	if !locationEqual(expected, actual) {
		return assert.Fail(t, fmt.Sprintf("Expected location '%s', got '%s'", expected, actual), msgAndArgs...)
	}
	return true
}

func locationEqual(a, b *time.Location) bool {
	// A note about comparring time.Locations:
	//   - can't only compare pointers
	//   - can't compare loc.String() because locations with the same
	//     name can have different offsets
	//   - can't use reflect.DeepEqual because time.Location has internal
	//     caches

	if a == b {
		return true
	} else if a == nil || b == nil {
		return false
	}

	// Check if they're equal by parsing times with a format that doesn't
	// include a timezone, which will interpret it as being a local time in
	// the given zone, and comparing the resulting local times.
	tA, err := time.ParseInLocation("2006-01-02", "2016-01-01", a)
	if err != nil {
		return false
	}

	tB, err := time.ParseInLocation("2006-01-02", "2016-01-01", b)
	if err != nil {
		return false
	}

	return tA.Equal(tB)
}
