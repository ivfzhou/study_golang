package main

import (
	"fmt"
	"unsafe"
)

func main() {
	to := 1
	switch to {
	case 1:
		x := byte(1)
		bitPrint(x)
		y := float32(7.3125)
		bitPrint(y)
		s := struct {
			x int32
			y float32
		}{x: int32(x), y: y}
		bitPrint(s)
	}
}

func bitPrint[T any](t T) {
	size := int64(unsafe.Sizeof(t))
	ptr := unsafe.Pointer(&t)
	var b byte
	for i := size - 1; i > 0; i-- {
		b = *(*byte)(unsafe.Add(ptr, i))
		for j := 7; j >= 0; j-- {
			fmt.Printf("%d", b>>j&0b1)
		}
		fmt.Print("_")
	}
	b = *(*byte)(ptr)
	for j := 7; j >= 0; j-- {
		fmt.Printf("%d", b>>j&0b1)
	}
	fmt.Println()
}
