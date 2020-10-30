package msgpack

// Codec is the interface that applies to data structures that can
// encode to and decode from the MessagPack format.
type Codec interface {
	Decode(decoder *Decoder) error
	Encode(encoder Writer) error
}
