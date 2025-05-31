// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cast"
)

type testCase struct {
	input       any
	expected    any
	expectError bool
}

func runTests[T cast.Basic](t *testing.T, testCases []testCase, to func(i any) T, toErr func(i any) (T, error)) {
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
					c.Assert(v, qt.Equals, testCase.expected)
				})

				t.Run("To", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v := cast.To[T](testCase.input)
					c.Assert(v, qt.Equals, testCase.expected)
				})

				t.Run("ToTypeE", func(t *testing.T) {
					t.Parallel()

					c := qt.New(t)

					v, err := toErr(testCase.input)
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

					v, err := cast.ToE[T](testCase.input)
					if testCase.expectError {
						c.Assert(err, qt.IsNotNil)
					} else {
						c.Assert(err, qt.IsNil)
						c.Assert(v, qt.Equals, testCase.expected)
					}
				})
			})

			t.Run("", func(t *testing.T) {
				t.Parallel()

				t.Run("Value", func(t *testing.T) {
					t.Run("ToType", func(t *testing.T) {
						t.Parallel()

						c := qt.New(t)

						v := to(&testCase.input)
						c.Assert(v, qt.Equals, testCase.expected)
					})

					t.Run("To", func(t *testing.T) {
						t.Parallel()

						c := qt.New(t)

						v := cast.To[T](&testCase.input)
						c.Assert(v, qt.Equals, testCase.expected)
					})

					t.Run("ToTypeE", func(t *testing.T) {
						t.Parallel()

						c := qt.New(t)

						v, err := toErr(&testCase.input)
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

						v, err := cast.ToE[T](&testCase.input)
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
}
