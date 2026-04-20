package cast

import (
	"fmt"
	"time"
)

// CustomBool demonstrates a custom type that implements ValueSetter
type CustomBool struct {
	value string
}

func (c CustomBool) SetValue(target any) error {
	switch t := target.(type) {
	case *bool:
		*t = c.value == "yes" || c.value == "true" || c.value == "1"
		return nil
	case *string:
		if c.value == "yes" || c.value == "true" || c.value == "1" {
			*t = "true"
		} else {
			*t = "false"
		}
		return nil
	default:
		return fmt.Errorf("unsupported target type: %T", target)
	}
}

// CustomTime demonstrates a custom time type
type CustomTime struct {
	timestamp int64
}

func (c CustomTime) SetValue(target any) error {
	switch t := target.(type) {
	case *time.Time:
		*t = time.Unix(c.timestamp, 0)
		return nil
	case *string:
		*t = time.Unix(c.timestamp, 0).Format(time.RFC3339)
		return nil
	default:
		return fmt.Errorf("unsupported target type: %T", target)
	}
}

// CustomInt demonstrates a custom integer type
type CustomInt struct {
	hexValue string
}

func (c CustomInt) SetValue(target any) error {
	switch t := target.(type) {
	case *int:
		if c.hexValue == "0xFF" {
			*t = 255
		} else {
			*t = 0
		}
		return nil
	case *string:
		*t = c.hexValue
		return nil
	default:
		return fmt.Errorf("unsupported target type: %T", target)
	}
}

func ExampleValueSetter_bool() {
	custom := CustomBool{value: "yes"}

	result, err := ToBoolE(custom)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Bool result: %v\n", result)
	// Output: Bool result: true
}

func ExampleValueSetter_string() {
	custom := CustomBool{value: "yes"}

	result, err := ToStringE(custom)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("String result: %s\n", result)
	// Output: String result: true
}

func ExampleValueSetter_time() {
	custom := CustomTime{timestamp: 1609459200} // 2021-01-01 00:00:00 UTC

	result, err := ToTimeE(custom)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Time result: %s\n", result.Format("2006-01-02"))
	// Output: Time result: 2021-01-01
}

func ExampleValueSetter_int() {
	custom := CustomInt{hexValue: "0xFF"}

	result, err := ToIntE(custom)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Int result: %d\n", result)
	// Output: Int result: 255
}
