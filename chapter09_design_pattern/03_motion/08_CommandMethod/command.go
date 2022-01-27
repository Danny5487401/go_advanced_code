package main

import "fmt"

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

// 调用者-遥控
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
// 客户端(就是我们人)拿起遥控器,瞅准了打开按钮(SetCommand),并且按钮该键(Do),按键被按下(Press),电视打开了(Open).
func main() {
	var tv TV
	openCommand := OpenCommand{tv} // 生成名利客户端
	invoker := Invoker{openCommand}
	invoker.Do()

	invoker.SetCommand(CloseCommand{tv})
	invoker.Do()
}
