package msgpack

// Codec is the interface that applies to data structures that can
// encode to and decode from the MessagPack format.
type Codec interface {
	Decode(decoder *Decoder) error
	Encode(encoder Writer) error
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
