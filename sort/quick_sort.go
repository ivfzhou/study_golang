package sort

func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr)-1)
}

func quickSort(arr []int, start, end int) {
	if start >= end {
		return
	}

	index := sortIndex(arr, start, end)
	quickSort(arr, start, index-1)
	quickSort(arr, index+1, end)
}

// 3, 7, 1, 5, 6, 4, 8, 9, 2
func sortIndex(arr []int, start, end int) int {
	tmp := arr[start]
	for start < end {

		for start < end && arr[end] >= tmp {
			end--
		}
		arr[start] = arr[end]

		for start < end && arr[start] <= tmp {
			start++
		}
		arr[end] = arr[start]
	}

	arr[start] = tmp

	return start
}
