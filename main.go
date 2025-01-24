package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

func main() {
	to := 3
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
	case 2:
		garbage()
	case 3:
		garbageNoFree()
	}
}

func garbage() {
	type S struct {
		Field *int
	}

	func() {
		f := new(int)
		s := new(S)
		s.Field = f
		runtime.SetFinalizer(s, func(s *S) {
			fmt.Println("s is freed")
		})
		runtime.SetFinalizer(f, func(s *int) {
			fmt.Println("f is freed")
		})
	}()
	for i := 0; i < 10; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}
}

func garbageNoFree() {
	type S struct {
		Field *S
	}

	func() {
		f := new(S)
		s := new(S)
		s.Field = f
		f.Field = s
		runtime.SetFinalizer(s, func(s *S) {
			fmt.Println("s is freed")
		})
		runtime.SetFinalizer(f, func(s *S) {
			fmt.Println("f is freed")
		})
	}()
	for i := 0; i < 10; i++ {
		runtime.GC()
		time.Sleep(time.Second)
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
