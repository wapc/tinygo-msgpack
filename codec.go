package msgpack

// Codec is the interface that applies to data structures that can
// encode to and decode from the MessagPack format.
type Codec interface {
	Decode(Reader) error
	Encode(Writer) error
}

// ToBytes creates a `[]byte` from `codec`.
func ToBytes(codec Codec) ([]byte, error) {
	var sizer Sizer
	if err := codec.Encode(&sizer); err != nil {
		return nil, err
	}
	buffer := make([]byte, sizer.Len())
	encoder := NewEncoder(buffer)
	if err := codec.Encode(&encoder); err != nil {
		return nil, err
	}
	return buffer, nil
}

// AnyToBytes creates a `[]byte` from `value`.
func AnyToBytes(value interface{}) ([]byte, error) {
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
