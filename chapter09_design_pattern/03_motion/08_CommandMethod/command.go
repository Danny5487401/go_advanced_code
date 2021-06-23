package main

import "fmt"

/*
命令模式：
	客户端通过调用者发送命令,命令调用接收者执行相应操作

形象描述：
	遥控器对应上面的角色就是调用者,电视就是接收者,命令呢?对应的就是遥控器上的按键,最后客户端就是我们人啦,当我们想打开电视的时候,
	就会通过遥控器(调用者)的电源按钮(命令)来打开电视(接收者),在这个过程中遥控器是不知道电视的,但是电源按钮是知道他要控制谁的什么操作.
*/

// 接收者,也就是一台大大的电视
type TV struct{}

// 定义接受者的功能
func (p TV) Open() {
	fmt.Println("play...")
}

func (p TV) Close() {
	fmt.Println("stop...")
}

// 按钮命令,按钮知道Tv的存在
// command
type Command interface {
	Press()
}

// 按钮一
type OpenCommand struct {
	tv TV // 知道Tv的存在
}

func (p OpenCommand) Press() {
	p.tv.Open()
}

// 按钮二
type CloseCommand struct {
	tv TV
}

func (p CloseCommand) Press() {
	p.tv.Close()
}

// 调用者
type Invoker struct {
	cmd Command
}

func (p *Invoker) SetCommand(cmd Command) {
	p.cmd = cmd
}

func (p Invoker) Do() {
	p.cmd.Press()
}

// 开始调用
func main() {
	var tv TV
	openCommand := OpenCommand{tv}
	invoker := Invoker{openCommand}
	invoker.Do()

	invoker.SetCommand(CloseCommand{tv})
	invoker.Do()
}
