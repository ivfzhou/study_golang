// 抽象工厂模式是一种创建型设计模式，它能创建一系列相关的对象，而无需指定其具体类。
// 如果代码需要与多个不同系列的相关产品交互，但是由于无法提前获取相关信息，或者出于对未来扩展性的考虑，你不希望代码基于产品的具体类进行构建，在这种情况下，你可以使用抽象工厂。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestAbstractFactory(t *testing.T) {
	adidasFactory, _ := GetSportsFactory("adidas")
	adidasShoe := adidasFactory.MakeShoe()
	adidasShirt := adidasFactory.MakeShirt()
	printShoeDetails(adidasShoe)
	printShirtDetails(adidasShirt)

	nikeFactory, _ := GetSportsFactory("nike")
	nikeShoe := nikeFactory.MakeShoe()
	nikeShirt := nikeFactory.MakeShirt()
	printShoeDetails(nikeShoe)
	printShirtDetails(nikeShirt)
}

func printShoeDetails(s IShoe) {
	fmt.Printf("Logo: %s", s.GetLogo())
	fmt.Printf("Size: %d", s.GetSize())
}

func printShirtDetails(s IShirt) {
	fmt.Printf("Logo: %s", s.GetLogo())
	fmt.Printf("Size: %d", s.GetSize())
}

// ===

type ISportsFactory interface {
	MakeShoe() IShoe
	MakeShirt() IShirt
}

func GetSportsFactory(brand string) (ISportsFactory, error) {
	if brand == "adidas" {
		return &Adidas{}, nil
	}

	if brand == "nike" {
		return &Nike{}, nil
	}

	return nil, fmt.Errorf("wrong brand type passed")
}

// =

type Adidas struct{}

func (a *Adidas) MakeShoe() IShoe {
	return &AdidasShoe{
		Shoe: Shoe{
			logo: "adidas",
			size: 14,
		},
	}
}

func (a *Adidas) MakeShirt() IShirt {
	return &AdidasShirt{
		Shirt: Shirt{
			logo: "adidas",
			size: 14,
		},
	}
}

// =

type Nike struct{}

func (n *Nike) MakeShoe() IShoe {
	return &NikeShoe{
		Shoe: Shoe{
			logo: "nike",
			size: 14,
		},
	}
}

func (n *Nike) MakeShirt() IShirt {
	return &NikeShirt{
		Shirt: Shirt{
			logo: "nike",
			size: 14,
		},
	}
}

// =

type IShoe interface {
	SetLogo(logo string)
	SetSize(size int)
	GetLogo() string
	GetSize() int
}

type Shoe struct {
	logo string
	size int
}

func (s *Shoe) SetLogo(logo string) {
	s.logo = logo
}

func (s *Shoe) GetLogo() string {
	return s.logo
}

func (s *Shoe) SetSize(size int) {
	s.size = size
}

func (s *Shoe) GetSize() int {
	return s.size
}

type AdidasShoe struct {
	Shoe
}

type NikeShoe struct {
	Shoe
}

// =

type IShirt interface {
	SetLogo(logo string)
	SetSize(size int)
	GetLogo() string
	GetSize() int
}

type Shirt struct {
	logo string
	size int
}

func (s *Shirt) SetLogo(logo string) {
	s.logo = logo
}

func (s *Shirt) GetLogo() string {
	return s.logo
}

func (s *Shirt) SetSize(size int) {
	s.size = size
}

func (s *Shirt) GetSize() int {
	return s.size
}

type AdidasShirt struct {
	Shirt
}

type NikeShirt struct {
	Shirt
}
