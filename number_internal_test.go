// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"strconv"
	"testing"

	qt "github.com/frankban/quicktest"
)

func BenchmarkToInt(b *testing.B) {
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

func TestTrimDecimal(t *testing.T) {
	c := qt.New(t)

	c.Assert(trimDecimal("10.0"), qt.Equals, "10")
	c.Assert(trimDecimal("10.00"), qt.Equals, "10")
	c.Assert(trimDecimal("10.010"), qt.Equals, "10")
	c.Assert(trimDecimal("0.0000000000"), qt.Equals, "0")
	c.Assert(trimDecimal("0.00000000001"), qt.Equals, "0")
}

func BenchmarkDecimal(b *testing.B) {
	testCases := []string{"10.0", "10.00", "10.010", "0.0000000000", "0.0000000001", "10.0000000000", "10.0000000001", "10000000000000.0000000000"}

	for _, testCase := range testCases {
		// TODO: remove after minimum Go version is >=1.22
		testCase := testCase

		b.Run(testCase, func(b *testing.B) {
			b.Run("ParseFloat", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					v, err := strconv.ParseFloat(testCase, 64)
					if err != nil {
						b.Fatal(err)
					}

					n := int64(v)
					_ = n
				}
			})

			b.Run("TrimDecimal", func(b *testing.B) {
				// TODO: use b.Loop() once updated to Go 1.24
				for i := 0; i < b.N; i++ {
					v, err := strconv.ParseInt(trimDecimal(testCase), 0, 0)
					if err != nil {
						b.Fatal(err)
					}

					_ = v
				}
			})
		})
	}
}
