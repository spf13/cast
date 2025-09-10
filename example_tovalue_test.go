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

// CustomBool demonstrates custom boolean conversion using ValueSetter
type CustomBool struct {
	Value bool
}

func (cb *CustomBool) SetValue(i any) error {
	if str, ok := i.(string); ok {
		switch str {
		case "yes", "Y", "on", "enabled":
			cb.Value = true
			return nil
		case "no", "N", "off", "disabled":
			cb.Value = false
			return nil
		}
	}

	// Fallback to standard conversion
	if v, err := cast.ToBoolE(i); err == nil {
		cb.Value = v
		return nil
	}

	return fmt.Errorf("cannot parse %v as custom boolean", i)
}

// ExampleToValue demonstrates using ToValue with custom boolean logic
func ExampleToValue() {
	var cb CustomBool

	// Standard conversion works normally
	cast.ToValue(&cb, true)
	fmt.Println(cb.Value) // true

	cast.ToValue(&cb, "true")
	fmt.Println(cb.Value) // true

	// Custom conversion kicks in for unsupported values
	cast.ToValue(&cb, "yes")
	fmt.Println(cb.Value) // true

	cast.ToValue(&cb, "enabled")
	fmt.Println(cb.Value) // true

	cast.ToValue(&cb, "no")
	fmt.Println(cb.Value) // false

	// Output:
	// true
	// true
	// true
	// true
	// false
}

// CustomTime demonstrates custom time conversion using ValueSetter
type CustomTime struct {
	Value time.Time
}

func (ct *CustomTime) SetValue(i any) error {
	if str, ok := i.(string); ok {
		// Try custom format: MM/DD/YYYY HH:MM:SS
		if t, err := time.Parse("01/02/2006 15:04:05", str); err == nil {
			ct.Value = t
			return nil
		}
	}

	// Fallback to standard conversion
	if v, err := cast.ToTimeE(i); err == nil {
		ct.Value = v
		return nil
	}

	return fmt.Errorf("cannot parse %v as custom time", i)
}

// Example_customTime demonstrates using ToValue with custom date formats
func Example_customTime() {
	var ct CustomTime

	// Standard conversion works normally
	standardTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	cast.ToValue(&ct, standardTime)
	fmt.Println(ct.Value.Year()) // 2023

	// Custom conversion for unsupported format
	cast.ToValue(&ct, "12/25/2023 14:30:00")
	fmt.Println(ct.Value.Month()) // December
	fmt.Println(ct.Value.Day())   // 25

	// Output:
	// 2023
	// December
	// 25
}

// ConfigInt demonstrates configuration with defaults using ValueSetter
type ConfigInt struct {
	Value   int
	Default int
}

func (ci *ConfigInt) SetValue(i any) error {
	// Try standard conversion first
	if v, err := cast.ToIntE(i); err == nil {
		ci.Value = v
		return nil
	}

	// Use default if conversion fails
	ci.Value = ci.Default
	return nil
}

// Example_config demonstrates using ToValue for configuration with defaults
func Example_config() {
	// Configuration with defaults
	timeout := &ConfigInt{Default: 30}
	retries := &ConfigInt{Default: 3}
	port := &ConfigInt{Default: 8080}

	// Valid values use standard conversion
	cast.ToValue(timeout, "60")
	fmt.Println(timeout.Value) // 60

	// Invalid values fall back to defaults
	cast.ToValue(retries, "invalid")
	fmt.Println(retries.Value) // 3

	cast.ToValue(port, make(chan int))
	fmt.Println(port.Value) // 8080

	// Output:
	// 60
	// 3
	// 8080
}

// Example_basicTypes demonstrates using ToValue with basic types
func Example_basicTypes() {
	// Basic type conversions
	var str string
	var num int
	var b bool

	cast.ToValue(&str, 123)
	fmt.Println(str) // "123"

	cast.ToValue(&num, "42")
	fmt.Println(num) // 42

	cast.ToValue(&b, "true")
	fmt.Println(b) // true

	// Output:
	// 123
	// 42
	// true
}
