package language_specification_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestSlice(t *testing.T) {
	s1 := make([]int, 8, 10)
	s2 := s1[3:5:5]
	fmt.Println(cap(s2)) // 2

	s1[3] = 1
	fmt.Println(s2[0]) // 1 change

	s2 = append(s2, 1) // allocate
	s1[3] = 2
	fmt.Println(s2[0]) // 1 not change
}

func TestSlice0(t *testing.T) {
	arr := new([5]int)
	s1 := arr[:3:5]
	s2 := s1[:2:4]
	s2 = append(s2, 1, 1)
	s2 = append(s2, 1)
	fmt.Println(arr) // [0, 0, 1, 1, 0]
}

func TestSlice1(t *testing.T) {
	s := make([]string, 0, 3)
	s = append(s, "1")
	v := (*reflect.SliceHeader)(unsafe.Pointer(reflect.ValueOf(&s).Elem().UnsafeAddr()))
	v0 := (*reflect.SliceHeader)(unsafe.Pointer(reflect.ValueOf(s).Pointer()))
	v1 := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	v2 := (*reflect.SliceHeader)(unsafe.Pointer(&s[0]))
	fmt.Println(unsafe.Slice(&s[0], 3))                                // ["1", "", ""]
	fmt.Println(v.Len, v0.Len, v1.Len, v2.Len)                         // 1 1 1 1
	fmt.Println(v.Cap, v0.Cap, v1.Cap, v2.Cap)                         // 3 0 3 0
	fmt.Printf("%#x %#x %#x %#x\n", v.Data, v0.Data, v1.Data, v2.Data) // <ptr>
	fmt.Println(cap(s))                                                // 3
}

func TestSlice2(t *testing.T) {
	s := []int{1, 2, 3}
	ss := unsafe.Slice(&s[2], 4)
	fmt.Println(ss) // [3, 0, 0, 0]
}
