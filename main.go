package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

func main() {
	to := 8
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
	case 4:
		pointerSize()
	case 5:
		var x, y uint8 = 15, 17
		fmt.Printf("isMultiplicationOverflow(%d, %d) is %v, and x*y = %d\n",
			x, y, isMultiplicationOverflow(x, y), x*y)
	case 6:
		x := uint8(3)
		fmt.Printf("isPowerOfTwo(%d) is %v\n", x, isPowerOfTwo(x))
	case 7:
		x, y := 6, 9
		gcd, lcm := greatestCommonDivisor(x, y)
		fmt.Printf("greatestCommonDivisor(%d, %d) is %d, %d\n", x, y, gcd, lcm)
	case 8:
		primeSieve()
	}
}

func primeSieve() {
	generate := func(ch chan<- int) {
		for i := 2; ; i++ {
			ch <- i
		}
	}
	filter := func(src <-chan int, dst chan<- int, prime int) {
		for i := range src {
			if i%prime != 0 {
				dst <- i
			}
		}
	}

	ch := make(chan int)
	go generate(ch)
	for {
		prime := <-ch
		fmt.Printf("%d\n", prime)
		time.Sleep(time.Millisecond * 100)
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}

// 两数乘积 = 两数的最大公约数 * 两数最小公倍数。
// 最小公倍数：两数共用的因数和各自的因数相乘。
func greatestCommonDivisor[T Integer](x, y T) (gcd T, leastCommonMultiple T) {
	mul := x * y

	for y != 0 {
		x, y = y, x%y
	}
	return x, mul / x
}

type Integer interface {
	int | int8 | int32 | int64 | UInteger
}

func isPowerOfTwo[T Integer](v T) bool {
	return (v-1)&v == 0
}

type UInteger interface {
	uint | uint8 | uint32 | uint64 | uintptr
}

func isMultiplicationOverflow[T UInteger](x, y T) bool {
	typeSize := unsafe.Sizeof(x)
	if uint64(x|y) < uint64(1<<4*typeSize) { // 该类型的一半位数量的最大值的平方不会大于该类型最大数。
		return false
	}
	return x > (^T(0))/y // 一个数超过了类型最大数的一半。
}

func pointerSize() {
	const PtrSize = 4 << (^uintptr(0) >> 63)
	fmt.Println("pointer size is", PtrSize) // 4 或者 8。
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
