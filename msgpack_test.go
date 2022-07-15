package msgpack_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	msgpack "github.com/wapc/tinygo-msgpack"
)

type Required struct {
	BoolValue   bool
	U8Value     uint8
	U16Value    uint16
	U32Value    uint32
	U64Value    uint64
	S8Value     int8
	S16Value    int16
	S32Value    int32
	S64Value    int64
	F32Value    float32
	F64Value    float64
	StringValue string
	BytesValue  []byte
	ArrayValue  []int64
	MapValue    map[string]int64
}

func DecodeRequiredNullable(decoder msgpack.Reader) (*Required, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeRequired(decoder)
	return &decoded, err
}

func DecodeRequired(decoder msgpack.Reader) (Required, error) {
	var o Required
	err := o.Decode(decoder)
	return o, err
}

func (o *Required) Decode(decoder msgpack.Reader) error {
	numFields, err := decoder.ReadMapSize()
	if err != nil {
		return err
	}

	for numFields > 0 {
		numFields--
		field, err := decoder.ReadString()
		if err != nil {
			return err
		}
		switch field {
		case "boolValue":
			o.BoolValue, err = decoder.ReadBool()
		case "u8Value":
			o.U8Value, err = decoder.ReadUint8()
		case "u16Value":
			o.U16Value, err = decoder.ReadUint16()
		case "u32Value":
			o.U32Value, err = decoder.ReadUint32()
		case "u64Value":
			o.U64Value, err = decoder.ReadUint64()
		case "s8Value":
			o.S8Value, err = decoder.ReadInt8()
		case "s16Value":
			o.S16Value, err = decoder.ReadInt16()
		case "s32Value":
			o.S32Value, err = decoder.ReadInt32()
		case "s64Value":
			o.S64Value, err = decoder.ReadInt64()
		case "f32Value":
			o.F32Value, err = decoder.ReadFloat32()
		case "f64Value":
			o.F64Value, err = decoder.ReadFloat64()
		case "stringValue":
			o.StringValue, err = decoder.ReadString()
		case "bytesValue":
			o.BytesValue, err = decoder.ReadByteArray()
		case "arrayValue":
			var isNil bool
			isNil, err = decoder.IsNextNil()
			if err == nil {
				if isNil {
					o.ArrayValue = nil
				} else {
					size, err := decoder.ReadArraySize()
					if err == nil {
						o.ArrayValue = make([]int64, size)
						i := 0
						for ; size > 0; i++ {
							size--
							o.ArrayValue[i], err = decoder.ReadInt64()
							if err != nil {
								break
							}
						}
					}
				}
			}
		case "mapValue":
			var isNil bool
			isNil, err = decoder.IsNextNil()
			if err == nil {
				if isNil {
					o.MapValue = nil
				} else {
					size, err := decoder.ReadMapSize()
					if err == nil {
						o.MapValue = make(map[string]int64, size)
						for size > 0 {
							size--
							var key string
							if key, err = decoder.ReadString(); err != nil {
								break
							}
							if o.MapValue[key], err = decoder.ReadInt64(); err != nil {
								break
							}
						}
					}
				}
			}
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Required) Encode(encoder msgpack.Writer) {
	if o == nil {
		encoder.WriteNil()
		return
	}
	encoder.WriteMapSize(16)

	// Add some nested data that must be skipped.
	// This tests skipping over unknown fields.
	encoder.WriteString("nested")
	encoder.WriteMapSize(2)
	encoder.WriteString("foo")
	encoder.WriteString("bar")
	encoder.WriteString("nested")
	encoder.WriteMapSize(3)
	encoder.WriteString("foo")
	encoder.WriteString("bar")
	encoder.WriteString("other")
	encoder.WriteString("value")
	encoder.WriteString("array")
	encoder.WriteArraySize(3)
	encoder.WriteInt64(1)
	encoder.WriteInt64(2)
	encoder.WriteInt64(3)

	encoder.WriteString("boolValue")
	encoder.WriteBool(o.BoolValue)
	encoder.WriteString("u8Value")
	encoder.WriteUint8(o.U8Value)
	encoder.WriteString("u16Value")
	encoder.WriteUint16(o.U16Value)
	encoder.WriteString("u32Value")
	encoder.WriteUint32(o.U32Value)
	encoder.WriteString("u64Value")
	encoder.WriteUint64(o.U64Value)
	encoder.WriteString("s8Value")
	encoder.WriteInt8(o.S8Value)
	encoder.WriteString("s16Value")
	encoder.WriteInt16(o.S16Value)
	encoder.WriteString("s32Value")
	encoder.WriteInt32(o.S32Value)
	encoder.WriteString("s64Value")
	encoder.WriteInt64(o.S64Value)
	encoder.WriteString("f32Value")
	encoder.WriteFloat32(o.F32Value)
	encoder.WriteString("f64Value")
	encoder.WriteFloat64(o.F64Value)
	encoder.WriteString("stringValue")
	encoder.WriteString(o.StringValue)
	encoder.WriteString("bytesValue")
	if o.BytesValue == nil {
		encoder.WriteNil()
	} else {
		encoder.WriteByteArray(o.BytesValue)
	}
	encoder.WriteString("arrayValue")
	if o.ArrayValue == nil {
		encoder.WriteNil()
	} else {
		encoder.WriteArraySize(uint32(len(o.ArrayValue)))
		for _, item := range o.ArrayValue {
			encoder.WriteInt64(item)
		}
	}
	encoder.WriteString("mapValue")
	if o.MapValue == nil {
		encoder.WriteNil()
	} else {
		encoder.WriteMapSize(uint32(len(o.MapValue)))
		for key, value := range o.MapValue {
			encoder.WriteString(key)
			encoder.WriteInt64(value)
		}
	}
}

func (o *Required) ToBuffer() []byte {
	var sizer msgpack.Sizer
	o.Encode(&sizer)
	buffer := make([]byte, sizer.Len())
	encoder := msgpack.NewEncoder(buffer)
	o.Encode(&encoder)
	return buffer
}

func TestEcho(t *testing.T) {
	expected := Required{
		BoolValue:   true,
		U8Value:     math.MaxUint8,
		U16Value:    math.MaxUint16,
		U32Value:    math.MaxUint32,
		U64Value:    math.MaxUint64,
		S8Value:     math.MinInt8,
		S16Value:    math.MinInt16,
		S32Value:    math.MinInt32,
		S64Value:    math.MinInt64,
		F32Value:    math.MaxFloat32,
		F64Value:    math.MaxFloat64,
		StringValue: "test",
		BytesValue:  []byte("test"),
		ArrayValue:  []int64{1, 2, 3, 4},
		MapValue: map[string]int64{
			"key": 1234,
		},
	}

	data := expected.ToBuffer()
	var actual Required
	decoder := msgpack.NewDecoder(data)
	err := actual.Decode(&decoder)
	require.NoError(t, err)
	assert.Equal(t, expected, actual, "mismatch with required fields")
}

func TestUint8Range(t *testing.T) {
	values := []uint8{}
	for i := int8(0); i <= 7; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxUint8)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteUint8(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteUint8(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadUint8()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch uint8 value")
	}
}

func TestUint16Range(t *testing.T) {
	values := []uint16{}
	for i := int16(0); i <= 15; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxUint16)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteUint16(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteUint16(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadUint16()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch uint16 value")
	}
}

func TestUint32Range(t *testing.T) {
	values := []uint32{}
	for i := int32(0); i <= 31; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxUint32)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteUint32(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteUint32(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadUint32()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch uint32 value")
	}
}

func TestUint64Range(t *testing.T) {
	values := []uint64{}
	for i := int32(0); i <= 31; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxUint64)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteUint64(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteUint64(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadUint64()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch uint64 value")
	}
}

func TestInt8Range(t *testing.T) {
	values := []int8{}
	values = append(values, math.MinInt8)
	for i := int8(6); i >= 0; i-- {
		values = append(values, -((1 << i) + 1))
		values = append(values, -(1 << i))
		values = append(values, -((1 << i) - 1))
	}
	for i := int8(0); i <= 6; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxInt8)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteInt8(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteInt8(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadInt8()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch int8 value")
	}
}

func TestInt16Range(t *testing.T) {
	values := []int16{}
	values = append(values, math.MinInt16)
	for i := int16(14); i >= 0; i-- {
		values = append(values, -((1 << i) + 1))
		values = append(values, -(1 << i))
		values = append(values, -((1 << i) - 1))
	}
	for i := int16(0); i <= 14; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxInt16)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteInt16(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteInt16(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadInt16()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch int16 value")
	}
}

func TestInt32Range(t *testing.T) {
	values := []int32{}
	values = append(values, math.MinInt32)
	for i := int32(30); i >= 0; i-- {
		values = append(values, -((1 << i) + 1))
		values = append(values, -(1 << i))
		values = append(values, -((1 << i) - 1))
	}
	for i := int32(0); i <= 30; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxInt32)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteInt32(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteInt32(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadInt32()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch int32 value")
	}
}

func TestInt64Range(t *testing.T) {
	values := []int64{}
	values = append(values, math.MinInt64)
	for i := int64(62); i >= 0; i-- {
		values = append(values, -((1 << i) + 1))
		values = append(values, -(1 << i))
		values = append(values, -((1 << i) - 1))
	}
	for i := int32(0); i <= 62; i++ {
		values = append(values, (1<<i)-1)
		values = append(values, (1 << i))
		values = append(values, (1<<i)+1)
	}
	values = append(values, math.MaxInt64)
	for _, expected := range values {
		var sizer msgpack.Sizer
		sizer.WriteInt64(expected)
		buffer := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buffer)
		encoder.WriteInt64(expected)
		decoder := msgpack.NewDecoder(buffer)
		actual, err := decoder.ReadInt64()
		require.NoError(t, err)
		assert.Equal(t, expected, actual, "mismatch int64 value")
	}
}
