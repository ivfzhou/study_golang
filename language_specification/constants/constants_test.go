package constants_test

import (
	"testing"
	"unsafe"
)

const (
	Len     = len("")
	Cap     = cap([1]int{})
	Imag    = imag(1 + 1i)
	Real    = real(1 + 1i)
	Bool    = false || true
	Iota    = iota
	Unsafe  = unsafe.Sizeof("")
	Complex = complex(1, 1)
)

func TestConstants(t *testing.T) {
	println(Len)
	println(Cap)
	println(Imag)
	println(Real)
	println(Bool)
	println(Iota)
	println(Unsafe)
	println(Complex)
}
