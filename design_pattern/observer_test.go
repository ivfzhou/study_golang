// 观察者是一种行为设计模式，允许一个对象将其状态的改变通知其他对象。
// 当一个对象状态的改变需要改变其他对象，或实际对象是事先未知的或动态变化的时，可使用观察者模式。
// 当应用中的一些对象必须观察其他对象时，可使用该模式。但仅能在有限时间内或特定情况下使用。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestObserver(t *testing.T) {
	shirtItem := NewItem("Nike Shirt")

	var observerFirst Observer = &Customer{id: "abc@gmail.com"}
	var observerSecond Observer = &Customer{id: "xyz@gmail.com"}

	shirtItem.Register(observerFirst)
	shirtItem.Register(observerSecond)

	shirtItem.UpdateAvailability()

	shirtItem.Deregister(observerSecond)
	shirtItem.UpdateAvailability()
}

// ===

type Subject interface {
	Register(Observer)
	Deregister(Observer)
	NotifyAll()
}

// =

type Item struct {
	observerList []Observer
	name         string
	inStock      bool
}

func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

func (i *Item) UpdateAvailability() {
	fmt.Printf("Item %s is now in stock\n", i.name)
	i.inStock = true
	i.NotifyAll()
}

func (i *Item) Register(o Observer) {
	i.observerList = append(i.observerList, o)
}

func (i *Item) Deregister(o Observer) {
	i.observerList = removeFromSlice(i.observerList, o)
}

func (i *Item) NotifyAll() {
	for _, observer := range i.observerList {
		observer.Update(i.name)
	}
}

// =

type Observer interface {
	Update(string)
	GetID() string
}

// =

type Customer struct {
	id string
}

func (c *Customer) Update(itemName string) {
	fmt.Printf("Sending email to customer %s for item %s\n", c.id, itemName)
}

func (c *Customer) GetID() string {
	return c.id
}

// =

func removeFromSlice(observerList []Observer, observerToRemove Observer) []Observer {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.GetID() == observer.GetID() {
			observerList[observerListLength-1], observerList[i] = nil, observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}
