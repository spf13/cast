// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast"
	"github.com/spf13/cast/internal"
)

func TestTime(t *testing.T) {
	var ptr *time.Time

	testCases := []testCase{
		{"2009-11-10 23:00:00 +0000 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // Time.String()
		{"Tue Nov 10 23:00:00 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},        // ANSIC
		{"Tue Nov 10 23:00:00 UTC 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},    // UnixDate
		{"Tue Nov 10 23:00:00 +0000 2009", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},  // RubyDate
		{"10 Nov 09 23:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},             // RFC822
		{"10 Nov 09 23:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},           // RFC822Z
		{"Tuesday, 10-Nov-09 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC850
		{"Tue, 10 Nov 2009 23:00:00 UTC", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},   // RFC1123
		{"Tue, 10 Nov 2009 23:00:00 +0000", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false}, // RFC1123Z
		{"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},            // RFC3339
		{"2018-10-21T23:21:29+0200", time.Date(2018, 10, 21, 21, 21, 29, 0, time.UTC), false},      // RFC3339 without timezone hh:mm colon
		{"2009-11-10T23:00:00Z", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), false},            // RFC3339Nano
		{"11:00PM", time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC), false},                              // Kitchen
		{"Nov 10 23:00:00", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},                    // Stamp
		{"Nov 10 23:00:00.000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},                // StampMilli
		{"Nov 10 23:00:00.000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},             // StampMicro
		{"Nov 10 23:00:00.000000000", time.Date(0, 11, 10, 23, 0, 0, 0, time.UTC), false},          // StampNano
		{"2016-03-06 15:28:01-00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},        // RFC3339 without T
		{"2016-03-06 15:28:01-0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},         // RFC3339 without T or timezone hh:mm colon
		{"2016-03-06 15:28:01", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 -0000", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 -00:00", time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 +0900", time.Date(2016, 3, 6, 6, 28, 1, 0, time.UTC), false},
		{"2016-03-06 15:28:01 +09:00", time.Date(2016, 3, 6, 6, 28, 1, 0, time.UTC), false},
		{"2006-01-02", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{"02 Jan 2006", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{1472574600, time.Date(2016, 8, 30, 16, 30, 0, 0, time.UTC), false},
		{int(1482597504), time.Date(2016, 12, 24, 16, 38, 24, 0, time.UTC), false},
		{int64(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{int32(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{uint(1482597504), time.Date(2016, 12, 24, 16, 38, 24, 0, time.UTC), false},
		{uint64(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{uint32(1234567890), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{json.Number("1234567890"), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},
		{time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC), false},

		{ptr, time.Time{}, false},

		// Failure cases
		{"2006", time.Time{}, true},
		{json.Number("123.4567890"), time.Time{}, true},
		{testing.T{}, time.Time{}, true},
	}

	runTests(t, testCases, cast.ToTime, cast.ToTimeE)
}

func TestDuration(t *testing.T) {
	type MyDuration time.Duration

	var ptr *time.Duration

	var expected time.Duration = 5

	testCases := []testCase{
		{time.Duration(5), expected, false},
		{int(5), expected, false},
		{int64(5), expected, false},
		{int32(5), expected, false},
		{int16(5), expected, false},
		{int8(5), expected, false},
		{uint(5), expected, false},
		{uint64(5), expected, false},
		{uint32(5), expected, false},
		{uint16(5), expected, false},
		{uint8(5), expected, false},
		{float64(5), expected, false},
		{float32(5), expected, false},
		{json.Number("5"), expected, false},
		{string("5"), expected, false},
		{string("5ns"), expected, false},
		{string("5us"), time.Microsecond * expected, false},
		{string("5µs"), time.Microsecond * expected, false},
		{string("5ms"), time.Millisecond * expected, false},
		{string("5s"), time.Second * expected, false},
		{string("5m"), time.Minute * expected, false},
		{string("5h"), time.Hour * expected, false},

		{0, time.Duration(0), false},
		{ptr, time.Duration(0), false},

		// Aliases
		{MyInt(5), expected, false},
		{MyString("5"), expected, false},
		{MyDuration(5), expected, false},

		// Failure cases
		{"test", time.Duration(0), true},
		{testing.T{}, time.Duration(0), true},
	}

	runTests(t, testCases, cast.ToDuration, cast.ToDurationE)
}

func TestTimeWithTimezones(t *testing.T) {
	c := qt.New(t)

	est, err := time.LoadLocation("EST")
	c.Assert(err, qt.IsNil)

	irn, err := time.LoadLocation("Iran")
	c.Assert(err, qt.IsNil)

	swd, err := time.LoadLocation("Europe/Stockholm")
	c.Assert(err, qt.IsNil)

	// Test same local time in different timezones
	utc2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	est2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, est)
	irn2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, irn)
	swd2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, swd)
	loc2016 := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local)

	for i, format := range internal.TimeFormats {
		format := format
		if format.Typ == internal.TimeFormatTimeOnly {
			continue
		}

		nameBase := fmt.Sprintf("%d;timeFormatType=%d;%s", i, format.Typ, format.Format)

		t.Run(path.Join(nameBase), func(t *testing.T) {
			est2016str := est2016.Format(format.Format)
			swd2016str := swd2016.Format(format.Format)

			t.Run("without default location", func(t *testing.T) {
				c := qt.New(t)
				converted, err := cast.ToTimeE(est2016str)
				c.Assert(err, qt.IsNil)
				if format.HasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, est2016, converted)
					assertLocationEqual(t, est, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in UTC.
					assertTimeEqual(t, utc2016, converted)
					assertLocationEqual(t, time.UTC, converted.Location())
				}
			})

			t.Run("local timezone without a default location", func(t *testing.T) {
				c := qt.New(t)
				converted, err := cast.ToTimeE(swd2016str)
				c.Assert(err, qt.IsNil)
				if format.HasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, swd2016, converted)
					assertLocationEqual(t, swd, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in UTC.
					assertTimeEqual(t, utc2016, converted)
					assertLocationEqual(t, time.UTC, converted.Location())
				}
			})

			t.Run("nil default location", func(t *testing.T) {
				c := qt.New(t)

				converted, err := cast.ToTimeInDefaultLocationE(est2016str, nil)
				c.Assert(err, qt.IsNil)
				if format.HasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, est2016, converted)
					assertLocationEqual(t, est, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in the local timezone.
					assertTimeEqual(t, loc2016, converted)
					assertLocationEqual(t, time.Local, converted.Location())
				}

			})

			t.Run("default location not UTC", func(t *testing.T) {
				c := qt.New(t)

				converted, err := cast.ToTimeInDefaultLocationE(est2016str, irn)
				c.Assert(err, qt.IsNil)
				if format.HasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, est2016, converted)
					assertLocationEqual(t, est, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in the given location.
					assertTimeEqual(t, irn2016, converted)
					assertLocationEqual(t, irn, converted.Location())
				}

			})

			t.Run("time in the local timezone default location not UTC", func(t *testing.T) {
				c := qt.New(t)

				converted, err := cast.ToTimeInDefaultLocationE(swd2016str, irn)
				c.Assert(err, qt.IsNil)
				if format.HasTimezone() {
					// Converting inputs with a timezone should preserve it
					assertTimeEqual(t, swd2016, converted)
					assertLocationEqual(t, swd, converted.Location())
				} else {
					// Converting inputs without a timezone should be interpreted
					// as a local time in the given location.
					assertTimeEqual(t, irn2016, converted)
					assertLocationEqual(t, irn, converted.Location())
				}

			})

		})

	}
}

func BenchmarkCommonTimeLayouts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, commonLayout := range []string{"2019-04-29", "2017-05-30T00:00:00Z"} {
			_, err := cast.StringToDateInDefaultLocation(commonLayout, time.UTC)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func assertTimeEqual(t *testing.T, expected, actual time.Time) {
	t.Helper()
	// Compare the dates using a numeric zone as there are cases where
	// time.Parse will assign a dummy location.
	qt.Assert(t, actual.Format(time.RFC1123Z), qt.Equals, expected.Format(time.RFC1123Z))
}

func assertLocationEqual(t *testing.T, expected, actual *time.Location) {
	t.Helper()
	qt.Assert(t, locationEqual(expected, actual), qt.IsTrue)
}

func locationEqual(a, b *time.Location) bool {
	// A note about comparring time.Locations:
	//   - can't only compare pointers
	//   - can't compare loc.String() because locations with the same
	//     name can have different offsets
	//   - can't use reflect.DeepEqual because time.Location has internal
	//     caches

	if a == b {
		return true
	} else if a == nil || b == nil {
		return false
	}

	// Check if they're equal by parsing times with a format that doesn't
	// include a timezone, which will interpret it as being a local time in
	// the given zone, and comparing the resulting local times.
	tA, err := time.ParseInLocation("2006-01-02", "2016-01-01", a)
	if err != nil {
		return false
	}

	tB, err := time.ParseInLocation("2006-01-02", "2016-01-01", b)
	if err != nil {
		return false
	}

	return tA.Equal(tB)
}

func TestStringToDate(t *testing.T) {
	c := qt.New(t)

	t.Run("Valid date formats in UTC", func(t *testing.T) {
		d, err := cast.StringToDate("2009-11-10 23:00:00 +0000 UTC")
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC))

		d2, err := cast.StringToDate("2009-11-10")
		c.Assert(err, qt.IsNil)
		c.Assert(d2, qt.Equals, time.Date(2009, 11, 10, 0, 0, 0, 0, time.UTC))

		d3, err := cast.StringToDate("2009-11-10T23:00:00Z")
		c.Assert(err, qt.IsNil)
		c.Assert(d3, qt.Equals, time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC))

		d4, err := cast.StringToDate("2016-03-06 15:28:01-00:00")
		c.Assert(err, qt.IsNil)
		c.Assert(d4.Equal(time.Date(2016, 3, 6, 15, 28, 1, 0, time.UTC)), qt.IsTrue)
	})

	t.Run("Invalid date returns error", func(t *testing.T) {
		_, err := cast.StringToDate("invalid-date")
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.Error(), qt.Equals, "unable to parse date: invalid-date")

		_, err = cast.StringToDate("")
		c.Assert(err, qt.Not(qt.IsNil))
	})
}

func TestToTimeInDefaultLocation(t *testing.T) {
	c := qt.New(t)

	loc, err := time.LoadLocation("America/New_York")
	c.Assert(err, qt.IsNil)

	t.Run("Success with UTC", func(t *testing.T) {
		d := cast.ToTimeInDefaultLocation("2020-01-01", time.UTC)
		c.Assert(d, qt.Equals, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	})

	t.Run("Success with custom location", func(t *testing.T) {
		d := cast.ToTimeInDefaultLocation("2020-01-01 10:00:00", loc)
		c.Assert(d, qt.Equals, time.Date(2020, 1, 1, 10, 0, 0, 0, loc))
	})

	t.Run("Invalid input returns zero time", func(t *testing.T) {
		d := cast.ToTimeInDefaultLocation("invalid-time-val", time.UTC)
		c.Assert(d, qt.Equals, time.Time{})
	})
}

type timeMockFloat64Provider struct {
	val float64
}

func (m timeMockFloat64Provider) Float64() float64 {
	return m.val
}

type timeMockFloat64EProvider struct {
	val float64
	err error
}

func (m timeMockFloat64EProvider) Float64() (float64, error) {
	return m.val, m.err
}

func TestDurationProvidersAndEdgeCases(t *testing.T) {
	c := qt.New(t)

	t.Run("float64Provider success", func(t *testing.T) {
		d, err := cast.ToDurationE(timeMockFloat64Provider{val: 5000})
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Duration(5000))
	})

	t.Run("float64EProvider success", func(t *testing.T) {
		d, err := cast.ToDurationE(timeMockFloat64EProvider{val: 2500, err: nil})
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Duration(2500))
	})

	t.Run("float64EProvider error branch", func(t *testing.T) {
		d, err := cast.ToDurationE(timeMockFloat64EProvider{val: 0, err: errors.New("provider failure")})
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(d, qt.Equals, time.Duration(0))
		c.Assert(err.Error(), qt.Contains, "time.Duration")
	})

	t.Run("String durations edge cases", func(t *testing.T) {
		d, err := cast.ToDurationE("-500")
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Duration(-500))

		d, err = cast.ToDurationE("-5s")
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, -5*time.Second)

		_, err = cast.ToDurationE("5m-invalid")
		c.Assert(err, qt.Not(qt.IsNil))

		_, err = cast.ToDurationE("")
		c.Assert(err, qt.Not(qt.IsNil))
	})

	t.Run("Alias unhandled type fallback", func(t *testing.T) {
		type MyBool bool
		_, err := cast.ToDurationE(MyBool(true))
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.Error(), qt.Contains, "unable to cast true of type bool to time.Duration")

		type unhandledStruct struct{}
		_, err = cast.ToDurationE(unhandledStruct{})
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.Error(), qt.Contains, "to time.Duration")
	})

	t.Run("Negative numbers to duration", func(t *testing.T) {
		d, err := cast.ToDurationE(int(-100))
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Duration(-100))

		d, err = cast.ToDurationE(float64(-200))
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Duration(-200))

		d, err = cast.ToDurationE(int64(-300))
		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, time.Duration(-300))
	})
}

func TestTimeEdgeCases(t *testing.T) {
	c := qt.New(t)

	t.Run("json.Number zero decimal trimming", func(t *testing.T) {
		tm, err := cast.ToTimeE(json.Number("1482597504.000000"))
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Unix(1482597504, 0))
	})

	t.Run("json.Number empty or invalid string", func(t *testing.T) {
		_, err := cast.ToTimeE(json.Number(""))
		c.Assert(err, qt.Not(qt.IsNil))

		_, err = cast.ToTimeE(json.Number("abc"))
		c.Assert(err, qt.Not(qt.IsNil))
	})

	t.Run("Negative timestamp integers", func(t *testing.T) {
		tm, err := cast.ToTimeE(int64(-1000))
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Unix(-1000, 0))

		tm, err = cast.ToTimeE(int(-1000))
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Unix(-1000, 0))

		tm, err = cast.ToTimeE(int32(-1000))
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Unix(-1000, 0))
	})

	t.Run("Pointer to int timestamp", func(t *testing.T) {
		val := 1482597504
		tm, err := cast.ToTimeE(&val)
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Unix(1482597504, 0))

		var nilInt *int
		tm, err = cast.ToTimeE(nilInt)
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Time{})
	})

	t.Run("Pointer to json.Number", func(t *testing.T) {
		jn := json.Number("1482597504")
		tm, err := cast.ToTimeE(&jn)
		c.Assert(err, qt.IsNil)
		c.Assert(tm, qt.Equals, time.Unix(1482597504, 0))
	})

	t.Run("Unsupported types", func(t *testing.T) {
		_, err := cast.ToTimeE(true)
		c.Assert(err, qt.Not(qt.IsNil))

		_, err = cast.ToTimeE(float64(1482597504))
		c.Assert(err, qt.Not(qt.IsNil))

		_, err = cast.ToTimeE([]byte("2020-01-01"))
		c.Assert(err, qt.Not(qt.IsNil))
	})
}

func TestTimeConcurrentSafety(t *testing.T) {
	c := qt.New(t)

	const goroutines = 50
	const iterations = 100
	done := make(chan bool, goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer func() { done <- true }()
			for i := 0; i < iterations; i++ {
				_, err1 := cast.ToTimeE("2023-05-15T10:30:00Z")
				_, err2 := cast.ToDurationE("5s")
				_, err3 := cast.StringToDate("2023-05-15")
				_, err4 := cast.ToDurationE(500)
				d5 := cast.ToTimeInDefaultLocation("2023-05-15 12:00:00", time.UTC)
				if err1 != nil || err2 != nil || err3 != nil || err4 != nil || d5.IsZero() {
					t.Errorf("unexpected error in concurrent test")
					return
				}
			}
		}(g)
	}

	for g := 0; g < goroutines; g++ {
		<-done
	}
	c.Assert(true, qt.IsTrue)
}
