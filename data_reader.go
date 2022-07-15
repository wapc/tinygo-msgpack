package msgpack

import (
	"encoding/binary"
	"errors"
	"math"
)

var ErrRange = errors.New("range error")

type DataReader struct {
	buffer     []byte
	byteOffset uint32
	err        error
}

func NewDataReader(buffer []byte) DataReader {
	return DataReader{
		buffer: buffer,
	}
}

func (d *DataReader) GetBytes(length uint32) ([]byte, error) {
	if err := d.checkBufferSize(length); err != nil {
		return nil, err
	}
	result := d.buffer[d.byteOffset : d.byteOffset+length]
	d.byteOffset += length
	return result, nil
}

func (d *DataReader) SetBytes(src []byte) error {
	srcLen := uint32(len(src))
	if err := d.checkBufferSize(srcLen); err != nil {
		return err
	}
	copy(d.buffer[d.byteOffset:], src)
	d.byteOffset += srcLen
	return nil
}

func (d *DataReader) PeekUint8() (uint8, error) {
	if err := d.checkBufferSize(1); err != nil {
		return 0, err
	}
	return d.buffer[d.byteOffset], nil
}

func (d *DataReader) Discard(length uint32) error {
	if err := d.checkBufferSize(length); err != nil {
		return err
	}
	d.byteOffset += length
	return nil
}

func (d *DataReader) GetFloat32() (float32, error) {
	if err := d.checkBufferSize(4); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(d.buffer[d.byteOffset:])
	d.byteOffset += 4
	return math.Float32frombits(v), nil
}

func (d *DataReader) GetFloat64() (float64, error) {
	if err := d.checkBufferSize(8); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint64(d.buffer[d.byteOffset:])
	d.byteOffset += 8
	return math.Float64frombits(v), nil
}

func (d *DataReader) GetInt8() (int8, error) {
	if err := d.checkBufferSize(1); err != nil {
		return 0, err
	}
	result := d.buffer[d.byteOffset]
	d.byteOffset++
	return int8(result), nil
}

func (d *DataReader) GetInt16() (int16, error) {
	if err := d.checkBufferSize(2); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint16(d.buffer[d.byteOffset:])
	d.byteOffset += 2
	return int16(result), nil
}

func (d *DataReader) GetInt32() (int32, error) {
	if err := d.checkBufferSize(4); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint32(d.buffer[d.byteOffset:])
	d.byteOffset += 4
	return int32(result), nil
}

func (d *DataReader) GetInt64() (int64, error) {
	if err := d.checkBufferSize(8); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint64(d.buffer[d.byteOffset:])
	d.byteOffset += 8
	return int64(result), nil
}

func (d *DataReader) GetUint8() (uint8, error) {
	if err := d.checkBufferSize(1); err != nil {
		return 0, err
	}
	result := d.buffer[d.byteOffset]
	d.byteOffset++
	return result, nil
}

func (d *DataReader) GetUint16() (uint16, error) {
	if err := d.checkBufferSize(2); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint16(d.buffer[d.byteOffset:])
	d.byteOffset += 2
	return result, nil
}

func (d *DataReader) GetUint32() (uint32, error) {
	if err := d.checkBufferSize(4); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint32(d.buffer[d.byteOffset:])
	d.byteOffset += 4
	return result, nil
}

func (d *DataReader) GetUint64() (uint64, error) {
	if err := d.checkBufferSize(8); err != nil {
		return 0, err
	}
	result := binary.BigEndian.Uint64(d.buffer[d.byteOffset:])
	d.byteOffset += 8
	return result, nil
}

func (d *DataReader) SetFloat32(value float32) error {
	if err := d.checkBufferSize(4); err != nil {
		return err
	}
	bits := math.Float32bits(value)
	binary.BigEndian.PutUint32(d.buffer[d.byteOffset:], bits)
	d.byteOffset += 4
	return nil
}

func (d *DataReader) SetFloat64(value float64) error {
	if err := d.checkBufferSize(8); err != nil {
		return err
	}
	bits := math.Float64bits(value)
	binary.BigEndian.PutUint64(d.buffer[d.byteOffset:], bits)
	d.byteOffset += 8
	return nil
}

func (d *DataReader) SetInt8(value int8) error {
	if err := d.checkBufferSize(1); err != nil {
		return err
	}
	d.buffer[d.byteOffset] = uint8(value)
	d.byteOffset++
	return nil
}

func (d *DataReader) SetInt16(value int16) error {
	if err := d.checkBufferSize(2); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(d.buffer[d.byteOffset:], uint16(value))
	d.byteOffset += 2
	return nil
}

func (d *DataReader) SetInt32(value int32) error {
	if err := d.checkBufferSize(4); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(d.buffer[d.byteOffset:], uint32(value))
	d.byteOffset += 4
	return nil
}

func (d *DataReader) SetInt64(value int64) error {
	if err := d.checkBufferSize(8); err != nil {
		return err
	}
	binary.BigEndian.PutUint64(d.buffer[d.byteOffset:], uint64(value))
	d.byteOffset += 8
	return nil
}

func (d *DataReader) SetUint8(value uint8) error {
	if err := d.checkBufferSize(1); err != nil {
		return err
	}
	d.buffer[d.byteOffset] = value
	d.byteOffset++
	return nil
}

func (d *DataReader) SetUint16(value uint16) error {
	if err := d.checkBufferSize(2); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(d.buffer[d.byteOffset:], value)
	d.byteOffset += 2
	return nil
}

func (d *DataReader) SetUint32(value uint32) error {
	if err := d.checkBufferSize(4); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(d.buffer[d.byteOffset:], value)
	d.byteOffset += 4
	return nil
}

func (d *DataReader) SetUint64(value uint64) error {
	if err := d.checkBufferSize(8); err != nil {
		return err
	}
	binary.BigEndian.PutUint64(d.buffer[d.byteOffset:], value)
	d.byteOffset += 8
	return nil
}

func (d *DataReader) checkBufferSize(length uint32) error {
	if d.err != nil {
		return d.err
	}
	if d.byteOffset+length > uint32(len(d.buffer)) {
		d.err = ErrRange
		return ErrRange
	}

	return nil
}

func (d *DataReader) Err() error {
	return d.err
}
