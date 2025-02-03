package expressions_test

import (
	"fmt"
	"testing"
)

var Arr = [4]int{1, 2, 3, 4}

func TestSlice(t *testing.T) {
	res := Arr[1:2] // addressable
	println(res[0])

	m := map[string][4]int{"": Arr}
	_ = m
	// res = m[""][1:2] // unaddressable

	fn := func() [4]int {
		return Arr
	}
	_ = fn
	// res = fn()[1:2] // unaddressable

	/*res = [4]int{1, 2, 3, 4}[1:2] // unaddressable*/

	m1 := map[string]*[4]int{"": &Arr}
	res = m1[""][1:2] // addressable
}

func TestSlice1(t *testing.T) {
	var s []int
	res := s[:]
	fmt.Println(res == nil) // true
	res = s[0:]
	fmt.Println(res == nil) // true
}

func TestSlice2(t *testing.T) {
	res := Arr[:0:4]
	fmt.Println(cap(res)) // 4

	res = Arr[:]

	res = Arr[1:]

	res = Arr[:1]
}
