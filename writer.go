package msgpack

// Writer is the interface for writing data using the MessagPack format.
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
}
