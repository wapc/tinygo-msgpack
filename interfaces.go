package msgpack

import (
	"time"
)

type Raw []byte

// Reader is the interface for reading data from the MessagePack format.
type Reader interface {
	IsNextNil() (bool, error)
	ReadBool() (bool, error)
	ReadNillableBool() (*bool, error)
	ReadInt8() (int8, error)
	ReadNillableInt8() (*int8, error)
	ReadInt16() (int16, error)
	ReadNillableInt16() (*int16, error)
	ReadInt32() (int32, error)
	ReadNillableInt32() (*int32, error)
	ReadInt64() (int64, error)
	ReadNillableInt64() (*int64, error)
	ReadUint8() (uint8, error)
	ReadNillableUint8() (*uint8, error)
	ReadUint16() (uint16, error)
	ReadNillableUint16() (*uint16, error)
	ReadUint32() (uint32, error)
	ReadNillableUint32() (*uint32, error)
	ReadUint64() (uint64, error)
	ReadNillableUint64() (*uint64, error)
	ReadFloat32() (float32, error)
	ReadNillableFloat32() (*float32, error)
	ReadFloat64() (float64, error)
	ReadNillableFloat64() (*float64, error)
	ReadString() (string, error)
	ReadNillableString() (*string, error)
	ReadTime() (time.Time, error)
	ReadNillableTime() (*time.Time, error)
	ReadByteArray() ([]byte, error)
	ReadNillableByteArray() ([]byte, error)
	ReadArraySize() (uint32, error)
	ReadMapSize() (uint32, error)
	ReadAny() (any, error)
	ReadRaw() (Raw, error)
	Skip() error
	Err() error
}

// Writer is the interface for writing data to the MessagePack format.
type Writer interface {
	WriteNil()
	WriteBool(value bool)
	WriteNillableBool(value *bool)
	WriteInt8(value int8)
	WriteNillableInt8(value *int8)
	WriteInt16(value int16)
	WriteNillableInt16(value *int16)
	WriteInt32(value int32)
	WriteNillableInt32(value *int32)
	WriteInt64(value int64)
	WriteNillableInt64(value *int64)
	WriteUint8(value uint8)
	WriteNillableUint8(value *uint8)
	WriteUint16(value uint16)
	WriteNillableUint16(value *uint16)
	WriteUint32(value uint32)
	WriteNillableUint32(value *uint32)
	WriteUint64(value uint64)
	WriteNillableUint64(value *uint64)
	WriteFloat32(value float32)
	WriteNillableFloat32(value *float32)
	WriteFloat64(value float64)
	WriteNillableFloat64(value *float64)
	WriteString(value string)
	WriteNillableString(value *string)
	WriteTime(value time.Time)
	WriteNillableTime(value *time.Time)
	WriteByteArray(value []byte)
	WriteNillableByteArray(value []byte)
	WriteArraySize(length uint32)
	WriteMapSize(length uint32)
	WriteAny(value any)
	WriteRaw(value Raw)
	Err() error
}
