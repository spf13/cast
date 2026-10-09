// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"encoding/json"
	"math"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast"
)

func TestHexStringToUint64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       any
		expected    uint64
		expectError bool
	}{
		// Happy path without 0x prefix
		{
			name:     "hex without prefix issue 334 example",
			input:    "882d5422d5fffff",
			expected: 613286979481632767,
		},
		{
			name:     "lowercase hex without prefix",
			input:    "1a",
			expected: 26,
		},
		{
			name:     "uppercase hex without prefix",
			input:    "1A",
			expected: 26,
		},
		{
			name:     "hex word beef",
			input:    "beef",
			expected: 48879,
		},
		{
			name:     "hex word uppercase BEEF",
			input:    "BEEF",
			expected: 48879,
		},
		{
			name:     "max uint64 hex without prefix",
			input:    "ffffffffffffffff",
			expected: math.MaxUint64,
		},
		{
			name:     "max uint64 uppercase hex without prefix",
			input:    "FFFFFFFFFFFFFFFF",
			expected: math.MaxUint64,
		},

		// Happy path with 0x and 0X prefix
		{
			name:     "hex with 0x prefix",
			input:    "0x1a",
			expected: 26,
		},
		{
			name:     "hex with 0X prefix",
			input:    "0X1A",
			expected: 26,
		},
		{
			name:     "issue 334 example with 0x prefix",
			input:    "0x882d5422d5fffff",
			expected: 613286979481632767,
		},
		{
			name:     "issue 334 example with 0X prefix uppercase",
			input:    "0X882D5422D5FFFFF",
			expected: 613286979481632767,
		},
		{
			name:     "hex zero with 0x prefix",
			input:    "0x0",
			expected: 0,
		},
		{
			name:     "max uint64 with 0x prefix",
			input:    "0xffffffffffffffff",
			expected: math.MaxUint64,
		},
		{
			name:     "max uint64 with 0X prefix uppercase",
			input:    "0XFFFFFFFFFFFFFFFF",
			expected: math.MaxUint64,
		},

		// Happy path with positive sign prefix
		{
			name:     "hex with positive sign and 0x prefix",
			input:    "+0x1a",
			expected: 26,
		},
		{
			name:     "hex with positive sign without prefix",
			input:    "+1a",
			expected: 26,
		},
		{
			name:     "issue 334 example with positive sign",
			input:    "+882d5422d5fffff",
			expected: 613286979481632767,
		},

		// Types and wrappers
		{
			name:     "json.Number with hex string",
			input:    json.Number("882d5422d5fffff"),
			expected: 613286979481632767,
		},
		{
			name:     "json.Number with 0x prefix",
			input:    json.Number("0x1a"),
			expected: 26,
		},
		{
			name:     "MyString alias with hex string",
			input:    MyString("882d5422d5fffff"),
			expected: 613286979481632767,
		},

		// Negative hex inputs (must return error)
		{
			name:        "negative hex without prefix",
			input:       "-882d5422d5fffff",
			expectError: true,
		},
		{
			name:        "negative hex with 0x prefix",
			input:       "-0x1a",
			expectError: true,
		},
		{
			name:        "negative hex short without prefix",
			input:       "-1a",
			expectError: true,
		},
		{
			name:        "negative hex zero with 0x prefix",
			input:       "-0x0",
			expectError: true,
		},

		// Invalid syntax cases
		{
			name:        "lone 0x prefix",
			input:       "0x",
			expectError: true,
		},
		{
			name:        "invalid hex characters with 0x",
			input:       "0xxyz",
			expectError: true,
		},
		{
			name:        "invalid hex characters without 0x",
			input:       "xyz",
			expectError: true,
		},
		{
			name:        "hex ending with invalid character",
			input:       "882d5422d5fffffz",
			expectError: true,
		},
		{
			name:        "hex with decimal point",
			input:       "882d5422d5fffff.5",
			expectError: true,
		},

		// Overflow cases
		{
			name:        "overflow with 0x prefix",
			input:       "0x10000000000000000",
			expectError: true,
		},
		{
			name:        "overflow without prefix",
			input:       "1ffffffffffffffff",
			expectError: true,
		},
		{
			name:        "all f overflow without prefix",
			input:       "fffffffffffffffff",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := qt.New(t)

			// Test ToUint64
			v := cast.ToUint64(tc.input)
			if tc.expectError {
				c.Assert(v, qt.Equals, uint64(0))
			} else {
				c.Assert(v, qt.Equals, tc.expected)
			}

			// Test ToUint64E
			vErr, err := cast.ToUint64E(tc.input)
			if tc.expectError {
				c.Assert(err, qt.IsNotNil)
				c.Assert(vErr, qt.Equals, uint64(0))
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(vErr, qt.Equals, tc.expected)
			}

			// Test generic To[uint64]
			vGen := cast.To[uint64](tc.input)
			if tc.expectError {
				c.Assert(vGen, qt.Equals, uint64(0))
			} else {
				c.Assert(vGen, qt.Equals, tc.expected)
			}

			// Test generic ToE[uint64]
			vGenErr, err := cast.ToE[uint64](tc.input)
			if tc.expectError {
				c.Assert(err, qt.IsNotNil)
				c.Assert(vGenErr, qt.Equals, uint64(0))
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(vGenErr, qt.Equals, tc.expected)
			}

			// Test generic ToNumber[uint64]
			vNum := cast.ToNumber[uint64](tc.input)
			if tc.expectError {
				c.Assert(vNum, qt.Equals, uint64(0))
			} else {
				c.Assert(vNum, qt.Equals, tc.expected)
			}

			// Test generic ToNumberE[uint64]
			vNumErr, err := cast.ToNumberE[uint64](tc.input)
			if tc.expectError {
				c.Assert(err, qt.IsNotNil)
				c.Assert(vNumErr, qt.Equals, uint64(0))
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(vNumErr, qt.Equals, tc.expected)
			}

			// Test pointer inputs
			inputPtr := tc.input
			vPtr := cast.ToUint64(&inputPtr)
			if tc.expectError {
				c.Assert(vPtr, qt.Equals, uint64(0))
			} else {
				c.Assert(vPtr, qt.Equals, tc.expected)
			}

			vPtrErr, err := cast.ToUint64E(&inputPtr)
			if tc.expectError {
				c.Assert(err, qt.IsNotNil)
				c.Assert(vPtrErr, qt.Equals, uint64(0))
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(vPtrErr, qt.Equals, tc.expected)
			}
		})
	}
}

func TestHexStringToOtherUintTypes(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	// ToUint
	c.Assert(cast.ToUint("1a"), qt.Equals, uint(26))
	c.Assert(cast.ToUint("0x1a"), qt.Equals, uint(26))
	uintVal, err := cast.ToUintE("1a")
	c.Assert(err, qt.IsNil)
	c.Assert(uintVal, qt.Equals, uint(26))

	// ToUint32
	c.Assert(cast.ToUint32("1a"), qt.Equals, uint32(26))
	c.Assert(cast.ToUint32("0x1a"), qt.Equals, uint32(26))
	uint32Val, err := cast.ToUint32E("1a")
	c.Assert(err, qt.IsNil)
	c.Assert(uint32Val, qt.Equals, uint32(26))

	// ToUint16
	c.Assert(cast.ToUint16("1a"), qt.Equals, uint16(26))
	c.Assert(cast.ToUint16("0x1a"), qt.Equals, uint16(26))
	uint16Val, err := cast.ToUint16E("1a")
	c.Assert(err, qt.IsNil)
	c.Assert(uint16Val, qt.Equals, uint16(26))

	// ToUint8
	c.Assert(cast.ToUint8("1a"), qt.Equals, uint8(26))
	c.Assert(cast.ToUint8("0x1a"), qt.Equals, uint8(26))
	uint8Val, err := cast.ToUint8E("1a")
	c.Assert(err, qt.IsNil)
	c.Assert(uint8Val, qt.Equals, uint8(26))
}
