package msgpack

import (
	"encoding/binary"
	"math"
	"strconv"
	"time"
)

// Decoder provides low-level MessagePack decoding functionality.
// It reads MessagePack data from a buffer.
type Decoder struct {
	reader DataReader
}

// Ensure `*Decoder` implements `Reader`.
var _ Reader = (*Decoder)(nil)

// NewDecoder creates a new Decoder that reads from the provided buffer.
func NewDecoder(buffer []byte) Decoder {
	return Decoder{
		reader: NewDataReader(buffer),
	}
}

// IsNextNil checks if the next value in the buffer is a MessagePack nil.
// It peeks at the next byte without advancing the read position.
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

// ReadBool reads a MessagePack boolean value.
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

// ReadNillableBool reads a MessagePack boolean value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableBool() (*bool, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadBool()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadInt8 reads a MessagePack int8 value.
// Performs range checking to ensure the value fits in 8 bits.
func (d *Decoder) ReadInt8() (int8, error) {
	v, err := d.ReadInt64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxInt8 && v >= math.MinInt8 {
		return int8(v), nil
	}
	return 0, ReadError{
		"integer overflow: value = " +
			strconv.FormatInt(v, 10) +
			"; bits = 8",
	}
}

// ReadNillableInt8 reads a MessagePack int8 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableInt8() (*int8, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadInt8()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadInt16 reads a MessagePack int16 value.
// Performs range checking to ensure the value fits in 16 bits.
func (d *Decoder) ReadInt16() (int16, error) {
	v, err := d.ReadInt64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxInt16 && v >= math.MinInt16 {
		return int16(v), nil
	}
	return 0, ReadError{
		"integer overflow: value = " +
			strconv.FormatInt(v, 10) +
			"; bits = 16",
	}
}

// ReadNillableInt16 reads a MessagePack int16 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableInt16() (*int16, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadInt16()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadInt32 reads a MessagePack int32 value.
// Performs range checking to ensure the value fits in 32 bits.
func (d *Decoder) ReadInt32() (int32, error) {
	v, err := d.ReadInt64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxInt32 && v >= math.MinInt32 {
		return int32(v), nil
	}
	return 0, ReadError{
		"integer overflow: value = " +
			strconv.FormatInt(v, 10) +
			"; bits = 32",
	}
}

// ReadNillableInt32 reads a MessagePack int32 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableInt32() (*int32, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadInt32()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadInt64 reads a MessagePack int64 value.
// Supports all MessagePack integer formats including fixint and typed integers.
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

// ReadNillableInt64 reads a MessagePack int64 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableInt64() (*int64, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadInt64()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadUint8 reads a MessagePack uint8 value.
// Performs range checking to ensure the value fits in 8 bits.
func (d *Decoder) ReadUint8() (uint8, error) {
	v, err := d.ReadUint64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxUint8 {
		return uint8(v), nil
	}
	return 0, ReadError{
		"integer overflow: value = " +
			strconv.FormatUint(v, 10) +
			"; bits = 8",
	}
}

// ReadNillableUint8 reads a MessagePack uint8 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableUint8() (*uint8, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadUint8()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadUint16 reads a MessagePack uint16 value.
// Performs range checking to ensure the value fits in 16 bits.
func (d *Decoder) ReadUint16() (uint16, error) {
	v, err := d.ReadUint64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxUint16 {
		return uint16(v), nil
	}
	return 0, ReadError{
		"integer overflow: value = " +
			strconv.FormatUint(v, 10) +
			"; bits = 16",
	}
}

// ReadNillableUint16 reads a MessagePack uint16 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableUint16() (*uint16, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadUint16()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadUint32 reads a MessagePack uint32 value.
// Performs range checking to ensure the value fits in 32 bits.
func (d *Decoder) ReadUint32() (uint32, error) {
	v, err := d.ReadUint64()
	if err != nil {
		return 0, err
	}
	if v <= math.MaxUint32 {
		return uint32(v), nil
	}
	return 0, ReadError{
		"integer overflow: value = " +
			strconv.FormatUint(v, 10) +
			"; bits = 32",
	}
}

// ReadNillableUint32 reads a MessagePack uint32 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableUint32() (*uint32, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadUint32()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadUint64 reads a MessagePack uint64 value.
// Supports all MessagePack unsigned integer formats including fixint and typed integers.
func (d *Decoder) ReadUint64() (uint64, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedInt(prefix) {
		return uint64(prefix), nil
	} else if isNegativeFixedInt(prefix) {
		v := int8(prefix)
		if v < 0 {
			return 0, ReadError{"bad prefix for uint"}
		}
		return uint64(v), err
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
	case FormatInt8:
		v, err := d.reader.GetInt8()
		if v < 0 {
			return 0, ReadError{"bad prefix for uint"}
		}
		return uint64(v), err
	case FormatInt16:
		v, err := d.reader.GetInt16()
		if v < 0 {
			return 0, ReadError{"bad prefix for uint"}
		}
		return uint64(v), err
	case FormatInt32:
		v, err := d.reader.GetInt32()
		if v < 0 {
			return 0, ReadError{"bad prefix for uint"}
		}
		return uint64(v), err
	case FormatInt64:
		v, err := d.reader.GetInt64()
		if v < 0 {
			return 0, ReadError{"bad prefix for uint"}
		}
		return uint64(v), err
	default:
		return 0, ReadError{"bad prefix for uint"}
	}
}

// ReadNillableUint64 reads a MessagePack uint64 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableUint64() (*uint64, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadUint64()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadFloat32 reads a MessagePack float32 value.
// Also accepts float64 values and converts them to float32.
func (d *Decoder) ReadFloat32() (float32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	switch prefix {
	case FormatFloat32:
		return d.reader.GetFloat32()
	case FormatFloat64:
		v, err := d.reader.GetFloat64()
		return float32(v), err
	}
	return 0, ReadError{"bad prefix for float32"}
}

// ReadNillableFloat32 reads a MessagePack float32 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableFloat32() (*float32, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadFloat32()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadFloat64 reads a MessagePack float64 value.
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

// ReadNillableFloat64 reads a MessagePack float64 value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableFloat64() (*float64, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadFloat64()
	if err != nil {
		return nil, err
	}
	return &val, err
}

func (d *Decoder) ReadTime() (time.Time, error) {
	prefix, err := d.reader.PeekUint8()
	if err != nil {
		return time.Time{}, err
	}

	if isString(prefix) {
		str, err := d.ReadString()
		if err != nil {
			return time.Time{}, err
		}
		return time.Parse(time.RFC3339Nano, str)
	}

	d.reader.Discard(1)
	extID, extLen, err := d.extHeader(prefix)
	if err != nil {
		return time.Time{}, err
	}

	// NodeJS seems to use extID 13.
	if extID != -1 && extID != 13 {
		return time.Time{}, ReadError{"msgpack: invalid time ext id=" + strconv.FormatUint(uint64(extID), 10)}
	}

	tm, err := d.decodeTime(extLen)
	if err != nil {
		return tm, err
	}

	if tm.IsZero() {
		// Zero time does not have timezone information.
		return tm.UTC(), nil
	}

	return tm, nil
}

func (d *Decoder) ReadNillableTime() (*time.Time, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadTime()
	if err != nil {
		return nil, err
	}
	return &val, err
}

func (d *Decoder) decodeTime(extLen uint32) (time.Time, error) {
	b, err := d.reader.GetBytes(extLen)
	if err != nil {
		return time.Time{}, err
	}

	switch len(b) {
	case 4:
		sec := binary.BigEndian.Uint32(b)
		return time.Unix(int64(sec), 0), nil
	case 8:
		sec := binary.BigEndian.Uint64(b)
		nsec := int64(sec >> 34)
		sec &= 0x00000003ffffffff
		return time.Unix(int64(sec), nsec), nil
	case 12:
		nsec := binary.BigEndian.Uint32(b)
		sec := binary.BigEndian.Uint64(b[4:])
		return time.Unix(int64(sec), int64(nsec)), nil
	default:
		return time.Time{}, ReadError{"msgpack: invalid time ext len=" + strconv.FormatUint(uint64(extLen), 10)}
	}
}

// ReadString reads a MessagePack string value.
// Supports both fixed and variable length string formats.
func (d *Decoder) ReadString() (string, error) {
	strLen, err := d.readStringLength()
	return d.readString(strLen, err)
}

// ReadNillableString reads a MessagePack string value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableString() (*string, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	val, err := d.ReadString()
	if err != nil {
		return nil, err
	}
	return &val, err
}

// readStringLength reads the length prefix for a MessagePack string.
// Supports fixed string, fixed array, and variable length string formats.
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
	case FormatString16, FormatArray16:
		v, err := d.reader.GetUint16()
		return uint32(v), err
	case FormatString32, FormatArray32:
		v, err := d.reader.GetUint32()
		return v, err
	}

	return 0, ReadError{"bad prefix for string length"}
}

// readString reads a string of the specified length from the buffer.
// Uses unsafe string conversion for performance.
func (d *Decoder) readString(strLen uint32, err error) (string, error) {
	if err != nil {
		return "", err
	}
	strBytes, err := d.reader.GetBytes(strLen)
	if err != nil {
		return "", err
	}
	return UnsafeString(strBytes), nil
}

// ReadByteArray reads a MessagePack binary array value.
// Returns the raw bytes without any conversion.
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

// ReadNillableByteArray reads a MessagePack binary array value or nil.
// Returns nil if the next value is MessagePack nil.
func (d *Decoder) ReadNillableByteArray() ([]byte, error) {
	isNil, err := d.IsNextNil()
	if isNil || err != nil {
		return nil, err
	}
	return d.ReadByteArray()
}

// readBinLength reads the length prefix for a MessagePack binary array.
// Supports fixed array and variable length binary formats.
func (d *Decoder) readBinLength() (uint32, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}

	if isFixedArray(prefix) {
		return uint32(prefix & FormatFourLeastSigBitsInByte), nil
	}
	switch prefix {
	case FormatNil:
		return 0, nil
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

// ReadArraySize reads the size prefix for a MessagePack array.
// Returns the number of elements in the array, or 0 for nil arrays.
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

// ReadMapSize reads the size prefix for a MessagePack map.
// Returns the number of key-value pairs in the map, or 0 for nil maps.
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

// ReadRaw reads the next MessagePack value as raw bytes without decoding.
// Returns the complete MessagePack representation of the value.
func (d *Decoder) ReadRaw() (Raw, error) {
	start := d.reader.byteOffset
	if err := d.Skip(); err != nil {
		return nil, err
	}
	end := d.reader.byteOffset
	bytes := d.reader.buffer[start:end]

	return bytes, nil
}

// Skip advances the decoder past the next MessagePack value.
// Recursively skips nested structures like arrays and maps.
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

// getSize determines the number of objects to skip for the current MessagePack value.
// Handles all MessagePack format types and returns the count of nested objects.
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

// ReadAny reads any MessagePack value and returns it as an interface{}.
// Supports all MessagePack types including primitive types, arrays, and maps.
func (d *Decoder) ReadAny() (any, error) {
	prefix, err := d.reader.GetUint8()
	if err != nil {
		return false, err
	}

	if isFixedInt(prefix) || isNegativeFixedInt(prefix) {
		return int64(int8(prefix)), nil
	}

	if isFixedString(prefix) {
		strLen := uint32(prefix & 0x1f)
		return d.readString(strLen, nil)
	}

	if isFixedArray(prefix) {
		aryLen := uint32(prefix & FormatFourLeastSigBitsInByte)
		ary := make([]any, aryLen)
		err := d.readArray(ary)
		return ary, err
	}

	if isFixedMap(prefix) {
		mapLen := uint32(prefix & FormatFourLeastSigBitsInByte)
		m := make(map[any]any, mapLen)
		err := d.readMap(m, mapLen)
		return m, err
	}

	switch prefix {
	case FormatNil:
		return nil, nil
	case FormatTrue:
		return true, nil
	case FormatFalse:
		return false, nil
	case FormatInt8:
		return d.reader.GetInt8()
	case FormatInt16:
		return d.reader.GetInt16()
	case FormatInt32:
		return d.reader.GetInt32()
	case FormatInt64:
		return d.reader.GetInt64()
	case FormatUint8:
		return d.reader.GetUint8()
	case FormatUint16:
		return d.reader.GetUint16()
	case FormatUint32:
		return d.reader.GetUint32()
	case FormatUint64:
		return d.reader.GetUint64()
	case FormatFloat32:
		return d.reader.GetFloat32()
	case FormatFloat64:
		return d.reader.GetFloat64()
	case FormatString8:
		v, err := d.reader.GetUint8()
		return d.readString(uint32(v), err)
	case FormatString16:
		v, err := d.reader.GetUint16()
		return d.readString(uint32(v), err)
	case FormatString32:
		v, err := d.reader.GetUint32()
		return d.readString(uint32(v), err)
	case FormatArray16:
		v, err := d.reader.GetUint16()
		if err != nil {
			return nil, err
		}
		ary := make([]any, v)
		err = d.readArray(ary)
		return ary, err
	case FormatArray32:
		v, err := d.reader.GetUint32()
		if err != nil {
			return nil, err
		}
		ary := make([]any, v)
		err = d.readArray(ary)
		return ary, err
	case FormatMap16:
		v, err := d.reader.GetUint16()
		if err != nil {
			return nil, err
		}
		m := make(map[any]any, v)
		err = d.readMap(m, uint32(v))
		return m, err
	case FormatMap32:
		v, err := d.reader.GetUint32()
		if err != nil {
			return nil, err
		}
		m := make(map[any]any, v)
		err = d.readMap(m, v)
		return m, err
	case FormatBin8:
		binLen, err := d.reader.GetUint8()
		if err != nil {
			return nil, err
		}
		return d.reader.GetBytes(uint32(binLen))
	case FormatBin16:
		binLen, err := d.reader.GetUint16()
		if err != nil {
			return nil, err
		}
		return d.reader.GetBytes(uint32(binLen))
	case FormatBin32:
		binLen, err := d.reader.GetUint32()
		if err != nil {
			return nil, err
		}
		return d.reader.GetBytes(binLen)
	}

	return nil, ReadError{"bad value for bool"}
}

// readArray populates an array with values read from the MessagePack stream.
// Reads the specified number of elements using ReadAny.
func (d *Decoder) readArray(array []any) error {
	for i := 0; i < len(array); i++ {
		value, err := d.ReadAny()
		if err != nil {
			return err
		}
		array[i] = value
	}
	return nil
}

// readMap populates a map with key-value pairs read from the MessagePack stream.
// Reads the specified number of pairs using ReadAny for both keys and values.
func (d *Decoder) readMap(m map[any]any, length uint32) error {
	for i := uint32(0); i < length; i++ {
		key, err := d.ReadAny()
		if err != nil {
			return err
		}
		value, err := d.ReadAny()
		if err != nil {
			return err
		}
		m[key] = value
	}
	return nil
}

// extHeader reads the extension header for MessagePack extension types.
// Returns the extension ID and length from the extension header.
func (d *Decoder) extHeader(c byte) (int8, uint32, error) {
	extLen, err := d.parseExtLen(c)
	if err != nil {
		return 0, 0, err
	}

	extID, err := d.readCode()
	if err != nil {
		return 0, 0, err
	}

	return int8(extID), extLen, nil
}

// readCode reads a single byte from the MessagePack stream.
// Used for reading extension type IDs.
func (d *Decoder) readCode() (byte, error) {
	c, err := d.reader.GetUint8()
	if err != nil {
		return 0, err
	}
	return c, nil
}

// parseExtLen determines the length of an extension based on the format byte.
// Supports all MessagePack extension length formats.
func (d *Decoder) parseExtLen(c byte) (uint32, error) {
	switch c {
	case FormatFixExt1:
		return 1, nil
	case FormatFixExt2:
		return 2, nil
	case FormatFixExt4:
		return 4, nil
	case FormatFixExt8:
		return 8, nil
	case FormatFixExt16:
		return 16, nil
	case FormatExt8:
		n, err := d.ReadUint8()
		return uint32(n), err
	case FormatExt16:
		n, err := d.ReadUint16()
		return uint32(n), err
	case FormatExt32:
		n, err := d.ReadUint32()
		return n, err
	default:
		return 0, ReadError{"msgpack: invalid code=" + strconv.FormatUint(uint64(c), 16) + " decoding ext len"}
	}
}

// Err returns any error that occurred during decoding.
// Returns the underlying error from the data reader.
func (d *Decoder) Err() error {
	return d.reader.Err()
}

// Decode is a generic function that decodes a MessagePack value into a type T.
// The type T must implement the Decodable interface via its pointer receiver.
func Decode[T any, PT interface {
	*T
	Decodable
}](decoder Reader) (T, error) {
	var inst T
	err := ((PT)(&inst)).Decode(decoder)
	return inst, err
}

// DecodeNillable is a generic function that decodes a MessagePack value into a type T or returns nil.
// The type T must implement the Decodable interface via its pointer receiver.
// Returns nil if the next value is MessagePack nil.
func DecodeNillable[T any, PT interface {
	*T
	Decodable
}](decoder Reader) (PT, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}

	codec := PT(new(T))
	err := codec.Decode(decoder)
	return codec, err
}

////////////////////

// isFixedInt checks if a byte represents a MessagePack fixed positive integer.
// Fixed positive integers have the most significant bit set to 0.
//
//go:inline
func isFixedInt(u byte) bool {
	return u>>7 == 0
}

// isNegativeFixedInt checks if a byte represents a MessagePack fixed negative integer.
// Fixed negative integers have the 3 most significant bits set to 111.
//
//go:inline
func isNegativeFixedInt(u byte) bool {
	return (u & 0xe0) == FormatNegativeFixInt
}

// isFixedMap checks if a byte represents a MessagePack fixed map.
// Fixed maps have the 4 most significant bits set to 1000.
//
//go:inline
func isFixedMap(u byte) bool {
	return (u & 0xf0) == FormatFixMap
}

// isFixedArray checks if a byte represents a MessagePack fixed array.
// Fixed arrays have the 4 most significant bits set to 1001.
//
//go:inline
func isFixedArray(u byte) bool {
	return (u & 0xf0) == FormatFixArray
}

// isFixedString checks if a byte represents a MessagePack fixed string.
// Fixed strings have the 3 most significant bits set to 101.
//
//go:inline
func isFixedString(u byte) bool {
	return (u & 0xe0) == FormatFixString
}

// isString checks if a byte represents any MessagePack string format.
// Includes fixed strings, variable length strings, and arrays (for compatibility).
//
//go:inline
func isString(u byte) bool {
	return isFixedString(u) ||
		u == FormatString8 ||
		u == FormatString16 ||
		u == FormatString32 ||
		isFixedArray(u) ||
		u == FormatArray16 ||
		u == FormatArray32
}

// ReadError represents an error that occurred during MessagePack decoding.
type ReadError struct {
	message string
}

// Error returns the error message for a ReadError.
func (e ReadError) Error() string {
	return e.message
}
