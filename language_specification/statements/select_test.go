package statements_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	var ch = make(chan int)
	close(ch)
	select { // 两个 case 都能进入。
	case i := <-ch:
		fmt.Println(i) // 0
		break
	case ch <- 1:
	}

	select {} // block forever
}

func TestSelect1(t *testing.T) {
	var ch chan int
	var ch1 chan int
	var ch2 chan int
	var ch3 chan int
	var m map[string]int

	f := func() map[string]int {
		println("map 返回")
		_ = m[""]
		return m
	}

	ri := func() int {
		println("return 1")
		return 1
	}

	chfn := func() chan int {
		println("return chan")
		return nil
	}

	go func() {
		time.Sleep(time.Second * 10)
		println("休眠结束")
	}()

	go func() {
		println("即将进入 select")
		select { // chan 锁定了值则修改 chan 对本次 select 无影响，右表达式会计算，选中 chan 后再计算左表达式。
		case m[""] = <-ch1:
		case f()[""] = <-ch2:
		case <-ch:
		case <-chfn():
		case ch3 <- ri():
		}
		println("已退出 select")
	}()

	runtime.Gosched()
	ch = make(chan int)
	println("ch 已修改完毕")
	ch <- 1
	println("finish")
}
