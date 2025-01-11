// 单例模式是一种创建型设计模式，让你能够保证一个类只有一个实例，并提供一个访问该实例的全局节点。
// 如果程序中的某个类对于所有客户端只有一个可用的实例，可以使用单例模式。
// 如果你需要更加严格地控制全局变量，可以使用单例模式。
package design_pattern_test

import (
	"fmt"
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
	for i := 0; i < 30; i++ {
		go GetInstance()
	}

	fmt.Scanln()
}

// ===

var lock = &sync.Mutex{}

type Single struct{}

var singleInstance *Single

func GetInstance() *Single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating Single instance now.")
			singleInstance = &Single{}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}
