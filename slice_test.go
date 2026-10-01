// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"errors"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast"
)

func runSliceTests[T cast.Basic | any](t *testing.T, testCases []testCase, to func(i any) []T, toErr func(i any) ([]T, error)) {
	for _, testCase := range testCases {
		// TODO: remove after minimum Go version is >=1.22
		testCase := testCase

		t.Run("", func(t *testing.T) {
			t.Parallel()

			t.Run("Value", func(t *testing.T) {
				t.Run("ToType", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v := to(testCase.input)
					if v == nil {
						return
					}

					c.Assert(v, qt.DeepEquals, testCase.expected)
				})

				// t.Run("To", func(t *testing.T) {
				// 	return

				// 	t.Parallel()

				// 	c := qt.New(t)

				// 	v := cast.To[T](testCase.input)
				// 	c.Assert(v, qt.DeepEquals, testCase.expected)
				// })

				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(v, qt.DeepEquals, testCase.expected)
					}
				})

				// t.Run("ToE", func(t *testing.T) {
				// 	return

				// 	t.Parallel()

				// 	c := qt.New(t)

				// 	v, err := cast.ToE[T](testCase.input)
				// 	if testCase.expectError {
				// 		c.Assert(err, qt.IsNotNil)
				// 	} else {
				// 		c.Assert(err, qt.IsNil)
				// 		c.Assert(v, qt.DeepEquals, testCase.expected)
				// 	}
				// })
			})

			t.Run("Pointer", func(t *testing.T) {
				t.Run("ToType", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v := to(&testCase.input)
					if v == nil {
						return
					}

					c.Assert(v, qt.DeepEquals, testCase.expected)
				})

				// t.Run("To", func(t *testing.T) {
				// 	return

				// 	t.Parallel()

				// 	c := qt.New(t)

				// 	v := cast.To[T](&testCase.input)
				// 	c.Assert(v, qt.DeepEquals, testCase.expected)
				// })

				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(&testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(v, qt.DeepEquals, testCase.expected)
					}
				})

				// t.Run("ToE", func(t *testing.T) {
				// 	return

				// 	t.Parallel()

				// 	c := qt.New(t)

				// 	v, err := cast.ToE[T](&testCase.input)
				// 	if testCase.expectError {
				// 		c.Assert(err, qt.IsNotNil)
				// 	} else {
				// 		c.Assert(err, qt.IsNil)
				// 		c.Assert(v, qt.DeepEquals, testCase.expected)
				// 	}
				// })
			})
		})
	}
}

func TestBoolSlice(t *testing.T) {
	testCases := []testCase{
		{[]bool{true, false, true}, []bool{true, false, true}, false},
		{[]any{true, false, true}, []bool{true, false, true}, false},
		{[]int{1, 0, 1}, []bool{true, false, true}, false},
		{[]string{"true", "false", "true"}, []bool{true, false, true}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToBoolSlice, cast.ToBoolSliceE)
}

func TestIntSlice(t *testing.T) {
	testCases := []testCase{
		{[]int{1, 3}, []int{1, 3}, false},
		{[]any{1.2, 3.2}, []int{1, 3}, false},
		{[]string{"2", "3"}, []int{2, 3}, false},
		{[2]string{"2", "3"}, []int{2, 3}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToIntSlice, cast.ToIntSliceE)
}

func TestInt64Slice(t *testing.T) {
	testCases := []testCase{
		{[]int{1, 3}, []int64{1, 3}, false},
		{[]any{1.2, 3.2}, []int64{1, 3}, false},
		{[]string{"2", "3"}, []int64{2, 3}, false},
		{[2]string{"2", "3"}, []int64{2, 3}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToInt64Slice, cast.ToInt64SliceE)
}

func TestFloat64Slice(t *testing.T) {
	testCases := []testCase{
		{[]int{1, 3}, []float64{1, 3}, false},
		{[]float64{1.2, 3.2}, []float64{1.2, 3.2}, false},
		{[]any{1.2, 3.2}, []float64{1.2, 3.2}, false},
		{[]string{"2", "3"}, []float64{2, 3}, false},
		{[]string{"1.2", "3.2"}, []float64{1.2, 3.2}, false},
		{[2]string{"2", "3"}, []float64{2, 3}, false},
		{[2]string{"1.2", "3.2"}, []float64{1.2, 3.2}, false},
		{[]int32{1, 3}, []float64{1.0, 3.0}, false},
		{[]int64{1, 3}, []float64{1.0, 3.0}, false},
		{[]bool{true, false}, []float64{1.0, 0.0}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToFloat64Slice, cast.ToFloat64SliceE)
}

func TestUintSlice(t *testing.T) {
	testCases := []testCase{
		{[]uint{1, 3}, []uint{1, 3}, false},
		{[]any{1, 3}, []uint{1, 3}, false},
		{[]string{"2", "3"}, []uint{2, 3}, false},
		{[]int{1, 3}, []uint{1, 3}, false},
		{[]int32{1, 3}, []uint{1, 3}, false},
		{[]int64{1, 3}, []uint{1, 3}, false},
		{[]float32{1.0, 3.0}, []uint{1, 3}, false},
		{[]float64{1.0, 3.0}, []uint{1, 3}, false},
		{[]bool{true, false}, []uint{1, 0}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToUintSlice, cast.ToUintSliceE)
}

func TestSlice(t *testing.T) {
	testCases := []testCase{
		{[]any{1, 3}, []any{1, 3}, false},
		{[]map[string]any{{"k1": 1}, {"k2": 2}}, []any{map[string]any{"k1": 1}, map[string]any{"k2": 2}}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToSlice, cast.ToSliceE)
}

func TestStringSlice(t *testing.T) {
	testCases := []testCase{
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
		{[]any{1, 3}, []string{"1", "3"}, false},
		{any(1), []string{"1"}, false},
		{[]error{errors.New("a"), errors.New("b")}, []string{"a", "b"}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToStringSlice, cast.ToStringSliceE)
}

func TestDurationSlice(t *testing.T) {
	testCases := []testCase{
		{[]string{"1s", "1m"}, []time.Duration{time.Second, time.Minute}, false},
		{[]int{1, 2}, []time.Duration{1, 2}, false},
		{[]any{1, 3}, []time.Duration{1, 3}, false},
		{[]time.Duration{1, 3}, []time.Duration{1, 3}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"invalid"}, nil, true},
	}

	runSliceTests(t, testCases, cast.ToDurationSlice, cast.ToDurationSliceE)
}

func runSliceETests[T cast.Basic](t *testing.T, testCases []testCase, toErr func(i any) ([]T, error)) {
	for _, testCase := range testCases {
		// TODO: remove after minimum Go version is >=1.22
		testCase := testCase

		t.Run("", func(t *testing.T) {
			t.Parallel()

			t.Run("Value", func(t *testing.T) {
				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(v, qt.DeepEquals, testCase.expected)
					}
				})
			})

			t.Run("Pointer", func(t *testing.T) {
				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(&testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(v, qt.DeepEquals, testCase.expected)
					}
				})
			})
		})
	}
}

func TestInt8Slice(t *testing.T) {
	testCases := []testCase{
		{[]int8{-128, 0, 127}, []int8{-128, 0, 127}, false},
		{[]any{int8(1), int(2), "3"}, []int8{1, 2, 3}, false},
		{[]int{1, 3}, []int8{1, 3}, false},
		{[]string{"1", "3"}, []int8{1, 3}, false},
		{[2]int{1, 3}, []int8{1, 3}, false},
		{[]float64{1.0, 3.0}, []int8{1, 3}, false},
		{[]bool{true, false}, []int8{1, 0}, false},
		{[]any{}, []int8{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]string{"9223372036854775808"}, nil, true},  // Overflow math.MaxInt64
		{[]string{"-9223372036854775809"}, nil, true}, // Underflow math.MinInt64
	}

	runSliceETests(t, testCases, cast.ToInt8SliceE)
}

func TestInt16Slice(t *testing.T) {
	testCases := []testCase{
		{[]int16{-32768, 0, 32767}, []int16{-32768, 0, 32767}, false},
		{[]any{int16(1), int(2), "3"}, []int16{1, 2, 3}, false},
		{[]int{1, 3}, []int16{1, 3}, false},
		{[]string{"1", "3"}, []int16{1, 3}, false},
		{[2]int{1, 3}, []int16{1, 3}, false},
		{[]float64{1.0, 3.0}, []int16{1, 3}, false},
		{[]bool{true, false}, []int16{1, 0}, false},
		{[]any{}, []int16{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]string{"9223372036854775808"}, nil, true},  // Overflow math.MaxInt64
		{[]string{"-9223372036854775809"}, nil, true}, // Underflow math.MinInt64
	}

	runSliceETests(t, testCases, cast.ToInt16SliceE)
}

func TestInt32Slice(t *testing.T) {
	testCases := []testCase{
		{[]int32{-2147483648, 0, 2147483647}, []int32{-2147483648, 0, 2147483647}, false},
		{[]any{int32(1), int(2), "3"}, []int32{1, 2, 3}, false},
		{[]int{1, 3}, []int32{1, 3}, false},
		{[]string{"1", "3"}, []int32{1, 3}, false},
		{[2]int{1, 3}, []int32{1, 3}, false},
		{[]float64{1.0, 3.0}, []int32{1, 3}, false},
		{[]bool{true, false}, []int32{1, 0}, false},
		{[]any{}, []int32{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]string{"9223372036854775808"}, nil, true},  // Overflow math.MaxInt64
		{[]string{"-9223372036854775809"}, nil, true}, // Underflow math.MinInt64
	}

	runSliceETests(t, testCases, cast.ToInt32SliceE)
}

func TestUint8Slice(t *testing.T) {
	testCases := []testCase{
		{[]uint8{0, 1, 255}, []uint8{0, 1, 255}, false},
		{[]any{uint8(1), int(2), "3"}, []uint8{1, 2, 3}, false},
		{[]int{1, 3}, []uint8{1, 3}, false},
		{[]string{"1", "3"}, []uint8{1, 3}, false},
		{[2]int{1, 3}, []uint8{1, 3}, false},
		{[]float64{1.0, 3.0}, []uint8{1, 3}, false},
		{[]bool{true, false}, []uint8{1, 0}, false},
		{[]any{}, []uint8{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]int{-1}, nil, true},                        // Negative value to uint8
		{[]string{"-1"}, nil, true},                   // Negative string to uint8
		{[]string{"18446744073709551616"}, nil, true}, // Overflow math.MaxUint64
	}

	runSliceETests(t, testCases, cast.ToUint8SliceE)
}

func TestUint16Slice(t *testing.T) {
	testCases := []testCase{
		{[]uint16{0, 1, 65535}, []uint16{0, 1, 65535}, false},
		{[]any{uint16(1), int(2), "3"}, []uint16{1, 2, 3}, false},
		{[]int{1, 3}, []uint16{1, 3}, false},
		{[]string{"1", "3"}, []uint16{1, 3}, false},
		{[2]int{1, 3}, []uint16{1, 3}, false},
		{[]float64{1.0, 3.0}, []uint16{1, 3}, false},
		{[]bool{true, false}, []uint16{1, 0}, false},
		{[]any{}, []uint16{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]int{-1}, nil, true},                        // Negative value to uint16
		{[]string{"-1"}, nil, true},                   // Negative string to uint16
		{[]string{"18446744073709551616"}, nil, true}, // Overflow math.MaxUint64
	}

	runSliceETests(t, testCases, cast.ToUint16SliceE)
}

func TestUint32Slice(t *testing.T) {
	testCases := []testCase{
		{[]uint32{0, 1, 4294967295}, []uint32{0, 1, 4294967295}, false},
		{[]any{uint32(1), int(2), "3"}, []uint32{1, 2, 3}, false},
		{[]int{1, 3}, []uint32{1, 3}, false},
		{[]string{"1", "3"}, []uint32{1, 3}, false},
		{[2]int{1, 3}, []uint32{1, 3}, false},
		{[]float64{1.0, 3.0}, []uint32{1, 3}, false},
		{[]bool{true, false}, []uint32{1, 0}, false},
		{[]any{}, []uint32{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]int{-1}, nil, true},                        // Negative value to uint32
		{[]string{"-1"}, nil, true},                   // Negative string to uint32
		{[]string{"18446744073709551616"}, nil, true}, // Overflow math.MaxUint64
	}

	runSliceETests(t, testCases, cast.ToUint32SliceE)
}

func TestUint64Slice(t *testing.T) {
	testCases := []testCase{
		{[]uint64{0, 1, 18446744073709551615}, []uint64{0, 1, 18446744073709551615}, false},
		{[]any{uint64(1), int(2), "3"}, []uint64{1, 2, 3}, false},
		{[]int{1, 3}, []uint64{1, 3}, false},
		{[]string{"1", "3"}, []uint64{1, 3}, false},
		{[2]int{1, 3}, []uint64{1, 3}, false},
		{[]float64{1.0, 3.0}, []uint64{1, 3}, false},
		{[]bool{true, false}, []uint64{1, 0}, false},
		{[]any{}, []uint64{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
		{[]int{-1}, nil, true},                        // Negative value to uint64
		{[]string{"18446744073709551616"}, nil, true}, // Overflow math.MaxUint64
		{[]string{"-1"}, nil, true},                   // Negative string to uint64
	}

	runSliceETests(t, testCases, cast.ToUint64SliceE)
}

func TestFloat32Slice(t *testing.T) {
	testCases := []testCase{
		{[]float32{1.5, 3.5}, []float32{1.5, 3.5}, false},
		{[]any{float32(1.5), int(2), "3.5"}, []float32{1.5, 2.0, 3.5}, false},
		{[]int{1, 3}, []float32{1.0, 3.0}, false},
		{[]string{"1.5", "3.5"}, []float32{1.5, 3.5}, false},
		{[2]float64{1.5, 3.5}, []float32{1.5, 3.5}, false},
		{[]bool{true, false}, []float32{1.0, 0.0}, false},
		{[]any{}, []float32{}, false},

		// Failure cases
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	runSliceETests(t, testCases, cast.ToFloat32SliceE)
}
