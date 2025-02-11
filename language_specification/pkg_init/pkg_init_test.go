package pkg_init_test

import "testing"

var x = I(T{}).ab()  // x has an undetected, hidden dependency on a and b
var _ = sideEffect() // unrelated to x, a, or b
var a = b
var b = 42

type I interface {
	ab() []int
}

type T struct{}

func (T) ab() []int {
	println(a, b)
	return []int{a, b}
}

func sideEffect() int {
	println("sideEffect")
	return 0
}

func TestI(t *testing.T) {
	println(x) // 42 42
}
