//go:build go1.27

// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast"
)

// casterNumberContexts mirrors numberContexts with Caster-based adapters.
//
// It is keyed by the same type names so TestCasterEquivalence can pair each
// entry with its existing-function counterpart.
var casterNumberContexts = map[string]struct {
	to    func(any) any
	toErr func(any) (any, error)
}{
	"int":     {to: toAny(cast.Default.To[int]), toErr: toAnyErr(cast.Default.ToE[int])},
	"int8":    {to: toAny(cast.Default.To[int8]), toErr: toAnyErr(cast.Default.ToE[int8])},
	"int16":   {to: toAny(cast.Default.To[int16]), toErr: toAnyErr(cast.Default.ToE[int16])},
	"int32":   {to: toAny(cast.Default.To[int32]), toErr: toAnyErr(cast.Default.ToE[int32])},
	"int64":   {to: toAny(cast.Default.To[int64]), toErr: toAnyErr(cast.Default.ToE[int64])},
	"uint":    {to: toAny(cast.Default.To[uint]), toErr: toAnyErr(cast.Default.ToE[uint])},
	"uint8":   {to: toAny(cast.Default.To[uint8]), toErr: toAnyErr(cast.Default.ToE[uint8])},
	"uint16":  {to: toAny(cast.Default.To[uint16]), toErr: toAnyErr(cast.Default.ToE[uint16])},
	"uint32":  {to: toAny(cast.Default.To[uint32]), toErr: toAnyErr(cast.Default.ToE[uint32])},
	"uint64":  {to: toAny(cast.Default.To[uint64]), toErr: toAnyErr(cast.Default.ToE[uint64])},
	"float32": {to: toAny(cast.Default.To[float32]), toErr: toAnyErr(cast.Default.ToE[float32])},
	"float64": {to: toAny(cast.Default.To[float64]), toErr: toAnyErr(cast.Default.ToE[float64])},
}

// casterKnownDivergences lists the sample slots (see numberContext.samples)
// where Caster intentionally differs from the existing functions.
//
// The existing parseInt/parseUint helpers call strconv with bitSize 0 for every
// integer width, so out-of-range strings for narrow types wrap silently
// (ToInt8E("128") == -128, nil). Caster passes the real bit size and returns a
// range error instead. Slot 9 is the underflow string, slot 10 the overflow
// string. Unsigned underflow ("-1") is a syntax error on both sides and is
// therefore equivalent, which is why it is not listed.
//
// If this test starts failing because the existing functions began returning
// errors, the divergence is gone and this table should be deleted.
var casterKnownDivergences = map[string][]int{
	"int8":   {9, 10},
	"int16":  {9, 10},
	"int32":  {9, 10},
	"uint8":  {10},
	"uint16": {10},
	"uint32": {10},
}

// casterToAnyErr is toAnyErr for every Basic type, not just Number.
func casterToAnyErr[T cast.Basic](fn func(any) (T, error)) func(any) (any, error) {
	return func(i any) (any, error) { return fn(i) }
}

// assertSameOutcome asserts that a Caster-based function and an existing
// package function return the same value and the same error (nil-ness and
// text) for the given input.
func assertSameOutcome(c *qt.C, input any, caster, existing func(any) (any, error)) {
	got, gotErr := caster(input)
	want, wantErr := existing(input)

	comment := qt.Commentf("input: %#v (%T)", input, input)

	if wantErr != nil {
		c.Assert(gotErr, qt.IsNotNil, comment)
		c.Assert(gotErr.Error(), qt.Equals, wantErr.Error(), comment)
	} else {
		c.Assert(gotErr, qt.IsNil, comment)
	}

	c.Assert(got, qt.Equals, want, comment)
}

// assertSameTime asserts that two time conversions agree on error, instant and
// zone. time.Time values are not compared with == because time.Parse allocates
// a fresh *Location for numeric offsets.
func assertSameTime(c *qt.C, input any, caster, existing func(any) (time.Time, error)) {
	got, gotErr := caster(input)
	want, wantErr := existing(input)

	comment := qt.Commentf("input: %#v (%T)", input, input)

	if wantErr != nil {
		c.Assert(gotErr, qt.IsNotNil, comment)
		c.Assert(gotErr.Error(), qt.Equals, wantErr.Error(), comment)
	} else {
		c.Assert(gotErr, qt.IsNil, comment)
	}

	c.Assert(got.Equal(want), qt.IsTrue, qt.Commentf("input: %#v (%T): got %v, want %v", input, input, got, want))
	c.Assert(got.Location().String(), qt.Equals, want.Location().String(), comment)

	gotName, gotOffset := got.Zone()
	wantName, wantOffset := want.Zone()

	c.Assert(gotName, qt.Equals, wantName, comment)
	c.Assert(gotOffset, qt.Equals, wantOffset, comment)
}

func contains(haystack []int, needle int) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func TestCasterEquivalence(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	c.Assert(len(casterNumberContexts), qt.Equals, len(numberContexts))

	for typeName, ctx := range numberContexts {
		casterCtx, ok := casterNumberContexts[typeName]
		c.Assert(ok, qt.IsTrue, qt.Commentf("missing caster context for %q", typeName))

		t.Run(typeName, func(t *testing.T) {
			t.Parallel()

			divergent := casterKnownDivergences[typeName]

			// Raw samples: typed values, aliases, min/max and the boundary strings.
			t.Run("Samples", func(t *testing.T) {
				t.Parallel()

				c := qt.New(t)

				for idx, sample := range ctx.samples {
					if sample == nil || contains(divergent, idx) {
						continue
					}

					assertSameOutcome(c, sample, casterCtx.toErr, ctx.toErr)
					assertSameOutcome(c, &sample, casterCtx.toErr, ctx.toErr)

					c.Assert(casterCtx.to(sample), qt.Equals, ctx.to(sample), qt.Commentf("sample %d: %#v", idx, sample))
					c.Assert(casterCtx.to(&sample), qt.Equals, ctx.to(&sample), qt.Commentf("sample %d: %#v", idx, sample))
				}
			})

			// The full corpus TestNumber runs the existing functions against:
			// strings, json.Number, bools, nil, aliases, decimals, failure cases.
			t.Run("Corpus", func(t *testing.T) {
				t.Parallel()

				c := qt.New(t)

				for _, testCase := range generateNumberTestCases(ctx.samples) {
					input := testCase.input

					assertSameOutcome(c, input, casterCtx.toErr, ctx.toErr)
					assertSameOutcome(c, &input, casterCtx.toErr, ctx.toErr)

					c.Assert(casterCtx.to(input), qt.Equals, ctx.to(input), qt.Commentf("input: %#v", input))
					c.Assert(casterCtx.to(&input), qt.Equals, ctx.to(&input), qt.Commentf("input: %#v", input))
				}
			})

			if len(divergent) == 0 {
				return
			}

			// Document the known divergence loudly rather than hiding it.
			t.Run("KnownDivergence", func(t *testing.T) {
				t.Parallel()

				c := qt.New(t)

				zero := ctx.samples[0]

				for _, idx := range divergent {
					sample := ctx.samples[idx]
					c.Assert(sample, qt.IsNotNil, qt.Commentf("divergent slot %d must hold a sample", idx))

					_, existingErr := ctx.toErr(sample)
					c.Assert(existingErr, qt.IsNil, qt.Commentf("existing %s function no longer wraps %#v; delete this divergence", typeName, sample))

					got, casterErr := casterCtx.toErr(sample)
					c.Assert(casterErr, qt.ErrorMatches, `.*value out of range`, qt.Commentf("sample %#v", sample))
					c.Assert(got, qt.Equals, zero, qt.Commentf("sample %#v", sample))
					c.Assert(casterCtx.to(sample), qt.Equals, zero, qt.Commentf("sample %#v", sample))
				}
			})
		})
	}
}

func TestCasterNonNumeric(t *testing.T) {
	t.Parallel()

	type MyDuration time.Duration

	var (
		nilIntPtr    *int
		nilStringPtr *string
		someInt      = 5
		someString   = "5"
		someBool     = true
		someDuration = 5 * time.Second
	)

	corpus := []any{
		nil,
		"", "true", "false", "T", "F", "1", "0", "5", "-5", "8.31", "5ns", "5s", "5m", "test",
		true, false,
		int(0), int(1), int(5), int8(-5), int16(5), int32(5), int64(5),
		uint(5), uint8(5), uint16(5), uint32(5), uint64(5),
		float32(8.31), float64(8.31), float64(0),
		json.Number("5"), json.Number("8.31"), json.Number(""), json.Number("true"),
		time.Duration(5), someDuration, time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC),
		time.Month(5), time.Weekday(5),
		MyString("5"), MyString("true"), MyBool(true), MyInt(5), MyInt64(-5), MyUint8(5), MyFloat64(8.31), MyDuration(5),
		nilIntPtr, nilStringPtr, &someInt, &someString, &someBool, &someDuration,
		[]byte("hello"), errors.New("boom"),
		testing.T{}, []string{"a"}, map[string]any{}, struct{}{},
	}

	contexts := map[string]struct {
		caster   func(any) (any, error)
		existing func(any) (any, error)
	}{
		"string":   {caster: casterToAnyErr(cast.Default.ToE[string]), existing: casterToAnyErr(cast.ToStringE)},
		"bool":     {caster: casterToAnyErr(cast.Default.ToE[bool]), existing: casterToAnyErr(cast.ToBoolE)},
		"duration": {caster: casterToAnyErr(cast.Default.ToE[time.Duration]), existing: casterToAnyErr(cast.ToDurationE)},
	}

	for typeName, ctx := range contexts {
		t.Run(typeName, func(t *testing.T) {
			t.Parallel()

			c := qt.New(t)

			for _, input := range corpus {
				assertSameOutcome(c, input, ctx.caster, ctx.existing)
				assertSameOutcome(c, &input, ctx.caster, ctx.existing)
			}
		})
	}

	t.Run("To", func(t *testing.T) {
		t.Parallel()

		c := qt.New(t)

		for _, input := range corpus {
			c.Assert(cast.Default.To[string](input), qt.Equals, cast.ToString(input), qt.Commentf("input: %#v", input))
			c.Assert(cast.Default.To[bool](input), qt.Equals, cast.ToBool(input), qt.Commentf("input: %#v", input))
			c.Assert(cast.Default.To[time.Duration](input), qt.Equals, cast.ToDuration(input), qt.Commentf("input: %#v", input))
		}
	})
}

func TestCasterTime(t *testing.T) {
	t.Parallel()

	var nilTimePtr *time.Time

	someTime := time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC)

	corpus := []any{
		"2009-11-10T23:00:00Z",          // RFC3339
		"2018-10-21T23:21:29+0200",      // RFC3339 without timezone colon
		"2016-03-06 15:28:01 +09:00",    // offset
		"Tue Nov 10 23:00:00 UTC 2009",  // named zone
		"Tue, 10 Nov 2009 23:00:00 UTC", // RFC1123
		"2006-01-02",                    // bare date, no zone
		"2016-03-06 15:28:01",           // datetime, no zone
		"02 Jan 2006",                   // no zone
		"11:00PM",                       // Kitchen, no zone
		"Nov 10 23:00:00",               // Stamp, no zone
		int(1482597504), int32(1234567890), int64(1234567890),
		uint(1482597504), uint32(1234567890), uint64(1234567890),
		json.Number("1234567890"), json.Number("1234567890.0"), json.Number("123.4567890"), json.Number(""),
		nil, nilTimePtr, someTime, &someTime,
		MyString("2006-01-02"),

		// Failure cases
		"2006", "test", "", testing.T{}, float64(1.5), true, int8(5),
	}

	t.Run("Default", func(t *testing.T) {
		t.Parallel()

		c := qt.New(t)

		for _, input := range corpus {
			assertSameTime(c, input, cast.Default.ToE[time.Time], cast.ToTimeE)
			assertSameTime(c, &input, cast.Default.ToE[time.Time], cast.ToTimeE)

			got := cast.Default.To[time.Time](input)
			want := cast.ToTime(input)
			c.Assert(got.Equal(want), qt.IsTrue, qt.Commentf("input: %#v: got %v, want %v", input, got, want))
			c.Assert(got.Location().String(), qt.Equals, want.Location().String(), qt.Commentf("input: %#v", input))
		}
	})

	locations := map[string]*time.Location{
		"UTC":   time.UTC,
		"Local": time.Local,
		"Fixed": time.FixedZone("UTC+9", 9*60*60),
	}

	for name, loc := range locations {
		t.Run("WithLocation/"+name, func(t *testing.T) {
			t.Parallel()

			c := qt.New(t)

			caster := cast.Default.WithLocation(loc)
			existing := func(i any) (time.Time, error) { return cast.ToTimeInDefaultLocationE(i, loc) }

			for _, input := range corpus {
				assertSameTime(c, input, caster.ToE[time.Time], existing)
				assertSameTime(c, &input, caster.ToE[time.Time], existing)
			}
		})
	}

	// ToTimeInDefaultLocationE treats a nil location as LOCAL time. The Caster
	// zero value must instead behave as UTC, otherwise Caster{} would not be
	// equivalent to ToTimeE. Location().String() distinguishes Local from UTC
	// even on machines whose local zone has a zero offset.
	t.Run("NilLocationIsUTC", func(t *testing.T) {
		t.Parallel()

		c := qt.New(t)

		const input = "2006-01-02 15:04:05"

		expected := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

		for name, caster := range map[string]cast.Caster{
			"Zero":               {},
			"Default":            cast.Default,
			"WithLocation(nil)":  cast.Default.WithLocation(nil),
			"ResetAfterLocal":    cast.Default.WithLocation(time.Local).WithLocation(nil),
			"ResetAfterFixed":    cast.Default.WithLocation(locations["Fixed"]).WithLocation(nil),
			"WithLocation(UTC)":  cast.Default.WithLocation(time.UTC),
			"WithBaseThenNilLoc": cast.Default.WithBase(16).WithLocation(nil),
		} {
			got, err := caster.ToE[time.Time](input)
			c.Assert(err, qt.IsNil, qt.Commentf("%s", name))
			c.Assert(got, qt.Equals, expected, qt.Commentf("%s: location must be UTC, got %q", name, got.Location()))
			c.Assert(got.Location().String(), qt.Equals, "UTC", qt.Commentf("%s", name))
			c.Assert(got.Location(), qt.Equals, time.UTC, qt.Commentf("%s", name))

			// Must match ToTimeE, not ToTimeInDefaultLocationE(_, nil).
			want, err := cast.ToTimeE(input)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.Equals, want, qt.Commentf("%s", name))
		}

		// Sanity check the trap this guards against: nil means Local downstream.
		local, err := cast.ToTimeInDefaultLocationE(input, nil)
		c.Assert(err, qt.IsNil)
		c.Assert(local.Location(), qt.Equals, time.Local)
	})
}

func TestCasterImmutability(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	const (
		octalLooking = "08" // invalid in base 0 (leading zero selects octal), 8 in base 10 and 16
		hex          = "ff"
		zoneless     = "2006-01-02 15:04:05"
	)

	fixed := time.FixedZone("UTC+3", 3*60*60)

	t.Run("ZeroValueIsDefault", func(t *testing.T) {
		t.Parallel()

		c := qt.New(t)

		c.Assert(cast.Caster{}, qt.Equals, cast.Default)

		for _, input := range []any{octalLooking, hex, "8", "8.31", 8, nil, true, "test", zoneless, "5s"} {
			assertSameOutcome(c, input, casterToAnyErr(cast.Caster{}.ToE[int]), casterToAnyErr(cast.Default.ToE[int]))
			assertSameOutcome(c, input, casterToAnyErr(cast.Caster{}.ToE[uint16]), casterToAnyErr(cast.Default.ToE[uint16]))
			assertSameOutcome(c, input, casterToAnyErr(cast.Caster{}.ToE[float64]), casterToAnyErr(cast.Default.ToE[float64]))
			assertSameOutcome(c, input, casterToAnyErr(cast.Caster{}.ToE[string]), casterToAnyErr(cast.Default.ToE[string]))
			assertSameOutcome(c, input, casterToAnyErr(cast.Caster{}.ToE[bool]), casterToAnyErr(cast.Default.ToE[bool]))
			assertSameOutcome(c, input, casterToAnyErr(cast.Caster{}.ToE[time.Duration]), casterToAnyErr(cast.Default.ToE[time.Duration]))
			assertSameTime(c, input, cast.Caster{}.ToE[time.Time], cast.Default.ToE[time.Time])
		}
	})

	t.Run("WithBase", func(t *testing.T) {
		t.Parallel()

		c := qt.New(t)

		original := cast.Default
		derived := original.WithBase(16)

		// The receiver is untouched and still equivalent to the package functions.
		c.Assert(original, qt.Equals, cast.Default)
		c.Assert(original, qt.Equals, cast.Caster{})
		c.Assert(original.To[int](octalLooking), qt.Equals, cast.ToInt(octalLooking))
		c.Assert(original.To[int](octalLooking), qt.Equals, 0)
		c.Assert(original.To[int](hex), qt.Equals, 0)

		// The copy carries the new base.
		c.Assert(derived, qt.Not(qt.Equals), original)
		c.Assert(derived.To[int](octalLooking), qt.Equals, 8)
		c.Assert(derived.To[int](hex), qt.Equals, 255)
		c.Assert(derived.To[uint8](hex), qt.Equals, uint8(255))

		// Chaining further options preserves the base.
		chained := derived.WithLocation(fixed)
		c.Assert(chained.To[int](hex), qt.Equals, 255)
		c.Assert(derived.To[int](hex), qt.Equals, 255)

		// Rebinding to a different base does not leak back.
		decimal := derived.WithBase(10)
		c.Assert(decimal.To[int](octalLooking), qt.Equals, 8)
		c.Assert(decimal.To[int](hex), qt.Equals, 0)
		c.Assert(derived.To[int](hex), qt.Equals, 255)

		// The package-level Default was never mutated.
		c.Assert(cast.Default, qt.Equals, cast.Caster{})
	})

	t.Run("WithLocation", func(t *testing.T) {
		t.Parallel()

		c := qt.New(t)

		original := cast.Default
		derived := original.WithLocation(fixed)

		c.Assert(original, qt.Equals, cast.Default)
		c.Assert(original, qt.Equals, cast.Caster{})
		c.Assert(original.To[time.Time](zoneless).Location().String(), qt.Equals, "UTC")

		c.Assert(derived, qt.Not(qt.Equals), original)
		c.Assert(derived.To[time.Time](zoneless).Location(), qt.Equals, fixed)
		c.Assert(derived.To[time.Time](zoneless), qt.Equals, time.Date(2006, 1, 2, 15, 4, 5, 0, fixed))

		// Chaining further options preserves the location.
		chained := derived.WithBase(16)
		c.Assert(chained.To[time.Time](zoneless).Location(), qt.Equals, fixed)
		c.Assert(chained.To[int](hex), qt.Equals, 255)
		c.Assert(derived.To[int](hex), qt.Equals, 0)

		c.Assert(cast.Default, qt.Equals, cast.Caster{})
	})

	c.Assert(cast.Default, qt.Equals, cast.Caster{})
}

func TestCasterInvalidBase(t *testing.T) {
	t.Parallel()

	for _, base := range []int{1, 37, -1} {
		caster := cast.Default.WithBase(base)

		t.Run(strconv.Itoa(base), func(t *testing.T) {
			t.Parallel()

			// String-like inputs reach strconv and surface its base error.
			t.Run("String", func(t *testing.T) {
				t.Parallel()

				c := qt.New(t)

				for _, input := range []any{"8", "-8", "+8", "0x08", json.Number("8"), MyString("8")} {
					var numErr *strconv.NumError

					comment := qt.Commentf("input: %#v", input)

					v, err := caster.ToE[int](input)
					c.Assert(err, qt.ErrorMatches, `unable to cast .* of type .* to int: strconv\.ParseInt: parsing .*: invalid base `+strconv.Itoa(base), comment)
					c.Assert(errors.As(err, &numErr), qt.IsTrue, comment)
					c.Assert(v, qt.Equals, 0, comment)
					c.Assert(caster.To[int](input), qt.Equals, 0, comment)

					// strconv validates the base before it looks at the digits or the
					// sign, so even "-8" reports an invalid base for unsigned targets.
					u, err := caster.ToE[uint64](input)
					c.Assert(err, qt.ErrorMatches, `unable to cast .* of type .* to uint64: strconv\.ParseUint: parsing .*: invalid base `+strconv.Itoa(base), comment)
					c.Assert(errors.As(err, &numErr), qt.IsTrue, comment)
					c.Assert(u, qt.Equals, uint64(0), comment)
					c.Assert(caster.To[uint64](input), qt.Equals, uint64(0), comment)
				}

				// Floats never consult the base.
				f, err := caster.ToE[float64]("8.31")
				c.Assert(err, qt.IsNil)
				c.Assert(f, qt.Equals, 8.31)
			})

			// Non-string inputs never reach strconv, so the base is simply unused.
			t.Run("NonString", func(t *testing.T) {
				t.Parallel()

				c := qt.New(t)

				for _, testCase := range []struct {
					input    any
					expected int
				}{
					{int(8), 8},
					{int8(-8), -8},
					{uint16(8), 8},
					{float64(8.31), 8},
					{true, 1},
					{false, 0},
					{nil, 0},
					{MyInt(8), 8},
					{time.Month(8), 8},
					{"", 0}, // empty string short-circuits before parsing
				} {
					comment := qt.Commentf("input: %#v", testCase.input)

					v, err := caster.ToE[int](testCase.input)
					c.Assert(err, qt.IsNil, comment)
					c.Assert(v, qt.Equals, testCase.expected, comment)
					c.Assert(caster.To[int](testCase.input), qt.Equals, testCase.expected, comment)
				}

				u, err := caster.ToE[uint](uint(8))
				c.Assert(err, qt.IsNil)
				c.Assert(u, qt.Equals, uint(8))

				// Non-numeric targets are unaffected too.
				s, err := caster.ToE[string](8)
				c.Assert(err, qt.IsNil)
				c.Assert(s, qt.Equals, "8")

				b, err := caster.ToE[bool]("true")
				c.Assert(err, qt.IsNil)
				c.Assert(b, qt.IsTrue)
			})
		})
	}
}

func BenchmarkCaster(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		b.Run("Existing", func(b *testing.B) {
			for b.Loop() {
				cast.ToInt("42")
			}
		})

		b.Run("Caster", func(b *testing.B) {
			for b.Loop() {
				cast.Default.To[int]("42")
			}
		})
	})

	b.Run("Int64", func(b *testing.B) {
		b.Run("Existing", func(b *testing.B) {
			for b.Loop() {
				cast.ToInt64("42")
			}
		})

		b.Run("Caster", func(b *testing.B) {
			for b.Loop() {
				cast.Default.To[int64]("42")
			}
		})
	})

	b.Run("Float64", func(b *testing.B) {
		b.Run("Existing", func(b *testing.B) {
			for b.Loop() {
				cast.ToFloat64("8.31")
			}
		})

		b.Run("Caster", func(b *testing.B) {
			for b.Loop() {
				cast.Default.To[float64]("8.31")
			}
		})
	})
}
