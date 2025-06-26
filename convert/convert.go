// Package convert provides utility functions for type conversion
// in MessagePack encoding and decoding operations.
package convert

import (
	"time"
)

// Package is the package name for the convert package.
const Package = "convert"

// Format creates a function that formats a value to a string.
// The returned function handles nil values by returning nil.
func Format[T any](fn func(*T) string) func(value *T) *string {
	return func(value *T) *string {
		if value == nil {
			return nil
		}

		str := fn(value)
		return &str
	}
}

// Parse creates a function that parses a value with error handling.
// The returned function preserves the original error if present.
func Parse[I, T any](fn func(I) (T, error)) func(I, error) (T, error) {
	return func(value I, decodeErr error) (val T, err error) {
		if decodeErr != nil {
			return val, decodeErr
		}
		return fn(value)
	}
}

// NillableParse creates a function that parses a nullable value with error handling.
// The returned function handles nil values and preserves original errors.
func NillableParse[I, T any](fn func(I) (T, error)) func(*I, error) (*T, error) {
	return func(value *I, decodeErr error) (val *T, err error) {
		if decodeErr != nil || value == nil {
			return val, decodeErr
		}
		var v T
		if v, err = fn(*value); err != nil {
			return val, err
		}
		return &v, nil
	}
}

// StringToTime converts a string to a time.Time using RFC3339Nano format.
// Returns the original error if present.
func StringToTime(value string, err error) (time.Time, error) {
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, value)
}

// StringToTimePtr converts a string pointer to a time.Time pointer using RFC3339Nano format.
// Returns nil if the input is nil or if there's an error.
func StringToTimePtr(value *string, err error) (*time.Time, error) {
	if value == nil || err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339Nano, *value)
	return &t, err
}

// TimeToString converts a time.Time to a string using RFC3339Nano format.
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

// TimeToStringPtr converts a time.Time pointer to a string pointer using RFC3339Nano format.
// Returns nil if the input is nil.
func TimeToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	val := t.Format(time.RFC3339Nano)
	return &val
}

// String converts a string to a typed string with error handling.
func String[T ~string](value string, err error) (T, error) {
	return T(value), err
}

// NillableString converts a string pointer to a typed string pointer with error handling.
func NillableString[T ~string](value *string, err error) (*T, error) {
	ret := T(*value)
	return &ret, err
}

// Bool converts a bool to a typed bool with error handling.
func Bool[T ~bool](value bool, err error) (T, error) {
	return T(value), err
}

// NillableBool converts a bool pointer to a typed bool pointer with error handling.
func NillableBool[T ~bool](value *bool, err error) (*T, error) {
	ret := T(*value)
	return &ret, err
}

// numberic is a constraint for numeric types that can be converted between each other.
type numberic interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Numeric converts a numeric value to another numeric type with error handling.
func Numeric[T, I numberic](value I, err error) (T, error) {
	return T(value), err
}

// NillableNumeric converts a numeric pointer to another numeric pointer with error handling.
func NillableNumeric[T, I numberic](value *I, err error) (*T, error) {
	ret := T(*value)
	return &ret, err
}

// ByteArray converts a byte slice to a typed byte slice with error handling.
func ByteArray[T, I ~[]byte](value I, err error) (T, error) {
	return T(value), err
}
