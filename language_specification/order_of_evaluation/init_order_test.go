package order_of_evaluation_test

import "testing"

func init() {
	println(2)
}

func init() {
	println(1)
}

var a, b, c = f() + v(), g(), sqr(u()) + v() // f v g u sqr v 2 1

func f() int {
	println("f")
	return 1
}
func u() int {
	println("u")
	return 1
}
func v() int {
	println("v")
	return 1
}
func g() int {
	println("g")
	return 1
}
func sqr(x int) int {
	println("sqr")
	return x * x
}

func TestInit(t *testing.T) {
	println(a, b, c) // 2 1 2
}
