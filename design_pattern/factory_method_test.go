// 工厂方法模式是一种创建型设计模式，其在父类中提供一个创建对象的方法，允许子类决定实例化对象的类型。
// 当你在编写代码的过程中，如果无法预知对象确切类别及其依赖关系时，可使用工厂方法。
// 如果你希望用户能扩展你软件库或框架的内部组件，可使用工厂方法。
// 如果你希望复用现有对象来节省系统资源，而不是每次都重新创建对象，可使用工厂方法。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestFactoryMethod(t *testing.T) {
	ak47, _ := GetGun("ak47")
	musket, _ := GetGun("Musket")

	printDetails(ak47)
	printDetails(musket)
}

func printDetails(g IGun) {
	fmt.Printf("Gun: %s", g.GetName())
	fmt.Println()
	fmt.Printf("Power: %d", g.GetPower())
	fmt.Println()
}

// ===

type IGun interface {
	SetName(name string)
	SetPower(power int)
	GetName() string
	GetPower() int
}

func GetGun(gunType string) (IGun, error) {
	if gunType == "ak47" {
		return NewAk47(), nil
	}
	if gunType == "Musket" {
		return NewMusket(), nil
	}
	return nil, fmt.Errorf("Wrong gun type passed")
}

// =

type Gun struct {
	name  string
	power int
}

func (g *Gun) SetName(name string) {
	g.name = name
}

func (g *Gun) GetName() string {
	return g.name
}

func (g *Gun) SetPower(power int) {
	g.power = power
}

func (g *Gun) GetPower() int {
	return g.power
}

// =

type Ak47 struct {
	Gun
}

func NewAk47() IGun {
	return &Ak47{
		Gun: Gun{
			name:  "AK47 gun",
			power: 4,
		},
	}
}

// =

type Musket struct {
	Gun
}

func NewMusket() IGun {
	return &Musket{
		Gun: Gun{
			name:  "Musket gun",
			power: 1,
		},
	}
}
