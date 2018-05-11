package main

//#cgo pkg-config: libalpm
//#include <alpm.h>
import "C"
import "unsafe"

func verCompIsNewer(old, new string) bool {
	ca := C.CString(old)
	defer C.free(unsafe.Pointer(ca))
	cb := C.CString(new)
	defer C.free(unsafe.Pointer(cb))
	ret := int(C.alpm_pkg_vercmp(ca, cb))
	return ret < 0
}
