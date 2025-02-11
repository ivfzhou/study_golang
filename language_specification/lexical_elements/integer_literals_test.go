package lexical_elements_test

import "testing"

const (
	// 十进制
	I1 int = 1   // 1
	I2 int = -1  // -1
	I3 int = 1_1 // 11

	// 二进制
	I4 int = 0b1   // 1
	I5 int = -0b1  // -1
	I6 int = 0b_1  // 1
	I7 int = 0b1_1 // 3

	// 八进制
	I8  int = 01   // 1
	I9  int = -01  // -1
	I10 int = 0o1  // 1
	I11 int = 0_1  // 1
	I12 int = 0o_1 // 1
	I13 int = 01_1 // 9

	// 十六进制
	I14 int = 0x1   // 1
	I15 int = -0x1  // -1
	I16 int = 0x_1  // 1
	I17 int = 0x1_1 // 17
)

func TestInteger(t *testing.T) {
	println(I1)
	println(I2)
	println(I3)
	println(I4)
	println(I5)
	println(I6)
	println(I7)
	println(I8)
	println(I9)
	println(I10)
	println(I11)
	println(I12)
	println(I13)
	println(I14)
	println(I15)
	println(I16)
	println(I17)
}
