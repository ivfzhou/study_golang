package expressions_test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	s := []int{1}
	a := &[1]int{1}
	fmt.Println(reflect.TypeOf(s).ConvertibleTo(reflect.TypeOf(a))) // true
	a = (*[1]int)(s)
}

func TestConvert1(t *testing.T) {
	type P struct {
		p string
	}

	type B struct {
		p string `json:"pp"`
	}

	var p P = P(B{"k"})
	val := reflect.ValueOf(p)
	field, _ := val.Type().FieldByName("p")
	fmt.Println(field.Tag) //
}

func TestConvert2(t *testing.T) {
	s := []int{1, 2}
	_ = (*[2]int)(s)
	_ = (*[1]int)(s)
	_ = (*[3]int)(s) // panic
}
