package msgpack

import (
	"time"
)

type Decodable interface {
	Decode(Reader) error
}

type Encodable interface {
	Encode(Writer) error
}

// Codable is the interface that applies to data structures that can
// encode to and decode from the MessagePack format.
type Codable interface {
	Decodable
	Encodable
}

// ToBytes creates a `[]byte` from `codec`.
func ToBytes(value Encodable) ([]byte, error) {
	var sizer Sizer
	if err := value.Encode(&sizer); err != nil {
		return nil, err
	}
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	if err := value.Encode(&encoder); err != nil {
		return nil, err
	}
	return buffer, nil
}

//go:inline
func Marshal(value any) ([]byte, error) {
	return AnyToBytes(value)
}

func Unmarshal(data []byte, value Decodable) error {
	decoder := NewDecoder(data)
	return value.Decode(&decoder)
}

func BytesToAny(data []byte) (any, error) {
	d := NewDecoder(data)
	return d.ReadAny()
}

// AnyToBytes creates a `[]byte` from `value`.
func AnyToBytes(value any) ([]byte, error) {
	var sizer Sizer
	sizer.WriteAny(value)
	if err := sizer.Err(); err != nil {
		return nil, err
	}
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteAny(value)
	if err := encoder.Err(); err != nil {
		return nil, err
	}
	return buffer, nil
}

// I8ToBytes creates a `[]byte` from `value`.
func I8ToBytes(value int8) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteInt8(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteInt8(value)
	return buffer, nil
}

// I16ToBytes creates a `[]byte` from `value`.
func I16ToBytes(value int16) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteInt16(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteInt16(value)
	return buffer, nil
}

// I32ToBytes creates a `[]byte` from `value`.
func I32ToBytes(value int32) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteInt32(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteInt32(value)
	return buffer, nil
}

// I64ToBytes creates a `[]byte` from `value`.
func I64ToBytes(value int64) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteInt64(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteInt64(value)
	return buffer, nil
}

// U8ToBytes creates a `[]byte` from `value`.
func U8ToBytes(value uint8) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteUint8(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteUint8(value)
	return buffer, nil
}

// U16ToBytes creates a `[]byte` from `value`.
func U16ToBytes(value uint16) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteUint16(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteUint16(value)
	return buffer, nil
}

// U32ToBytes creates a `[]byte` from `value`.
func U32ToBytes(value uint32) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteUint32(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteUint32(value)
	return buffer, nil
}

// U64ToBytes creates a `[]byte` from `value`.
func U64ToBytes(value uint64) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteUint64(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteUint64(value)
	return buffer, nil
}

// F32ToBytes creates a `[]byte` from `value`.
func F32ToBytes(value float32) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteFloat32(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteFloat32(value)
	return buffer, nil
}

// F64ToBytes creates a `[]byte` from `value`.
func F64ToBytes(value float64) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteFloat64(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteFloat64(value)
	return buffer, nil
}

// BoolToBytes creates a `[]byte` from `value`.
func BoolToBytes(value bool) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteBool(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteBool(value)
	return buffer, nil
}

// StringToBytes creates a `[]byte` from `value`.
func StringToBytes(value string) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteString(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteString(value)
	return buffer, nil
}

// BytesToBytes creates a `[]byte` from `value`.
func BytesToBytes(value []byte) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteByteArray(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteByteArray(value)
	return buffer, nil
}

// TimeToBytes creates a `[]byte` from `value`.
func TimeToBytes(value time.Time) ([]byte, error) {
	sizer := NewSizer()
	sizer.WriteTime(value)
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	encoder.WriteTime(value)
	return buffer, nil
}

func SliceToBytes[T any](values []T, valF func(Writer, T)) ([]byte, error) {
	sizer := NewSizer()
	if err := WriteSlice(&sizer, values, valF); err != nil {
		return nil, err
	}
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	if err := WriteSlice(&encoder, values, valF); err != nil {
		return nil, err
	}
	return buffer, nil
}

func MapToBytes[K comparable, V any](values map[K]V, keyF func(Writer, K), valF func(Writer, V)) ([]byte, error) {
	sizer := NewSizer()
	if err := WriteMap(&sizer, values, keyF, valF); err != nil {
		return nil, err
	}
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	if err := WriteMap(&encoder, values, keyF, valF); err != nil {
		return nil, err
	}
	return buffer, nil
}
