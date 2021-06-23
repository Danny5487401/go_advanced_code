package main

import "fmt"

/*
命令模式：
	客户端通过调用者发送命令,命令调用接收者执行相应操作

形象描述：
	遥控器对应上面的角色就是调用者,电视就是接收者,命令呢?对应的就是遥控器上的按键,最后客户端就是我们人啦,当我们想打开电视的时候,
	就会通过遥控器(调用者)的电源按钮(命令)来打开电视(接收者),在这个过程中遥控器是不知道电视的,但是电源按钮是知道他要控制谁的什么操作.
在命令模式中有如下几个角色:
	Command: 命令
	Invoker: 调用者
	Receiver: 接受者
	Client: 客户端

协作过程：
	1.Client创建一个ConcreteCommand对象并指定它的Receiver对象
	2.某Invoker对象存储该ConcreteCommand对象
	3.该Invoker通过调用Command对象的Excute操作来提交一个请求。若该命令是可撤消的，ConcreteCommand就在执行Excute操作之前存储当前状态以用于取消该命令
	4.ConcreteCommand对象对调用它的Receiver的一些操作以执行该请求
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

/*
go-redis源码分析
	// commands.go

	// 所有的命令
	type Cmdable interface {
		  Pipeline() Pipeliner
		  Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
		  TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
		  TxPipeline() Pipeliner
		  Command(ctx context.Context) *CommandsInfoCmd
		  ClientGetName(ctx context.Context) *StringCmd
		  // ...
		  // 和所有Redis命令的相关方法
	}

	// cmdable实现了Cmdable接口
	type cmdable func(ctx context.Context, cmd Cmder) error
	func (c cmdable) Echo(ctx context.Context, message interface{}) *StringCmd {
		cmd := NewStringCmd(ctx, "echo", message)
		_ = c(ctx, cmd)
		return cmd
	}
	//这里值得一提的是cmdable是一个函数类型，func(ctx context.Context, cmd Cmder) error
	//并且每个cmdable方法里都会有_ = c(ctx, cmd)

	type Client struct {
		  *baseClient
		  cmdable
		  hooks
		  ctx context.Context
	}

	func NewClient(opt *Options) *Client {
		  opt.init()

		  c := Client{
				baseClient: newBaseClient(opt, newConnPool(opt)),
				ctx:        context.Background(),
		  }
		  c.cmdable = c.Process //划线

		  return &c
	}


*/
