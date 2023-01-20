package msgpack

import (
	"encoding/binary"
	"math"
	"time"
)

type Encoder struct {
	reader DataReader
}

func NewEncoder(buffer []byte) Encoder {
	return Encoder{
		reader: NewDataReader(buffer),
	}
}

func (e *Encoder) WriteNil() {
	e.reader.SetUint8(FormatNil)
}

func (e *Encoder) WriteBool(value bool) {
	if value {
		e.reader.SetUint8(FormatTrue)
	} else {
		e.reader.SetUint8(FormatFalse)
	}
}

func (e *Encoder) WriteNillableBool(value *bool) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteBool(*value)
	}
}

func (e *Encoder) WriteInt8(value int8) {
	e.WriteInt64(int64(value))
}

func (e *Encoder) WriteNillableInt8(value *int8) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt8(*value)
	}
}

func (e *Encoder) WriteInt16(value int16) {
	e.WriteInt64(int64(value))
}

func (e *Encoder) WriteNillableInt16(value *int16) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt16(*value)
	}
}

func (e *Encoder) WriteInt32(value int32) {
	e.WriteInt64(int64(value))
}

func (e *Encoder) WriteNillableInt32(value *int32) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt32(*value)
	}
}

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

func (e *Encoder) WriteNillableInt64(value *int64) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteInt64(*value)
	}
}

func (e *Encoder) WriteUint8(value uint8) {
	e.WriteUint64(uint64(value))
}

func (e *Encoder) WriteNillableUint8(value *uint8) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint8(*value)
	}
}

func (e *Encoder) WriteUint16(value uint16) {
	e.WriteUint64(uint64(value))
}

func (e *Encoder) WriteNillableUint16(value *uint16) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint16(*value)
	}
}

func (e *Encoder) WriteUint32(value uint32) {
	e.WriteUint64(uint64(value))
}

func (e *Encoder) WriteNillableUint32(value *uint32) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint32(*value)
	}
}

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

func (e *Encoder) WriteNillableUint64(value *uint64) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteUint64(*value)
	}
}

func (e *Encoder) WriteFloat32(value float32) {
	e.reader.SetUint8(FormatFloat32)
	e.reader.SetFloat32(value)
}

func (e *Encoder) WriteNillableFloat32(value *float32) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteFloat32(*value)
	}
}

func (e *Encoder) WriteFloat64(value float64) {
	e.reader.SetUint8(FormatFloat64)
	e.reader.SetFloat64(value)
}

func (e *Encoder) WriteNillableFloat64(value *float64) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteFloat64(*value)
	}
}

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

func (e *Encoder) WriteString(value string) {
	valueBytes := UnsafeBytes(value)
	e.writeStringLength(uint32(len(valueBytes)))
	e.reader.SetBytes(valueBytes)
}

func (e *Encoder) WriteNillableString(value *string) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteString(*value)
	}
}

func (e *Encoder) WriteTime(tm time.Time) {
	var timeBuf [12]byte
	b := e.encodeTime(tm, timeBuf[:])
	e.encodeExtLen(len(b))
	e.reader.SetInt8(-1)
	e.reader.SetBytes(b)
}

func (e *Encoder) WriteNillableTime(value *time.Time) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteTime(*value)
	}
}

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

func (e *Encoder) WriteNillableByteArray(value []byte) {
	if value == nil {
		e.WriteNil()
	} else {
		e.WriteByteArray(value)
	}
}

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

func (e *Encoder) WriteAny(value any) {
	if value == nil {
		e.WriteNil()
	}
	switch v := value.(type) {
	case nil:
		e.WriteNil()
	case Codec:
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

func (e *Encoder) write1(code byte, n uint8) error {
	var buf [2]byte
	buf[0] = code
	buf[1] = n
	return e.reader.SetBytes(buf[:])
}

func (e *Encoder) write2(code byte, n uint16) error {
	var buf [3]byte
	buf[0] = code
	buf[1] = byte(n >> 8)
	buf[2] = byte(n)
	return e.reader.SetBytes(buf[:])
}

func (e *Encoder) write4(code byte, n uint32) error {
	var buf [5]byte
	buf[0] = code
	buf[1] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[3] = byte(n >> 8)
	buf[4] = byte(n)
	return e.reader.SetBytes(buf[:])
}

func (e *Encoder) Err() error {
	return e.reader.Err()
}
