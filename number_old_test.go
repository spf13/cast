// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

type testStep struct {
	input  interface{}
	expect interface{}
	iserr  bool
}

func createNumberTestSteps(zero, one, eight, eightnegative, eightpoint31, eightpoint31negative interface{}) []testStep {
	var jeight, jminuseight, jfloateight json.Number
	_ = json.Unmarshal([]byte("8"), &jeight)
	_ = json.Unmarshal([]byte("-8"), &jminuseight)
	_ = json.Unmarshal([]byte("8.0"), &jfloateight)

	kind := reflect.TypeOf(zero).Kind()
	isUint := kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64

	// Some precision is lost when converting from float64 to float32.
	eightpoint31_32 := eightpoint31
	eightpoint31negative_32 := eightpoint31negative
	if kind == reflect.Float64 {
		eightpoint31_32 = float64(float32(eightpoint31.(float64)))
		eightpoint31negative_32 = float64(float32(eightpoint31negative.(float64)))
	}

	return []testStep{
		{int(8), eight, false},
		{int8(8), eight, false},
		{int16(8), eight, false},
		{int32(8), eight, false},
		{int64(8), eight, false},
		{time.Weekday(8), eight, false},
		{time.Month(8), eight, false},
		{uint(8), eight, false},
		{uint8(8), eight, false},
		{uint16(8), eight, false},
		{uint32(8), eight, false},
		{uint64(8), eight, false},
		{float32(8.31), eightpoint31_32, false},
		{float64(8.31), eightpoint31, false},
		{true, one, false},
		{false, zero, false},
		{"8", eight, false},
		{nil, zero, false},
		{int(-8), eightnegative, isUint},
		{int8(-8), eightnegative, isUint},
		{int16(-8), eightnegative, isUint},
		{int32(-8), eightnegative, isUint},
		{int64(-8), eightnegative, isUint},
		{float32(-8.31), eightpoint31negative_32, isUint},
		{float64(-8.31), eightpoint31negative, isUint},
		{"-8", eightnegative, isUint},
		{jeight, eight, false},
		{jminuseight, eightnegative, isUint},
		{jfloateight, eight, false},
		{"test", zero, true},
		{testing.T{}, zero, true},
	}
}

// Maybe Go 1.18 generics will make this less ugly?
func runNumberTest(c *qt.C, tests []testStep, tove func(interface{}) (interface{}, error), tov func(interface{}) interface{}) {
	c.Helper()

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i)

		v, err := tove(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
			continue
		}
		c.Assert(err, qt.IsNil, errmsg)
		c.Assert(v, qt.Equals, test.expect, errmsg)

		// Non-E test:
		v = tov(test.input)
		c.Assert(v, qt.Equals, test.expect, errmsg)
	}
}

func TestToUintE(t *testing.T) {
	tests := createNumberTestSteps(uint(0), uint(1), uint(8), uint(0), uint(8), uint(8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToUintE(v) },
		func(v interface{}) interface{} { return ToUint(v) },
	)
}

func TestToUint64E(t *testing.T) {
	tests := createNumberTestSteps(uint64(0), uint64(1), uint64(8), uint64(0), uint64(8), uint64(8))

	// Maximum value of uint64
	tests = append(tests,
		testStep{"18446744073709551615", uint64(18446744073709551615), false},
		testStep{"18446744073709551616", uint64(0), true},
	)

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToUint64E(v) },
		func(v interface{}) interface{} { return ToUint64(v) },
	)
}

func TestToUint32E(t *testing.T) {
	tests := createNumberTestSteps(uint32(0), uint32(1), uint32(8), uint32(0), uint32(8), uint32(8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToUint32E(v) },
		func(v interface{}) interface{} { return ToUint32(v) },
	)
}

func TestToUint16E(t *testing.T) {
	tests := createNumberTestSteps(uint16(0), uint16(1), uint16(8), uint16(0), uint16(8), uint16(8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToUint16E(v) },
		func(v interface{}) interface{} { return ToUint16(v) },
	)
}

func TestToUint8E(t *testing.T) {
	tests := createNumberTestSteps(uint8(0), uint8(1), uint8(8), uint8(0), uint8(8), uint8(8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToUint8E(v) },
		func(v interface{}) interface{} { return ToUint8(v) },
	)
}
func TestToIntE(t *testing.T) {
	tests := createNumberTestSteps(int(0), int(1), int(8), int(-8), int(8), int(-8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToIntE(v) },
		func(v interface{}) interface{} { return ToInt(v) },
	)
}

func TestToInt64E(t *testing.T) {
	tests := createNumberTestSteps(int64(0), int64(1), int64(8), int64(-8), int64(8), int64(-8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToInt64E(v) },
		func(v interface{}) interface{} { return ToInt64(v) },
	)
}

func TestToInt32E(t *testing.T) {
	tests := createNumberTestSteps(int32(0), int32(1), int32(8), int32(-8), int32(8), int32(-8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToInt32E(v) },
		func(v interface{}) interface{} { return ToInt32(v) },
	)
}

func TestToInt16E(t *testing.T) {
	tests := createNumberTestSteps(int16(0), int16(1), int16(8), int16(-8), int16(8), int16(-8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToInt16E(v) },
		func(v interface{}) interface{} { return ToInt16(v) },
	)
}

func TestToInt8E(t *testing.T) {
	tests := createNumberTestSteps(int8(0), int8(1), int8(8), int8(-8), int8(8), int8(-8))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToInt8E(v) },
		func(v interface{}) interface{} { return ToInt8(v) },
	)
}

func TestToFloat64E(t *testing.T) {
	tests := createNumberTestSteps(float64(0), float64(1), float64(8), float64(-8), float64(8.31), float64(-8.31))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToFloat64E(v) },
		func(v interface{}) interface{} { return ToFloat64(v) },
	)
}

func TestToFloat32E(t *testing.T) {
	tests := createNumberTestSteps(float32(0), float32(1), float32(8), float32(-8), float32(8.31), float32(-8.31))

	runNumberTest(
		qt.New(t),
		tests,
		func(v interface{}) (interface{}, error) { return ToFloat32E(v) },
		func(v interface{}) interface{} { return ToFloat32(v) },
	)
}

func BenchmarkTooInt(b *testing.B) {
	convert := func(num52 interface{}) {
		if v := ToInt(num52); v != 52 {
			b.Fatalf("ToInt returned wrong value, got %d, want %d", v, 32)
		}
	}
	for i := 0; i < b.N; i++ {
		convert("52")
		convert(52.0)
		convert(uint64(52))
	}
}

func BenchmarkTrimZeroDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		trimZeroDecimal("")
		trimZeroDecimal("123")
		trimZeroDecimal("120")
		trimZeroDecimal("120.00")
	}
}

func TestTrimZeroDecimal(t *testing.T) {
	c := qt.New(t)

	c.Assert(trimZeroDecimal("10.0"), qt.Equals, "10")
	c.Assert(trimZeroDecimal("10.00"), qt.Equals, "10")
	c.Assert(trimZeroDecimal("10.010"), qt.Equals, "10.010")
	c.Assert(trimZeroDecimal("0.0000000000"), qt.Equals, "0")
	c.Assert(trimZeroDecimal("0.00000000001"), qt.Equals, "0.00000000001")
}
