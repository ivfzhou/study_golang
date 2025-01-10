package sort_test

import (
	"math/rand"
	"testing"

	"gitee.com/ivfzhou/study_golang/sort"
)

func TestDichotomySearch(t *testing.T) {
	for i := 0; i < 100; i++ {
		arr := randomArr()
		if !checkOrder(arr) {
			t.Error("failure: ", arr)
		}
	}
}

func checkOrder(arr []int) bool {
	length := len(arr) - 1
	for i := 0; i < length; i++ {
		if arr[i] > arr[i+1] {
			return false
		}
	}

	return true
}

func randomArr() []int {
	arr := make([]int, 0, rand.Intn(100))
	for i := 0; i < 100; i++ {
		arr = sort.OrderInsert(arr, rand.Intn(100))
	}
	return arr
}
