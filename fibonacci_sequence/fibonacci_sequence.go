package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(fibonacci1(30), fibonacci2(30), fibonacci3(30))

	t := time.Now()
	for i := 0; i < 10000; i++ {
		fibonacci1(30)
	}
	fmt.Println(time.Since(t))

	t = time.Now()
	for i := 0; i < 10000; i++ {
		fibonacci2(30)
	}
	fmt.Println(time.Since(t))

	t = time.Now()
	for i := 0; i < 10000; i++ {
		fibonacci3(30)
	}
	fmt.Println(time.Since(t))
}

func fibonacci1(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	return fibonacci1(n-1) + fibonacci1(n-2)
}

func fibonacci2(n int) int {
	arr := make([]int, n+1)
	arr[1] = 1
	arr[2] = 1
	for i := 3; i <= n; i++ {
		arr[i] = arr[i-1] + arr[i-2]
	}
	return arr[n]
}

func fibonacci3(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	a, b := 1, 1
	for i := 3; i <= n; i++ {
		a, b = a+b, a
	}
	return a
}
