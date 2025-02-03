package expressions_test

import "testing"

type V struct{}

func (V) M1() {}

func (*V) M2() {}

func TestMethodValue1(t *testing.T) {
	fn := V.M1
	fn(V{})
	V.M1(V{})

	fn1 := (*V).M2
	fn1(&V{})
	// fn2 := V.M2
	// fn2(V{})

	//V.M2(V{})
	(*V).M2(&V{})
}
