// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"errors"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

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
