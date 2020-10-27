// +build wasm

package msgpack

import (
	"reflect"
	"unsafe"
)

// UnsafeString returns the byte slice as a volatile string
// THIS SHOULD ONLY BE USED BY THE CODE GENERATOR.
// THIS IS EVIL CODE.
// YOU HAVE BEEN WARNED.
func UnsafeString(b []byte) string {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: sh.Data, Len: sh.Len}))
}

// UnsafeBytes returns the string as a byte slice
// THIS SHOULD ONLY BE USED BY THE CODE GENERATOR.
// THIS IS EVIL CODE.
// YOU HAVE BEEN WARNED.
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  uintptr(len(s)),
		Cap:  uintptr(len(s)),
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
	}))
}
