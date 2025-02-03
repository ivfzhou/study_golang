package language_specification_test

import (
	"fmt"
	"testing"
)

func TestTypeAssert(t *testing.T) {
	var i interface{}
	fmt.Println(i) // nil
	v, ok := i.(struct{})
	fmt.Println(ok) // false
	fmt.Println(v)  // {}
}
