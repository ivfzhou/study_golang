package statements_test

import (
	"fmt"
	"testing"
)

func TestDefer(t *testing.T) {
	fn := func() func() {
		return nil
	}
	defer fn()()
	fmt.Println("ok")
}

func TestDefer1(t *testing.T) {
	fn := func() {
		fmt.Println(1)
	}
	defer fn()
	fn = func() {
		fmt.Println(2)
	}
	fn = nil
	// 1
}

func TestDefer2(t *testing.T) {
	fn := func(i *int) {
		fmt.Println(*i)
	}
	i := 3
	defer fn(&i)
	i = 4
}
