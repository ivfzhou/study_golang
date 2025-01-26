package main

import (
	"fmt"
	"sync"
)

func main() {
	ChannelLock(10)
}

func ChannelLock(incr int) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			AddValue(incr)
		}()
		go func() {
			defer wg.Done()
			AddValue(-incr)
		}()
	}
	wg.Wait()

	fmt.Println(QueryValue())
}

var (
	ch  = make(chan struct{}, 1)
	val int
)

func AddValue(a int) {
	ch <- struct{}{}
	val += a
	<-ch
}

func QueryValue() interface{} {
	ch <- struct{}{}
	v := val
	<-ch
	return v
}
