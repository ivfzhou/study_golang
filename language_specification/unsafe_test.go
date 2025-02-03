package language_specification_test

import (
	"fmt"
	"testing"
	"unsafe"
)

type UnsafeStruct struct {
	a bool  // 占 8 字节，实用 1 字节。
	b []int // 占 24 字节。
	c int16 // 占 8 字节，实用 2 字节。
}

func TestUnsafe(t *testing.T) {
	s := UnsafeStruct{}
	fmt.Println(unsafe.Sizeof(s), unsafe.Alignof(s))                           // 40 8
	fmt.Println(unsafe.Sizeof(s.a), unsafe.Alignof(s.a), unsafe.Offsetof(s.a)) // 1 1 0
	fmt.Println(unsafe.Sizeof(s.b), unsafe.Alignof(s.b), unsafe.Offsetof(s.b)) // 24 8 8
	fmt.Println(unsafe.Sizeof(s.c), unsafe.Alignof(s.c), unsafe.Offsetof(s.c)) // 2 2 32
}

func TestUnsafe0(t *testing.T) {
	f := 1.1
	a := *(*uint64)(unsafe.Pointer(&f))
	fmt.Println("float64 -> uint64:", a)
	fmt.Printf("%b\n", a)
}

func TestUnsafe1(t *testing.T) {
	s := UnsafeStruct{}
	b := (*[]int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.b)))
	*b = []int{1, 2, 3}
	fmt.Println(s.b) // [1, 2, 3]
}

type UnsafeStruct0 struct {
	// *UnsafeStruct
	UnsafeStruct
	d int
}

func TestUnsafe2(t *testing.T) {
	s := &UnsafeStruct0{}
	fmt.Println(unsafe.Offsetof(s.d)) // 40
}

func TestUnsafe3(t *testing.T) {
	var (
		i   interface{}
		s   []int
		str string
		in  int
		f   float64
		m   map[string]string
		ch  chan int
		b   bool
		fn  func()
		ui  uintptr
	)
	fmt.Println(unsafe.Sizeof(i))   // 16 字节数
	fmt.Println(unsafe.Sizeof(s))   // 24
	fmt.Println(unsafe.Sizeof(str)) // 16
	fmt.Println(unsafe.Sizeof(in))  //  8
	fmt.Println(unsafe.Sizeof(f))   //  8
	fmt.Println(unsafe.Sizeof(m))   //  8
	fmt.Println(unsafe.Sizeof(ch))  //  8
	fmt.Println(unsafe.Sizeof(b))   //  1
	fmt.Println(unsafe.Sizeof(fn))  //  8
	fmt.Println(unsafe.Sizeof(ui))  //  8
}

func TestUnsafe4(t *testing.T) {
	s := UnsafeStruct0{}
	scp := unsafe.Add(unsafe.Pointer(&s.b), 24)
	sc := (*int16)(scp)
	*sc = 12
	fmt.Println(s.c) // 12

	sap := unsafe.Add(scp, -32)
	sa := (*int8)(sap)
	*sa = 2
	fmt.Println(s.a) // true
}
