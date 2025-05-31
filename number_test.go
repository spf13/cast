// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cast"
)

type numberContext struct {
	to          func(any) any
	toErr       func(any) (any, error)
	specific    func(any) any
	generic     func(any) any
	specificErr func(any) (any, error)
	genericErr  func(any) (any, error)

	// Order of samples:
	// zero, one, 8, -8, 8.3, -8.3, min, max, underflow string, overflow string
	samples []any
}

func toAny[T cast.Number](fn func(any) T) func(i any) any {
	return func(i any) any { return fn(i) }
}

func toAnyErr[T cast.Number](fn func(any) (T, error)) func(i any) (any, error) {
	return func(i any) (any, error) { return fn(i) }
}

var numberContexts = map[string]numberContext{
	"int": {
		to:          toAny(cast.To[int]),
		toErr:       toAnyErr(cast.ToE[int]),
		specific:    toAny(cast.ToInt),
		generic:     toAny(cast.ToNumber[int]),
		specificErr: toAnyErr(cast.ToIntE),
		genericErr:  toAnyErr(cast.ToNumberE[int]),
		samples:     []any{int(0), int(1), int(8), int(-8), int(8), int(-8), math.MinInt, math.MaxInt, nil, nil},
	},
	"int8": {
		to:          toAny(cast.To[int8]),
		toErr:       toAnyErr(cast.ToE[int8]),
		specific:    toAny(cast.ToInt8),
		generic:     toAny(cast.ToNumber[int8]),
		specificErr: toAnyErr(cast.ToInt8E),
		genericErr:  toAnyErr(cast.ToNumberE[int8]),
		samples:     []any{int8(0), int8(1), int8(8), int8(-8), int8(8), int8(-8), int8(math.MinInt8), int8(math.MaxInt8), "-129", "128"},
	},
	"int16": {
		to:          toAny(cast.To[int16]),
		toErr:       toAnyErr(cast.ToE[int16]),
		specific:    toAny(cast.ToInt16),
		generic:     toAny(cast.ToNumber[int16]),
		specificErr: toAnyErr(cast.ToInt16E),
		genericErr:  toAnyErr(cast.ToNumberE[int16]),
		samples:     []any{int16(0), int16(1), int16(8), int16(-8), int16(8), int16(-8), int16(math.MinInt16), int16(math.MaxInt16), "-32769", "32768"},
	},
	"int32": {
		to:          toAny(cast.To[int32]),
		toErr:       toAnyErr(cast.ToE[int32]),
		specific:    toAny(cast.ToInt32),
		generic:     toAny(cast.ToNumber[int32]),
		specificErr: toAnyErr(cast.ToInt32E),
		genericErr:  toAnyErr(cast.ToNumberE[int32]),
		samples:     []any{int32(0), int32(1), int32(8), int32(-8), int32(8), int32(-8), int32(math.MinInt32), int32(math.MaxInt32), "-2147483649", "2147483648"},
	},
	"int64": {
		to:          toAny(cast.To[int64]),
		toErr:       toAnyErr(cast.ToE[int64]),
		specific:    toAny(cast.ToInt64),
		generic:     toAny(cast.ToNumber[int64]),
		specificErr: toAnyErr(cast.ToInt64E),
		genericErr:  toAnyErr(cast.ToNumberE[int64]),
		samples:     []any{int64(0), int64(1), int64(8), int64(-8), int64(8), int64(-8), int64(math.MinInt64), int64(math.MaxInt64), "-9223372036854775809", "9223372036854775808"},
	},
	"uint": {
		to:          toAny(cast.To[uint]),
		toErr:       toAnyErr(cast.ToE[uint]),
		specific:    toAny(cast.ToUint),
		generic:     toAny(cast.ToNumber[uint]),
		specificErr: toAnyErr(cast.ToUintE),
		genericErr:  toAnyErr(cast.ToNumberE[uint]),
		samples:     []any{uint(0), uint(1), uint(8), uint(0), uint(8), uint(0), uint(0), uint(math.MaxUint), nil, nil},
	},
	"uint8": {
		to:          toAny(cast.To[uint8]),
		toErr:       toAnyErr(cast.ToE[uint8]),
		specific:    toAny(cast.ToUint8),
		generic:     toAny(cast.ToNumber[uint8]),
		specificErr: toAnyErr(cast.ToUint8E),
		genericErr:  toAnyErr(cast.ToNumberE[uint8]),
		samples:     []any{uint8(0), uint8(1), uint8(8), uint8(0), uint8(8), uint8(0), uint8(0), uint8(math.MaxUint8), "-1", "256"},
	},
	"uint16": {
		to:          toAny(cast.To[uint16]),
		toErr:       toAnyErr(cast.ToE[uint16]),
		specific:    toAny(cast.ToUint16),
		generic:     toAny(cast.ToNumber[uint16]),
		specificErr: toAnyErr(cast.ToUint16E),
		genericErr:  toAnyErr(cast.ToNumberE[uint16]),
		samples:     []any{uint16(0), uint16(1), uint16(8), uint16(0), uint16(8), uint16(0), uint16(0), uint16(math.MaxUint16), "-1", "65536"},
	},
	"uint32": {
		to:          toAny(cast.To[uint32]),
		toErr:       toAnyErr(cast.ToE[uint32]),
		specific:    toAny(cast.ToUint32),
		generic:     toAny(cast.ToNumber[uint32]),
		specificErr: toAnyErr(cast.ToUint32E),
		genericErr:  toAnyErr(cast.ToNumberE[uint32]),
		samples:     []any{uint32(0), uint32(1), uint32(8), uint32(0), uint32(8), uint32(0), uint32(0), uint32(math.MaxUint32), "-1", "4294967296"},
	},
	"uint64": {
		to:          toAny(cast.To[uint64]),
		toErr:       toAnyErr(cast.ToE[uint64]),
		specific:    toAny(cast.ToUint64),
		generic:     toAny(cast.ToNumber[uint64]),
		specificErr: toAnyErr(cast.ToUint64E),
		genericErr:  toAnyErr(cast.ToNumberE[uint64]),
		samples:     []any{uint64(0), uint64(1), uint64(8), uint64(0), uint64(8), uint64(0), uint64(0), uint64(math.MaxUint64), "-1", "18446744073709551616"},
	},
	"float32": {
		to:          toAny(cast.To[float32]),
		toErr:       toAnyErr(cast.ToE[float32]),
		specific:    toAny(cast.ToFloat32),
		generic:     toAny(cast.ToNumber[float32]),
		specificErr: toAnyErr(cast.ToFloat32E),
		genericErr:  toAnyErr(cast.ToNumberE[float32]),
		samples:     []any{float32(0), float32(1), float32(8), float32(-8), float32(8.31), float32(-8.31), float32(-math.MaxFloat32), float32(math.MaxFloat32), nil, nil},
	},
	"float64": {
		to:          toAny(cast.To[float64]),
		toErr:       toAnyErr(cast.ToE[float64]),
		specific:    toAny(cast.ToFloat64),
		generic:     toAny(cast.ToNumber[float64]),
		specificErr: toAnyErr(cast.ToFloat64E),
		genericErr:  toAnyErr(cast.ToNumberE[float64]),
		samples:     []any{float64(0), float64(1), float64(8), float64(-8), float64(8.31), float64(-8.31), float64(-math.MaxFloat64), float64(math.MaxFloat64), nil, nil},
	},
}

func generateNumberTestCases(samples []any) []testCase {
	zero := samples[0]
	one := samples[1]
	eight := samples[2]
	eightNegative := samples[3]
	eightPoint31 := samples[4]
	eightPoint31Negative := samples[5]
	min := samples[6]
	max := samples[7]
	underflowString := samples[8]
	overflowString := samples[9]

	_ = min
	_ = max
	_ = underflowString
	_ = overflowString

	var jsonEight, jsonEightNegative, jsonEightPointZero json.Number
	_ = json.Unmarshal([]byte("8"), &jsonEight)
	_ = json.Unmarshal([]byte("-8"), &jsonEightNegative)
	_ = json.Unmarshal([]byte("8.0"), &jsonEightPointZero)

	kind := reflect.TypeOf(zero).Kind()
	isUint := kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64

	// Some precision is lost when converting from float64 to float32.
	eightPoint31_32 := eightPoint31
	eightPoint31Negative_32 := eightPoint31Negative
	if kind == reflect.Float64 {
		eightPoint31_32 = float64(float32(eightPoint31.(float64)))
		eightPoint31Negative_32 = float64(float32(eightPoint31Negative.(float64)))
	}

	testCases := []testCase{
		// Positive numbers
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
		{float32(8.31), eightPoint31_32, false},
		{float64(8.31), eightPoint31, false},
		{cast.ToString(max), max, false},

		// Negative numbers
		{int(-8), eightNegative, isUint},
		{int8(-8), eightNegative, isUint},
		{int16(-8), eightNegative, isUint},
		{int32(-8), eightNegative, isUint},
		{int64(-8), eightNegative, isUint},
		{float32(-8.31), eightPoint31Negative_32, isUint},
		{float64(-8.31), eightPoint31Negative, isUint},
		{"-8", eightNegative, isUint},

		// Other basic types
		{true, one, false},
		{false, zero, false},
		{"8", eight, false},
		{nil, zero, false},

		// JSON
		{jsonEight, eight, false},
		{jsonEightNegative, eightNegative, isUint},
		{jsonEightPointZero, eight, false},

		// Failure cases
		{"test", zero, true},
		{testing.T{}, zero, true},
	}

	if isUint && underflowString != nil {
		testCases = append(testCases, testCase{underflowString, zero, true})
	}

	if kind == reflect.Uint64 && isUint && overflowString != nil {
		testCases = append(testCases, testCase{overflowString, zero, true})
	}

	return testCases
}

func TestToNumber(t *testing.T) {
	t.Parallel()

	for typeName, ctx := range numberContexts {
		// TODO: remove after minimum Go version is >=1.22
		typeName := typeName
		ctx := ctx

		t.Run(typeName, func(t *testing.T) {
			t.Parallel()

			testCases := generateNumberTestCases(ctx.samples)

			for _, testCase := range testCases {
				// TODO: remove after minimum Go version is >=1.22
				testCase := testCase

				t.Run("", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					{
						v := ctx.specific(testCase.input)
						c.Assert(v, qt.Equals, testCase.expected)
					}

					{
						v := ctx.generic(testCase.input)
						c.Assert(v, qt.Equals, testCase.expected)
					}

					{
						v, err := ctx.specificErr(testCase.input)
						if testCase.expectError {
							c.Assert(err, qt.IsNotNil)
						} else {
							c.Assert(err, qt.IsNil)
							c.Assert(v, qt.Equals, testCase.expected)
						}
					}

					{
						v, err := ctx.genericErr(testCase.input)
						if testCase.expectError {
							c.Assert(err, qt.IsNotNil)
						} else {
							c.Assert(err, qt.IsNil)
							c.Assert(v, qt.Equals, testCase.expected)
						}
					}
				})
			}
		})
	}
}

func BenchmarkNumber(b *testing.B) {
	type testCase struct {
		name     string
		input    any
		specific func(any) (any, error)
		generic  func(any) (any, error)
	}

	var cases []testCase

	// TODO: sort keys before iterating (once Go version is updated)
	for typeName, ctx := range numberContexts {
		cases = append(
			cases,
			testCase{
				name:     typeName,
				input:    "123",
				specific: ctx.specificErr,
				generic:  ctx.genericErr,
			},
			testCase{
				name:     typeName,
				input:    "1234567890123",
				specific: ctx.specificErr,
				generic:  ctx.genericErr,
			},
			testCase{
				name:     typeName,
				input:    "-123",
				specific: ctx.specificErr,
				generic:  ctx.genericErr,
			},
			testCase{
				name:     typeName,
				input:    "-1234567890123",
				specific: ctx.specificErr,
				generic:  ctx.genericErr,
			},
			testCase{
				name:     typeName,
				input:    "0000000000123",
				specific: ctx.specificErr,
				generic:  ctx.genericErr,
			},
			testCase{
				name:     typeName,
				input:    "00000000001234567890123",
				specific: ctx.specificErr,
				generic:  ctx.genericErr,
			},
		)
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.Run("Specific", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					_, _ = c.specific(c.input)
				}
			})

			b.Run("Generic", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					_, _ = c.generic(c.input)
				}
			})
		})
	}
}
