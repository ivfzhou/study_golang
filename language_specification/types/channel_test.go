package types

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	var ch chan struct{}
	go func() {
		ch = make(chan struct{})
		runtime.Gosched()
		fmt.Println("make success")
		ch <- struct{}{}
		go func() {
			time.Sleep(time.Hour)
		}()
	}()
	fmt.Println("listen")
	<-ch // 一旦监听，无法再修改其值。
	fmt.Println("OK")
}
