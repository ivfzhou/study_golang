package language_specification_test

import (
	"fmt"
	"testing"
)

func TestOperandConvert(t *testing.T) {
	const (
		x         = 10
		y float64 = 1.1
	)
	fmt.Println(x + y) // 11.1

	const (
		m int = 10
		n     = 1.0
	)
	fmt.Println(m + n) // 11
}

func TestShift(t *testing.T) {
	s := uint(32)
	var j uint32 = 1 << s // 1 强转成 uint32。
	fmt.Println(j)

	var jj = uint8(1 << s) // 1 强转成 uint8。
	fmt.Println(jj)

	var jjj uint8 = 1.0 << s // 1.0 强转成 uint8。
	fmt.Println(jjj)

	var jjjj = 1<<s != uint8(1) // 1 强转成 uint8。
	fmt.Println(jjjj)

	arr := []int{1, 2, 3}
	var j5 = arr[1.0<<(s+32)] // 1.0 强转成 int。
	fmt.Println(j5)
}

func TestShift1(t *testing.T) {
	fn := func() int {
		return -1
	}

	var i = uint8(fn() << 1)
	fmt.Println(i)
}
