// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/cast"
)

// TestPlusFunctions tests the "plus functions" that provide fallback logic
func TestPlusFunctions(t *testing.T) {
	t.Parallel()

	t.Run("ToBoolP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion - fallback should not be called
		result := cast.ToBoolP(func(i any) (bool, error) {
			t.Error("Fallback should not be called for successful conversion")
			return false, nil
		}, true)
		c.Assert(result, qt.Equals, true)

		// Test failed conversion - fallback should be called
		fallbackCalled := false
		result = cast.ToBoolP(func(i any) (bool, error) {
			fallbackCalled = true
			return true, nil
		}, "invalid")
		c.Assert(result, qt.Equals, true)
		c.Assert(fallbackCalled, qt.Equals, true)

		// Test nil fallback function
		result = cast.ToBoolP(nil, "invalid")
		c.Assert(result, qt.Equals, false) // Should return zero value
	})

	t.Run("ToStringP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		result := cast.ToStringP(func(i any) (string, error) {
			t.Error("Fallback should not be called for successful conversion")
			return "", nil
		}, 123)
		c.Assert(result, qt.Equals, "123")

		// Test failed conversion with fallback
		fallbackCalled := false
		result = cast.ToStringP(func(i any) (string, error) {
			fallbackCalled = true
			return "fallback", nil
		}, make(chan int))
		c.Assert(result, qt.Equals, "fallback")
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToIntP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		result := cast.ToIntP(func(i any) (int, error) {
			t.Error("Fallback should not be called for successful conversion")
			return 0, nil
		}, "42")
		c.Assert(result, qt.Equals, 42)

		// Test failed conversion with fallback
		fallbackCalled := false
		result = cast.ToIntP(func(i any) (int, error) {
			fallbackCalled = true
			return 999, nil
		}, "invalid")
		c.Assert(result, qt.Equals, 999)
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToTimeP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		now := time.Now()
		result := cast.ToTimeP(func(i any) (time.Time, error) {
			t.Error("Fallback should not be called for successful conversion")
			return time.Time{}, nil
		}, now)
		c.Assert(result.UTC(), qt.Equals, now.UTC())

		// Test failed conversion with fallback
		fallbackTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		fallbackCalled := false
		result = cast.ToTimeP(func(i any) (time.Time, error) {
			fallbackCalled = true
			return fallbackTime, nil
		}, "invalid")
		c.Assert(result.UTC(), qt.Equals, fallbackTime.UTC())
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToTimeInDefaultLocationP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		loc := time.FixedZone("TEST", 3600)
		
		// Test successful conversion
		result := cast.ToTimeInDefaultLocationP(func(i any, location *time.Location) (time.Time, error) {
			t.Error("Fallback should not be called for successful conversion")
			return time.Time{}, nil
		}, "2023-01-01", loc)
		c.Assert(result.Year(), qt.Equals, 2023)

		// Test failed conversion with fallback
		fallbackTime := time.Date(2024, 1, 1, 0, 0, 0, 0, loc)
		fallbackCalled := false
		result = cast.ToTimeInDefaultLocationP(func(i any, location *time.Location) (time.Time, error) {
			fallbackCalled = true
			c.Assert(location, qt.Equals, loc)
			return fallbackTime, nil
		}, "invalid", loc)
		c.Assert(result, qt.Equals, fallbackTime)
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToDurationP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		result := cast.ToDurationP(func(i any) (time.Duration, error) {
			t.Error("Fallback should not be called for successful conversion")
			return 0, nil
		}, "5s")
		c.Assert(result, qt.Equals, 5*time.Second)

		// Test failed conversion with fallback
		fallbackDuration := 10 * time.Minute
		fallbackCalled := false
		result = cast.ToDurationP(func(i any) (time.Duration, error) {
			fallbackCalled = true
			return fallbackDuration, nil
		}, "invalid")
		c.Assert(result, qt.Equals, fallbackDuration)
		c.Assert(fallbackCalled, qt.Equals, true)
	})
}

// TestGenericPlusFunctions tests the generic ToP and ToPE functions
func TestGenericPlusFunctions(t *testing.T) {
	t.Parallel()

	t.Run("ToP[int]", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		result := cast.ToP[int](func(i any) (int, error) {
			t.Error("Fallback should not be called for successful conversion")
			return 0, nil
		}, "42")
		c.Assert(result, qt.Equals, 42)

		// Test failed conversion with fallback
		fallbackCalled := false
		result = cast.ToP[int](func(i any) (int, error) {
			fallbackCalled = true
			return 999, nil
		}, "invalid")
		c.Assert(result, qt.Equals, 999)
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToPE[string]", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		result, err := cast.ToPE[string](func(i any) (string, error) {
			t.Error("Fallback should not be called for successful conversion")
			return "", nil
		}, 123)
		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, "123")

		// Test failed conversion with fallback
		fallbackCalled := false
		result, err = cast.ToPE[string](func(i any) (string, error) {
			fallbackCalled = true
			return "fallback", nil
		}, make(chan int))
		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, "fallback")
		c.Assert(fallbackCalled, qt.Equals, true)

		// Test fallback returning error
		result, err = cast.ToPE[string](func(i any) (string, error) {
			return "", fmt.Errorf("fallback error")
		}, make(chan int))
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.Error(), qt.Equals, "fallback error")
	})

	t.Run("ToPE[bool] with nil fallback", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test with nil fallback - should return original error
		result, err := cast.ToPE[bool](nil, "invalid")
		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.Equals, false)
	})
}

// TestMapPlusFunctions tests the map-related plus functions
func TestMapPlusFunctions(t *testing.T) {
	t.Parallel()

	t.Run("ToStringMapStringP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		input := map[string]string{"key": "value"}
		result := cast.ToStringMapStringP(func(i any) (map[string]string, error) {
			t.Error("Fallback should not be called for successful conversion")
			return nil, nil
		}, input)
		c.Assert(result, qt.DeepEquals, input)

		// Test failed conversion with fallback
		fallbackMap := map[string]string{"fallback": "value"}
		fallbackCalled := false
		result = cast.ToStringMapStringP(func(i any) (map[string]string, error) {
			fallbackCalled = true
			return fallbackMap, nil
		}, "invalid")
		c.Assert(result, qt.DeepEquals, fallbackMap)
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToStringMapIntP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		input := map[string]int{"key": 42}
		result := cast.ToStringMapIntP(func(i any) (map[string]int, error) {
			t.Error("Fallback should not be called for successful conversion")
			return nil, nil
		}, input)
		c.Assert(result, qt.DeepEquals, input)

		// Test failed conversion with fallback
		fallbackMap := map[string]int{"fallback": 999}
		fallbackCalled := false
		result = cast.ToStringMapIntP(func(i any) (map[string]int, error) {
			fallbackCalled = true
			return fallbackMap, nil
		}, "invalid")
		c.Assert(result, qt.DeepEquals, fallbackMap)
		c.Assert(fallbackCalled, qt.Equals, true)
	})
}

// TestSlicePlusFunctions tests the slice-related plus functions
func TestSlicePlusFunctions(t *testing.T) {
	t.Parallel()

	t.Run("ToSliceP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		input := []any{1, 2, 3}
		result := cast.ToSliceP(func(i any) ([]any, error) {
			t.Error("Fallback should not be called for successful conversion")
			return nil, nil
		}, input)
		c.Assert(result, qt.DeepEquals, input)

		// Test failed conversion with fallback
		fallbackSlice := []any{"fallback"}
		fallbackCalled := false
		result = cast.ToSliceP(func(i any) ([]any, error) {
			fallbackCalled = true
			return fallbackSlice, nil
		}, "invalid")
		c.Assert(result, qt.DeepEquals, fallbackSlice)
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToStringSliceP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		input := []string{"a", "b", "c"}
		result := cast.ToStringSliceP(func(i any) ([]string, error) {
			t.Error("Fallback should not be called for successful conversion")
			return nil, nil
		}, input)
		c.Assert(result, qt.DeepEquals, input)

		// Test failed conversion with fallback - use a type that cannot be converted
		fallbackSlice := []string{"fallback"}
		fallbackCalled := false
		result = cast.ToStringSliceP(func(i any) ([]string, error) {
			fallbackCalled = true
			return fallbackSlice, nil
		}, make(chan int)) // channels cannot be converted to []string
		c.Assert(result, qt.DeepEquals, fallbackSlice)
		c.Assert(fallbackCalled, qt.Equals, true)
	})

	t.Run("ToIntSliceP", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Test successful conversion
		input := []int{1, 2, 3}
		result := cast.ToIntSliceP(func(i any) ([]int, error) {
			t.Error("Fallback should not be called for successful conversion")
			return nil, nil
		}, input)
		c.Assert(result, qt.DeepEquals, input)

		// Test failed conversion with fallback
		fallbackSlice := []int{999}
		fallbackCalled := false
		result = cast.ToIntSliceP(func(i any) ([]int, error) {
			fallbackCalled = true
			return fallbackSlice, nil
		}, "invalid")
		c.Assert(result, qt.DeepEquals, fallbackSlice)
		c.Assert(fallbackCalled, qt.Equals, true)
	})
}

// TestRealWorldScenarios tests practical use cases for plus functions
func TestRealWorldScenarios(t *testing.T) {
	t.Parallel()

	t.Run("CustomDateFormat", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Custom date format that cast doesn't support by default
		customDateStr := "01/02/2006 15:04:05"
		inputDate := "12/25/2023 14:30:00"

		result := cast.ToTimeP(func(i any) (time.Time, error) {
			if str, ok := i.(string); ok {
				return time.Parse(customDateStr, str)
			}
			return time.Time{}, fmt.Errorf("cannot parse %v as custom date", i)
		}, inputDate)

		expected := time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
		c.Assert(result.UTC(), qt.Equals, expected)
	})

	t.Run("CustomBooleanValues", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Custom boolean values that cast doesn't support by default
		customBoolFallback := func(i any) (bool, error) {
			if str, ok := i.(string); ok {
				switch str {
				case "yes", "Y", "on", "enabled":
					return true, nil
				case "no", "N", "off", "disabled":
					return false, nil
				}
			}
			return false, fmt.Errorf("cannot parse %v as custom boolean", i)
		}

		testCases := []struct {
			input    string
			expected bool
		}{
			{"yes", true},
			{"Y", true},
			{"on", true},
			{"enabled", true},
			{"no", false},
			{"N", false},
			{"off", false},
			{"disabled", false},
		}

		for _, tc := range testCases {
			result := cast.ToBoolP(customBoolFallback, tc.input)
			c.Assert(result, qt.Equals, tc.expected, qt.Commentf("input: %s", tc.input))
		}
	})

	t.Run("LegacyDataMigration", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Simulate migrating legacy data with custom string format
		legacyStringFallback := func(i any) (string, error) {
			// Handle legacy format: "LEGACY:actual_value"
			if str, ok := i.(string); ok && len(str) > 7 && str[:7] == "LEGACY:" {
				return str[7:], nil
			}
			return fmt.Sprintf("migrated_%v", i), nil
		}

		// Use types that cast cannot convert to string to trigger fallback
		testCases := []struct {
			input    any
			expected string
		}{
			{make(chan int), "migrated_0x"}, // channels cannot be converted to string
			{func() {}, "migrated_0x"},      // functions cannot be converted to string
		}

		for _, tc := range testCases {
			result := cast.ToStringP(legacyStringFallback, tc.input)
			// For channels and functions, the result will contain memory addresses
			// so we just check that it starts with "migrated_"
			c.Assert(len(result) > 9, qt.Equals, true, qt.Commentf("input: %v", tc.input))
			c.Assert(result[:9], qt.Equals, "migrated_", qt.Commentf("input: %v", tc.input))
		}

		// Test the legacy format handling with a type that cast can't handle
		type customType struct{ value string }
		legacyInput := customType{value: "LEGACY:test_value"}
		result := cast.ToStringP(legacyStringFallback, legacyInput)
		c.Assert(len(result) > 9, qt.Equals, true)
		c.Assert(result[:9], qt.Equals, "migrated_")
	})

	t.Run("ConfigurationDefaults", func(t *testing.T) {
		t.Parallel()
		c := qt.New(t)

		// Simulate configuration parsing with defaults
		configDefaults := map[string]any{
			"timeout":    30,
			"retries":    3,
			"debug":      false,
			"server":     "localhost",
			"port":       8080,
		}

		getConfigInt := func(key string) func(any) (int, error) {
			return func(i any) (int, error) {
				if defaultVal, ok := configDefaults[key].(int); ok {
					return defaultVal, nil
				}
				return 0, fmt.Errorf("no default for key %s", key)
			}
		}

		// Test with invalid config values falling back to defaults
		timeout := cast.ToIntP(getConfigInt("timeout"), "invalid")
		c.Assert(timeout, qt.Equals, 30)

		// Use a type that cannot be converted to int to trigger fallback
		retries := cast.ToIntP(getConfigInt("retries"), make(chan int))
		c.Assert(retries, qt.Equals, 3)

		// Test with valid config values (no fallback)
		port := cast.ToIntP(getConfigInt("port"), "9000")
		c.Assert(port, qt.Equals, 9000) // Should use the provided value, not default
	})
}
