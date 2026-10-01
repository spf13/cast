// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"encoding/json"
	"errors"
	"math"
	"math/big"
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
		samples:     []any{int(0), int(1), int(8), int(-8), int(8), int(-8), MyInt(8), math.MinInt, math.MaxInt, new(big.Int).Sub(big.NewInt(math.MinInt), big.NewInt(1)).String(), new(big.Int).Add(big.NewInt(math.MaxInt), big.NewInt(1)).String()},
	},
	"int8": {
		to:          toAny(cast.To[int8]),
		toErr:       toAnyErr(cast.ToE[int8]),
		specific:    toAny(cast.ToInt8),
		generic:     toAny(cast.ToNumber[int8]),
		specificErr: toAnyErr(cast.ToInt8E),
		genericErr:  toAnyErr(cast.ToNumberE[int8]),
		samples:     []any{int8(0), int8(1), int8(8), int8(-8), int8(8), int8(-8), MyInt8(8), int8(math.MinInt8), int8(math.MaxInt8), "-129", "128"},
	},
	"int16": {
		to:          toAny(cast.To[int16]),
		toErr:       toAnyErr(cast.ToE[int16]),
		specific:    toAny(cast.ToInt16),
		generic:     toAny(cast.ToNumber[int16]),
		specificErr: toAnyErr(cast.ToInt16E),
		genericErr:  toAnyErr(cast.ToNumberE[int16]),
		samples:     []any{int16(0), int16(1), int16(8), int16(-8), int16(8), int16(-8), MyInt16(8), int16(math.MinInt16), int16(math.MaxInt16), "-32769", "32768"},
	},
	"int32": {
		to:          toAny(cast.To[int32]),
		toErr:       toAnyErr(cast.ToE[int32]),
		specific:    toAny(cast.ToInt32),
		generic:     toAny(cast.ToNumber[int32]),
		specificErr: toAnyErr(cast.ToInt32E),
		genericErr:  toAnyErr(cast.ToNumberE[int32]),
		samples:     []any{int32(0), int32(1), int32(8), int32(-8), int32(8), int32(-8), MyInt32(8), int32(math.MinInt32), int32(math.MaxInt32), "-2147483649", "2147483648"},
	},
	"int64": {
		to:          toAny(cast.To[int64]),
		toErr:       toAnyErr(cast.ToE[int64]),
		specific:    toAny(cast.ToInt64),
		generic:     toAny(cast.ToNumber[int64]),
		specificErr: toAnyErr(cast.ToInt64E),
		genericErr:  toAnyErr(cast.ToNumberE[int64]),
		samples:     []any{int64(0), int64(1), int64(8), int64(-8), int64(8), int64(-8), MyInt64(8), int64(math.MinInt64), int64(math.MaxInt64), "-9223372036854775809", "9223372036854775808"},
	},
	"uint": {
		to:          toAny(cast.To[uint]),
		toErr:       toAnyErr(cast.ToE[uint]),
		specific:    toAny(cast.ToUint),
		generic:     toAny(cast.ToNumber[uint]),
		specificErr: toAnyErr(cast.ToUintE),
		genericErr:  toAnyErr(cast.ToNumberE[uint]),
		samples:     []any{uint(0), uint(1), uint(8), uint(0), uint(8), uint(0), MyUint(8), uint(0), uint(math.MaxUint), nil, nil},
	},
	"uint8": {
		to:          toAny(cast.To[uint8]),
		toErr:       toAnyErr(cast.ToE[uint8]),
		specific:    toAny(cast.ToUint8),
		generic:     toAny(cast.ToNumber[uint8]),
		specificErr: toAnyErr(cast.ToUint8E),
		genericErr:  toAnyErr(cast.ToNumberE[uint8]),
		samples:     []any{uint8(0), uint8(1), uint8(8), uint8(0), uint8(8), uint8(0), MyUint8(8), uint8(0), uint8(math.MaxUint8), "-1", "256"},
	},
	"uint16": {
		to:          toAny(cast.To[uint16]),
		toErr:       toAnyErr(cast.ToE[uint16]),
		specific:    toAny(cast.ToUint16),
		generic:     toAny(cast.ToNumber[uint16]),
		specificErr: toAnyErr(cast.ToUint16E),
		genericErr:  toAnyErr(cast.ToNumberE[uint16]),
		samples:     []any{uint16(0), uint16(1), uint16(8), uint16(0), uint16(8), uint16(0), MyUint16(8), uint16(0), uint16(math.MaxUint16), "-1", "65536"},
	},
	"uint32": {
		to:          toAny(cast.To[uint32]),
		toErr:       toAnyErr(cast.ToE[uint32]),
		specific:    toAny(cast.ToUint32),
		generic:     toAny(cast.ToNumber[uint32]),
		specificErr: toAnyErr(cast.ToUint32E),
		genericErr:  toAnyErr(cast.ToNumberE[uint32]),
		samples:     []any{uint32(0), uint32(1), uint32(8), uint32(0), uint32(8), uint32(0), MyUint32(8), uint32(0), uint32(math.MaxUint32), "-1", "4294967296"},
	},
	"uint64": {
		to:          toAny(cast.To[uint64]),
		toErr:       toAnyErr(cast.ToE[uint64]),
		specific:    toAny(cast.ToUint64),
		generic:     toAny(cast.ToNumber[uint64]),
		specificErr: toAnyErr(cast.ToUint64E),
		genericErr:  toAnyErr(cast.ToNumberE[uint64]),
		samples:     []any{uint64(0), uint64(1), uint64(8), uint64(0), uint64(8), uint64(0), MyUint64(8), uint64(0), uint64(math.MaxUint64), "-1", "18446744073709551616"},
	},
	"float32": {
		to:          toAny(cast.To[float32]),
		toErr:       toAnyErr(cast.ToE[float32]),
		specific:    toAny(cast.ToFloat32),
		generic:     toAny(cast.ToNumber[float32]),
		specificErr: toAnyErr(cast.ToFloat32E),
		genericErr:  toAnyErr(cast.ToNumberE[float32]),
		samples:     []any{float32(0), float32(1), float32(8), float32(-8), float32(8.31), float32(-8.31), MyFloat32(8), float32(-math.MaxFloat32), float32(math.MaxFloat32), nil, nil},
	},
	"float64": {
		to:          toAny(cast.To[float64]),
		toErr:       toAnyErr(cast.ToE[float64]),
		specific:    toAny(cast.ToFloat64),
		generic:     toAny(cast.ToNumber[float64]),
		specificErr: toAnyErr(cast.ToFloat64E),
		genericErr:  toAnyErr(cast.ToNumberE[float64]),
		samples:     []any{float64(0), float64(1), float64(8), float64(-8), float64(8.31), float64(-8.31), MyFloat64(8), float64(-math.MaxFloat64), float64(math.MaxFloat64), nil, nil},
	},
}

// TODO: separate test and failure cases?
// Kinda hard to track cases right now.
func generateNumberTestCases(samples []any) []testCase {
	zero := samples[0]
	one := samples[1]
	eight := samples[2]
	eightNegative := samples[3]
	eightPoint31 := samples[4]
	eightPoint31Negative := samples[5]
	aliasEight := samples[6]
	min := samples[7]
	max := samples[8]
	underflowString := samples[9]
	overflowString := samples[10]

	_ = min
	_ = max
	_ = underflowString
	_ = overflowString

	kind := reflect.TypeOf(zero).Kind()
	isSint := kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64
	isUint := kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64
	isInt := isSint || isUint

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

		// Other basic types
		{true, one, false},
		{false, zero, false},
		{"8", eight, false},
		{"-8", eightNegative, isUint},
		{"8.31", eightPoint31, false},
		{"-8.31", eightPoint31Negative, isUint},
		{"", zero, false},
		{nil, zero, false},

		// JSON
		{json.Number("8"), eight, false},
		{json.Number("-8"), eightNegative, isUint},
		{json.Number("8.31"), eightPoint31, false},
		{json.Number("-8.31"), eightPoint31Negative, isUint},
		{json.Number(""), zero, false},

		// Aliases
		{aliasEight, eight, false},
		{MyString("8"), eight, false},
		{MyBool(true), one, false},

		// Failure cases
		{"test", zero, true},
		{testing.T{}, zero, true},

		{"10...17", zero, true},
		{"10.foobar", zero, true},
		{"10.0i", zero, true},
	}

	if isInt {
		testCases = append(
			testCases,

			testCase{".5", zero, false},
			testCase{"+8.", eight, false},
			testCase{"+.25", zero, false},
			testCase{"-.25", zero, isUint},

			testCase{"10.0E9", zero, true},
		)
	} else if kind == reflect.Float32 {
		testCases = append(testCases, testCase{"10.0E9", float32(10000000000.000000), false})
	} else if kind == reflect.Float64 {
		testCases = append(testCases, testCase{"10.0E9", float64(10000000000.000000), false})
	}

	if isUint && underflowString != nil {
		testCases = append(testCases, testCase{underflowString, zero, true})
	}

	if kind == reflect.Uint64 && isUint && overflowString != nil {
		testCases = append(testCases, testCase{overflowString, zero, true})
	}

	return testCases
}

func TestNumber(t *testing.T) {
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

					t.Run("Value", func(t *testing.T) {
						t.Run("ToType", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v := ctx.specific(testCase.input)
							c.Assert(v, qt.Equals, testCase.expected)
						})

						t.Run("To", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v := ctx.generic(testCase.input)
							c.Assert(v, qt.Equals, testCase.expected)
						})

						t.Run("ToBasic", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v := ctx.to(testCase.input)
							c.Assert(v, qt.Equals, testCase.expected)
						})

						t.Run("ToTypeE", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v, err := ctx.specificErr(testCase.input)
							if testCase.expectError {
								c.Assert(err, qt.IsNotNil)
							} else {
								c.Assert(err, qt.IsNil)
								c.Assert(v, qt.Equals, testCase.expected)
							}
						})

						t.Run("ToE", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v, err := ctx.genericErr(testCase.input)
							if testCase.expectError {
								c.Assert(err, qt.IsNotNil)
							} else {
								c.Assert(err, qt.IsNil)
								c.Assert(v, qt.Equals, testCase.expected)
							}
						})

						t.Run("ToBasicE", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v, err := ctx.toErr(testCase.input)
							if testCase.expectError {
								c.Assert(err, qt.IsNotNil)
							} else {
								c.Assert(err, qt.IsNil)
								c.Assert(v, qt.Equals, testCase.expected)
							}
						})

						t.Run("Pointer", func(t *testing.T) {
							t.Run("ToType", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v := ctx.specific(&testCase.input)
								c.Assert(v, qt.Equals, testCase.expected)
							})

							t.Run("To", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v := ctx.generic(&testCase.input)
								c.Assert(v, qt.Equals, testCase.expected)
							})

							t.Run("ToBasic", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v := ctx.to(&testCase.input)
								c.Assert(v, qt.Equals, testCase.expected)
							})

							t.Run("ToTypeE", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v, err := ctx.specificErr(&testCase.input)
								if testCase.expectError {
									c.Assert(err, qt.IsNotNil)
								} else {
									c.Assert(err, qt.IsNil)
									c.Assert(v, qt.Equals, testCase.expected)
								}
							})

							t.Run("ToE", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v, err := ctx.genericErr(&testCase.input)
								if testCase.expectError {
									c.Assert(err, qt.IsNotNil)
								} else {
									c.Assert(err, qt.IsNil)
									c.Assert(v, qt.Equals, testCase.expected)
								}
							})

							t.Run("ToBasicE", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v, err := ctx.toErr(&testCase.input)
								if testCase.expectError {
									c.Assert(err, qt.IsNotNil)
								} else {
									c.Assert(err, qt.IsNil)
									c.Assert(v, qt.Equals, testCase.expected)
								}
							})
						})
					})
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

	var testCases []testCase

	// TODO: sort keys before iterating (once Go version is updated)
	for typeName, ctx := range numberContexts {
		testCases = append(
			testCases,
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

	for _, testCase := range testCases {
		// TODO: remove after minimum Go version is >=1.22
		testCase := testCase

		b.Run(testCase.name, func(b *testing.B) {
			b.Run("Specific", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					_, _ = testCase.specific(testCase.input)
				}
			})

			b.Run("Generic", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					_, _ = testCase.generic(testCase.input)
				}
			})
		})
	}
}

type mockFloat64Provider struct {
	val float64
}

func (m mockFloat64Provider) Float64() float64 {
	return m.val
}

type mockFloat64EProvider struct {
	val float64
	err error
}

func (m mockFloat64EProvider) Float64() (float64, error) {
	return m.val, m.err
}

func TestFloat64Providers_Public(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// mockFloat64Provider success on ToFloat64 / ToFloat64E
	pValid := mockFloat64Provider{val: 42.5}
	c.Assert(cast.ToFloat64(pValid), qt.Equals, 42.5)
	vF64, err := cast.ToFloat64E(pValid)
	c.Assert(err, qt.IsNil)
	c.Assert(vF64, qt.Equals, 42.5)

	// Converting float64Provider to non-float64 number fails:
	_, err = cast.ToIntE(pValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToInt64E(pValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToFloat32E(pValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToUintE(pValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToUint64E(pValid)
	c.Assert(err, qt.IsNotNil)

	// mockFloat64EProvider success
	pEValid := mockFloat64EProvider{val: 84.25, err: nil}
	vFE64, err := cast.ToFloat64E(pEValid)
	c.Assert(err, qt.IsNil)
	c.Assert(vFE64, qt.Equals, 84.25)
	c.Assert(cast.ToFloat64(pEValid), qt.Equals, 84.25)

	// mockFloat64EProvider error
	pEErr := mockFloat64EProvider{err: errors.New("provider failed")}
	_, err = cast.ToFloat64E(pEErr)
	c.Assert(err, qt.IsNotNil)
	c.Assert(cast.ToFloat64(pEErr), qt.Equals, float64(0))

	// Converting float64EProvider to non-float64 number fails:
	_, err = cast.ToIntE(pEValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToInt64E(pEValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToFloat32E(pEValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToUintE(pEValid)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToUint64E(pEValid)
	c.Assert(err, qt.IsNotNil)
}

func TestToUnsignedNumber_NegativeTime(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// Negative time.Weekday
	v1, err := cast.ToUintE(time.Weekday(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v1, qt.Equals, uint(0))
	c.Assert(cast.ToUint(time.Weekday(-1)), qt.Equals, uint(0))

	v1_8, err := cast.ToUint8E(time.Weekday(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v1_8, qt.Equals, uint8(0))

	v1_16, err := cast.ToUint16E(time.Weekday(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v1_16, qt.Equals, uint16(0))

	v1_32, err := cast.ToUint32E(time.Weekday(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v1_32, qt.Equals, uint32(0))

	v1_64, err := cast.ToUint64E(time.Weekday(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v1_64, qt.Equals, uint64(0))

	// Negative time.Month
	v2, err := cast.ToUintE(time.Month(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v2, qt.Equals, uint(0))
	c.Assert(cast.ToUint(time.Month(-1)), qt.Equals, uint(0))

	v2_8, err := cast.ToUint8E(time.Month(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v2_8, qt.Equals, uint8(0))

	v2_16, err := cast.ToUint16E(time.Month(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v2_16, qt.Equals, uint16(0))

	v2_32, err := cast.ToUint32E(time.Month(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v2_32, qt.Equals, uint32(0))

	v2_64, err := cast.ToUint64E(time.Month(-1))
	c.Assert(err, qt.IsNotNil)
	c.Assert(v2_64, qt.Equals, uint64(0))
}

func TestJSONNumber_Errors(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	invalidJSON := json.Number("not_a_valid_number")

	// Errors in integer casting
	_, err := cast.ToIntE(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	c.Assert(cast.ToInt(invalidJSON), qt.Equals, 0)

	_, err = cast.ToInt64E(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	c.Assert(cast.ToInt64(invalidJSON), qt.Equals, int64(0))

	_, err = cast.ToInt32E(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToInt16E(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToInt8E(invalidJSON)
	c.Assert(err, qt.IsNotNil)

	// Errors in float casting
	_, err = cast.ToFloat64E(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	c.Assert(cast.ToFloat64(invalidJSON), qt.Equals, float64(0))

	_, err = cast.ToFloat32E(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	c.Assert(cast.ToFloat32(invalidJSON), qt.Equals, float32(0))

	// Error in basic.go:ToBoolE (covers basic.go:59!)
	b, err := cast.ToBoolE(invalidJSON)
	c.Assert(err, qt.IsNotNil)
	c.Assert(b, qt.Equals, false)
	c.Assert(cast.ToBool(invalidJSON), qt.Equals, false)

	// Valid json.Number with ToBoolE
	b, err = cast.ToBoolE(json.Number("1"))
	c.Assert(err, qt.IsNil)
	c.Assert(b, qt.Equals, true)

	b, err = cast.ToBoolE(json.Number("0"))
	c.Assert(err, qt.IsNil)
	c.Assert(b, qt.Equals, false)
}

func TestNumber_FormatsAndOverflow(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// Hexadecimal
	c.Assert(cast.ToInt("0x10"), qt.Equals, 16)
	c.Assert(cast.ToUint("0x10"), qt.Equals, uint(16))
	c.Assert(cast.ToInt("-0x10"), qt.Equals, -16)
	c.Assert(cast.ToUint("+0x10"), qt.Equals, uint(16))
	_, err := cast.ToUintE("-0x10")
	c.Assert(err, qt.IsNotNil)

	// Octal (0o and legacy 0 prefix)
	c.Assert(cast.ToInt("0o77"), qt.Equals, 63)
	c.Assert(cast.ToUint("0o77"), qt.Equals, uint(63))
	c.Assert(cast.ToInt("-0o77"), qt.Equals, -63)
	c.Assert(cast.ToInt("077"), qt.Equals, 63)
	c.Assert(cast.ToUint("077"), qt.Equals, uint(63))

	// Binary
	c.Assert(cast.ToInt("0b1010"), qt.Equals, 10)
	c.Assert(cast.ToUint("0b1010"), qt.Equals, uint(10))
	c.Assert(cast.ToInt("-0b1010"), qt.Equals, -10)
	_, err = cast.ToUintE("-0b1010")
	c.Assert(err, qt.IsNotNil)

	// Float overflow (ErrRange)
	_, err = cast.ToFloat32E("1e40")
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToFloat64E("1e309")
	c.Assert(err, qt.IsNotNil)

	// Float hex notation
	c.Assert(cast.ToFloat64("0x10p0"), qt.Equals, float64(16))
	c.Assert(cast.ToFloat64("0x1p-2"), qt.Equals, float64(0.25))

	// Whitespace rejects
	_, err = cast.ToIntE(" 123 ")
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToUintE(" 123 ")
	c.Assert(err, qt.IsNotNil)
	_, err = cast.ToFloat64E(" 123.45 ")
	c.Assert(err, qt.IsNotNil)

	// Decimal trim special cases
	c.Assert(cast.ToInt(".5"), qt.Equals, 0)
	c.Assert(cast.ToInt("+8."), qt.Equals, 8)
	c.Assert(cast.ToInt("+0.25"), qt.Equals, 0)
	c.Assert(cast.ToInt("-0.25"), qt.Equals, 0)
	c.Assert(cast.ToUint("+8."), qt.Equals, uint(8))
}

type customStringer struct {
	text string
}

func (s customStringer) String() string {
	return s.text
}

func TestPointersAndInterfaces(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// Nil pointers
	var pi *int
	c.Assert(cast.ToInt(pi), qt.Equals, 0)
	vInt, err := cast.ToIntE(pi)
	c.Assert(err, qt.IsNil)
	c.Assert(vInt, qt.Equals, 0)

	var pu *uint
	c.Assert(cast.ToUint(pu), qt.Equals, uint(0))
	vUint, err := cast.ToUintE(pu)
	c.Assert(err, qt.IsNil)
	c.Assert(vUint, qt.Equals, uint(0))

	var pf *float64
	c.Assert(cast.ToFloat64(pf), qt.Equals, float64(0))
	vFloat, err := cast.ToFloat64E(pf)
	c.Assert(err, qt.IsNil)
	c.Assert(vFloat, qt.Equals, float64(0))

	var ps *string
	c.Assert(cast.ToString(ps), qt.Equals, "")
	vStr, err := cast.ToStringE(ps)
	c.Assert(err, qt.IsNil)
	c.Assert(vStr, qt.Equals, "")

	var pb *bool
	c.Assert(cast.ToBool(pb), qt.Equals, false)
	vBool, err := cast.ToBoolE(pb)
	c.Assert(err, qt.IsNil)
	c.Assert(vBool, qt.Equals, false)

	// Multi-level pointers
	val := 99
	p1 := &val
	p2 := &p1
	p3 := &p2
	c.Assert(cast.ToInt(p3), qt.Equals, 99)
	vMulti, err := cast.ToIntE(p3)
	c.Assert(err, qt.IsNil)
	c.Assert(vMulti, qt.Equals, 99)

	// Multi-level nil pointer
	var pNil *int
	pNil2 := &pNil
	c.Assert(cast.ToInt(pNil2), qt.Equals, 0)

	// Stringer in ToStringE
	cs := customStringer{text: "hello world"}
	c.Assert(cast.ToString(cs), qt.Equals, "hello world")
	sRes, err := cast.ToStringE(cs)
	c.Assert(err, qt.IsNil)
	c.Assert(sRes, qt.Equals, "hello world")

	// Stringer passed to ToIntE fails (cast does not parse Stringer to number)
	_, err = cast.ToIntE(cs)
	c.Assert(err, qt.IsNotNil)

	// error interface in ToStringE
	myErr := errors.New("something went wrong")
	c.Assert(cast.ToString(myErr), qt.Equals, "something went wrong")
	eRes, err := cast.ToStringE(myErr)
	c.Assert(err, qt.IsNil)
	c.Assert(eRes, qt.Equals, "something went wrong")
}
