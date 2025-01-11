// 策略模式是一种行为设计模式，它能让你定义一系列算法，并将每种算法分别放入独立的类中，以使算法的对象能够相互替换。
// 当你想使用对象中各种不同的算法变体，并希望能在运行时切换算法时，可使用策略模式。
// 当你有许多仅在执行某些行为时略有不同的相似类时，可使用策略模式。
// 如果算法在上下文的逻辑中不是特别重要，使用该模式能将类的业务逻辑与其算法实现细节隔离开来。
// 当类中使用了复杂条件运算符以在同一算法的不同变体中切换时，可使用该模式。
package design_pattern_test

import "testing"

func TestStrategy(t *testing.T) {
	c := &Cache{}
	var lru EvictionAlgo = &Lru{}
	c.SetEvictionStrategy(lru)
	c.Evict()
}

// ===

type EvictionAlgo interface {
	Evict(*Cache)
}

// =

type Cache struct {
	e EvictionAlgo
}

func (c *Cache) SetEvictionStrategy(e EvictionAlgo) {
	c.e = e
}

func (c *Cache) Evict() {
	c.e.Evict(c)
}

// =

type Fifo struct{}

func (*Fifo) Evict(c *Cache) {
	println("fifo")
}

// =

type Lfu struct{}

func (*Lfu) Evict(c *Cache) {
	println("lfu")
}

// =

type Lru struct{}

func (*Lru) Evict(c *Cache) {
	println("lru")
}
