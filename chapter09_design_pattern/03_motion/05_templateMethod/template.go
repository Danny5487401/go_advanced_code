package main

import (
	"fmt"
)

/*
模版模式
	定义一个算法骨架，将一些步骤延迟到子类进行。模板模式使得子类可以不改变一个算法的结构，即可重新定义该算法的某些特定步骤
优点
	封装不变的部分，扩展可变部分
	提取公共代码，便于维护
	行为由父类控制，子类实现
缺点
	每一个不同的实例，都需要一个子类来实现，导致类的个数增加。
需求：
	Game的父类，作为模板类。然后实现football和basketball两个子类，

*/

//Game 模板基类
type Game struct {
	Initialize func() //初始化
	StartPlay  func() // 开始
	EndPlay    func() // 结束
}

//Play 模板基类的Play方法
func (g Game) Play() {
	g.Initialize()
	g.StartPlay()
	g.EndPlay()
}

// 由于Go语言中，没有继承的概念，所以我们使用匿名组合的方式，来实现继承
//FootBall 子类，继承ame类
type FootBall struct {
	Game
}

//NewFootBall 实例化football子类
func NewFootBall() *FootBall {
	ft := new(FootBall)
	ft.Game.Initialize = ft.Initialize
	ft.Game.StartPlay = ft.StartPlay
	ft.Game.EndPlay = ft.EndPlay
	return ft
}

//Initialize 子类的Initialize方法
func (ft *FootBall) Initialize() {
	fmt.Println("Football game initialize")
}

//StartPlay 子类的StartPlay方法
func (ft *FootBall) StartPlay() {
	fmt.Println("Football game started.")
}

//EndPlay 子类的EndPlay方法
func (ft *FootBall) EndPlay() {
	fmt.Println("Football game Finished!")
}

// 开始调用
func main() {
	football := NewFootBall()
	football.Play()
	fmt.Println("-------------------")

}
