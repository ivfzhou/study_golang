package language_specification_test

import (
	"fmt"
	"testing"
)

func TestSilentlyDefer(t *testing.T) {
	Caller()
}

func Caller() {
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Println("恐慌了，恐慌恢复。")
		}
	}()

	res, err := Callee("")
	if err != nil {
		println("请求失败。")
	}

	fmt.Println("请求结果：" + res)
}

func Callee(nonEmpty string) (string, error) {
	if nonEmpty == "" {
		panic(nil)
	}

	return "ok", nil
}

func TestDeferOrder(t *testing.T) {
	defer func() {
		fmt.Println("OK4")
	}()
	defer func() {
		fmt.Println("OK3")
		panic("OK7")
	}()
	defer func() {
		fmt.Println("OK2")
		panic("OK6")
	}()
	fmt.Println("OK1")
	panic("OK5")
}

func TestDeferOrder0(t *testing.T) {
	defer func() {
		fmt.Println("OK3")
		fmt.Println(recover())
	}()
	defer func() {
		fmt.Println("OK2")
		panic("OK4")
	}()
	fmt.Println("OK1")
	panic(nil)
}
