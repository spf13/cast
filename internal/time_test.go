package internal_test

import (
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast/internal"
)

func TestTimeFormatHasTimezone(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		name     string
		format   internal.TimeFormat
		expected bool
	}{
		{
			name:     "NoTimezone",
			format:   internal.TimeFormat{Format: "2006-01-02", Typ: internal.TimeFormatNoTimezone},
			expected: false,
		},
		{
			name:     "NamedTimezone",
			format:   internal.TimeFormat{Format: time.RFC1123, Typ: internal.TimeFormatNamedTimezone},
			expected: false,
		},
		{
			name:     "NumericTimezone",
			format:   internal.TimeFormat{Format: time.RFC3339, Typ: internal.TimeFormatNumericTimezone},
			expected: true,
		},
		{
			name:     "NumericAndNamedTimezone",
			format:   internal.TimeFormat{Format: "2006-01-02 15:04:05.999999999 -0700 MST", Typ: internal.TimeFormatNumericAndNamedTimezone},
			expected: true,
		},
		{
			name:     "TimeOnly",
			format:   internal.TimeFormat{Format: time.Kitchen, Typ: internal.TimeFormatTimeOnly},
			expected: false,
		},
		{
			name:     "NegativeBoundary",
			format:   internal.TimeFormat{Format: "custom", Typ: internal.TimeFormatType(-1)},
			expected: false,
		},
		{
			name:     "UpperBoundBoundary",
			format:   internal.TimeFormat{Format: "custom", Typ: internal.TimeFormatType(5)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Assert(tt.format.HasTimezone(), qt.Equals, tt.expected)
		})
	}

	for i, tf := range internal.TimeFormats {
		t.Run(fmt.Sprintf("PredefinedFormat_%d_%s", i, tf.Format), func(t *testing.T) {
			expected := tf.Typ == internal.TimeFormatNumericTimezone || tf.Typ == internal.TimeFormatNumericAndNamedTimezone
			c.Assert(tf.HasTimezone(), qt.Equals, expected)
		})
	}
}

func TestTimeFormatTypeString(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		typ      internal.TimeFormatType
		expected string
	}{
		{internal.TimeFormatNoTimezone, "TimeFormatNoTimezone"},
		{internal.TimeFormatNamedTimezone, "TimeFormatNamedTimezone"},
		{internal.TimeFormatNumericTimezone, "TimeFormatNumericTimezone"},
		{internal.TimeFormatNumericAndNamedTimezone, "TimeFormatNumericAndNamedTimezone"},
		{internal.TimeFormatTimeOnly, "TimeFormatTimeOnly"},
		{internal.TimeFormatType(-1), "TimeFormatType(-1)"},
		{internal.TimeFormatType(5), "TimeFormatType(5)"},
		{internal.TimeFormatType(100), "TimeFormatType(100)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			c.Assert(tt.typ.String(), qt.Equals, tt.expected)
		})
	}
}

func TestParseDateWith(t *testing.T) {
	c := qt.New(t)

	ny, err := time.LoadLocation("America/New_York")
	c.Assert(err, qt.IsNil)

	t.Run("NoTimezone with nil location defaults to Local", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: "2006-01-02", Typ: internal.TimeFormatNoTimezone},
		}
		res, err := internal.ParseDateWith("2023-05-15", nil, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Location(), qt.Equals, time.Local)
		c.Assert(res.Year(), qt.Equals, 2023)
		c.Assert(res.Month(), qt.Equals, time.Month(5))
		c.Assert(res.Day(), qt.Equals, 15)
	})

	t.Run("NoTimezone with specified location", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: "2006-01-02 15:04:05", Typ: internal.TimeFormatNoTimezone},
		}
		res, err := internal.ParseDateWith("2023-05-15 10:30:00", ny, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Location(), qt.Equals, ny)
		c.Assert(res.Hour(), qt.Equals, 10)
		c.Assert(res.Minute(), qt.Equals, 30)
	})

	t.Run("NamedTimezone with nil location defaults to Local", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: time.RFC1123, Typ: internal.TimeFormatNamedTimezone},
		}
		res, err := internal.ParseDateWith("Mon, 15 May 2023 10:30:00 UTC", nil, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Location(), qt.Equals, time.Local)
		c.Assert(res.Year(), qt.Equals, 2023)
	})

	t.Run("NamedTimezone with specified location", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: time.RFC1123, Typ: internal.TimeFormatNamedTimezone},
		}
		res, err := internal.ParseDateWith("Mon, 15 May 2023 10:30:00 UTC", time.UTC, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Location(), qt.Equals, time.UTC)
	})

	t.Run("NumericTimezone preserves parsed location", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: time.RFC3339, Typ: internal.TimeFormatNumericTimezone},
		}
		res, err := internal.ParseDateWith("2023-05-15T10:30:00+04:00", time.UTC, formats)
		c.Assert(err, qt.IsNil)
		_, offset := res.Zone()
		c.Assert(offset, qt.Equals, 4*3600)
	})

	t.Run("NumericAndNamedTimezone preserves parsed location", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: "2006-01-02 15:04:05.999999999 -0700 MST", Typ: internal.TimeFormatNumericAndNamedTimezone},
		}
		res, err := internal.ParseDateWith("2023-05-15 10:30:00.123456789 -0700 MST", time.UTC, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Nanosecond(), qt.Equals, 123456789)
		_, offset := res.Zone()
		c.Assert(offset, qt.Equals, -7*3600)
	})

	t.Run("TimeOnly format", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: time.Kitchen, Typ: internal.TimeFormatTimeOnly},
		}
		res, err := internal.ParseDateWith("3:04PM", time.UTC, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Hour(), qt.Equals, 15)
		c.Assert(res.Minute(), qt.Equals, 4)
	})

	t.Run("Multiple formats fallback", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: "2006-01-02", Typ: internal.TimeFormatNoTimezone},
			{Format: time.RFC3339, Typ: internal.TimeFormatNumericTimezone},
		}
		res, err := internal.ParseDateWith("2023-05-15T10:30:00Z", time.UTC, formats)
		c.Assert(err, qt.IsNil)
		c.Assert(res.Year(), qt.Equals, 2023)
	})

	t.Run("Parse failure unmatched format", func(t *testing.T) {
		formats := []internal.TimeFormat{
			{Format: "2006-01-02", Typ: internal.TimeFormatNoTimezone},
		}
		_, err := internal.ParseDateWith("not-a-date", time.UTC, formats)
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.Error(), qt.Equals, "unable to parse date: not-a-date")
	})

	t.Run("Empty formats slice", func(t *testing.T) {
		_, err := internal.ParseDateWith("2023-05-15", time.UTC, nil)
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.Error(), qt.Equals, "unable to parse date: 2023-05-15")
	})

	t.Run("Empty string input", func(t *testing.T) {
		_, err := internal.ParseDateWith("", time.UTC, internal.TimeFormats)
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.Error(), qt.Equals, "unable to parse date: ")
	})
}
