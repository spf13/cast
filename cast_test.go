// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestStringMapStringSliceE(t *testing.T) {
	c := qt.New(t)

	// ToStringMapString inputs/outputs
	var stringMapString = map[string]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var stringMapInterface = map[string]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapString = map[interface{}]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapInterface = map[interface{}]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}

	// ToStringMapStringSlice inputs/outputs
	var stringMapStringSlice = map[string][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapInterfaceSlice = map[string][]interface{}{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapInterfaceInterfaceSlice = map[string]interface{}{"key 1": []interface{}{"value 1", "value 2", "value 3"}, "key 2": []interface{}{"value 1", "value 2", "value 3"}, "key 3": []interface{}{"value 1", "value 2", "value 3"}}
	var stringMapStringSingleSliceFieldsResult = map[string][]string{"key 1": {"value", "1"}, "key 2": {"value", "2"}, "key 3": {"value", "3"}}
	var interfaceMapStringSlice = map[interface{}][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var interfaceMapInterfaceSlice = map[interface{}][]interface{}{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}

	var stringMapStringSliceMultiple = map[string][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapStringSliceSingle = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2"}, "key 3": {"value 3"}}

	var stringMapInterface1 = map[string]interface{}{"key 1": []string{"value 1"}, "key 2": []string{"value 2"}}
	var stringMapInterfaceResult1 = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2"}}

	var jsonStringMapString = `{"key 1": "value 1", "key 2": "value 2"}`
	var jsonStringMapStringArray = `{"key 1": ["value 1"], "key 2": ["value 2", "value 3"]}`
	var jsonStringMapStringArrayResult = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2", "value 3"}}

	type Key struct {
		k string
	}

	tests := []struct {
		input  interface{}
		expect map[string][]string
		iserr  bool
	}{
		{stringMapStringSlice, stringMapStringSlice, false},
		{stringMapInterfaceSlice, stringMapStringSlice, false},
		{stringMapInterfaceInterfaceSlice, stringMapStringSlice, false},
		{stringMapStringSliceMultiple, stringMapStringSlice, false},
		{stringMapStringSliceMultiple, stringMapStringSlice, false},
		{stringMapString, stringMapStringSliceSingle, false},
		{stringMapInterface, stringMapStringSliceSingle, false},
		{stringMapInterface1, stringMapInterfaceResult1, false},
		{interfaceMapStringSlice, stringMapStringSlice, false},
		{interfaceMapInterfaceSlice, stringMapStringSlice, false},
		{interfaceMapString, stringMapStringSingleSliceFieldsResult, false},
		{interfaceMapInterface, stringMapStringSingleSliceFieldsResult, false},
		{jsonStringMapStringArray, jsonStringMapStringArrayResult, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{map[interface{}]interface{}{"foo": testing.T{}}, nil, true},
		{map[interface{}]interface{}{Key{"foo"}: "bar"}, nil, true}, // ToStringE(Key{"foo"}) should fail
		{jsonStringMapString, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringMapStringSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
			continue
		}

		c.Assert(err, qt.IsNil, errmsg)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringMapStringSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToStringMapE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect map[string]interface{}
		iserr  bool
	}{
		{map[interface{}]interface{}{"tag": "tags", "group": "groups"}, map[string]interface{}{"tag": "tags", "group": "groups"}, false},
		{map[string]interface{}{"tag": "tags", "group": "groups"}, map[string]interface{}{"tag": "tags", "group": "groups"}, false},
		{`{"tag": "tags", "group": "groups"}`, map[string]interface{}{"tag": "tags", "group": "groups"}, false},
		{`{"tag": "tags", "group": true}`, map[string]interface{}{"tag": "tags", "group": true}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringMapE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
			continue
		}

		c.Assert(err, qt.IsNil, errmsg)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringMap(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToStringMapBoolE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect map[string]bool
		iserr  bool
	}{
		{map[interface{}]interface{}{"v1": true, "v2": false}, map[string]bool{"v1": true, "v2": false}, false},
		{map[string]interface{}{"v1": true, "v2": false}, map[string]bool{"v1": true, "v2": false}, false},
		{map[string]bool{"v1": true, "v2": false}, map[string]bool{"v1": true, "v2": false}, false},
		{`{"v1": true, "v2": false}`, map[string]bool{"v1": true, "v2": false}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringMapBoolE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
			continue
		}

		c.Assert(err, qt.IsNil, errmsg)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringMapBool(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToStringMapIntE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect map[string]int
		iserr  bool
	}{
		{map[interface{}]interface{}{"v1": 1, "v2": 222}, map[string]int{"v1": 1, "v2": 222}, false},
		{map[string]interface{}{"v1": 342, "v2": 5141}, map[string]int{"v1": 342, "v2": 5141}, false},
		{map[string]int{"v1": 33, "v2": 88}, map[string]int{"v1": 33, "v2": 88}, false},
		{map[string]int32{"v1": int32(33), "v2": int32(88)}, map[string]int{"v1": 33, "v2": 88}, false},
		{map[string]uint16{"v1": uint16(33), "v2": uint16(88)}, map[string]int{"v1": 33, "v2": 88}, false},
		{map[string]float64{"v1": float64(8.22), "v2": float64(43.32)}, map[string]int{"v1": 8, "v2": 43}, false},
		{`{"v1": 67, "v2": 56}`, map[string]int{"v1": 67, "v2": 56}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringMapIntE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
			continue
		}

		c.Assert(err, qt.IsNil, errmsg)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringMapInt(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToStringMapInt64E(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect map[string]int64
		iserr  bool
	}{
		{map[interface{}]interface{}{"v1": int32(8), "v2": int32(888)}, map[string]int64{"v1": int64(8), "v2": int64(888)}, false},
		{map[string]interface{}{"v1": int64(45), "v2": int64(67)}, map[string]int64{"v1": 45, "v2": 67}, false},
		{map[string]int64{"v1": 33, "v2": 88}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]int{"v1": 33, "v2": 88}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]int32{"v1": int32(33), "v2": int32(88)}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]uint16{"v1": uint16(33), "v2": uint16(88)}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]float64{"v1": float64(8.22), "v2": float64(43.32)}, map[string]int64{"v1": 8, "v2": 43}, false},
		{`{"v1": 67, "v2": 56}`, map[string]int64{"v1": 67, "v2": 56}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringMapInt64E(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringMapInt64(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToStringMapStringE(t *testing.T) {
	c := qt.New(t)

	var stringMapString = map[string]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var stringMapInterface = map[string]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapString = map[interface{}]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapInterface = map[interface{}]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var jsonString = `{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}`
	var invalidJsonString = `{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"`
	var emptyString = ""

	tests := []struct {
		input  interface{}
		expect map[string]string
		iserr  bool
	}{
		{stringMapString, stringMapString, false},
		{stringMapInterface, stringMapString, false},
		{interfaceMapString, stringMapString, false},
		{interfaceMapInterface, stringMapString, false},
		{jsonString, stringMapString, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{invalidJsonString, nil, true},
		{emptyString, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringMapStringE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringMapString(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToBoolSliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []bool
		iserr  bool
	}{
		{[]bool{true, false, true}, []bool{true, false, true}, false},
		{[]interface{}{true, false, true}, []bool{true, false, true}, false},
		{[]int{1, 0, 1}, []bool{true, false, true}, false},
		{[]string{"true", "false", "true"}, []bool{true, false, true}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToBoolSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToBoolSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToIntSliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []int
		iserr  bool
	}{
		{[]int{1, 3}, []int{1, 3}, false},
		{[]interface{}{1.2, 3.2}, []int{1, 3}, false},
		{[]string{"2", "3"}, []int{2, 3}, false},
		{[2]string{"2", "3"}, []int{2, 3}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToIntSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToIntSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToInt64SliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []int64
		iserr  bool
	}{
		{[]int{1, 3}, []int64{1, 3}, false},
		{[]interface{}{1.2, 3.2}, []int64{1, 3}, false},
		{[]string{"2", "3"}, []int64{2, 3}, false},
		{[2]string{"2", "3"}, []int64{2, 3}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToInt64SliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToInt64Slice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToFloat64SliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []float64
		iserr  bool
	}{
		{[]int{1, 3}, []float64{1, 3}, false},
		{[]float64{1.2, 3.2}, []float64{1.2, 3.2}, false},
		{[]interface{}{1.2, 3.2}, []float64{1.2, 3.2}, false},
		{[]string{"2", "3"}, []float64{2, 3}, false},
		{[]string{"1.2", "3.2"}, []float64{1.2, 3.2}, false},
		{[2]string{"2", "3"}, []float64{2, 3}, false},
		{[2]string{"1.2", "3.2"}, []float64{1.2, 3.2}, false},
		{[]int32{1, 3}, []float64{1.0, 3.0}, false},
		{[]int64{1, 3}, []float64{1.0, 3.0}, false},
		{[]bool{true, false}, []float64{1.0, 0.0}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToFloat64SliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToFloat64Slice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToUintSliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []uint
		iserr  bool
	}{
		{[]uint{1, 3}, []uint{1, 3}, false},
		{[]interface{}{1, 3}, []uint{1, 3}, false},
		{[]string{"2", "3"}, []uint{2, 3}, false},
		{[]int{1, 3}, []uint{1, 3}, false},
		{[]int32{1, 3}, []uint{1, 3}, false},
		{[]int64{1, 3}, []uint{1, 3}, false},
		{[]float32{1.0, 3.0}, []uint{1, 3}, false},
		{[]float64{1.0, 3.0}, []uint{1, 3}, false},
		{[]bool{true, false}, []uint{1, 0}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToUintSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToUintSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToSliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []interface{}
		iserr  bool
	}{
		{[]interface{}{1, 3}, []interface{}{1, 3}, false},
		{[]map[string]interface{}{{"k1": 1}, {"k2": 2}}, []interface{}{map[string]interface{}{"k1": 1}, map[string]interface{}{"k2": 2}}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToStringSliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []string
		iserr  bool
	}{
		{[]int{1, 2}, []string{"1", "2"}, false},
		{[]int8{int8(1), int8(2)}, []string{"1", "2"}, false},
		{[]int32{int32(1), int32(2)}, []string{"1", "2"}, false},
		{[]int64{int64(1), int64(2)}, []string{"1", "2"}, false},
		{[]uint{uint(1), uint(2)}, []string{"1", "2"}, false},
		{[]uint8{uint8(1), uint8(2)}, []string{"1", "2"}, false},
		{[]uint32{uint32(1), uint32(2)}, []string{"1", "2"}, false},
		{[]uint64{uint64(1), uint64(2)}, []string{"1", "2"}, false},
		{[]float32{float32(1.01), float32(2.01)}, []string{"1.01", "2.01"}, false},
		{[]float64{float64(1.01), float64(2.01)}, []string{"1.01", "2.01"}, false},
		{[]string{"a", "b"}, []string{"a", "b"}, false},
		{[]interface{}{1, 3}, []string{"1", "3"}, false},
		{interface{}(1), []string{"1"}, false},
		{[]error{errors.New("a"), errors.New("b")}, []string{"a", "b"}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToStringSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToStringSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToDurationSliceE(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input  interface{}
		expect []time.Duration
		iserr  bool
	}{
		{[]string{"1s", "1m"}, []time.Duration{time.Second, time.Minute}, false},
		{[]int{1, 2}, []time.Duration{1, 2}, false},
		{[]interface{}{1, 3}, []time.Duration{1, 3}, false},
		{[]time.Duration{1, 3}, []time.Duration{1, 3}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"invalid"}, nil, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToDurationSliceE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)

		// Non-E test
		v = ToDurationSlice(test.input)
		c.Assert(v, qt.DeepEquals, test.expect, errmsg)
	}
}

func TestToBoolE(t *testing.T) {
	c := qt.New(t)

	var jf, jt, je json.Number
	_ = json.Unmarshal([]byte("0"), &jf)
	_ = json.Unmarshal([]byte("1"), &jt)
	_ = json.Unmarshal([]byte("1.0"), &je)
	tests := []struct {
		input  interface{}
		expect bool
		iserr  bool
	}{
		{0, false, false},
		{int64(0), false, false},
		{int32(0), false, false},
		{int16(0), false, false},
		{int8(0), false, false},
		{uint(0), false, false},
		{uint64(0), false, false},
		{uint32(0), false, false},
		{uint16(0), false, false},
		{uint8(0), false, false},
		{float64(0), false, false},
		{float32(0), false, false},
		{time.Duration(0), false, false},
		{jf, false, false},
		{nil, false, false},
		{"false", false, false},
		{"FALSE", false, false},
		{"False", false, false},
		{"f", false, false},
		{"F", false, false},
		{false, false, false},

		{"true", true, false},
		{"TRUE", true, false},
		{"True", true, false},
		{"t", true, false},
		{"T", true, false},
		{1, true, false},
		{int64(1), true, false},
		{int32(1), true, false},
		{int16(1), true, false},
		{int8(1), true, false},
		{uint(1), true, false},
		{uint64(1), true, false},
		{uint32(1), true, false},
		{uint16(1), true, false},
		{uint8(1), true, false},
		{float64(1), true, false},
		{float32(1), true, false},
		{time.Duration(1), true, false},
		{jt, true, false},
		{je, true, false},
		{true, true, false},
		{-1, true, false},
		{int64(-1), true, false},
		{int32(-1), true, false},
		{int16(-1), true, false},
		{int8(-1), true, false},

		// errors
		{"test", false, true},
		{testing.T{}, false, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToBoolE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, test.expect, errmsg)

		// Non-E test
		v = ToBool(test.input)
		c.Assert(v, qt.Equals, test.expect, errmsg)
	}
}

func BenchmarkTooBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !ToBool(true) {
			b.Fatal("ToBool returned false")
		}
	}
}

func BenchmarkCommonTimeLayouts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, commonLayout := range []string{"2019-04-29", "2017-05-30T00:00:00Z"} {
			_, err := StringToDateInDefaultLocation(commonLayout, time.UTC)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func TestIndirectPointers(t *testing.T) {
	c := qt.New(t)

	x := 13
	y := &x
	z := &y

	c.Assert(ToInt(y), qt.Equals, 13)
	c.Assert(ToInt(z), qt.Equals, 13)

}

func TestToTime(t *testing.T) {
	c := qt.New(t)

	var jntime, jnetime json.Number
	_ = json.Unmarshal([]byte("1234567890"), &jntime)
	_ = json.Unmarshal([]byte("123.4567890"), &jnetime)
	tests := []struct {
		input  interface{}
		expect time.Time
		iserr  bool
	}{
		{"2009-11-10 23:00:00 +0000 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // Time.String()
		{"Tue Nov 10 23:00:00 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},        // ANSIC
		{"Tue Nov 10 23:00:00 UTC 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},    // UnixDate
		{"Tue Nov 10 23:00:00 +0000 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},  // RubyDate
		{"10 Nov 09 23:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},             // RFC822
		{"10 Nov 09 23:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},           // RFC822Z
		{"Tuesday, 10-Nov-09 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC850
		{"Tue, 10 Nov 2009 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // RFC1123
		{"Tue, 10 Nov 2009 23:00:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC1123Z
		{"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},            // RFC3339
		{"2018-10-21T23:21:29+0200", time.Date(2018, 10, 21, 21, 21, 29, 0, time.UTC), false},      // RFC3339 without timezone hh:mm colon
		{"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},            // RFC3339Nano
		{"11:00PM", time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC), false},                              // Kitchen
		{"Nov 10 23:00:00", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},                    // Stamp
		{"Nov 10 23:00:00.000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},                // StampMilli
		{"Nov 10 23:00:00.000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},             // StampMicro
		{"Nov 10 23:00:00.000000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},          // StampNano
		{"2016-03-06 15:28:01-00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},        // RFC3339 without T
		{"2016-03-06 15:28:01-0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},         // RFC3339 without T or timezone hh:mm colon
		{"2016-03-06 15:28:01", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 -0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 -00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 +0900", time.Date(2016, 3, 6, 6, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 +09:00", time.Date(2016, 3, 6, 6, 28, 1, 0, time.UTC), false},
		{"2006-01-02", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{"02 Jan 2006", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{1472574600, time.Date(2016, 8, 30, 16, 30, 0, 0, time.UTC), false},
		{int(1482597504), time.Date(2016, 12, 24, 16, 38, 24, 0, time.UTC), false},
		{int64(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{int32(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{uint(1482597504), time.Date(2016, 12, 24, 16, 38, 24, 0, time.UTC), false},
		{uint64(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{uint32(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{jntime, time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		// errors
		{"2006", time.Time{}, true},
		{jnetime, time.Time{}, true},
		{testing.T{}, time.Time{}, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToTimeE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v.UTC(), qt.Equals, test.expect, errmsg)

		// Non-E test
		v = ToTime(test.input)
		c.Assert(v.UTC(), qt.Equals, test.expect, errmsg)
	}
}

func TestToDurationE(t *testing.T) {
	c := qt.New(t)

	var td time.Duration = 5
	var jn json.Number
	_ = json.Unmarshal([]byte("5"), &jn)

	tests := []struct {
		input  interface{}
		expect time.Duration
		iserr  bool
	}{
		{time.Duration(5), td, false},
		{int(5), td, false},
		{int64(5), td, false},
		{int32(5), td, false},
		{int16(5), td, false},
		{int8(5), td, false},
		{uint(5), td, false},
		{uint64(5), td, false},
		{uint32(5), td, false},
		{uint16(5), td, false},
		{uint8(5), td, false},
		{float64(5), td, false},
		{float32(5), td, false},
		{jn, td, false},
		{string("5"), td, false},
		{string("5ns"), td, false},
		{string("5us"), time.Microsecond * td, false},
		{string("5µs"), time.Microsecond * td, false},
		{string("5ms"), time.Millisecond * td, false},
		{string("5s"), time.Second * td, false},
		{string("5m"), time.Minute * td, false},
		{string("5h"), time.Hour * td, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToDurationE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, test.expect, errmsg)

		// Non-E test
		v = ToDuration(test.input)
		c.Assert(v, qt.Equals, test.expect, errmsg)
	}
}

func TestToTimeWithTimezones(t *testing.T) {
	c := qt.New(t)

	est, err := time.LoadLocation("EST")
	c.Assert(err, qt.IsNil)

	irn, err := time.LoadLocation("Iran")
	c.Assert(err, qt.IsNil)

	swd, err := time.LoadLocation("Europe/Stockholm")
	c.Assert(err, qt.IsNil)

	// Test same local time in different timezones
	utc2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	est2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, est)
	irn2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, irn)
	swd2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, swd)
	loc2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local)

	for i, format := range timeFormats {
		format := format
		if format.typ == timeFormatTimeOnly {
			continue
		}

		nameBase := fmt.Sprintf("%d;timeFormatType=%d;%s", i, format.typ, format.format)

		t.Run(path.Join(nameBase), func(t *testing.T) {
			est2016str := est2016.Format(format.format)
			swd2016str := swd2016.Format(format.format)

			t.Run("without default location", func(t *testing.T) {
				c := qt.New(t)
				converted, err := ToTimeE(est2016str)
				c.Assert(err, qt.IsNil)
				if format.hasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, est2016, converted)
					assertLocationEqual(t, est, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in UTC.
					assertTimeEqual(t, utc2016, converted)
					assertLocationEqual(t, time.UTC, converted.Location())
				}
			})

			t.Run("local timezone without a default location", func(t *testing.T) {
				c := qt.New(t)
				converted, err := ToTimeE(swd2016str)
				c.Assert(err, qt.IsNil)
				if format.hasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, swd2016, converted)
					assertLocationEqual(t, swd, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in UTC.
					assertTimeEqual(t, utc2016, converted)
					assertLocationEqual(t, time.UTC, converted.Location())
				}
			})

			t.Run("nil default location", func(t *testing.T) {
				c := qt.New(t)

				converted, err := ToTimeInDefaultLocationE(est2016str, nil)
				c.Assert(err, qt.IsNil)
				if format.hasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, est2016, converted)
					assertLocationEqual(t, est, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in the local timezone.
					assertTimeEqual(t, loc2016, converted)
					assertLocationEqual(t, time.Local, converted.Location())
				}

			})

			t.Run("default location not UTC", func(t *testing.T) {
				c := qt.New(t)

				converted, err := ToTimeInDefaultLocationE(est2016str, irn)
				c.Assert(err, qt.IsNil)
				if format.hasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, est2016, converted)
					assertLocationEqual(t, est, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in the given location.
					assertTimeEqual(t, irn2016, converted)
					assertLocationEqual(t, irn, converted.Location())
				}

			})

			t.Run("time in the local timezone default location not UTC", func(t *testing.T) {
				c := qt.New(t)

				converted, err := ToTimeInDefaultLocationE(swd2016str, irn)
				c.Assert(err, qt.IsNil)
				if format.hasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, swd2016, converted)
					assertLocationEqual(t, swd, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in the given location.
					assertTimeEqual(t, irn2016, converted)
					assertLocationEqual(t, irn, converted.Location())
				}

			})

		})

	}
}

func assertTimeEqual(t *testing.T, expected, actual time.Time) {
	t.Helper()
	// Compare the dates using a numeric zone as there are cases where
	// time.Parse will assign a dummy location.
	qt.Assert(t, actual.Format(time.RFC1123Z), qt.Equals, expected.Format(time.RFC1123Z))
}

func assertLocationEqual(t *testing.T, expected, actual *time.Location) {
	t.Helper()
	qt.Assert(t, locationEqual(expected, actual), qt.IsTrue)
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
