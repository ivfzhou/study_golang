package types_test

import (
	"fmt"
	"testing"
)

type I interface {
	M1()
	M2()
}

type S struct {
	T
}

type T struct{}

func (T) M1() {}

func (*T) M2() {}

// S 拥有 T 的方法 M1，但没有 M2。
// *S 则拥有 T 的所有方法。
// 若匿名字段是 *T，则 S 和 *S 都拥有 T 的所有方法。
func TestMethodSet(t *testing.T) {
	var i interface{} = T{}
	_, ok := i.(I)
	fmt.Println(ok) // false

	i = &T{}
	_, ok = i.(I) // true
	fmt.Println(ok)
}
