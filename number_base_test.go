// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

type baseTestCase struct {
	input       any
	base        int
	expected    any
	expectError bool
}

func generateNumberBaseTestCases(samples []any) []baseTestCase {
	zero := samples[0]
	// one := samples[1]
	eight := samples[2]
	// eightNegative := samples[3]
	// eightPoint31 := samples[4]
	// eightPoint31Negative := samples[5]
	// aliasEight := samples[6]
	min := samples[7]
	max := samples[8]
	underflowString := samples[9]
	overflowString := samples[10]

	_ = min
	_ = max
	_ = underflowString
	_ = overflowString

	// kind := reflect.TypeOf(zero).Kind()
	// isSint := kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64
	// isUint := kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64
	// isInt := isSint || isUint

	// Some precision is lost when converting from float64 to float32.
	// eightPoint31_32 := eightPoint31
	// eightPoint31Negative_32 := eightPoint31Negative
	// if kind == reflect.Float64 {
	// 	eightPoint31_32 = float64(float32(eightPoint31.(float64)))
	// 	eightPoint31Negative_32 = float64(float32(eightPoint31Negative.(float64)))
	// }

	testCases := []baseTestCase{
		{"08", 10, eight, false},
		{"0008", 10, eight, false},
		{"010", 8, eight, false},
		{"08", 16, eight, false},

		{"0x08", 10, zero, true},
	}

	return testCases
}

func TestNumberBase(t *testing.T) {
	t.Parallel()

	for typeName, ctx := range numberContexts {
		// TODO: remove after minimum Go version is >=1.22
		typeName := typeName
		ctx := ctx

		if typeName == "float32" || typeName == "float64" {
			continue
		}

		t.Run(typeName, func(t *testing.T) {
			t.Parallel()

			testCases := generateNumberBaseTestCases(ctx.samples)

			for _, testCase := range testCases {
				// TODO: remove after minimum Go version is >=1.22
				testCase := testCase

				t.Run("", func(t *testing.T) {
					t.Parallel()

					t.Run("Value", func(t *testing.T) {
						t.Run("To", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v := ctx.base(testCase.input, testCase.base)
							c.Assert(v, qt.Equals, testCase.expected)
						})

						t.Run("ToE", func(t *testing.T) {
							t.Parallel()

							c := qt.New(t)

							v, err := ctx.baseErr(testCase.input, testCase.base)
							if testCase.expectError {
								c.Assert(err, qt.IsNotNil)
							} else {
								c.Assert(err, qt.IsNil)
								c.Assert(v, qt.Equals, testCase.expected)
							}
						})

						t.Run("Pointer", func(t *testing.T) {
							t.Run("To", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v := ctx.base(&testCase.input, testCase.base)
								c.Assert(v, qt.Equals, testCase.expected)
							})

							t.Run("ToE", func(t *testing.T) {
								t.Parallel()

								c := qt.New(t)

								v, err := ctx.baseErr(&testCase.input, testCase.base)
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
