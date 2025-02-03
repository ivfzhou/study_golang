package main_test

import "testing"

func TestAppend(t *testing.T) {
	println(append([]byte{}, "你好"...))
}
