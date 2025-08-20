# Plus Functions

This document describes the "plus functions" feature added to the cast library. Plus functions provide a way to extend cast's type conversion capabilities with custom fallback logic.

## Overview

Plus functions are variants of cast's standard conversion functions that accept a fallback function. When the standard cast conversion fails, the fallback function is called to provide an alternative conversion strategy.

## Function Naming Convention

Plus functions follow the naming pattern: `To{Type}P` and `To{Type}PE`

- `To{Type}P`: Returns the converted value, using fallback on failure
- `To{Type}PE`: Returns the converted value and error, using fallback on failure

## Generic Plus Functions

### `ToP[T Basic](fn func(any) (T, error), i any) T`
Generic plus function for any Basic type (string, bool, numbers, time.Time, time.Duration).

### `ToPE[T Basic](fn func(any) (T, error), i any) (T, error)`
Generic plus function with error return.

## Specific Plus Functions

### Basic Types
- `ToBoolP(fn func(any) (bool, error), i any) bool`
- `ToStringP(fn func(any) (string, error), i any) string`

### Numeric Types
- `ToIntP(fn func(any) (int, error), i any) int`
- `ToInt8P`, `ToInt16P`, `ToInt32P`, `ToInt64P`
- `ToUintP`, `ToUint8P`, `ToUint16P`, `ToUint32P`, `ToUint64P`
- `ToFloat32P`, `ToFloat64P`

### Time Types
- `ToTimeP(fn func(any) (time.Time, error), i any) time.Time`
- `ToTimeInDefaultLocationP(fn func(any, *time.Location) (time.Time, error), i any, location *time.Location) time.Time`
- `ToDurationP(fn func(any) (time.Duration, error), i any) time.Duration`

### Map Types
- `ToStringMapStringP(fn func(any) (map[string]string, error), i any) map[string]string`
- `ToStringMapIntP(fn func(any) (map[string]int, error), i any) map[string]int`
- `ToStringMapBoolP(fn func(any) (map[string]bool, error), i any) map[string]bool`
- And more...

### Slice Types
- `ToSliceP(fn func(any) ([]any, error), i any) []any`
- `ToStringSliceP(fn func(any) ([]string, error), i any) []string`
- `ToIntSliceP(fn func(any) ([]int, error), i any) []int`
- And more...

## Usage Examples

### Custom Boolean Values
```go
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

result := cast.ToBoolP(customBoolFallback, "yes") // true
```

### Custom Date Formats
```go
customDateFallback := func(i any) (time.Time, error) {
    if str, ok := i.(string); ok {
        return time.Parse("01/02/2006 15:04:05", str)
    }
    return time.Time{}, fmt.Errorf("cannot parse %v as custom date", i)
}

result := cast.ToTimeP(customDateFallback, "12/25/2023 14:30:00")
```

### Configuration with Defaults
```go
getDefault := func(defaultVal int) func(any) (int, error) {
    return func(i any) (int, error) {
        return defaultVal, nil
    }
}

timeout := cast.ToP[int](getDefault(30), "invalid") // 30
```

### Generic Usage
```go
// Works with any Basic type
port := cast.ToP[int](getDefault(8080), invalidConfig)
debug := cast.ToP[bool](getDefault(false), invalidConfig)
host := cast.ToP[string](getDefault("localhost"), invalidConfig)
```

## Benefits

1. **Extends Coverage**: Handle edge cases that cast doesn't support natively
2. **Backward Compatible**: Doesn't change existing cast behavior
3. **Type Safe**: Leverages Go's generics for type safety
4. **Flexible**: Supports custom conversion logic for any scenario
5. **Consistent**: Follows cast's existing API patterns

## When to Use

- Parsing legacy data formats
- Handling domain-specific string representations
- Providing configuration defaults
- Supporting custom boolean/date/number formats
- Migrating from other type conversion libraries
- Adding application-specific conversion rules

## Performance

Plus functions have minimal overhead:
- If standard conversion succeeds, fallback is never called
- Only when conversion fails does the fallback function execute
- No performance impact on existing cast usage
