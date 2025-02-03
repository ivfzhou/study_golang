package language_specification_test

import (
	"fmt"
	"testing"
)

type I interface {
	M()
}

type T struct {
	A string
}

func (t T) M() {
	fmt.Println("T.M()")
}

type T0 struct {
	B string
}

type T1 struct {
	T
	A string
	B string
}

func (t *T1) M() {
	fmt.Println("*T1.M()")
}

type T2 struct {
	T0
	T1
}

func TestStructSelector(t *testing.T) {
	v := T2{
		T0{"T0.B"},
		T1{T{"T.A"}, "T1.A", "T1.B"},
	}
	fmt.Println(v.A) // T1.A
	v.M()            // *T1.M()
	// fmt.Println(v.B) ambiguous selector
}

func TestStructSelector0(t *testing.T) {
	var i interface{} = T2{}
	_, ok := i.(I)
	fmt.Println(ok) // false

	i = &T2{}
	im, ok := i.(I)
	fmt.Println(ok) // true
	im.M()          // *T1.M()
}
