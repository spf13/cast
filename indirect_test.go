// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestIndirectPointers(t *testing.T) {
	c := qt.New(t)

	x := 13
	y := &x
	z := &y

	c.Assert(ToInt(y), qt.Equals, 13)
	c.Assert(ToInt(z), qt.Equals, 13)
}
