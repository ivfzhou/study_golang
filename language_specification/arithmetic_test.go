package language_specification_test

import (
	"fmt"
	"math/rand"
	"testing"
)

/*
x/y = 除数取整。
x%y = 余数。
除数取整*y+余数 = x
*/

func TestArithmetic(t *testing.T) {
	intMax := 1<<63 - 1
	intMin := -1 << 63
	for x := intMin; x <= intMax; x++ {
		y := -1
		quotient := x / y
		remainder := x % y
		if !(x == quotient*y+remainder) {
			fmt.Println(x)
		}
	}
}

/*
x/2^n = x>>n
x%2^n = x%2^n&x
*/

func TestArithmetic0(t *testing.T) {
	intMax := 1<<63 - 1
	for x := 0; x <= intMax; x++ {
		bit := rand.Intn(64)
		y := 1 << bit
		if !(x/y == x>>bit && x%y == x%y&x) {
			fmt.Println(x)
		}
	}
}
