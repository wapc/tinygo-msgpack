// Package msgpack provides a high-performance MessagePack implementation
// optimized for TinyGo and WebAssembly environments.
package msgpack

import (
	"reflect"
)

// WriteSlice writes a slice of values to MessagePack format.
// It uses the provided function to encode each individual value.
func WriteSlice[T any](encoder Writer, values []T, valF func(Writer, T)) error {
	encoder.WriteArraySize(uint32(len(values)))
	for _, item := range values {
		valF(encoder, item)
	}

	return encoder.Err()
}

// ReadSlice reads a slice of values from MessagePack format.
// It uses the provided function to decode each individual value.
func ReadSlice[T any](reader Reader, valF func(reader Reader) (T, error)) ([]T, error) {
	listSize, err := reader.ReadArraySize()
	if err != nil {
		return nil, err
	}
	request := make([]T, 0, listSize)
	for listSize > 0 {
		listSize--
		item, err := valF(reader)
		if err != nil {
			return nil, err
		}
		request = append(request, item)
	}

	return request, nil
}

// WriteMap writes a map to MessagePack format.
// It uses the provided functions to encode keys and values.
func WriteMap[K comparable, V any](writer Writer,
	m map[K]V, keyF func(Writer, K),
	valF func(Writer, V)) error {
	writer.WriteMapSize(uint32(len(m)))
	for key, val := range m {
		keyF(writer, key)
		valF(writer, val)
	}

	return writer.Err()
}

// ReadMap reads a map from MessagePack format.
// It uses the provided functions to decode keys and values.
func ReadMap[K comparable, V any](reader Reader,
	keyF func(reader Reader) (K, error),
	valF func(reader Reader) (V, error)) (map[K]V, error) {
	mapSize, err := reader.ReadMapSize()
	if err != nil {
		return nil, err
	}
	m := make(map[K]V, mapSize)
	for mapSize > 0 {
		mapSize--
		key, err := keyF(reader)
		if err != nil {
			return nil, err
		}
		value, err := valF(reader)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}

	return m, nil
}

// ReadAny reads any value from MessagePack format.
// It's a convenience function that creates a decoder and calls ReadAny.
func ReadAny(data []byte) (any, error) {
	d := NewDecoder(data)
	return d.ReadAny()
}

// isNil checks if a value is nil, including nil pointers.
// It uses reflection to handle interface types properly.
func isNil(val any) bool {
	return val == nil ||
		(reflect.ValueOf(val).Kind() == reflect.Ptr &&
			reflect.ValueOf(val).IsNil())
}
