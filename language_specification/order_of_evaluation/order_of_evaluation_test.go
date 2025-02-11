package order_of_evaluation_test

import (
	"fmt"
	"testing"
)

func TestArrayOrder(t *testing.T) {
	for {
		a := 1
		f := func() int { a++; return a }
		x := []int{a, f()} // 先 f()
		if !(x[0] == 2 && x[1] == 2) {
			fmt.Println(x)
			break
		}
	}
}

func TestMapOrder(t *testing.T) {
	for {
		a := 2
		m := map[int]int{a: 1, a: 2}
		if m[2] == 1 {
			println(m)
			break
		}
	}
}

func TestMapOrder0(t *testing.T) {
	for {
		a := 1
		f := func() int { a++; return a }
		n := map[int]int{a: f()}
		if _, ok := n[1]; ok {
			println(n)
			break
		}
	}
}

func TestOrder(t *testing.T) {
	y := [1]int{}
	f := func() int {
		println("f1")
		return 0
	}
	g := func(h, ix, c int) int {
		println("g5")
		return 1
	}
	h := func() int {
		println("h2")
		return 1
	}
	i := func() int {
		println("i3")
		return 1
	}
	j := func() int {
		println("j4")
		return 0
	}
	x := [1]int{1}
	c := make(chan int)
	close(c)
	k := func() bool {
		println("k6")
		return true
	}

	ok := false
	y[f()], ok = g(h(), i()+x[j()], <-c), k() // left to right f1 h2 i3 j4 g5 k6
	println(ok)
}
