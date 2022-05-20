package msgpack

import "math"

type Sizer struct {
	length uint32
}

func NewSizer() Sizer {
	return Sizer{}
}

func (s *Sizer) Len() uint32 {
	return s.length
}

func (s *Sizer) WriteNil() {
	s.length++
}

func (s *Sizer) WriteString(value string) {
	buf := []byte(value)
	length := uint32(len(buf))
	s.writeStringLength(length)
	s.length += length
}

func (s *Sizer) writeStringLength(length uint32) {
	if length < 32 {
		s.length++
	} else if length <= math.MaxUint8 {
		s.length += 2
	} else if length <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

func (s *Sizer) WriteBool(value bool) {
	s.length++
}

func (s *Sizer) WriteArraySize(length uint32) {
	if length < 16 {
		s.length++
	} else if length <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

func (s *Sizer) writeBinLength(length uint32) {
	if length < math.MaxUint8 {
		s.length += 1
	} else if length <= math.MaxUint16 {
		s.length += 2
	} else {
		s.length += 4
	}
}

func (s *Sizer) WriteByteArray(value []byte) {
	length := uint32(len(value))
	if length == 0 {
		s.length++
		return
	}
	s.writeBinLength(length)
	s.length += length + 1
}

func (s *Sizer) WriteMapSize(length uint32) {
	if length < 16 {
		s.length++
	} else if length <= math.MaxUint16 {
		s.length += 3
	} else {
		s.length += 5
	}
}

func (s *Sizer) WriteInt8(value int8) {
	s.WriteInt64(int64(value))
}
func (s *Sizer) WriteInt16(value int16) {
	s.WriteInt64(int64(value))
}
func (s *Sizer) WriteInt32(value int32) {
	s.WriteInt64(int64(value))
}
func (s *Sizer) WriteInt64(value int64) {
	if value >= -(1<<5) && value < 1<<7 {
		s.length++
	} else if value < 1<<7 && value >= -(1<<7) {
		s.length += 2
	} else if value < 1<<15 && value >= -(1<<15) {
		s.length += 3
	} else if value < 1<<31 && value >= -(1<<31) {
		s.length += 5
	} else {
		s.length += 9
	}
}

func (s *Sizer) WriteUint8(value uint8) {
	s.WriteUint64(uint64(value))
}
func (s *Sizer) WriteUint16(value uint16) {
	s.WriteUint64(uint64(value))
}
func (s *Sizer) WriteUint32(value uint32) {
	s.WriteUint64(uint64(value))
}
func (s *Sizer) WriteUint64(value uint64) {
	if value < 1<<7 {
		s.length++
	} else if value < 1<<8 {
		s.length += 2
	} else if value < 1<<16 {
		s.length += 3
	} else if value < 1<<32 {
		s.length += 5
	} else {
		s.length += 9
	}
}

func (s *Sizer) WriteFloat32(value float32) {
	s.length += 5
}
func (s *Sizer) WriteFloat64(value float64) {
	s.length += 9
}
