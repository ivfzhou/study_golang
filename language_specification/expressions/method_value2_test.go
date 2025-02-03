package expressions_test

import (
	"fmt"
	"testing"
)

type II interface {
	M1()
}

type s struct{}

func (*s) M1() {
	fmt.Println("ok")
}

func TestMethodValue2(t *testing.T) {
	fn := II.M1
	fn(&s{}) // ok
}
