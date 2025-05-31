// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestToBoolE(t *testing.T) {
	c := qt.New(t)

	var jf, jt, je json.Number
	_ = json.Unmarshal([]byte("0"), &jf)
	_ = json.Unmarshal([]byte("1"), &jt)
	_ = json.Unmarshal([]byte("1.0"), &je)
	tests := []struct {
		input  interface{}
		expect bool
		iserr  bool
	}{
		{0, false, false},
		{int64(0), false, false},
		{int32(0), false, false},
		{int16(0), false, false},
		{int8(0), false, false},
		{uint(0), false, false},
		{uint64(0), false, false},
		{uint32(0), false, false},
		{uint16(0), false, false},
		{uint8(0), false, false},
		{float64(0), false, false},
		{float32(0), false, false},
		{time.Duration(0), false, false},
		{jf, false, false},
		{nil, false, false},
		{"false", false, false},
		{"FALSE", false, false},
		{"False", false, false},
		{"f", false, false},
		{"F", false, false},
		{false, false, false},

		{"true", true, false},
		{"TRUE", true, false},
		{"True", true, false},
		{"t", true, false},
		{"T", true, false},
		{1, true, false},
		{int64(1), true, false},
		{int32(1), true, false},
		{int16(1), true, false},
		{int8(1), true, false},
		{uint(1), true, false},
		{uint64(1), true, false},
		{uint32(1), true, false},
		{uint16(1), true, false},
		{uint8(1), true, false},
		{float64(1), true, false},
		{float32(1), true, false},
		{time.Duration(1), true, false},
		{jt, true, false},
		{je, true, false},
		{true, true, false},
		{-1, true, false},
		{int64(-1), true, false},
		{int32(-1), true, false},
		{int16(-1), true, false},
		{int8(-1), true, false},

		// errors
		{"test", false, true},
		{testing.T{}, false, true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		v, err := ToBoolE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, test.expect, errmsg)

		// Non-E test
		v = ToBool(test.input)
		c.Assert(v, qt.Equals, test.expect, errmsg)
	}
}

func BenchmarkTooBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !ToBool(true) {
			b.Fatal("ToBool returned false")
		}
	}
}

func TestIndirectPointers(t *testing.T) {
	c := qt.New(t)

	x := 13
	y := &x
	z := &y

	c.Assert(ToInt(y), qt.Equals, 13)
	c.Assert(ToInt(z), qt.Equals, 13)

}
