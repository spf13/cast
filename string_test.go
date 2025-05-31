// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"html/template"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestToStringE(t *testing.T) {
	c := qt.New(t)

	var jn json.Number
	_ = json.Unmarshal([]byte("8"), &jn)
	type Key struct {
		k string
	}
	key := &Key{"foo"}

	tests := []struct {
		input  interface{}
		expect string
		iserr  bool
	}{
		{int(8), "8", false},
		{int8(8), "8", false},
		{int16(8), "8", false},
		{int32(8), "8", false},
		{int64(8), "8", false},
		{uint(8), "8", false},
		{uint8(8), "8", false},
		{uint16(8), "8", false},
		{uint32(8), "8", false},
		{uint64(8), "8", false},
		{float32(8.31), "8.31", false},
		{float64(8.31), "8.31", false},
		{jn, "8", false},
		{true, "true", false},
		{false, "false", false},
		{nil, "", false},
		{[]byte("one time"), "one time", false},
		{"one more time", "one more time", false},
		{template.HTML("one time"), "one time", false},
		{template.URL("http://somehost.foo"), "http://somehost.foo", false},
		{template.JS("(1+2)"), "(1+2)", false},
		{template.CSS("a"), "a", false},
		{template.HTMLAttr("a"), "a", false},
		// errors
		{testing.T{}, "", true},
		{key, "", true},
	}

	for i, test := range tests {
		errmsg := qt.Commentf("i = %d", i) // assert helper message

		// Non-E test
		v := ToString(test.input)
		c.Assert(v, qt.Equals, test.expect, errmsg)

		// Non-pointer test
		v, err := ToStringE(test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
		} else {
			c.Assert(err, qt.IsNil, errmsg)
			c.Assert(v, qt.Equals, test.expect, errmsg)
		}

		// Pointer test
		v, err = ToStringE(&test.input)
		if test.iserr {
			c.Assert(err, qt.IsNotNil, errmsg)
		} else {
			c.Assert(err, qt.IsNil, errmsg)
			c.Assert(v, qt.Equals, test.expect, errmsg)
		}
	}
}

type foo struct {
	val string
}

func (x foo) String() string {
	return x.val
}

func TestStringerToString(t *testing.T) {
	c := qt.New(t)

	var x foo
	x.val = "bar"
	c.Assert(ToString(x), qt.Equals, "bar", qt.Commentf("non-pointer test"))
	c.Assert(ToString(&x), qt.Equals, "bar", qt.Commentf("pointer test"))
}

type fu struct {
	val string
}

func (x fu) Error() string {
	return x.val
}

func TestErrorToString(t *testing.T) {
	c := qt.New(t)

	var x fu
	x.val = "bar"
	c.Assert(ToString(x), qt.Equals, "bar", qt.Commentf("non-pointer test"))
	c.Assert(ToString(&x), qt.Equals, "bar", qt.Commentf("pointer test"))
}
