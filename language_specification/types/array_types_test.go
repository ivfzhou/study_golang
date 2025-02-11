package types_test

import (
	"fmt"
	"testing"
)

var (
	arr0 = []string{1: "hello", 3: "hello"}
	arr1 = [...]string{1: "hello", 3: "hello"}
	arr2 = [4]string{1: "hello", 3: "hello"}
)

/*
var Index = 2
var Arr3 = [Index]string{} // must be a constant expression
*/

const Index1 = 2

var arr4 = [Index1]int{}
var arr5 = [Index1 + 1]int{}

func TestArray(t *testing.T) {
	fmt.Println(len(arr0), cap(arr0))
	fmt.Println(len(arr1), cap(arr1))
	fmt.Println(len(arr2), cap(arr2))

	fmt.Println([1]string{""} == [1]string{""}) // true
}

func TestArray0(t *testing.T) {
	arr := new([2]int)
	fmt.Println(arr[0])
}
