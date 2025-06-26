// Package msgpack provides a high-performance MessagePack implementation
// optimized for TinyGo and WebAssembly environments.
//
// This library offers both high-level convenience functions and low-level
// streaming APIs for maximum control over MessagePack encoding and decoding.
package msgpack

// The following examples demonstrate common usage patterns.

// ExampleMarshal demonstrates basic MessagePack marshaling.
func ExampleMarshal() {
	// Marshal any value to MessagePack
	data, err := Marshal(map[string]interface{}{
		"name":   "Bob",
		"age":    25,
		"active": true,
		"scores": []int{95, 87, 92},
	})
	if err != nil {
		// Handle error
	}
	_ = data // Use the marshaled data
}

// ExampleUnmarshal demonstrates basic MessagePack unmarshaling.
func ExampleUnmarshal() {
	data := []byte{0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa4, 0x42, 0x6f, 0x62, 0x00}
	result, err := BytesToAny(data)
	if err != nil {
		// Handle error
	}
	_ = result // Use the unmarshaled data
}

// ExampleEncoder demonstrates low-level encoding.
func ExampleEncoder() {
	// Pre-allocate buffer
	buffer := make([]byte, 1024)
	encoder := NewEncoder(buffer)

	// Write MessagePack data
	encoder.WriteMapSize(2)
	encoder.WriteString("key1")
	encoder.WriteString("value1")
	encoder.WriteString("key2")
	encoder.WriteInt64(42)

	// Check for errors
	if err := encoder.Err(); err != nil {
		// Handle error
	}
}

// ExampleDecoder demonstrates low-level decoding.
func ExampleDecoder() {
	data := []byte{0x82, 0xa4, 0x6b, 0x65, 0x79, 0x31, 0xa6, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31}
	decoder := NewDecoder(data)

	// Read MessagePack data
	mapSize, err := decoder.ReadMapSize()
	if err != nil {
		// Handle error
	}

	for i := uint32(0); i < mapSize; i++ {
		key, err := decoder.ReadString()
		if err != nil {
			// Handle error
		}

		// Read value based on key
		switch key {
		case "key1":
			value, err := decoder.ReadString()
			if err != nil {
				// Handle error
			}
			_ = value
		}
	}
}

// ExampleSizer demonstrates size calculation.
func ExampleSizer() {
	// Calculate size before encoding
	sizer := NewSizer()
	sizer.WriteMapSize(2)
	sizer.WriteString("name")
	sizer.WriteString("Alice")
	sizer.WriteString("age")
	sizer.WriteInt32(30)

	bufferSize := sizer.Len()

	// Pre-allocate buffer
	buffer := make([]byte, bufferSize)
	encoder := NewEncoder(buffer)

	// Encode to the pre-allocated buffer
	encoder.WriteMapSize(2)
	encoder.WriteString("name")
	encoder.WriteString("Alice")
	encoder.WriteString("age")
	encoder.WriteInt32(30)
}

// ExampleNullable demonstrates nullable type usage.
func ExampleNullable() {
	// Reading nullable values
	decoder := NewDecoder([]byte{0xc0}) // nil value

	isNil, err := decoder.IsNextNil()
	if err != nil {
		// Handle error
	}
	if isNil {
		// Value is nil
	} else {
		value, err := decoder.ReadString()
		if err != nil {
			// Handle error
		}
		_ = value
	}

	// Writing nullable values
	buffer := make([]byte, 10)
	encoder := NewEncoder(buffer)

	var str *string = nil
	encoder.WriteNillableString(str) // Writes nil

	str2 := "hello"
	encoder.WriteNillableString(&str2) // Writes the string
}

// ExampleCustomType demonstrates custom type encoding/decoding.
func ExampleCustomType() {
	// Define a custom type
	type Person struct {
		Name string
		Age  int32
	}

	// Implement encoding
	person := &Person{Name: "Alice", Age: 30}
	data, err := AnyToBytes(person)
	if err != nil {
		// Handle error
	}

	// Implement decoding
	result, err := BytesToAny(data)
	if err != nil {
		// Handle error
	}
	_ = result
}

// ExampleSlice demonstrates slice encoding/decoding.
func ExampleSlice() {
	// Using helper functions for slices
	buffer := make([]byte, 100)
	encoder := NewEncoder(buffer)
	values := []string{"a", "b", "c"}

	err := WriteSlice(&encoder, values, func(w Writer, v string) {
		w.WriteString(v)
	})
	if err != nil {
		// Handle error
	}

	// Reading slices
	data := []byte{0x93, 0xa1, 0x61, 0xa1, 0x62, 0xa1, 0x63}
	decoder := NewDecoder(data)
	result, err := ReadSlice(&decoder, func(r Reader) (string, error) {
		return r.ReadString()
	})
	if err != nil {
		// Handle error
	}
	_ = result
}

// ExampleMap demonstrates map encoding/decoding.
func ExampleMap() {
	// Using helper functions for maps
	buffer := make([]byte, 100)
	encoder := NewEncoder(buffer)
	values := map[string]int{"a": 1, "b": 2, "c": 3}

	err := WriteMap(&encoder, values,
		func(w Writer, k string) { w.WriteString(k) },
		func(w Writer, v int) { w.WriteInt64(int64(v)) },
	)
	if err != nil {
		// Handle error
	}

	// Reading maps
	data := []byte{0x83, 0xa1, 0x61, 0x01, 0xa1, 0x62, 0x02, 0xa1, 0x63, 0x03}
	decoder := NewDecoder(data)
	result, err := ReadMap(&decoder,
		func(r Reader) (string, error) { return r.ReadString() },
		func(r Reader) (int, error) { v, err := r.ReadInt64(); return int(v), err },
	)
	if err != nil {
		// Handle error
	}
	_ = result
}
