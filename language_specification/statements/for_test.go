package statements_test

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {
	var arr *[2]int
	for i, v := range arr { // 取值 panic
		fmt.Println(i, v)
	}
}

func TestFor1(t *testing.T) {
	var ch chan int
	for range ch {

	}
}
