// 备忘录是一种行为设计模式，允许生成对象状态的快照并在以后将其还原。
// 当你需要创建对象状态快照来恢复其之前的状态时，可以使用备忘录模式。
// 当直接访问对象的成员变量、获取器或设置器将导致封装被突破时，可以使用该模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestMemento(t *testing.T) {
	caretaker := &Caretaker{
		mementoArray: make([]*Memento, 0),
	}

	originator := &Originator{
		state: "A",
	}

	fmt.Printf("Originator Current State: %s\n", originator.GetState())
	caretaker.AddMemento(originator.CreateMemento())

	originator.SetState("B")
	fmt.Printf("Originator Current State: %s\n", originator.GetState())
	caretaker.AddMemento(originator.CreateMemento())

	originator.SetState("C")
	fmt.Printf("Originator Current State: %s\n", originator.GetState())
	caretaker.AddMemento(originator.CreateMemento())

	originator.RestoreMemento(caretaker.GetMemento(1))
	fmt.Printf("Restored to State: %s\n", originator.GetState())

	originator.RestoreMemento(caretaker.GetMemento(0))
	fmt.Printf("Restored to State: %s\n", originator.GetState())
}

// ===

type Originator struct {
	state string
}

func (e *Originator) CreateMemento() *Memento {
	return &Memento{state: e.state}
}

func (e *Originator) RestoreMemento(m *Memento) {
	e.state = m.GetSavedState()
}

func (e *Originator) SetState(state string) {
	e.state = state
}

func (e *Originator) GetState() string {
	return e.state
}

// =

type Memento struct {
	state string
}

func (m *Memento) GetSavedState() string {
	return m.state
}

// =

type Caretaker struct {
	mementoArray []*Memento
}

func (c *Caretaker) AddMemento(m *Memento) {
	c.mementoArray = append(c.mementoArray, m)
}

func (c *Caretaker) GetMemento(index int) *Memento {
	return c.mementoArray[index]
}
