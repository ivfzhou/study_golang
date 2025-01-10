package sort_test

import (
	"math/rand"
	"testing"
	"time"

	"gitee.com/ivfzhou/study_golang/sort"
)

func TestQuickSort(t *testing.T) {
	arr := make([]int, 100)
	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < 100; i++ {
		arr[i] = rand.Intn(100)
	}
	sort.QuickSort(arr)
	t.Log(order(arr))
	t.Log(arr)
}

func order(arr []int) bool {
	for i := 0; i < len(arr)-1; i++ {
		if arr[i] > arr[i+1] {
			return false
		}
	}
	return true
}
