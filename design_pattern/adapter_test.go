// 适配器模式是一种结构型设计模式，它能使接口不兼容的对象能够相互合作。
// 当你希望使用某个类，但是其接口与其他代码不兼容时，可以使用适配器类。
// 如果您需要复用这样一些类，他们处于同一个继承体系，并且他们又有了额外的一些共同的方法，
// 但是这些共同的方法不是所有在这一继承体系中的子类所具有的共性。
package design_pattern_test

import (
	"fmt"
	"testing"
)

/*
有不同类型的电脑，他们都兼容一个产品，所以产品能连接到这些电脑上。
*/

func TestAdapter(t *testing.T) {
	client := &Client{}

	var mac Computer_ = &Mac_{}
	client.InsertLightningConnectorIntoComputer(mac)

	windowsMachine := &Windows_{}
	windowsMachineAdapter := &WindowsAdapter{
		windowMachine: windowsMachine,
	}

	client.InsertLightningConnectorIntoComputer(windowsMachineAdapter)
}

// ===

type Client struct{}

func (c *Client) InsertLightningConnectorIntoComputer(com Computer_) {
	fmt.Println("Client inserts Lightning connector into computer.")
	com.InsertIntoLightningPort()
}

// =

type Computer_ interface {
	InsertIntoLightningPort()
}

// =

type Mac_ struct{}

func (m *Mac_) InsertIntoLightningPort() {
	fmt.Println("Lightning connector is plugged into mac machine.")
}

// =

type WindowsAdapter struct {
	windowMachine *Windows_
}

func (w *WindowsAdapter) InsertIntoLightningPort() {
	fmt.Println("Adapter converts Lightning signal to USB.")
	w.windowMachine.InsertIntoUSBPort()
}

// =

type Windows_ struct{}

func (w *Windows_) InsertIntoUSBPort() {
	fmt.Println("USB connector is plugged into windows machine.")
}
