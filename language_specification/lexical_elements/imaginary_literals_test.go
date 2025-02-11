package lexical_elements_test

import "testing"

const (
	// 十进制
	Im1 complex128 = 1i   // 1i
	Im2 complex128 = -1i  // -1i
	Im3 complex128 = 1e1i // 10i

	// 二进制
	Im4 complex128 = 0b1i // 1i

	// 八进制
	Im5 complex128 = 0o1i // 1i
	Im6 complex128 = 01i   // 1i

	// 十六进制
	Im7 complex128 = 0x1i   // 1i
	Im8 complex128 = 0x1p1i // 2i
)

func TestImaginary(t *testing.T) {
	println(Im1)
	println(Im2)
	println(Im3)
	println(Im4)
	println(Im5)
	println(Im6)
	println(Im7)
	println(Im8)
}
