// 生成器模式是一种创建型设计模式，使你能够分步骤创建复杂对象。该模式允许你使用相同地创建代码生成不同类型和形式的对象。
// 使用生成器模式可避免“重叠构造函数 （telescoping constructor）”的出现。
// 当你希望使用代码创建不同形式的产品（例如石头或木头房屋）时，可使用生成器模式。
// 使用生成器构造组合树或其他复杂对象。
package design_pattern_test

import (
	"fmt"
	"testing"
)

/*
建筑师可以使用不同的建筑方式造房子。
*/

func TestBuilder(*testing.T) {
	var normalBuilder IBuilder = GetBuilder("normal")
	var iglooBuilder IBuilder = GetBuilder("igloo")

	director := NewDirector(normalBuilder)
	normalHouse := director.BuildHouse()

	fmt.Printf("Normal House Door Type: %s\n", normalHouse.doorType)
	fmt.Printf("Normal House Window Type: %s\n", normalHouse.windowType)
	fmt.Printf("Normal House Num Floor: %d\n", normalHouse.floor)

	director.SetBuilder(iglooBuilder)
	iglooHouse := director.BuildHouse()

	fmt.Printf("Igloo House Door Type: %s\n", iglooHouse.doorType)
	fmt.Printf("Igloo House Window Type: %s\n", iglooHouse.windowType)
	fmt.Printf("Igloo House Num Floor: %d\n", iglooHouse.floor)
}

// ===

type IBuilder interface {
	SetWindowType()
	SetDoorType()
	SetNumFloor()
	GetHouse() House
}

type House struct {
	windowType string
	doorType   string
	floor      int
}

func GetBuilder(builderType string) IBuilder {
	if builderType == "normal" {
		return NewNormalBuilder()
	}

	if builderType == "igloo" {
		return NewIglooBuilder()
	}
	return nil
}

// =

type NormalBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func NewNormalBuilder() *NormalBuilder {
	return &NormalBuilder{}
}

func (b *NormalBuilder) SetWindowType() {
	b.windowType = "Wooden Window"
}

func (b *NormalBuilder) SetDoorType() {
	b.doorType = "Wooden Door"
}

func (b *NormalBuilder) SetNumFloor() {
	b.floor = 2
}

func (b *NormalBuilder) GetHouse() House {
	return House{
		doorType:   b.doorType,
		windowType: b.windowType,
		floor:      b.floor,
	}
}

// =

type IglooBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func NewIglooBuilder() *IglooBuilder {
	return &IglooBuilder{}
}

func (b *IglooBuilder) SetWindowType() {
	b.windowType = "Snow Window"
}

func (b *IglooBuilder) SetDoorType() {
	b.doorType = "Snow Door"
}

func (b *IglooBuilder) SetNumFloor() {
	b.floor = 1
}

func (b *IglooBuilder) GetHouse() House {
	return House{
		doorType:   b.doorType,
		windowType: b.windowType,
		floor:      b.floor,
	}
}

// =

type Director struct {
	builder IBuilder
}

func NewDirector(b IBuilder) *Director {
	return &Director{
		builder: b,
	}
}

func (d *Director) SetBuilder(b IBuilder) {
	d.builder = b
}

func (d *Director) BuildHouse() House {
	d.builder.SetDoorType()
	d.builder.SetWindowType()
	d.builder.SetNumFloor()
	return d.builder.GetHouse()
}
