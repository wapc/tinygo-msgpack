# TinyGo MessagePack

A high-performance MessagePack implementation for TinyGo and Go, designed for WebAssembly (WASM) and embedded environments.

## Features

- **TinyGo Compatible**: Optimized for TinyGo compilation to WebAssembly
- **Zero Allocations**: Efficient memory usage with minimal allocations
- **Type Safety**: Strongly typed encoding/decoding with compile-time safety
- **Comprehensive Support**: Full MessagePack specification support including:
  - All primitive types (int8-64, uint8-64, float32/64, bool, string)
  - Arrays and maps
  - Binary data
  - Time values
  - Nil values
  - Extension types
- **High-level API**: Convenient Marshal/Unmarshal functions
- **Size Calculation**: Built-in sizer for pre-allocating buffers
- **Nullable Types**: Support for optional/nullable values

## Installation

```bash
go get github.com/wapc/tinygo-msgpack
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/wapc/tinygo-msgpack"
)

type Person struct {
    Name string
    Age  int32
    Tags []string
}

func (p *Person) Encode(encoder msgpack.Writer) {
    encoder.WriteMapSize(3)
    encoder.WriteString("name")
    encoder.WriteString(p.Name)
    encoder.WriteString("age")
    encoder.WriteInt32(p.Age)
    encoder.WriteString("tags")
    if p.Tags == nil {
        encoder.WriteNil()
    } else {
        encoder.WriteArraySize(uint32(len(p.Tags)))
        for _, tag := range p.Tags {
            encoder.WriteString(tag)
        }
    }
}

func (p *Person) Decode(decoder msgpack.Reader) error {
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
        case "name":
            p.Name, err = decoder.ReadString()
        case "age":
            p.Age, err = decoder.ReadInt32()
        case "tags":
            isNil, err := decoder.IsNextNil()
            if err != nil {
                return err
            }
            if isNil {
                p.Tags = nil
            } else {
                size, err := decoder.ReadArraySize()
                if err != nil {
                    return err
                }
                p.Tags = make([]string, size)
                for i := uint32(0); i < size; i++ {
                    p.Tags[i], err = decoder.ReadString()
                    if err != nil {
                        return err
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

func main() {
    person := &Person{
        Name: "Alice",
        Age:  30,
        Tags: []string{"developer", "golang"},
    }

    // Encode to bytes
    data, err := msgpack.ToBytes(person)
    if err != nil {
        log.Fatal(err)
    }

    // Decode from bytes
    var decoded Person
    err = msgpack.Unmarshal(data, &decoded)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Decoded: %+v\n", decoded)
}
```

### High-level API

```go
// Marshal any value to MessagePack
data, err := msgpack.Marshal(map[string]interface{}{
    "name": "Bob",
    "age":  25,
    "active": true,
    "scores": []int{95, 87, 92},
})

// Unmarshal to any value
var result map[string]interface{}
err = msgpack.Unmarshal(data, &result)
```

### Streaming API

```go
// Encoding
buffer := make([]byte, 1024)
encoder := msgpack.NewEncoder(buffer)

encoder.WriteMapSize(2)
encoder.WriteString("key1")
encoder.WriteString("value1")
encoder.WriteString("key2")
encoder.WriteInt64(42)

// Decoding
decoder := msgpack.NewDecoder(buffer)
mapSize, err := decoder.ReadMapSize()
if err != nil {
    log.Fatal(err)
}

for i := uint32(0); i < mapSize; i++ {
    key, err := decoder.ReadString()
    if err != nil {
        log.Fatal(err)
    }
    
    // Read value based on key
    switch key {
    case "key1":
        value, err := decoder.ReadString()
    case "key2":
        value, err := decoder.ReadInt64()
    }
}
```

### Size Calculation

```go
// Calculate size before encoding
sizer := msgpack.NewSizer()
person.Encode(&sizer)
bufferSize := sizer.Len()

// Pre-allocate buffer
buffer := make([]byte, bufferSize)
encoder := msgpack.NewEncoder(buffer)
person.Encode(&encoder)
```

## API Reference

### Core Interfaces

#### `Reader`
Interface for reading MessagePack data:
```go
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
    ReadTime() (time.Time, error)
    ReadByteArray() ([]byte, error)
    ReadArraySize() (uint32, error)
    ReadMapSize() (uint32, error)
    ReadAny() (any, error)
    ReadRaw() (Raw, error)
    Skip() error
    Err() error
}
```

#### `Writer`
Interface for writing MessagePack data:
```go
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
    WriteTime(value time.Time)
    WriteByteArray(value []byte)
    WriteArraySize(length uint32)
    WriteMapSize(length uint32)
    WriteAny(value any)
    WriteRaw(value Raw)
    Err() error
}
```

### Main Types

#### `Encoder`
Low-level encoder for writing MessagePack data:
```go
func NewEncoder(buffer []byte) Encoder
```

#### `Decoder`
Low-level decoder for reading MessagePack data:
```go
func NewDecoder(buffer []byte) Decoder
```

#### `Sizer`
Calculates the size of encoded data:
```go
func NewSizer() Sizer
func (s *Sizer) Len() uint32
```

### High-level Functions

#### Encoding
```go
func Marshal(value any) ([]byte, error)
func ToBytes(value Encodable) ([]byte, error)
func AnyToBytes(value any) ([]byte, error)
```

#### Decoding
```go
func Unmarshal(data []byte, value Decodable) error
func BytesToAny(data []byte) (any, error)
```

#### Type-specific Encoding
```go
func I8ToBytes(value int8) ([]byte, error)
func I16ToBytes(value int16) ([]byte, error)
func I32ToBytes(value int32) ([]byte, error)
func I64ToBytes(value int64) ([]byte, error)
func U8ToBytes(value uint8) ([]byte, error)
func U16ToBytes(value uint16) ([]byte, error)
func U32ToBytes(value uint32) ([]byte, error)
func U64ToBytes(value uint64) ([]byte, error)
func F32ToBytes(value float32) ([]byte, error)
func F64ToBytes(value float64) ([]byte, error)
func BoolToBytes(value bool) ([]byte, error)
func StringToBytes(value string) ([]byte, error)
func BytesToBytes(value []byte) ([]byte, error)
func TimeToBytes(value time.Time) ([]byte, error)
```

### Nullable Types

All primitive types have nullable variants:
```go
// Reading nullable values
value, err := decoder.ReadNillableBool()
value, err := decoder.ReadNillableInt8()
value, err := decoder.ReadNillableString()
// ... etc

// Writing nullable values
encoder.WriteNillableBool(&value)
encoder.WriteNillableInt8(&value)
encoder.WriteNillableString(&value)
// ... etc
```

## TinyGo Compatibility

This library is specifically designed for TinyGo and WebAssembly environments:

- No reflection usage
- Minimal memory allocations
- Compatible with TinyGo's limited standard library
- Optimized for WASM performance

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
