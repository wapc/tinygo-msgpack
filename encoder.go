package msgpack

import (
	"encoding/binary"
	"math"
	"time"
)

// Encoder provides low-level MessagePack encoding functionality.
// It writes MessagePack data to a pre-allocated buffer.
type Encoder struct {
	reader DataReader
}

// Ensure `*Encoder` implements `Writer`.
var _ Writer = (*Encoder)(nil)

// NewEncoder creates a new Encoder that writes to the provided buffer.
// The buffer must be large enough to hold all encoded data.
func NewEncoder(buffer []byte) Encoder {
	return Encoder{
		reader: NewDataReader(buffer),
	}
}

// WriteNil writes a MessagePack nil value.
func (e *Encoder) WriteNil() {
	e.reader.SetUint8(FormatNil)
}

// WriteBool writes a MessagePack boolean value.
func (e *Encoder) WriteBool(value bool) {
	if value {
		e.reader.SetUint8(FormatTrue)
	} else {
		e.reader.SetUint8(FormatFalse)
	}
}

// WriteNillableBool writes a MessagePack boolean value or nil if the pointer is nil.
func (e *Encoder) WriteNillableBool(value *bool) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteBool(*value)
	}
}

// WriteInt8 writes a MessagePack int8 value.
// The value is encoded using the most efficient MessagePack integer format.
func (e *Encoder) WriteInt8(value int8) {
	e.WriteInt64(int64(value))
}

// WriteNillableInt8 writes a MessagePack int8 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableInt8(value *int8) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt8(*value)
	}
}

// WriteInt16 writes a MessagePack int16 value.
// The value is encoded using the most efficient MessagePack integer format.
func (e *Encoder) WriteInt16(value int16) {
	e.WriteInt64(int64(value))
}

// WriteNillableInt16 writes a MessagePack int16 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableInt16(value *int16) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt16(*value)
	}
}

// WriteInt32 writes a MessagePack int32 value.
// The value is encoded using the most efficient MessagePack integer format.
func (e *Encoder) WriteInt32(value int32) {
	e.WriteInt64(int64(value))
}

// WriteNillableInt32 writes a MessagePack int32 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableInt32(value *int32) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt32(*value)
	}
}

// WriteInt64 writes a MessagePack int64 value using the most efficient format:
// - Positive fixint (0x00-0x7f) for values 0-127
// - Negative fixint (0xe0-0xff) for values -32 to -1
// - int8, int16, int32, or int64 format for larger values
func (e *Encoder) WriteInt64(value int64) {
	if value >= 0 && value < 1<<7 {
		e.reader.SetUint8(uint8(value))
	} else if value < 0 && value >= -(1<<5) {
		e.reader.SetUint8(uint8(value) | FormatNegativeFixInt)
	} else if value <= math.MaxInt8 && value >= math.MinInt8 {
		e.reader.SetUint8(FormatInt8)
		e.reader.SetInt8(int8(value))
	} else if value <= math.MaxInt16 && value >= math.MinInt16 {
		e.reader.SetUint8(FormatInt16)
		e.reader.SetInt16(int16(value))
	} else if value <= math.MaxInt32 && value >= math.MinInt32 {
		e.reader.SetUint8(FormatInt32)
		e.reader.SetInt32(int32(value))
	} else {
		e.reader.SetUint8(FormatInt64)
		e.reader.SetInt64(value)
	}
}

// WriteNillableInt64 writes a MessagePack int64 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableInt64(value *int64) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt64(*value)
	}
}

// WriteUint8 writes a MessagePack uint8 value.
// The value is encoded using the most efficient MessagePack integer format.
func (e *Encoder) WriteUint8(value uint8) {
	e.WriteUint64(uint64(value))
}

// WriteNillableUint8 writes a MessagePack uint8 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableUint8(value *uint8) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint8(*value)
	}
}

// WriteUint16 writes a MessagePack uint16 value.
// The value is encoded using the most efficient MessagePack integer format.
func (e *Encoder) WriteUint16(value uint16) {
	e.WriteUint64(uint64(value))
}

// WriteNillableUint16 writes a MessagePack uint16 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableUint16(value *uint16) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint16(*value)
	}
}

// WriteUint32 writes a MessagePack uint32 value.
// The value is encoded using the most efficient MessagePack integer format.
func (e *Encoder) WriteUint32(value uint32) {
	e.WriteUint64(uint64(value))
}

// WriteNillableUint32 writes a MessagePack uint32 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableUint32(value *uint32) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint32(*value)
	}
}

// WriteUint64 writes a MessagePack uint64 value using the most efficient format:
// - Positive fixint (0x00-0x7f) for values 0-127
// - uint8, uint16, uint32, or uint64 format for larger values
func (e *Encoder) WriteUint64(value uint64) {
	if value < 1<<7 {
		e.reader.SetUint8(uint8(value))
	} else if value <= math.MaxUint8 {
		e.reader.SetUint8(FormatUint8)
		e.reader.SetUint8(uint8(value))
	} else if value <= math.MaxUint16 {
		e.reader.SetUint8(FormatUint16)
		e.reader.SetUint16(uint16(value))
	} else if value <= math.MaxUint32 {
		e.reader.SetUint8(FormatUint32)
		e.reader.SetUint32(uint32(value))
	} else {
		e.reader.SetUint8(FormatUint64)
		e.reader.SetUint64(value)
	}
}

// WriteNillableUint64 writes a MessagePack uint64 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableUint64(value *uint64) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint64(*value)
	}
}

// WriteFloat32 writes a MessagePack float32 value.
func (e *Encoder) WriteFloat32(value float32) {
	e.reader.SetUint8(FormatFloat32)
	e.reader.SetFloat32(value)
}

// WriteNillableFloat32 writes a MessagePack float32 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableFloat32(value *float32) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteFloat32(*value)
	}
}

// WriteFloat64 writes a MessagePack float64 value.
func (e *Encoder) WriteFloat64(value float64) {
	e.reader.SetUint8(FormatFloat64)
	e.reader.SetFloat64(value)
}

// WriteNillableFloat64 writes a MessagePack float64 value or nil if the pointer is nil.
func (e *Encoder) WriteNillableFloat64(value *float64) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteFloat64(*value)
	}
}

// writeStringLength writes the MessagePack string length header using the most efficient format:
// - Fixstr (0xa0-0xbf) for lengths 0-31
// - str8, str16, or str32 format for longer strings
func (e *Encoder) writeStringLength(length uint32) {
	if length < 32 {
		e.reader.SetUint8(uint8(length) | FormatFixString)
	} else if length <= math.MaxUint8 {
		e.reader.SetUint8(FormatString8)
		e.reader.SetUint8(uint8(length))
	} else if length <= math.MaxUint16 {
		e.reader.SetUint8(FormatString16)
		e.reader.SetUint16(uint16(length))
	} else {
		e.reader.SetUint8(FormatString32)
		e.reader.SetUint32(length)
	}
}

// WriteString writes a MessagePack string value.
func (e *Encoder) WriteString(value string) {
	valueBytes := UnsafeBytes(value)
	e.writeStringLength(uint32(len(valueBytes)))
	e.reader.SetBytes(valueBytes)
}

// WriteNillableString writes a MessagePack string value or nil if the pointer is nil.
func (e *Encoder) WriteNillableString(value *string) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteString(*value)
	}
}

// WriteTime writes a MessagePack time value as an extension type.
// Times are encoded using extension type -1 with a compact binary format.
func (e *Encoder) WriteTime(tm time.Time) {
	var timeBuf [12]byte
	b := e.encodeTime(tm, timeBuf[:])
	e.encodeExtLen(len(b))
	e.reader.SetInt8(-1)
	e.reader.SetBytes(b)
}

// WriteNillableTime writes a MessagePack time value or nil if the pointer is nil.
func (e *Encoder) WriteNillableTime(value *time.Time) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteTime(*value)
	}
}

// encodeTime encodes a time.Time value into a compact binary format.
// The format optimizes for common time ranges to minimize size.
func (e *Encoder) encodeTime(tm time.Time, timeBuf []byte) []byte {
	secs := uint64(tm.Unix())
	if secs>>34 == 0 {
		data := uint64(tm.Nanosecond())<<34 | secs

		if data&0xffffffff00000000 == 0 {
			b := timeBuf[:4]
			binary.BigEndian.PutUint32(b, uint32(data))
			return b
		}

		b := timeBuf[:8]
		binary.BigEndian.PutUint64(b, data)
		return b
	}

	b := timeBuf[:12]
	binary.BigEndian.PutUint32(b, uint32(tm.Nanosecond()))
	binary.BigEndian.PutUint64(b[4:], secs)
	return b
}

// writeBinLength writes the MessagePack binary length header using the most efficient format:
// - bin8, bin16, or bin32 format based on length
func (e *Encoder) writeBinLength(length uint32) {
	if length <= math.MaxUint8 {
		e.reader.SetUint8(FormatBin8)
		e.reader.SetUint8(uint8(length))
	} else if length <= math.MaxUint16 {
		e.reader.SetUint8(FormatBin16)
		e.reader.SetUint16(uint16(length))
	} else {
		e.reader.SetUint8(FormatBin32)
		e.reader.SetUint32(length)
	}
}

// WriteByteArray writes a MessagePack binary value.
func (e *Encoder) WriteByteArray(value []byte) {
	valueLen := uint32(len(value))
	if valueLen == 0 {
		e.reader.SetUint8(FormatBin8)
		e.reader.SetUint8(0)
		return
	}
	e.writeBinLength(valueLen)
	e.reader.SetBytes(value)
}

// WriteNillableByteArray writes a MessagePack binary value or nil if the slice is nil.
func (e *Encoder) WriteNillableByteArray(value []byte) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteByteArray(value)
	}
}

// WriteArraySize writes a MessagePack array size header using the most efficient format:
// - Fixarray (0x90-0x9f) for sizes 0-15
// - array16 or array32 format for larger arrays
func (e *Encoder) WriteArraySize(length uint32) {
	if length < 16 {
		e.reader.SetUint8(uint8(length) | FormatFixArray)
	} else if length <= math.MaxUint16 {
		e.reader.SetUint8(FormatArray16)
		e.reader.SetUint16(uint16(length))
	} else {
		e.reader.SetUint8(FormatArray32)
		e.reader.SetUint32(length)
	}
}

// WriteMapSize writes a MessagePack map size header using the most efficient format:
// - Fixmap (0x80-0x8f) for sizes 0-15
// - map16 or map32 format for larger maps
func (e *Encoder) WriteMapSize(length uint32) {
	if length < 16 {
		e.reader.SetUint8(uint8(length) | FormatFixMap)
	} else if length <= math.MaxUint16 {
		e.reader.SetUint8(FormatMap16)
		e.reader.SetUint16(uint16(length))
	} else {
		e.reader.SetUint8(FormatMap32)
		e.reader.SetUint32(length)
	}
}

// WriteAny writes any value to MessagePack format.
// It uses type assertions to determine the appropriate encoding method.
func (e *Encoder) WriteAny(value any) {
	if isNil(value) {
		e.WriteNil()
		return
	}

	switch v := value.(type) {
	case nil:
		e.WriteNil()
	case Encodable:
		v.Encode(e)
	case int:
		e.WriteInt64(int64(v))
	case int8:
		e.WriteInt8(v)
	case int16:
		e.WriteInt16(v)
	case int32:
		e.WriteInt32(v)
	case int64:
		e.WriteInt64(v)
	case uint:
		e.WriteUint64(uint64(v))
	case uint8:
		e.WriteUint8(v)
	case uint16:
		e.WriteUint16(v)
	case uint32:
		e.WriteUint32(v)
	case uint64:
		e.WriteUint64(v)
	case bool:
		e.WriteBool(v)
	case float32:
		e.WriteFloat32(v)
	case float64:
		e.WriteFloat64(v)
	case string:
		e.WriteString(v)
	case []byte:
		e.WriteByteArray(v)
	case []interface{}:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteAny(v)
		}
	case []string:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteString(v)
		}
	case []bool:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteBool(v)
		}
	case []int:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteInt64(int64(v))
		}
	case []int8:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteInt8(v)
		}
	case []int16:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteInt16(v)
		}
	case []int32:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteInt32(v)
		}
	case []int64:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteInt64(v)
		}

	case []uint:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteUint64(uint64(v))
		}
	case []uint16:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteUint16(v)
		}
	case []uint32:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteUint32(v)
		}
	case []uint64:
		size := uint32(len(v))
		e.WriteArraySize(size)
		for _, v := range v {
			e.WriteUint64(v)
		}

	case map[string]string:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteString(k)
			e.WriteString(v)
		}
	case map[string]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteString(k)
			e.WriteAny(v)
		}
	case map[int]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteInt64(int64(k))
			e.WriteAny(v)
		}
	case map[int8]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteInt8(k)
			e.WriteAny(v)
		}
	case map[int16]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteInt16(k)
			e.WriteAny(v)
		}
	case map[int32]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteInt32(k)
			e.WriteAny(v)
		}
	case map[int64]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteInt64(k)
			e.WriteAny(v)
		}
	case map[uint]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteUint64(uint64(k))
			e.WriteAny(v)
		}
	case map[uint8]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteUint8(k)
			e.WriteAny(v)
		}
	case map[uint16]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteUint16(k)
			e.WriteAny(v)
		}
	case map[uint32]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteUint32(k)
			e.WriteAny(v)
		}
	case map[uint64]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteUint64(k)
			e.WriteAny(v)
		}
	case map[interface{}]interface{}:
		size := uint32(len(v))
		e.WriteMapSize(size)
		for k, v := range v {
			e.WriteAny(k)
			e.WriteAny(v)
		}
	}
}

// WriteRaw writes raw MessagePack bytes without any processing.
func (e *Encoder) WriteRaw(value Raw) {
	e.reader.SetBytes(value)
}

// encodeExtLen writes the MessagePack extension length header using the most efficient format:
// - Fixext1, Fixext2, Fixext4, Fixext8, Fixext16 for specific lengths
// - ext8, ext16, or ext32 format for other lengths
func (e *Encoder) encodeExtLen(l int) error {
	switch l {
	case 1:
		return e.reader.SetUint8(FormatFixExt1)
	case 2:
		return e.reader.SetUint8(FormatFixExt2)
	case 4:
		return e.reader.SetUint8(FormatFixExt4)
	case 8:
		return e.reader.SetUint8(FormatFixExt8)
	case 16:
		return e.reader.SetUint8(FormatFixExt16)
	}
	if l <= math.MaxUint8 {
		return e.write1(FormatExt8, uint8(l))
	}
	if l <= math.MaxUint16 {
		return e.write2(FormatExt16, uint16(l))
	}
	return e.write4(FormatExt32, uint32(l))
}

// write1 writes a 1-byte extension header with the given code and length.
func (e *Encoder) write1(code byte, n uint8) error {
	var buf [2]byte
	buf[0] = code
	buf[1] = n
	return e.reader.SetBytes(buf[:])
}

// write2 writes a 2-byte extension header with the given code and length.
func (e *Encoder) write2(code byte, n uint16) error {
	var buf [3]byte
	buf[0] = code
	buf[1] = byte(n >> 8)
	buf[2] = byte(n)
	return e.reader.SetBytes(buf[:])
}

// write4 writes a 4-byte extension header with the given code and length.
func (e *Encoder) write4(code byte, n uint32) error {
	var buf [5]byte
	buf[0] = code
	buf[1] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[3] = byte(n >> 8)
	buf[4] = byte(n)
	return e.reader.SetBytes(buf[:])
}

// Err returns any error that occurred during encoding.
func (e *Encoder) Err() error {
	return e.reader.Err()
}
