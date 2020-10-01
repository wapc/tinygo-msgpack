package msgpack

const (
	FormatError                  = 0
	FormatFourBytes              = 0xffffffff
	FormatFourLeastSigBitsInByte = 0x0f
	FormatFourSigBitsInByte      = 0xf0
	FormatPositiveFixInt         = 0x00
	FormatFixMap                 = 0x80
	FormatFixArray               = 0x90
	FormatFixString              = 0xa0
	FormatNil                    = 0xc0
	FormatFalse                  = 0xc2
	FormatTrue                   = 0xc3
	FormatBin8                   = 0xc4
	FormatBin16                  = 0xc5
	FormatBin32                  = 0xc6
	FormatExt8                   = 0xc7
	FormatExt16                  = 0xc8
	FormatExt32                  = 0xc9
	FormatFloat32                = 0xca
	FormatFloat64                = 0xcb
	FormatUint8                  = 0xcc
	FormatUint16                 = 0xcd
	FormatUint32                 = 0xce
	FormatUint64                 = 0xcf
	FormatInt8                   = 0xd0
	FormatInt16                  = 0xd1
	FormatInt32                  = 0xd2
	FormatInt64                  = 0xd3
	FormatFixExt1                = 0xd4
	FormatFixExt2                = 0xd5
	FormatFixExt4                = 0xd6
	FormatFixExt8                = 0xd7
	FormatFixExt16               = 0xd8
	FormatString8                = 0xd9
	FormatString16               = 0xda
	FormatString32               = 0xdb
	FormatArray16                = 0xdc
	FormatArray32                = 0xdd
	FormatMap16                  = 0xde
	FormatMap32                  = 0xdf
	FormatNegativeFixInt         = 0xe0
)
