package generic_test

import (
	"strings"
	"testing"

	"gitee.com/ivfzhou/study_golang/language_specification/generic"
)

func TestGeneric(t *testing.T) {
	t.Log(generic.Sum[int, int](map[int]int{1: 2}))
	list := &generic.ArrayList[uint]{}
	list.Add(1)
	list.Add(2)
	t.Log(list.Get(1))
	builder := strings.Builder{}
	builder.WriteString("kkk")
	t.Log(generic.ToString[*strings.Builder](&builder))
	_ = make([]generic.ArrayList[int], 0)
}

func TestImplements(t *testing.T) {
	var c generic.Car[string] = nil
	var bm *generic.BaoMa[string] = &generic.BaoMa[string]{}
	bm.SetName("BaoMa")
	c = bm
	t.Log(c.Name())
	var byd = &generic.BYD{}
	byd.SetName("BYD")
	c = byd
	t.Log(c.Name())
}

func TestAFunc(t *testing.T) {
	generic.AFunc([]int{1, 2, 3})
}
