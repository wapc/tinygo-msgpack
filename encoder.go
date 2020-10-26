package msgpack

import (
	"math"
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

func (e *Encoder) WriteInt8(value int8) {
	e.WriteInt64(int64(value))
}

func (e *Encoder) WriteInt16(value int16) {
	e.WriteInt64(int64(value))
}

func (e *Encoder) WriteInt32(value int32) {
	e.WriteInt64(int64(value))
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

func (e *Encoder) WriteUint8(value uint8) {
	e.WriteUint64(uint64(value))
}

func (e *Encoder) WriteUint16(value uint16) {
	e.WriteUint64(uint64(value))
}

func (e *Encoder) WriteUint32(value uint32) {
	e.WriteUint64(uint64(value))
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

func (e *Encoder) WriteFloat32(value float32) {
	e.reader.SetUint8(FormatFloat32)
	e.reader.SetFloat32(value)
}

func (e *Encoder) WriteFloat64(value float64) {
	e.reader.SetUint8(FormatFloat64)
	e.reader.SetFloat64(value)
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
		e.WriteNil()
		return
	}
	e.writeBinLength(valueLen)
	e.reader.SetBytes(value)
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
