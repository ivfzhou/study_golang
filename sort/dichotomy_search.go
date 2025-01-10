package sort

// DichotomySearch 二分查找法，arr 为升序切片，将 elem 插入 arr 返回值位置之后，arr 仍保持有序。
// 返回 -1 表示插入 arr 开头。
func DichotomySearch(arr []int, elem int) int {
	// 空切片返回 -1。
	if len(arr) == 0 {
		return -1
	}

	left, right, isBigger := 0, len(arr)-1, false
	currentPosition := (left + right) / 2

	for {
		val := arr[currentPosition]
		if elem > val { // 指针位置值偏大，left 设为当前指针位。
			left = currentPosition
			isBigger = true
		} else if elem < val { // 指针位置值偏大，right 设为当前指针位。
			right = currentPosition
			isBigger = false
		} else { // 相等返回该位置。
			return currentPosition
		}

		newPosition := (left + right) / 2
		if currentPosition == newPosition { // 相等表明位置已确定。
			if isBigger {
				if elem > arr[right] {
					return right // 返回指针位右边位。
				}
				return left // 返回指针位。
			}
			return left - 1 // 返回指针位左边位。
		}
		// 缩小范围继续循环。
		currentPosition = newPosition
	}
}

// OrderInsert arr 为升序切片，按顺序插入 elem。
func OrderInsert(arr []int, elem int) []int {
	index := DichotomySearch(arr, elem)
	if index == -1 { // 插入前面。
		return append([]int{elem}, arr...)
	} else if index == len(arr)-1 { // 插入最后。
		return append(arr, elem)
	} else { // 中间插入。
		newArr := make([]int, 0, len(arr)+1)
		return append(append(append(newArr, arr[:index+1]...), elem), arr[index+1:]...)
	}
}
