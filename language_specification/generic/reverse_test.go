package generic_test

import (
	"testing"

	"gitee.com/ivfzhou/study_golang/language_specification/generic"
)

func TestReverse(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	generic.Reverse(arr)
	t.Log(arr)
}
