package sb

import (
	"reflect"
	"unsafe"
)

func B2S(bs []byte) string {
	if nil == bs {
		return ""
	}
	return *(*string)(unsafe.Pointer(&bs))
}

func S2B(s *string) []byte {

	if nil == s {
		return nil
	}
	var bs []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	bh.Data, bh.Cap, bh.Len = sh.Data, sh.Len, sh.Len

	return bs
}
