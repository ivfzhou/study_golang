// 享元模式是一种结构型设计模式，它摒弃了在每个对象中保存所有数据的方式，通过共享多个对象所共有的相同状态，让你能在有限的内存容量中载入更多对象。
// 仅在程序必须支持大量对象且没有足够的内存容量时使用享元模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

/*
游戏中有角色衣服，衣服由衣服工厂提供。同一件衣服可以给角色共用。
*/

func TestFlyweight(t *testing.T) {
	game := NewGame()

	// Add Terrorist
	game.AddTerrorist(TerroristDressType)
	game.AddTerrorist(TerroristDressType)
	game.AddTerrorist(TerroristDressType)
	game.AddTerrorist(TerroristDressType)

	// Add CounterTerrorist
	game.AddCounterTerrorist(CounterTerroristDressType)
	game.AddCounterTerrorist(CounterTerroristDressType)
	game.AddCounterTerrorist(CounterTerroristDressType)

	dressFactoryInstance := GetDressFactorySingleInstance()

	for dressType, dress := range dressFactoryInstance.dressMap {
		fmt.Printf("DressColorType: %s\nDressColor: %s\n", dressType, dress.GetColor())
	}
}

// ===

const (
	// TerroristDressType terrorist dress type
	TerroristDressType = "tDress"
	// CounterTerroristDressType terrorist dress type
	CounterTerroristDressType = "ctDress"
)

type Dress interface {
	GetColor() string
}

// =

type DressFactory struct {
	dressMap map[string]Dress
}

var dressFactorySingleInstance = &DressFactory{
	dressMap: make(map[string]Dress),
}

func GetDressFactorySingleInstance() *DressFactory {
	return dressFactorySingleInstance
}

func (d *DressFactory) GetDressByType(dressType string) (Dress, error) {
	if d.dressMap[dressType] != nil {
		return d.dressMap[dressType], nil
	}

	if dressType == TerroristDressType {
		d.dressMap[dressType] = NewTerroristDress()
		return d.dressMap[dressType], nil
	}
	if dressType == CounterTerroristDressType {
		d.dressMap[dressType] = NewCounterTerroristDress()
		return d.dressMap[dressType], nil
	}

	return nil, fmt.Errorf("Wrong dress type passed")
}

// =

type TerroristDress struct {
	color string
}

func NewTerroristDress() *TerroristDress {
	return &TerroristDress{color: "red"}
}

func (t *TerroristDress) GetColor() string {
	return t.color
}

// =

type CounterTerroristDress struct {
	color string
}

func NewCounterTerroristDress() *CounterTerroristDress {
	return &CounterTerroristDress{color: "green"}
}

func (c *CounterTerroristDress) GetColor() string {
	return c.color
}

// =

type Player struct {
	dress      Dress
	playerType string
	lat        int
	long       int
}

func NewPlayer(playerType, dressType string) *Player {
	dress, _ := GetDressFactorySingleInstance().GetDressByType(dressType)
	return &Player{
		playerType: playerType,
		dress:      dress,
	}
}

func (p *Player) newLocation(lat, long int) {
	p.lat = lat
	p.long = long
}

// =

type Game struct {
	terrorists        []*Player
	counterTerrorists []*Player
}

func NewGame() *Game {
	return &Game{
		terrorists:        make([]*Player, 1),
		counterTerrorists: make([]*Player, 1),
	}
}

func (c *Game) AddTerrorist(dressType string) {
	player := NewPlayer("T", dressType)
	c.terrorists = append(c.terrorists, player)
	return
}

func (c *Game) AddCounterTerrorist(dressType string) {
	player := NewPlayer("CT", dressType)
	c.counterTerrorists = append(c.counterTerrorists, player)
	return
}
