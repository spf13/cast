// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ToTimeE casts an interface to a time.Time type.
func ToTimeE(i interface{}) (tim time.Time, err error) {
	return ToTimeInDefaultLocationE(i, time.UTC)
}

// ToTimeInDefaultLocationE casts an empty interface to time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func ToTimeInDefaultLocationE(i interface{}, location *time.Location) (tim time.Time, err error) {
	i = indirect(i)

	switch v := i.(type) {
	case time.Time:
		return v, nil
	case string:
		return StringToDateInDefaultLocation(v, location)
	case json.Number:
		s, err1 := ToInt64E(v)
		if err1 != nil {
			return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
		}
		return time.Unix(s, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case uint32:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
	}
}

// ToDurationE casts an interface to a time.Duration type.
func ToDurationE(i interface{}) (d time.Duration, err error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(ToInt64(s))
		return
	case float32, float64:
		d = time.Duration(ToFloat64(s))
		return
	case string:
		if strings.ContainsAny(s, "nsuµmh") {
			d, err = time.ParseDuration(s)
		} else {
			d, err = time.ParseDuration(s + "ns")
		}
		return
	case float64EProvider:
		var v float64
		v, err = s.Float64()
		d = time.Duration(v)
		return
	case float64Provider:
		d = time.Duration(s.Float64())
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	}
}

// ToDurationSliceE casts an interface to a []time.Duration type.
func ToDurationSliceE(i interface{}) ([]time.Duration, error) {
	return toSliceE(i, ToDurationE)
}

// StringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, time.UTC, timeFormats)
}

// StringToDateInDefaultLocation casts an empty interface to a time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func StringToDateInDefaultLocation(s string, location *time.Location) (time.Time, error) {
	return parseDateWith(s, location, timeFormats)
}

type timeFormatType int

const (
	timeFormatNoTimezone timeFormatType = iota
	timeFormatNamedTimezone
	timeFormatNumericTimezone
	timeFormatNumericAndNamedTimezone
	timeFormatTimeOnly
)

type timeFormat struct {
	format string
	typ    timeFormatType
}

func (f timeFormat) hasTimezone() bool {
	// We don't include the formats with only named timezones, see
	// https://github.com/golang/go/issues/19694#issuecomment-289103522
	return f.typ >= timeFormatNumericTimezone && f.typ <= timeFormatNumericAndNamedTimezone
}

var timeFormats = []timeFormat{
	// Keep common formats at the top.
	{"2006-01-02", timeFormatNoTimezone},
	{time.RFC3339, timeFormatNumericTimezone},
	{"2006-01-02T15:04:05", timeFormatNoTimezone}, // iso8601 without timezone
	{time.RFC1123Z, timeFormatNumericTimezone},
	{time.RFC1123, timeFormatNamedTimezone},
	{time.RFC822Z, timeFormatNumericTimezone},
	{time.RFC822, timeFormatNamedTimezone},
	{time.RFC850, timeFormatNamedTimezone},
	{"2006-01-02 15:04:05.999999999 -0700 MST", timeFormatNumericAndNamedTimezone}, // Time.String()
	{"2006-01-02T15:04:05-0700", timeFormatNumericTimezone},                        // RFC3339 without timezone hh:mm colon
	{"2006-01-02 15:04:05Z0700", timeFormatNumericTimezone},                        // RFC3339 without T or timezone hh:mm colon
	{"2006-01-02 15:04:05", timeFormatNoTimezone},
	{time.ANSIC, timeFormatNoTimezone},
	{time.UnixDate, timeFormatNamedTimezone},
	{time.RubyDate, timeFormatNumericTimezone},
	{"2006-01-02 15:04:05Z07:00", timeFormatNumericTimezone},
	{"02 Jan 2006", timeFormatNoTimezone},
	{"2006-01-02 15:04:05 -07:00", timeFormatNumericTimezone},
	{"2006-01-02 15:04:05 -0700", timeFormatNumericTimezone},
	{time.Kitchen, timeFormatTimeOnly},
	{time.Stamp, timeFormatTimeOnly},
	{time.StampMilli, timeFormatTimeOnly},
	{time.StampMicro, timeFormatTimeOnly},
	{time.StampNano, timeFormatTimeOnly},
}

func parseDateWith(s string, location *time.Location, formats []timeFormat) (d time.Time, e error) {
	for _, format := range formats {
		if d, e = time.Parse(format.format, s); e == nil {

			// Some time formats have a zone name, but no offset, so it gets
			// put in that zone name (not the default one passed in to us), but
			// without that zone's offset. So set the location manually.
			if format.typ <= timeFormatNamedTimezone {
				if location == nil {
					location = time.Local
				}
				year, month, day := d.Date()
				hour, min, sec := d.Clock()
				d = time.Date(year, month, day, hour, min, sec, d.Nanosecond(), location)
			}

			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}
