package iterations

import (
	"unsafe"
)

func Repeat(s string, num int) string {
	length := len(s) * num

	repeated := make([]byte, len(s), length)
	copy(repeated, s)

	for len(repeated) < length {
		if len(repeated) <= length/2 {
			repeated = append(repeated, repeated...)
		} else {
			repeated = append(repeated, repeated[:length-len(repeated)]...)
		}
	}

	return *(*string)(unsafe.Pointer(&repeated))
}
