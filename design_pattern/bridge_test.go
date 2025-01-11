// 桥接模式是一种结构型设计模式，可将一个大类或一系列紧密相关的类拆分为抽象和实现两个独立的层次结构，从而能在开发时分别使用。
// 如果你想要拆分或重组一个具有多重功能的庞杂类（例如能与多个数据库服务器进行交互的类），可以使用桥接模式。
// 如果你希望在几个独立维度上扩展一个类，可使用该模式。
// 如果你需要在运行时切换不同实现方法，可使用桥接模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestBridge(t *testing.T) {
	var hpPrinter Printer = &Hp{}
	var epsonPrinter Printer = &Epson{}

	var macComputer Computer = &Mac{}

	macComputer.SetPrinter(hpPrinter)
	macComputer.Print()
	fmt.Println()

	macComputer.SetPrinter(epsonPrinter)
	macComputer.Print()
	fmt.Println()

	var winComputer Computer = &Windows{}

	winComputer.SetPrinter(hpPrinter)
	winComputer.Print()
	fmt.Println()

	winComputer.SetPrinter(epsonPrinter)
	winComputer.Print()
	fmt.Println()
}

// ===

type Computer interface {
	Print()
	SetPrinter(Printer)
}

// =

type Mac struct {
	printer Printer
}

func (m *Mac) Print() {
	fmt.Println("Print request for mac")
	m.printer.PrintFile()
}

func (m *Mac) SetPrinter(p Printer) {
	m.printer = p
}

// =

type Windows struct {
	printer Printer
}

func (w *Windows) Print() {
	fmt.Println("Print request for windows")
	w.printer.PrintFile()
}

func (w *Windows) SetPrinter(p Printer) {
	w.printer = p
}

// =

type Printer interface {
	PrintFile()
}

// =

type Epson struct{}

func (p *Epson) PrintFile() {
	fmt.Println("Printing by a EPSON Printer")
}

// =

type Hp struct{}

func (p *Hp) PrintFile() {
	fmt.Println("Printing by a HP Printer")
}
