package functools

import (
	"reflect"
	"unsafe"
)

func ByteSliceToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func StringToByteSlice(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}
