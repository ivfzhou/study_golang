// 状态模式是一种行为设计模式，让你能在一个对象的内部状态变化时改变其行为，使其看上去就像改变了自身所属的类一样。
// 如果对象需要根据自身当前状态进行不同行为，同时状态的数量非常多且与状态相关的代码会频繁变更的话，可使用状态模式。
// 如果某个类需要根据成员变量的当前值改变自身行为，从而需要使用大量的条件语句时，可使用该模式。
// 当相似状态和基于条件的状态机转换中存在许多重复代码时，可使用状态模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

/*
一个提供货物的机器，它有有货、无货、待投币、已投币四种状态。
四个状态用四个类实现，调用机器功能时，状态在这些类中切换。
*/

func TestState(t *testing.T) {
	vendingMachine := NewVendingMachine(1, 10)

	err := vendingMachine.RequestItem()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = vendingMachine.InsertMoney(10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = vendingMachine.DispenseItem()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = vendingMachine.AddItem(2)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = vendingMachine.RequestItem()
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = vendingMachine.InsertMoney(10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = vendingMachine.DispenseItem()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// ===

type State interface {
	AddItem(int) error
	RequestItem() error
	InsertMoney(money int) error
	DispenseItem() error
}

// =

type VendingMachine struct {
	hasItem       State
	itemRequested State
	hasMoney      State
	noItem        State

	currentState State

	itemCount int
	itemPrice int
}

func NewVendingMachine(itemCount, itemPrice int) *VendingMachine {
	v := &VendingMachine{
		itemCount: itemCount,
		itemPrice: itemPrice,
	}
	hasItemState := &HasItemState{
		vendingMachine: v,
	}
	itemRequestedState := &ItemRequestedState{
		vendingMachine: v,
	}
	hasMoneyState := &HasMoneyState{
		vendingMachine: v,
	}
	noItemState := &NoItemState{
		vendingMachine: v,
	}

	v.setState(hasItemState)
	v.hasItem = hasItemState
	v.itemRequested = itemRequestedState
	v.hasMoney = hasMoneyState
	v.noItem = noItemState
	return v
}

func (v *VendingMachine) RequestItem() error {
	return v.currentState.RequestItem()
}

func (v *VendingMachine) AddItem(count int) error {
	return v.currentState.AddItem(count)
}

func (v *VendingMachine) InsertMoney(money int) error {
	return v.currentState.InsertMoney(money)
}

func (v *VendingMachine) DispenseItem() error {
	return v.currentState.DispenseItem()
}

func (v *VendingMachine) setState(s State) {
	v.currentState = s
}

func (v *VendingMachine) incrementItemCount(count int) {
	fmt.Printf("Adding %d items\n", count)
	v.itemCount = v.itemCount + count
}

// =

type NoItemState struct {
	vendingMachine *VendingMachine
}

func (i *NoItemState) RequestItem() error {
	return fmt.Errorf("Item out of stock")
}

func (i *NoItemState) AddItem(count int) error {
	i.vendingMachine.incrementItemCount(count)
	i.vendingMachine.setState(i.vendingMachine.hasItem)
	return nil
}

func (i *NoItemState) InsertMoney(int) error {
	return fmt.Errorf("Item out of stock")
}

func (i *NoItemState) DispenseItem() error {
	return fmt.Errorf("Item out of stock")
}

// =

type HasItemState struct {
	vendingMachine *VendingMachine
}

func (i *HasItemState) RequestItem() error {
	if i.vendingMachine.itemCount == 0 {
		i.vendingMachine.setState(i.vendingMachine.noItem)
		return fmt.Errorf("No item present")
	}
	fmt.Printf("Item requestd\n")
	i.vendingMachine.setState(i.vendingMachine.itemRequested)
	return nil
}

func (i *HasItemState) AddItem(count int) error {
	fmt.Printf("%d items added\n", count)
	i.vendingMachine.incrementItemCount(count)
	return nil
}

func (i *HasItemState) InsertMoney(int) error {
	return fmt.Errorf("Please select item first")
}

func (i *HasItemState) DispenseItem() error {
	return fmt.Errorf("Please select item first")
}

// =

type ItemRequestedState struct {
	vendingMachine *VendingMachine
}

func (i *ItemRequestedState) RequestItem() error {
	return fmt.Errorf("Item already requested")
}

func (i *ItemRequestedState) AddItem(int) error {
	return fmt.Errorf("Item Dispense in progress")
}

func (i *ItemRequestedState) InsertMoney(money int) error {
	if money < i.vendingMachine.itemPrice {
		return fmt.Errorf("Inserted money is less. Please insert %d", i.vendingMachine.itemPrice)
	}
	fmt.Println("Money entered is ok")
	i.vendingMachine.setState(i.vendingMachine.hasMoney)
	return nil
}

func (i *ItemRequestedState) DispenseItem() error {
	return fmt.Errorf("Please insert money first")
}

// =

type HasMoneyState struct {
	vendingMachine *VendingMachine
}

func (i *HasMoneyState) RequestItem() error {
	return fmt.Errorf("Item dispense in progress")
}

func (i *HasMoneyState) AddItem(int) error {
	return fmt.Errorf("Item dispense in progress")
}

func (i *HasMoneyState) InsertMoney(int) error {
	return fmt.Errorf("Item out of stock")
}

func (i *HasMoneyState) DispenseItem() error {
	fmt.Println("Dispensing Item")
	i.vendingMachine.itemCount = i.vendingMachine.itemCount - 1
	if i.vendingMachine.itemCount == 0 {
		i.vendingMachine.setState(i.vendingMachine.noItem)
	} else {
		i.vendingMachine.setState(i.vendingMachine.hasItem)
	}
	return nil
}
