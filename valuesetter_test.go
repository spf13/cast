package cast

import (
	"testing"
	"time"
)

// TestValueSetterIntegration tests that ValueSetter interface works correctly
// with all conversion functions in their default branches
func TestValueSetterIntegration(t *testing.T) {
	// Test with CustomBool
	customBool := CustomBool{value: "yes"}

	// Test ToBoolE
	result, err := ToBoolE(customBool)
	if err != nil {
		t.Errorf("ToBoolE failed: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got %v", result)
	}

	// Test ToStringE
	strResult, err := ToStringE(customBool)
	if err != nil {
		t.Errorf("ToStringE failed: %v", err)
	}
	if strResult != "true" {
		t.Errorf("Expected 'true', got %v", strResult)
	}

	// Test with CustomTime
	customTime := CustomTime{timestamp: 1609459200} // 2021-01-01 00:00:00 UTC

	timeResult, err := ToTimeE(customTime)
	if err != nil {
		t.Errorf("ToTimeE failed: %v", err)
	}
	expected := time.Unix(1609459200, 0)
	if !timeResult.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, timeResult)
	}

	// Test with CustomInt
	customInt := CustomInt{hexValue: "0xFF"}

	intResult, err := ToIntE(customInt)
	if err != nil {
		t.Errorf("ToIntE failed: %v", err)
	}
	if intResult != 255 {
		t.Errorf("Expected 255, got %v", intResult)
	}
}

// TestStandardConversionPriority ensures standard conversions work normally
// and ValueSetter is only used as fallback
func TestStandardConversionPriority(t *testing.T) {
	// Test that standard string conversion works normally
	result, err := ToStringE("hello")
	if err != nil {
		t.Errorf("ToStringE failed: %v", err)
	}
	if result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}

	// Test that standard bool conversion works normally
	boolResult, err := ToBoolE(true)
	if err != nil {
		t.Errorf("ToBoolE failed: %v", err)
	}
	if !boolResult {
		t.Errorf("Expected true, got %v", boolResult)
	}

	// Test that standard int conversion works normally
	intResult, err := ToIntE(42)
	if err != nil {
		t.Errorf("ToIntE failed: %v", err)
	}
	if intResult != 42 {
		t.Errorf("Expected 42, got %v", intResult)
	}
}
