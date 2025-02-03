package language_specification_test

import (
	"fmt"
	"testing"
)

func TestRound(t *testing.T) {
	// 0000 0101
	x := uint8(5)
	// 0000 0110
	y := uint8(10)
	// 0000 1010 10
	// 1111 0110 -10
	// 0000 0101 5
	// 1111 1011 补码
	// 1000 0101 原码
	// -5
	z := x - y
	fmt.Printf("%b %d %d\n", z, z, int8(z)) // 11111011 251 -5

	x = 0b0000_0001
	y = 0b1000_0010
	fmt.Printf("%d %d\n", x, y) // 1 130
	z = x - y
	fmt.Printf("%b %d %d\n", z, z, int8(z)) // 01111111 127
	// 两值相差 128 以上则强转值错误。
}
