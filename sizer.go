// Package msgpack provides a high-performance MessagePack implementation
// optimized for TinyGo and WebAssembly environments.
package msgpack

import (
	"math"
	"time"
)

// Sizer calculates the size of MessagePack encoded data without actually encoding.
// It implements the Writer interface to track the size of data that would be written.
type Sizer struct {
	length uint32
}

// Ensure `*Sizer` implements `Writer`.
var _ Writer = (*Sizer)(nil)

// NewSizer creates a new Sizer with zero length.
func NewSizer() Sizer {
	return Sizer{}
}

// Len returns the calculated length of the MessagePack data.
func (s *Sizer) Len() uint32 {
	return s.length
}

// WriteNil adds 1 byte for a MessagePack nil value.
func (s *Sizer) WriteNil() {
	s.length++
}

// WriteString calculates the size needed for a MessagePack string value.
func (s *Sizer) WriteString(value string) {
	buf := UnsafeBytes(value)
	length := uint32(len(buf))
	s.writeStringLength(length)
	s.length += length
}

// WriteNillableString calculates the size for a string or nil value.
func (s *Sizer) WriteNillableString(value *string) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteString(*value)
	}
}

// writeStringLength calculates the size of the string length header.
func (s *Sizer) writeStringLength(length uint32) {
	if length < 32 {
		s.length++
	} else if length <= math.MaxUint8 {
		s.length += 2
	} else if length <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

// WriteTime calculates the size needed for a MessagePack time value.
func (s *Sizer) WriteTime(value time.Time) {
	l := s.encodeTime(value)
	s.encodeExtLen(l)
	s.length += 1 + uint32(l)
}

// encodeTime calculates the size of the encoded time data.
func (s *Sizer) encodeTime(tm time.Time) int {
	secs := uint64(tm.Unix())
	if secs>>34 == 0 {
		data := uint64(tm.Nanosecond())<<34 | secs

		if data&0xffffffff00000000 == 0 {
			return 4
		}

		return 8
	}

	return 12
}

// encodeExtLen calculates the size of the extension length header.
func (s *Sizer) encodeExtLen(l int) {
	switch l {
	case 1, 2, 4, 8, 16:
		s.length++
		return
	}
	if l <= math.MaxUint8 {
		s.length += 2
	} else if l <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

// WriteNillableTime calculates the size for a time or nil value.
func (s *Sizer) WriteNillableTime(value *time.Time) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteTime(*value)
	}
}

// WriteBool adds 1 byte for a MessagePack boolean value.
func (s *Sizer) WriteBool(value bool) {
	s.length++
}

// WriteNillableBool calculates the size for a boolean or nil value.
func (s *Sizer) WriteNillableBool(value *bool) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteBool(*value)
	}
}

// WriteArraySize calculates the size of the array size header.
func (s *Sizer) WriteArraySize(length uint32) {
	if length < 16 {
		s.length++
	} else if length <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

// writeBinLength calculates the size of the binary length header.
func (s *Sizer) writeBinLength(length uint32) {
	if length < math.MaxUint8 {
		s.length += 1
	} else if length <= math.MaxUint16 {
		s.length += 2
	} else {
		s.length += 4
	}
}

// WriteByteArray calculates the size needed for a MessagePack binary value.
func (s *Sizer) WriteByteArray(value []byte) {
	length := uint32(len(value))
	if length == 0 {
		s.length += 2
		return
	}
	s.writeBinLength(length)
	s.length += length + 1
}

// WriteNillableByteArray calculates the size for a byte array or nil value.
func (s *Sizer) WriteNillableByteArray(value []byte) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteByteArray(value)
	}
}

// WriteMapSize calculates the size of the map size header.
func (s *Sizer) WriteMapSize(length uint32) {
	if length < 16 {
		s.length++
	} else if length <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

// WriteInt8 calculates the size needed for an int8 value.
func (s *Sizer) WriteInt8(value int8) {
	s.WriteInt64(int64(value))
}

// WriteNillableInt8 calculates the size for an int8 or nil value.
func (s *Sizer) WriteNillableInt8(value *int8) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteInt8(*value)
	}
}

// WriteInt16 calculates the size needed for an int16 value.
func (s *Sizer) WriteInt16(value int16) {
	s.WriteInt64(int64(value))
}

// WriteNillableInt16 calculates the size for an int16 or nil value.
func (s *Sizer) WriteNillableInt16(value *int16) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteInt16(*value)
	}
}

// WriteInt32 calculates the size needed for an int32 value.
func (s *Sizer) WriteInt32(value int32) {
	s.WriteInt64(int64(value))
}

// WriteNillableInt32 calculates the size for an int32 or nil value.
func (s *Sizer) WriteNillableInt32(value *int32) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteInt32(*value)
	}
}

// WriteInt64 calculates the size needed for an int64 value.
func (s *Sizer) WriteInt64(value int64) {
	if value >= -(1<<5) && value < 1<<7 {
		s.length++
	} else if value < 1<<7 && value >= -(1<<7) {
		s.length += 2
	} else if value < 1<<15 && value >= -(1<<15) {
		s.length += 3
	} else if value < 1<<31 && value >= -(1<<31) {
		s.length += 5
	} else {
		s.length += 9
	}
}

// WriteNillableInt64 calculates the size for an int64 or nil value.
func (s *Sizer) WriteNillableInt64(value *int64) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteInt64(*value)
	}
}

// WriteUint8 calculates the size needed for a uint8 value.
func (s *Sizer) WriteUint8(value uint8) {
	s.WriteUint64(uint64(value))
}

// WriteNillableUint8 calculates the size for a uint8 or nil value.
func (s *Sizer) WriteNillableUint8(value *uint8) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteUint8(*value)
	}
}

// WriteUint16 calculates the size needed for a uint16 value.
func (s *Sizer) WriteUint16(value uint16) {
	s.WriteUint64(uint64(value))
}

// WriteNillableUint16 calculates the size for a uint16 or nil value.
func (s *Sizer) WriteNillableUint16(value *uint16) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteUint16(*value)
	}
}

// WriteUint32 calculates the size needed for a uint32 value.
func (s *Sizer) WriteUint32(value uint32) {
	s.WriteUint64(uint64(value))
}

// WriteNillableUint32 calculates the size for a uint32 or nil value.
func (s *Sizer) WriteNillableUint32(value *uint32) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteUint32(*value)
	}
}

// WriteUint64 calculates the size needed for a uint64 value.
func (s *Sizer) WriteUint64(value uint64) {
	if value < 1<<7 {
		s.length++
	} else if value < 1<<8 {
		s.length += 2
	} else if value < 1<<16 {
		s.length += 3
	} else if value < 1<<32 {
		s.length += 5
	} else {
		s.length += 9
	}
}

// WriteNillableUint64 calculates the size for a uint64 or nil value.
func (s *Sizer) WriteNillableUint64(value *uint64) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteUint64(*value)
	}
}

// WriteFloat32 adds 5 bytes for a MessagePack float32 value.
func (s *Sizer) WriteFloat32(value float32) {
	s.length += 5
}

// WriteNillableFloat32 calculates the size for a float32 or nil value.
func (s *Sizer) WriteNillableFloat32(value *float32) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteFloat32(*value)
	}
}

// WriteFloat64 adds 9 bytes for a MessagePack float64 value.
func (s *Sizer) WriteFloat64(value float64) {
	s.length += 9
}

// WriteNillableFloat64 calculates the size for a float64 or nil value.
func (s *Sizer) WriteNillableFloat64(value *float64) {
	if value == nil {
		s.WriteNil()
	} else {
		s.WriteFloat64(*value)
	}
}

// WriteAny calculates the size needed for any value.
// It uses type assertions to determine the appropriate sizing method.
func (s *Sizer) WriteAny(value any) {
	if value == nil {
		s.WriteNil()
	}
	switch v := value.(type) {
	case nil:
		s.WriteNil()
	case Encodable:
		v.Encode(s)
	case int:
		s.WriteInt64(int64(v))
	case int8:
		s.WriteInt8(v)
	case int16:
		s.WriteInt16(v)
	case int32:
		s.WriteInt32(v)
	case int64:
		s.WriteInt64(v)
	case uint:
		s.WriteUint64(uint64(v))
	case uint8:
		s.WriteUint8(v)
	case uint16:
		s.WriteUint16(v)
	case uint32:
		s.WriteUint32(v)
	case uint64:
		s.WriteUint64(v)
	case bool:
		s.WriteBool(v)
	case float32:
		s.WriteFloat32(v)
	case float64:
		s.WriteFloat64(v)
	case string:
		s.WriteString(v)
	case time.Time:
		s.WriteTime(v)
	case []byte:
		s.WriteByteArray(v)
	case []any:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteAny(v)
		}
	case []string:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteString(v)
		}
	case []time.Time:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteTime(v)
		}
	case []bool:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteBool(v)
		}
	case []int:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteInt64(int64(v))
		}
	case []int8:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteInt8(v)
		}
	case []int16:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteInt16(v)
		}
	case []int32:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteInt32(v)
		}
	case []int64:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteInt64(v)
		}
	case []uint:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteUint64(uint64(v))
		}
	// case []uint8:
	// 	size := uint32(len(v))
	// 	s.WriteArraySize(size)
	// 	for _, v := range v {
	// 		s.WriteUint8(v)
	// 	}
	case []uint16:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteUint16(v)
		}
	case []uint32:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteUint32(v)
		}
	case []uint64:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteUint64(v)
		}
	case []float32:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteFloat32(v)
		}
	case []float64:
		size := uint32(len(v))
		s.WriteArraySize(size)
		for _, v := range v {
			s.WriteFloat64(v)
		}
	case map[string]string:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteString(k)
			s.WriteString(v)
		}
	case map[string]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteString(k)
			s.WriteAny(v)
		}
	case map[int]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteInt64(int64(k))
			s.WriteAny(v)
		}
	case map[int8]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteInt8(k)
			s.WriteAny(v)
		}
	case map[int16]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteInt16(k)
			s.WriteAny(v)
		}
	case map[int32]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteInt32(k)
			s.WriteAny(v)
		}
	case map[int64]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteInt64(k)
			s.WriteAny(v)
		}
	case map[uint]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteUint64(uint64(k))
			s.WriteAny(v)
		}
	case map[uint8]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteUint8(k)
			s.WriteAny(v)
		}
	case map[uint16]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteUint16(k)
			s.WriteAny(v)
		}
	case map[uint32]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteUint32(k)
			s.WriteAny(v)
		}
	case map[uint64]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteUint64(k)
			s.WriteAny(v)
		}
	case map[any]any:
		size := uint32(len(v))
		s.WriteMapSize(size)
		for k, v := range v {
			s.WriteAny(k)
			s.WriteAny(v)
		}
	}
}

// WriteRaw adds the size of raw MessagePack bytes.
func (s *Sizer) WriteRaw(value Raw) {
	s.length += uint32(len(value))
}

// Err returns any error that occurred during sizing.
// Sizer operations don't typically produce errors, so this usually returns nil.
func (s *Sizer) Err() error {
	return nil
}
