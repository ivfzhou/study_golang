package main_test

import "testing"

func TestCopy(t *testing.T) {
	var b = make([]byte, 5)
	println(copy(b, "你好"))
	println(b)
}
