package statements_test

import (
	"fmt"
	"testing"
)

func TestAssign(t *testing.T) {
	x := []int{1, 2, 3}
	i := 0
	i, x[i] = 1, 2 // set i = 1, x[0] = 2
	t.Log(i, x)

	i = 2
	x = []int{3, 5, 7}
	for i, x[i] = range x {
		// set i, x[2] = 0, x[0]  3,5,3
		// set i, x[0] = 1, x[1]  5,5,3
		// set i, x[1] = 2, x[2]  5,3,3
	}
	fmt.Println(i, x)
}
