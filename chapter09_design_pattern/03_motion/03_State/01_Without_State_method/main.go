package main

import "fmt"

/*
需求：
	当今社会，论坛贴吧很多，我们也会加入感兴趣的论坛，偶尔进行发言，但有时却会发现不能发帖了，原来是昨天的某个帖子引发了口水战，被举报了。
	假设有三种状态，normal(正常），restricted(受限)，closed(封号)，判断依据是一个健康值（这里只是假设）。不同的状态有不同的权限
实现方式一：
	未使用状态模式
*/

// 定义状态
type AccountState int

const (
	NORMAL     AccountState = iota //正常0
	RESTRICTED                     //受限
	CLOSED                         //封号
)

// 通过健康值判断状态
type Account struct {
	State       AccountState
	HealthValue int
}

// 账户初始化
func NewAccount(health int) *Account {
	a := &Account{
		HealthValue: health,
	}
	// 设置状态
	a.changeState()
	return a
}

// 通过健康值判断
func (a *Account) changeState() {
	if a.HealthValue <= -10 {
		a.State = CLOSED
	} else if a.HealthValue > -10 && a.HealthValue <= 0 {
		a.State = RESTRICTED
	} else if a.HealthValue > 0 {
		a.State = NORMAL
	}
}

// 不同的状态值拥有不同的权限，三种权限--看帖-评论-发帖
//看帖
func (a *Account) View() {
	if a.State == NORMAL || a.State == RESTRICTED {
		fmt.Println("正常看帖")
	} else if a.State == CLOSED {
		fmt.Println("账号被封，无法看帖")
	}

}

//评论
func (a *Account) Comment() {
	if a.State == NORMAL || a.State == RESTRICTED {
		fmt.Println("正常评论")
	} else if a.State == CLOSED {
		fmt.Println("抱歉，你的健康值小于-10，不能评论")
	}

}

//发帖
func (a *Account) Post() {
	if a.State == NORMAL {
		fmt.Println("正常发帖")
	} else if a.State == RESTRICTED || a.State == CLOSED {
		fmt.Println("抱歉，你的健康值小于0，不能发帖")
	}
}

//给账户设定健康值
func (a *Account) SetHealth(value int) {
	a.HealthValue = value
	a.changeState()

}

// 具体调用
func main() {
	account := NewAccount(10)
	account.View()
	account.Comment()
	account.Post()
	account.SetHealth(-10)
	account.View()
	account.Comment()
	account.Post()

}

/*
缺点：
	1。看帖和发帖方法中都包含状态判断语句，以判断在该状态下是否具有该方法以及在特定状态下该方法如何实现，导致代码非常冗长，可维护性较差；
	2。系统扩展性较差，如果需要增加一种新的状态，如hot状态（活跃用户，该状态用户发帖积分增加更多），需要对原有代码进行大量修改，
		扩展起来非常麻烦
*/
