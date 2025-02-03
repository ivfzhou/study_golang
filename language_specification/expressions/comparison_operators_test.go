package expressions_test

import (
	"fmt"
	"testing"
)

type I interface {
	M()
}

type TT struct {
	x int
	Y int
}

func (TT) M() {}

func TestComparison(t *testing.T) {
	var tt = TT{}
	var i I = TT{}
	fmt.Println(tt == i) // true
	tt.x = 1
	fmt.Println(tt == i) // false
	tt.Y = 1
	tt.x = 0
	fmt.Println(tt == i) // false
}
