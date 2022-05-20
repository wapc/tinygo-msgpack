package msgpack

import (
	"math"
	"strconv"
)

type Decoder struct {
	reader DataReader
}

func NewDecoder(buffer []byte) Decoder {
	return Decoder{
		reader: NewDataReader(buffer),
	}
}

func (d *Decoder) IsNextNil() (bool, error) {
	prefix, err := d.reader.PeekUint8()
	if err != nil {
		return false, err
	}
	if prefix == FormatNil {
		d.reader.Discard(1)
		return true, nil
	}
	return false, nil
}

func (d *Decoder) ReadBool() (bool, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return false, err
	}
	if prefix == FormatTrue {
		return true, nil
	} else if prefix == FormatFalse {
		return false, nil
	}
	return false, ReadError{"bad value for bool"}
}

func (d *Decoder) ReadInt8() (int8, error) {
	v, err := d.ReadInt64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxInt8 && v >= math.MinInt8 {
		return int8(v), nil
	}
	return 0, ReadError{
		"interger overflow: value = " +
			strconv.FormatInt(v, 10) +
			"; bits = 8",
	}
}

func (d *Decoder) ReadInt16() (int16, error) {
	v, err := d.ReadInt64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxInt16 && v >= math.MinInt16 {
		return int16(v), nil
	}
	return 0, ReadError{
		"interger overflow: value = " +
			strconv.FormatInt(v, 10) +
			"; bits = 16",
	}
}

func (d *Decoder) ReadInt32() (int32, error) {
	v, err := d.ReadInt64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxInt32 && v >= math.MinInt32 {
		return int32(v), nil
	}
	return 0, ReadError{
		"interger overflow: value = " +
			strconv.FormatInt(v, 10) +
			"; bits = 32",
	}
}

func (d *Decoder) ReadInt64() (int64, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedInt(prefix) || isNegativeFixedInt(prefix) {
		return int64(int8(prefix)), nil
	}
	switch prefix {
	case FormatInt8:
		v, err := d.reader.GetInt8()
		return int64(v), err
	case FormatInt16:
		v, err := d.reader.GetInt16()
		return int64(v), err
	case FormatInt32:
		v, err := d.reader.GetInt32()
		return int64(v), err
	case FormatInt64:
		v, err := d.reader.GetInt64()
		return int64(v), err
	default:
		return 0, ReadError{"bad prefix for int64"}
	}
}

func (d *Decoder) ReadUint8() (uint8, error) {
	v, err := d.ReadUint64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxUint8 {
		return uint8(v), nil
	}
	return 0, ReadError{
		"interger overflow: value = " +
			strconv.FormatUint(v, 64) +
			"; bits = 8",
	}
}

func (d *Decoder) ReadUint16() (uint16, error) {
	v, err := d.ReadUint64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxUint16 {
		return uint16(v), nil
	}
	return 0, ReadError{
		"interger overflow: value = " +
			strconv.FormatUint(v, 64) +
			"; bits = 16",
	}
}

func (d *Decoder) ReadUint32() (uint32, error) {
	v, err := d.ReadUint64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxUint32 {
		return uint32(v), nil
	}
	return 0, ReadError{
		"interger overflow: value = " +
			strconv.FormatUint(v, 64) +
			"; bits = 32",
	}
}

func (d *Decoder) ReadUint64() (uint64, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedInt(prefix) {
		return uint64(prefix), nil
	} else if isNegativeFixedInt(prefix) {
		return 0, ReadError{"bad prefix for uint64"}
	}
	switch prefix {
	case FormatUint8:
		v, err := d.reader.GetUint8()
		return uint64(v), err
	case FormatUint16:
		v, err := d.reader.GetUint16()
		return uint64(v), err
	case FormatUint32:
		v, err := d.reader.GetUint32()
		return uint64(v), err
	case FormatUint64:
		v, err := d.reader.GetUint64()
		return uint64(v), err
	default:
		return 0, ReadError{"bad prefix for uint64"}
	}
}

func (d *Decoder) ReadFloat32() (float32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if prefix == FormatFloat32 {
		return d.reader.GetFloat32()
	}
	return 0, ReadError{"bad prefix for float32"}
}

func (d *Decoder) ReadFloat64() (float64, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if prefix == FormatFloat64 {
		return d.reader.GetFloat64()
	}
	return 0, ReadError{"bad prefix for float64"}
}

func (d *Decoder) ReadString() (string, error) {
	strLen, err := d.readStringLength()
	if err != nil {
		return "", err
	}
	strBytes, err := d.reader.GetBytes(strLen)
	if err != nil {
		return "", err
	}
	return string(strBytes), nil
}

func (d *Decoder) readStringLength() (uint32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedString(prefix) {
		return uint32(prefix & 0x1f), nil
	}
	if isFixedArray(prefix) {
		return uint32(prefix & FormatFourLeastSigBitsInByte), nil
	}
	switch prefix {
	case FormatString8:
		v, err := d.reader.GetUint8()
		return uint32(v), err
	case FormatString16:
		v, err := d.reader.GetUint16()
		return uint32(v), err
	case FormatString32:
		v, err := d.reader.GetUint32()
		return v, err
	}
	return 0, ReadError{"bad prefix for string length"}
}

func (d *Decoder) ReadByteArray() ([]byte, error) {
	binLen, err := d.readBinLength()
	if err != nil {
		return nil, err
	}
	binBytes, err := d.reader.GetBytes(binLen)
	if err != nil {
		return nil, err
	}
	return binBytes, nil
}

func (d *Decoder) readBinLength() (uint32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedArray(prefix) {
		return uint32(prefix & FormatFourLeastSigBitsInByte), nil
	}
	switch prefix {
	case FormatBin8:
		v, err := d.reader.GetUint8()
		return uint32(v), err
	case FormatBin16:
		v, err := d.reader.GetUint16()
		return uint32(v), err
	case FormatBin32:
		v, err := d.reader.GetUint32()
		return v, err
	}
	return 0, ReadError{"bad prefix for binary length"}
}

func (d *Decoder) ReadArraySize() (uint32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedArray(prefix) {
		return uint32(prefix & FormatFourLeastSigBitsInByte), nil
	} else if prefix == FormatArray16 {
		v, err := d.reader.GetUint16()
		return uint32(v), err
	} else if prefix == FormatArray32 {
		v, err := d.reader.GetUint32()
		return v, err
	} else if prefix == FormatNil {
		return 0, nil
	}
	return 0, ReadError{"bad prefix for array length"}
}

func (d *Decoder) ReadMapSize() (uint32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedMap(prefix) {
		return uint32(prefix & FormatFourLeastSigBitsInByte), nil
	} else if prefix == FormatMap16 {
		v, err := d.reader.GetUint16()
		return uint32(v), err
	} else if prefix == FormatMap32 {
		v, err := d.reader.GetUint32()
		return v, err
	} else if prefix == FormatNil {
		return 0, nil
	}
	return 0, ReadError{"bad prefix for map length"}
}

func (d *Decoder) Skip() error {
	numberOfObjectsToDiscard, err := d.getSize()
	if err != nil {
		return err
	}

	for numberOfObjectsToDiscard > 0 {
		err = d.Skip() // Skip recursively
		if err != nil {
			return err
		}
		numberOfObjectsToDiscard--
	}
	return nil
}

func (d *Decoder) getSize() (uint32, error) {
	leadByte, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}
	var objectsToDiscard uint32 = 0

	if isNegativeFixedInt(leadByte) || isFixedInt(leadByte) {
		// Noop, will just discard the leadbyte
	} else if isFixedString(leadByte) {
		strLen := uint32(leadByte & 0x1f)
		d.reader.Discard(strLen)
	} else if isFixedArray(leadByte) {
		objectsToDiscard = uint32(leadByte & FormatFourLeastSigBitsInByte)
	} else if isFixedMap(leadByte) {
		objectsToDiscard = 2 * uint32(leadByte&FormatFourLeastSigBitsInByte)
	} else {
		switch leadByte {
		case FormatNil, FormatTrue, FormatFalse:
		case FormatString8, FormatBin8:
			length, err := d.reader.GetUint8()
			if err != nil {
				return 0, err
			}
			err = d.reader.Discard(uint32(length))
			if err != nil {
				return 0, err
			}
		case FormatString16, FormatBin16:
			length, err := d.reader.GetUint16()
			if err != nil {
				return 0, err
			}
			err = d.reader.Discard(uint32(length))
			if err != nil {
				return 0, err
			}
		case FormatString32, FormatBin32:
			length, err := d.reader.GetUint32()
			if err != nil {
				return 0, err
			}
			err = d.reader.Discard(length)
			if err != nil {
				return 0, err
			}
		case FormatFloat32:
			d.reader.Discard(4)
		case FormatFloat64:
			d.reader.Discard(8)
		case FormatUint8, FormatInt8:
			d.reader.Discard(1)
		case FormatUint16, FormatInt16:
			d.reader.Discard(2)
		case FormatUint32, FormatInt32:
			d.reader.Discard(4)
		case FormatUint64, FormatInt64:
			d.reader.Discard(8)
		case FormatFixExt1:
			d.reader.Discard(1)
		case FormatFixExt2:
			d.reader.Discard(3)
		case FormatFixExt4:
			d.reader.Discard(5)
		case FormatFixExt8:
			d.reader.Discard(9)
		case FormatFixExt16:
			d.reader.Discard(17)
		case FormatArray16:
			v, err := d.reader.GetUint16()
			if err != nil {
				return 0, err
			}
			objectsToDiscard = uint32(v)
		case FormatArray32:
			v, err := d.reader.GetUint32()
			if err != nil {
				return 0, err
			}
			objectsToDiscard = v
		case FormatMap16:
			v, err := d.reader.GetUint16()
			if err != nil {
				return 0, err
			}
			objectsToDiscard = 2 * uint32(v)
		case FormatMap32:
			v, err := d.reader.GetUint32()
			if err != nil {
				return 0, err
			}
			objectsToDiscard = 2 * v
		default:
			return 0, ReadError{"bad prefix"}
		}
	}

	return objectsToDiscard, nil
}

////////////////////

//go:inline
func isFixedInt(u byte) bool {
	return u>>7 == 0
}

//go:inline
func isNegativeFixedInt(u byte) bool {
	return (u & 0xe0) == FormatNegativeFixInt
}

//go:inline
func isFixedMap(u byte) bool {
	return (u & 0xf0) == FormatFixMap
}

//go:inline
func isFixedArray(u byte) bool {
	return (u & 0xf0) == FormatFixArray
}

//go:inline
func isFixedString(u byte) bool {
	return (u & 0xe0) == FormatFixString
}

type ReadError struct {
	message string
}

func (e ReadError) Error() string {
	return e.message
}
