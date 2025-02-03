package expressions_test

import (
	"fmt"
	"testing"
)

type T struct {
	A int
}

func (t T) M() {
	fmt.Println(t.A)
}

func (t *T) M1() {
	fmt.Println(t.A)
}

func TestMethodValue(t *testing.T) {
	v := T{}
	fn := v.M
	fn() // 0
	v.A = 1
	fn() // 0

	v1 := &T{}
	fn = v1.M
	fn() // 0
	v1.A = 1
	fn() // 0

	v2 := T{}
	fn = v2.M1
	fn() // 0
	v2.A = 1
	fn() // 1

	fn = func() {}
	fn() //
}
