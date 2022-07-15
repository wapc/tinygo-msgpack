package msgpack

// Reader is the interface for reading data from the MessagePack format.
type Reader interface {
	IsNextNil() (bool, error)
	ReadBool() (bool, error)
	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadUint8() (uint8, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint64() (uint64, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	ReadString() (string, error)
	ReadByteArray() ([]byte, error)
	ReadArraySize() (uint32, error)
	ReadMapSize() (uint32, error)
	Skip() error
	Err() error
}

// Writer is the interface for writing data to the MessagePack format.
type Writer interface {
	WriteNil()
	WriteBool(value bool)
	WriteInt8(value int8)
	WriteInt16(value int16)
	WriteInt32(value int32)
	WriteInt64(value int64)
	WriteUint8(value uint8)
	WriteUint16(value uint16)
	WriteUint32(value uint32)
	WriteUint64(value uint64)
	WriteFloat32(value float32)
	WriteFloat64(value float64)
	WriteString(value string)
	WriteByteArray(value []byte)
	WriteArraySize(length uint32)
	WriteMapSize(length uint32)
	Err() error
}
