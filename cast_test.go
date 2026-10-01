// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast"
)

type testCase struct {
	input       any
	expected    any
	expectError bool
}

func runTests[T cast.Basic](t *testing.T, testCases []testCase, to func(i any) T, toErr func(i any) (T, error)) {
	var typ T
	_, isTime := any(typ).(time.Time)

	res := func(i any) any {
		return i
	}

	if isTime {
		res = func(i any) any {
			return i.(time.Time).UTC()
		}
	}

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
					c.Assert(res(v), qt.Equals, testCase.expected)
				})

				t.Run("To", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v := cast.To[T](testCase.input)
					c.Assert(res(v), qt.Equals, testCase.expected)
				})

				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(res(v), qt.Equals, testCase.expected)
					}
				})

				t.Run("ToE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := cast.ToE[T](testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(res(v), qt.Equals, testCase.expected)
					}
				})
			})

			t.Run("Pointer", func(t *testing.T) {
				t.Run("ToType", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v := to(&testCase.input)
					c.Assert(res(v), qt.Equals, testCase.expected)
				})

				t.Run("To", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v := cast.To[T](&testCase.input)
					c.Assert(res(v), qt.Equals, testCase.expected)
				})

				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(&testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(res(v), qt.Equals, testCase.expected)
					}
				})

				t.Run("ToE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := cast.ToE[T](&testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(res(v), qt.Equals, testCase.expected)
					}
				})
			})
		})
	}
}

func Example() {
	// Cast a value to another type
	{
		v, err := cast.ToIntE("1234")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%T(%v)\n", v, v)
	}

	// Alternatively, you can use the generic [ToE] function for [Basic] types
	{
		v, err := cast.ToE[int]("4321")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%T(%v)\n", v, v)
	}

	// You can suppress errors by using the non-error versions
	{
		v := cast.ToInt("9876")
		fmt.Printf("%T(%v)\n", v, v)
	}

	// Similarly, there is a generic [To] function for [Basic] types
	{
		v := cast.To[int]("6789")
		fmt.Printf("%T(%v)\n", v, v)
	}

	// Finally, you can use [Must] to panic if there is an error.
	{
		v := cast.Must[int](cast.ToE[int]("5555"))
		fmt.Printf("%T(%v)\n", v, v)
	}

	// Output:
	// int(1234)
	// int(4321)
	// int(9876)
	// int(6789)
	// int(5555)
}

func BenchmarkCast(b *testing.B) {
	b.Run("Bool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cast.ToBool("true")
		}
	})

	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cast.ToString(123456789)
		}
	})

	b.Run("Number", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cast.ToNumber[int]("123456789")
		}
	})

	b.Run("Int64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cast.ToInt64("123456789")
		}
	})

	b.Run("BoolSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cast.ToBoolSlice([]string{"true", "false", "TRUE", "false"})
		}
	})

	b.Run("StringSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cast.ToStringSlice([]int{123456789, 123456789, 123456789, 123456789})
		}
	})
}

// Alias types for alias testing
type MyString string
type MyBool bool
type MyInt int
type MyInt8 int8
type MyInt16 int16
type MyInt32 int32
type MyInt64 int64
type MyUint uint
type MyUint8 uint8
type MyUint16 uint16
type MyUint32 uint32
type MyUint64 uint64
type MyFloat32 float32
type MyFloat64 float64

func TestMust(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// Valid values pass through
	c.Assert(cast.Must[int](123, nil), qt.Equals, 123)
	c.Assert(cast.Must[string]("hello", nil), qt.Equals, "hello")
	c.Assert(cast.Must[bool](true, nil), qt.Equals, true)

	// Panic occurs when error is non-nil
	c.Assert(func() {
		cast.Must[int](0, errors.New("cannot cast"))
	}, qt.PanicMatches, "cannot cast")

	c.Assert(func() {
		cast.Must[string]("", errors.New("string error"))
	}, qt.PanicMatches, "string error")

	c.Assert(func() {
		cast.Must[int](cast.ToIntE("invalid_number"))
	}, qt.PanicMatches, "unable to cast .*")
}

func TestToAndToE_AllBasicTypes(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	// string
	{
		v, err := cast.ToE[string](123)
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, "123")
		c.Assert(cast.To[string](123), qt.Equals, "123")

		type uncastableStruct struct{ A int }
		_, err = cast.ToE[string](uncastableStruct{A: 1})
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[string](uncastableStruct{A: 1}), qt.Equals, "")
	}

	// bool
	{
		v, err := cast.ToE[bool]("true")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, true)
		c.Assert(cast.To[bool]("true"), qt.Equals, true)

		_, err = cast.ToE[bool]("invalid")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[bool]("invalid"), qt.Equals, false)
	}

	// int
	{
		v, err := cast.ToE[int]("123")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, 123)
		c.Assert(cast.To[int]("123"), qt.Equals, 123)

		_, err = cast.ToE[int]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[int]("bad"), qt.Equals, 0)
	}

	// int8
	{
		v, err := cast.ToE[int8]("12")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, int8(12))
		c.Assert(cast.To[int8]("12"), qt.Equals, int8(12))

		_, err = cast.ToE[int8]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[int8]("bad"), qt.Equals, int8(0))
	}

	// int16
	{
		v, err := cast.ToE[int16]("1234")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, int16(1234))
		c.Assert(cast.To[int16]("1234"), qt.Equals, int16(1234))

		_, err = cast.ToE[int16]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[int16]("bad"), qt.Equals, int16(0))
	}

	// int32
	{
		v, err := cast.ToE[int32]("12345")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, int32(12345))
		c.Assert(cast.To[int32]("12345"), qt.Equals, int32(12345))

		_, err = cast.ToE[int32]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[int32]("bad"), qt.Equals, int32(0))
	}

	// int64
	{
		v, err := cast.ToE[int64]("123456")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, int64(123456))
		c.Assert(cast.To[int64]("123456"), qt.Equals, int64(123456))

		_, err = cast.ToE[int64]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[int64]("bad"), qt.Equals, int64(0))
	}

	// uint
	{
		v, err := cast.ToE[uint]("123")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, uint(123))
		c.Assert(cast.To[uint]("123"), qt.Equals, uint(123))

		_, err = cast.ToE[uint]("-1")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[uint]("-1"), qt.Equals, uint(0))
	}

	// uint8
	{
		v, err := cast.ToE[uint8]("12")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, uint8(12))
		c.Assert(cast.To[uint8]("12"), qt.Equals, uint8(12))

		_, err = cast.ToE[uint8]("-1")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[uint8]("-1"), qt.Equals, uint8(0))
	}

	// uint16
	{
		v, err := cast.ToE[uint16]("1234")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, uint16(1234))
		c.Assert(cast.To[uint16]("1234"), qt.Equals, uint16(1234))

		_, err = cast.ToE[uint16]("-1")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[uint16]("-1"), qt.Equals, uint16(0))
	}

	// uint32
	{
		v, err := cast.ToE[uint32]("12345")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, uint32(12345))
		c.Assert(cast.To[uint32]("12345"), qt.Equals, uint32(12345))

		_, err = cast.ToE[uint32]("-1")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[uint32]("-1"), qt.Equals, uint32(0))
	}

	// uint64
	{
		v, err := cast.ToE[uint64]("123456")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, uint64(123456))
		c.Assert(cast.To[uint64]("123456"), qt.Equals, uint64(123456))

		_, err = cast.ToE[uint64]("-1")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[uint64]("-1"), qt.Equals, uint64(0))
	}

	// float32
	{
		v, err := cast.ToE[float32]("12.34")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, float32(12.34))
		c.Assert(cast.To[float32]("12.34"), qt.Equals, float32(12.34))

		_, err = cast.ToE[float32]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[float32]("bad"), qt.Equals, float32(0))
	}

	// float64
	{
		v, err := cast.ToE[float64]("12.34")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, float64(12.34))
		c.Assert(cast.To[float64]("12.34"), qt.Equals, float64(12.34))

		_, err = cast.ToE[float64]("bad")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[float64]("bad"), qt.Equals, float64(0))
	}

	// time.Time
	{
		v, err := cast.ToE[time.Time]("2020-01-01T00:00:00Z")
		c.Assert(err, qt.IsNil)
		c.Assert(v.UTC(), qt.Equals, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
		c.Assert(cast.To[time.Time]("2020-01-01T00:00:00Z").UTC(), qt.Equals, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))

		_, err = cast.ToE[time.Time]("invalid")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[time.Time]("invalid"), qt.Equals, time.Time{})
	}

	// time.Duration
	{
		v, err := cast.ToE[time.Duration]("5s")
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, 5*time.Second)
		c.Assert(cast.To[time.Duration]("5s"), qt.Equals, 5*time.Second)

		_, err = cast.ToE[time.Duration]("invalid")
		c.Assert(err, qt.IsNotNil)
		c.Assert(cast.To[time.Duration]("invalid"), qt.Equals, time.Duration(0))
	}
}
