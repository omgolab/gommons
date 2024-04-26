// The above code is a Go package that provides functions for retrieving environment variables with
// type conversion and fallback values.
package genv

import (
	"os"
	"strconv"
	"time"
)

// The above type defines an interface that can be implemented by any integer type in Go.
type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// The above type defines an interface that can be implemented by any unsigned integer type in Go.
type UintType interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// The FloatType interface represents either a float32 or a float64 type.
type FloatType interface {
	~float32 | ~float64
}

// The type EnvValueType represents a value that can be an integer, unsigned integer, floating-point
// number, boolean, string, or time.
type EnvValueType interface {
	IntType | UintType | FloatType | ~bool | ~string | time.Time
}

type EnvConverter[T any] func(string, bool) T

// The `defaultConverter` function is a generic function that converts environment variable values to a
// specified type, with a fallback value if the variable is not found.
func defaultConverter[T EnvValueType](fallback T) EnvConverter[T] {
	return func(val string, found bool) T {
		if found {
			return any(val).(T)
		}

		switch any(fallback).(type) {
		case int:
			return intConvert(val, fallback, 10).(T)
		case int8:
			return intConvert(val, fallback, 8).(T)
		case int16:
			return intConvert(val, fallback, 16).(T)
		case int32:
			return intConvert(val, fallback, 32).(T)
		case int64:
			return intConvert(val, fallback, 64).(T)
		case uint:
			return uintConvert(val, fallback, 10).(T)
		case uint8:
			return uintConvert(val, fallback, 8).(T)
		case uint16:
			return uintConvert(val, fallback, 16).(T)
		case uint32:
			return uintConvert(val, fallback, 32).(T)
		case uint64:
			return uintConvert(val, fallback, 64).(T)
		case float32:
			return floatConvert(val, fallback, 32).(T)
		case float64:
			return floatConvert(val, fallback, 64).(T)
		case bool:
			return boolConvert(val, fallback).(T)
		case string:
			return any(val).(T)
		case time.Time:
			return timeConvert(val, fallback).(T)
		default:
			return fallback
		}
	}
}

// The function `EnvWithConverter` retrieves an environment variable value and converts it to a
// specified type using a provided converter function.
func EnvWithConverter[T EnvValueType](key string, converter EnvConverter[T]) T {
	val, found := os.LookupEnv(key)
	if converter != nil {
		return converter(val, found)
	}
	return any(val).(T)
}

// The `Env` function retrieves the value of an environment variable and provides a fallback value if
// the variable is not set.
func Env[T EnvValueType](key string, fallback T) T {
	return EnvWithConverter(key, defaultConverter[T](fallback))
}

// The function intConvert converts a string to an integer with a specified bit size, returning a
// fallback value if the conversion fails.
func intConvert(val string, fallback any, bitSize int) any {
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err != nil {
		return fallback
	}
	return intVal
}

// The function `uintConvert` converts a string to an unsigned integer of a specified bit size,
// returning a fallback value if the conversion fails.
func uintConvert(val string, fallback any, bitSize int) any {
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err != nil {
		return fallback
	}
	return uintVal
}

// The function "floatConvert" converts a string to a float value with a specified bit size, and
// returns a fallback value if the conversion fails.
func floatConvert(val string, fallback any, bitSize int) any {
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err != nil {
		return fallback
	}
	return floatVal
}

// The function `boolConvert` converts a string value to a boolean and returns a fallback value if the
// conversion fails.
func boolConvert(val string, fallback any) any {
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return boolVal
}

// The timeConvert function parses a string into a time value using the "2006-01-02" format and returns
// the parsed time value or a fallback value if parsing fails.
func timeConvert(val string, fallback any) any {
	timeVal, err := time.Parse("2006-01-02", val)
	if err != nil {
		return fallback
	}
	return timeVal
}
