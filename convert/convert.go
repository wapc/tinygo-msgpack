package convert

import (
	"time"
)

const Package = "convert"

func Format[T any](fn func(*T) string) func(value *T) *string {
	return func(value *T) *string {
		if value == nil {
			return nil
		}

		str := fn(value)
		return &str
	}
}

func Parse[I, T any](fn func(I) (T, error)) func(I, error) (T, error) {
	return func(value I, decodeErr error) (val T, err error) {
		if decodeErr != nil {
			return val, decodeErr
		}
		return fn(value)
	}
}

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

func StringToTime(value string, err error) (time.Time, error) {
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, value)
}

func StringToTimePtr(value *string, err error) (*time.Time, error) {
	if value == nil || err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339Nano, *value)
	return &t, err
}

func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func TimeToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	val := t.Format(time.RFC3339Nano)
	return &val
}

func String[T ~string](value string, err error) (T, error) {
	return T(value), err
}

func NillableString[T ~string](value *string, err error) (*T, error) {
	ret := T(*value)
	return &ret, err
}

func Bool[T ~bool](value bool, err error) (T, error) {
	return T(value), err
}

func NillableBool[T ~bool](value *bool, err error) (*T, error) {
	ret := T(*value)
	return &ret, err
}

type numberic interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func Numeric[T, I numberic](value I, err error) (T, error) {
	return T(value), err
}

func NillableNumeric[T, I numberic](value *I, err error) (*T, error) {
	ret := T(*value)
	return &ret, err
}

func ByteArray[T, I ~[]byte](value I, err error) (T, error) {
	return T(value), err
}
