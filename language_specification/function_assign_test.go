package language_specification_test

import (
	"fmt"
	"strings"
	"testing"
)

func TestFuncAssign(t *testing.T) {
	FuncTest0(FuncTest1("a.b.c"))    // 1
	FuncTest0(FuncTest2("a.b.c"))    // 1
	FuncTest0(FuncTest2("a.b.c")...) // 3

	FuncTest3(FuncTest4("a.b.c"))                               // a.b.c 5
	FuncTest3(FuncIntType.MethodTest0(FuncIntType(0), "a.b.c")) // a.b.c 5
}

func FuncTest0(params ...interface{}) {
	fmt.Println(len(params))
}

func FuncTest1(s string) []string {
	return strings.Split(s, ".")
}

func FuncTest2(s string) (is []interface{}) {
	for _, v := range strings.Split(s, ".") {
		is = append(is, v)
	}
	return
}

func FuncTest3(s string, i int) {
	fmt.Println(s, i)
}

func FuncTest4(s string) (string, int) {
	return s, len(s)
}

type FuncIntType int

func (FuncIntType) MethodTest0(s string) (string, int) {
	return s, len(s)
}
