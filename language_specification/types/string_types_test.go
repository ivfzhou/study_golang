package types_test

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	s := "hello 你好"
	for i, v := range s {
		/*
			0 h
			1 e
			2 l
			3 l
			4 o
			5
			6 你
			9 好
		*/
		fmt.Println(i, string(v))
	}

	fmt.Println(s[7]) // 189
}

func TestString0(t *testing.T) {
	// s := "你好"

	/*p :=&s[0] // 无法取指。*/

	/*s[0] = 1 // 无法设值。*/
}
