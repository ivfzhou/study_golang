// 命令是一种行为设计模式，它可将请求或简单操作转换为一个对象。
// 如果你需要通过操作来参数化对象，可使用命令模式。
// 如果你想要将操作放入队列中、操作的执行或者远程执行操作，可使用命令模式。
// 如果你想要实现操作回滚功能，可使用命令模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestCommand(t *testing.T) {
	var tv Device = &TV{}

	var onCommand Command = &OnCommand{
		device: tv,
	}

	var offCommand Command = &OffCommand{
		device: tv,
	}

	onButton := &Button{
		command: onCommand,
	}
	onButton.Press()

	offButton := &Button{
		command: offCommand,
	}
	offButton.Press()
}

// ===

type Button struct {
	command Command
}

func (b *Button) Press() {
	b.command.Execute()
}

// =

type Command interface {
	Execute()
}

// =

type OnCommand struct {
	device Device
}

func (c *OnCommand) Execute() {
	c.device.On()
}

// =

type OffCommand struct {
	device Device
}

func (c *OffCommand) Execute() {
	c.device.Off()
}

// =

type Device interface {
	On()
	Off()
}

// =

type TV struct {
	isRunning bool
}

func (t *TV) On() {
	t.isRunning = true
	fmt.Println("Turning tv On")
}

func (t *TV) Off() {
	t.isRunning = false
	fmt.Println("Turning tv Off")
}
