package lexical_elements_test

import "testing"

const (
	// 十进制
	F1 float64 = 1.      // 1
	F2 float64 = -1.     // -1
	F3 float64 = 1_1.1_1 // 11.11
	F4 float64 = 08.     // 8
	F5 float64 = 0_8.    // 8
	F6 float64 = 1e-1    // 0.1
	F7 float64 = .1e1    // 1
	F8 float64 = 1e1 + 1 // 11
	F9 float64 = 1_1e1_1 // 1.1 * 10^12

	// 十六进制，e p 后面是十进制。
	F10 float64 = 0x1       // 1
	F11 float64 = 0x_1      // 1
	F12 float64 = -0x1p1    // 2
	F13 float64 = 0x1p-1    // 0.5
	F14 float64 = 0x1_1p1_1 // 34816
	F15 float64 = 0x1e_1    // 481
	F16 float64 = 0x1e      // 30
)

func TestFloat(t *testing.T) {
	println(F1)
	println(F2)
	println(F3)
	println(F4)
	println(F5)
	println(F6)
	println(F7)
	println(F8)
	println(F9)
	println(F10)
	println(F11)
	println(F12)
	println(F13)
	println(F14)
	println(F15)
	println(F16)
}
