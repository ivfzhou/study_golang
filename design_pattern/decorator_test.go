// 装饰模式是一种结构型设计模式，允许你通过将对象放入包含行为的特殊封装对象中来为原对象绑定新的行为。
// 如果你希望在无需修改代码的情况下即可使用对象，且希望在运行时为对象新增额外的行为，可以使用装饰模式。
// 如果用继承来扩展对象行为的方案难以实现或者根本不可行，你可以使用该模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestDecorator(t *testing.T) {
	var pizza IPizza = &VeggieMania{}

	// Add cheese topping
	var pizzaWithCheese IPizza = &CheeseTopping{
		pizza: pizza,
	}

	// Add tomato topping
	var pizzaWithCheeseAndTomato IPizza = &TomatoTopping{
		pizza: pizzaWithCheese,
	}

	fmt.Printf("Price of veggeMania with tomato and cheese topping is %d\n", pizzaWithCheeseAndTomato.GetPrice())
}

// ===

type IPizza interface {
	GetPrice() int
}

// =

type VeggieMania struct{}

func (p *VeggieMania) GetPrice() int {
	return 15
}

// =

type TomatoTopping struct {
	pizza IPizza
}

func (c *TomatoTopping) GetPrice() int {
	pizzaPrice := c.pizza.GetPrice()
	return pizzaPrice + 7
}

// =

type CheeseTopping struct {
	pizza IPizza
}

func (c *CheeseTopping) GetPrice() int {
	pizzaPrice := c.pizza.GetPrice()
	return pizzaPrice + 10
}
