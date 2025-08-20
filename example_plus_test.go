// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast_test

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

// ExampleToBoolP demonstrates using ToBoolP with custom boolean logic
func ExampleToBoolP() {
	// Custom boolean conversion for "yes"/"no" strings
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

	// Standard conversion works normally
	fmt.Println(cast.ToBoolP(customBoolFallback, true))   // true
	fmt.Println(cast.ToBoolP(customBoolFallback, "true")) // true

	// Custom conversion kicks in for unsupported values
	fmt.Println(cast.ToBoolP(customBoolFallback, "yes"))     // true
	fmt.Println(cast.ToBoolP(customBoolFallback, "enabled")) // true
	fmt.Println(cast.ToBoolP(customBoolFallback, "no"))      // false

	// Output:
	// true
	// true
	// true
	// true
	// false
}

// ExampleToTimeP demonstrates using ToTimeP with custom date formats
func ExampleToTimeP() {
	// Custom date format fallback
	customDateFallback := func(i any) (time.Time, error) {
		if str, ok := i.(string); ok {
			// Try custom format: MM/DD/YYYY HH:MM:SS
			if t, err := time.Parse("01/02/2006 15:04:05", str); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse %v as custom date", i)
	}

	// Standard conversion works normally
	standardTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	result1 := cast.ToTimeP(customDateFallback, standardTime)
	fmt.Println(result1.Year()) // 2023

	// Custom conversion for unsupported format
	result2 := cast.ToTimeP(customDateFallback, "12/25/2023 14:30:00")
	fmt.Println(result2.Month()) // December
	fmt.Println(result2.Day())   // 25

	// Output:
	// 2023
	// December
	// 25
}

// ExampleToP demonstrates using the generic ToP function
func ExampleToP() {
	// Configuration with defaults
	configDefaults := map[string]int{
		"timeout": 30,
		"retries": 3,
		"port":    8080,
	}

	getDefault := func(key string) func(any) (int, error) {
		return func(i any) (int, error) {
			if val, ok := configDefaults[key]; ok {
				return val, nil
			}
			return 0, fmt.Errorf("no default for %s", key)
		}
	}

	// Valid values use standard conversion
	timeout := cast.ToP[int](getDefault("timeout"), "60")
	fmt.Println(timeout) // 60

	// Invalid values fall back to defaults
	retries := cast.ToP[int](getDefault("retries"), "invalid")
	fmt.Println(retries) // 3

	port := cast.ToP[int](getDefault("port"), make(chan int))
	fmt.Println(port) // 8080

	// Output:
	// 60
	// 3
	// 8080
}

// ExampleToStringMapStringP demonstrates using map plus functions
func ExampleToStringMapStringP() {
	// Fallback that provides default configuration
	defaultConfig := func(i any) (map[string]string, error) {
		return map[string]string{
			"host": "localhost",
			"port": "8080",
			"env":  "development",
		}, nil
	}

	// Valid map works normally
	validMap := map[string]string{"host": "example.com", "port": "9000"}
	result1 := cast.ToStringMapStringP(defaultConfig, validMap)
	fmt.Println(result1["host"]) // example.com

	// Invalid input falls back to defaults
	result2 := cast.ToStringMapStringP(defaultConfig, "invalid")
	fmt.Println(result2["host"]) // localhost
	fmt.Println(result2["env"])  // development

	// Output:
	// example.com
	// localhost
	// development
}

// ExampleToStringSliceP demonstrates using slice plus functions
func ExampleToStringSliceP() {
	// Fallback that provides default values
	defaultSliceFallback := func(i any) ([]string, error) {
		return []string{"default", "values"}, nil
	}

	// Valid slice works normally
	validSlice := []string{"a", "b", "c"}
	result1 := cast.ToStringSliceP(defaultSliceFallback, validSlice)
	fmt.Println(len(result1)) // 3
	fmt.Println(result1[0])   // a

	// Invalid input uses fallback (use a type that can't be converted)
	result2 := cast.ToStringSliceP(defaultSliceFallback, make(chan int))
	fmt.Println(len(result2))  // 2
	fmt.Println(result2[0])    // default

	// Output:
	// 3
	// a
	// 2
	// default
}
