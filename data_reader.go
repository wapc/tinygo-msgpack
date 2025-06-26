package msgpack

import (
	"encoding/binary"
	"errors"
	"math"
)

// ErrRange is returned when attempting to read beyond the buffer boundaries.
var ErrRange = errors.New("range error")

// DataReader provides low-level buffer reading and writing operations.
// It manages a byte buffer with an offset for sequential access.
type DataReader struct {
	buffer     []byte
	byteOffset uint32
	err        error
}

// NewDataReader creates a new DataReader with the provided buffer.
func NewDataReader(buffer []byte) DataReader {
	return DataReader{
		buffer: buffer,
	}
}

// GetBytes reads a specified number of bytes from the buffer.
// Advances the read position by the number of bytes read.
func (d *DataReader) GetBytes(length uint32) ([]byte, error) {
	if err := d.checkBufferSize(length); err != nil {
		return nil, err
	}
	result := d.buffer[d.byteOffset : d.byteOffset+length]
	d.byteOffset += length
	return result, nil
}

// SetBytes writes bytes to the buffer at the current position.
// Advances the write position by the number of bytes written.
func (d *DataReader) SetBytes(src []byte) error {
	srcLen := uint32(len(src))
	if err := d.checkBufferSize(srcLen); err != nil {
		return err
	}
	copy(d.buffer[d.byteOffset:], src)
	d.byteOffset += srcLen
	return nil
}

// PeekUint8 reads a uint8 value without advancing the read position.
func (d *DataReader) PeekUint8() (uint8, error) {
	if err := d.checkBufferSize(1); err != nil {
		return 0, err
	}
	return d.buffer[d.byteOffset], nil
}

// Discard advances the read position by the specified number of bytes.
func (d *DataReader) Discard(length uint32) error {
	if err := d.checkBufferSize(length); err != nil {
		return err
	}
	d.byteOffset += length
	return nil
}

// GetFloat32 reads a float32 value from the buffer in big-endian format.
func (d *DataReader) GetFloat32() (float32, error) {
	if err := d.checkBufferSize(4); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(d.buffer[d.byteOffset:])
	d.byteOffset += 4
	return math.Float32frombits(v), nil
}

// GetFloat64 reads a float64 value from the buffer in big-endian format.
func (d *DataReader) GetFloat64() (float64, error) {
	if err := d.checkBufferSize(8); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint64(d.buffer[d.byteOffset:])
	d.byteOffset += 8
	return math.Float64frombits(v), nil
}

// GetInt8 reads an int8 value from the buffer.
func (d *DataReader) GetInt8() (int8, error) {
	if err := d.checkBufferSize(1); err != nil {
		return 0, err
	}
	result := d.buffer[d.byteOffset]
	d.byteOffset++
	return int8(result), nil
}

// GetInt16 reads an int16 value from the buffer in big-endian format.
func (d *DataReader) GetInt16() (int16, error) {
	if err := d.checkBufferSize(2); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint16(d.buffer[d.byteOffset:])
	d.byteOffset += 2
	return int16(result), nil
}

// GetInt32 reads an int32 value from the buffer in big-endian format.
func (d *DataReader) GetInt32() (int32, error) {
	if err := d.checkBufferSize(4); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint32(d.buffer[d.byteOffset:])
	d.byteOffset += 4
	return int32(result), nil
}

// GetInt64 reads an int64 value from the buffer in big-endian format.
func (d *DataReader) GetInt64() (int64, error) {
	if err := d.checkBufferSize(8); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint64(d.buffer[d.byteOffset:])
	d.byteOffset += 8
	return int64(result), nil
}

// GetUint8 reads a uint8 value from the buffer.
func (d *DataReader) GetUint8() (uint8, error) {
	if err := d.checkBufferSize(1); err != nil {
		return 0, err
	}
	result := d.buffer[d.byteOffset]
	d.byteOffset++
	return result, nil
}

// GetUint16 reads a uint16 value from the buffer in big-endian format.
func (d *DataReader) GetUint16() (uint16, error) {
	if err := d.checkBufferSize(2); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint16(d.buffer[d.byteOffset:])
	d.byteOffset += 2
	return result, nil
}

// GetUint32 reads a uint32 value from the buffer in big-endian format.
func (d *DataReader) GetUint32() (uint32, error) {
	if err := d.checkBufferSize(4); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint32(d.buffer[d.byteOffset:])
	d.byteOffset += 4
	return result, nil
}

// GetUint64 reads a uint64 value from the buffer in big-endian format.
func (d *DataReader) GetUint64() (uint64, error) {
	if err := d.checkBufferSize(8); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint64(d.buffer[d.byteOffset:])
	d.byteOffset += 8
	return result, nil
}

// SetFloat32 writes a float32 value to the buffer in big-endian format.
func (d *DataReader) SetFloat32(value float32) error {
	if err := d.checkBufferSize(4); err != nil {
		return err
	}
	bits := math.Float32bits(value)
	binary.BigEndian.PutUint32(d.buffer[d.byteOffset:], bits)
	d.byteOffset += 4
	return nil
}

// SetFloat64 writes a float64 value to the buffer in big-endian format.
func (d *DataReader) SetFloat64(value float64) error {
	if err := d.checkBufferSize(8); err != nil {
		return err
	}
	bits := math.Float64bits(value)
	binary.BigEndian.PutUint64(d.buffer[d.byteOffset:], bits)
	d.byteOffset += 8
	return nil
}

// SetInt8 writes an int8 value to the buffer.
func (d *DataReader) SetInt8(value int8) error {
	if err := d.checkBufferSize(1); err != nil {
		return err
	}
	d.buffer[d.byteOffset] = uint8(value)
	d.byteOffset++
	return nil
}

// SetInt16 writes an int16 value to the buffer in big-endian format.
func (d *DataReader) SetInt16(value int16) error {
	if err := d.checkBufferSize(2); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(d.buffer[d.byteOffset:], uint16(value))
	d.byteOffset += 2
	return nil
}

// SetInt32 writes an int32 value to the buffer in big-endian format.
func (d *DataReader) SetInt32(value int32) error {
	if err := d.checkBufferSize(4); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(d.buffer[d.byteOffset:], uint32(value))
	d.byteOffset += 4
	return nil
}

// SetInt64 writes an int64 value to the buffer in big-endian format.
func (d *DataReader) SetInt64(value int64) error {
	if err := d.checkBufferSize(8); err != nil {
		return err
	}
	binary.BigEndian.PutUint64(d.buffer[d.byteOffset:], uint64(value))
	d.byteOffset += 8
	return nil
}

// SetUint8 writes a uint8 value to the buffer.
func (d *DataReader) SetUint8(value uint8) error {
	if err := d.checkBufferSize(1); err != nil {
		return err
	}
	d.buffer[d.byteOffset] = value
	d.byteOffset++
	return nil
}

// SetUint16 writes a uint16 value to the buffer in big-endian format.
func (d *DataReader) SetUint16(value uint16) error {
	if err := d.checkBufferSize(2); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(d.buffer[d.byteOffset:], value)
	d.byteOffset += 2
	return nil
}

// SetUint32 writes a uint32 value to the buffer in big-endian format.
func (d *DataReader) SetUint32(value uint32) error {
	if err := d.checkBufferSize(4); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(d.buffer[d.byteOffset:], value)
	d.byteOffset += 4
	return nil
}

// SetUint64 writes a uint64 value to the buffer in big-endian format.
func (d *DataReader) SetUint64(value uint64) error {
	if err := d.checkBufferSize(8); err != nil {
		return err
	}
	binary.BigEndian.PutUint64(d.buffer[d.byteOffset:], value)
	d.byteOffset += 8
	return nil
}

// checkBufferSize verifies that the requested number of bytes is available.
// Returns ErrRange if the buffer is too small.
func (d *DataReader) checkBufferSize(length uint32) error {
	if d.err != nil {
		return d.err
	}
	if d.byteOffset+length > uint32(len(d.buffer)) {
		d.err = ErrRange
		return d.err
	}

	return nil
}

// Err returns any error that occurred during buffer operations.
func (d *DataReader) Err() error {
	return d.err
}
