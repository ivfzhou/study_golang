package generic

func Reverse[E any](arr []E) {
	last := len(arr) - 1
	first := 0
	for last > first {
		arr[first], arr[last] = arr[last], arr[first]
		first++
		last--
	}
}
