package main

import "fmt"

// Strategy 模式应用场景，在我们写的程序中，大多有if else的条件语句基本上都适合
// Strategy 模式，但是if else 条件的情况是不变的，则不适合此模式，例如一周7天
// Strategy 及其子类为组建提供了一系列可重用的算法，从而使得类型在运行时方便的根据
// 需要在各个算法之间进行切换

//一般的做法

/*
type API struct {

}
// 这个主题时变化不稳定的，没有可扩展性
func (a *API) Recognition(name string){
	if name == "ali"{
		fmt.Println("ali api调用")
	}else if name == "baidu"{
		fmt.Println("baidu api调用")
	}else if name == "xunfei"{
		fmt.Println("xunfei api调用")
	}
	//...如果有其他新的在此添加
}
*/
// ====================下面是用Strategy设计模式=====================

// ==============稳定===============
// 定义一个api接口，添加一个抽象方法 Recognition()
type IAPI interface {
	Recognition()
}

type API struct {
	// 这里个人认为不是继承而是组合，在重构关键技法中这也是一种提倡做法，继承——>组合
	iapi IAPI
}

func (a *API) OnProgress(){
	// 运行时动态改变
	a.iapi.Recognition()
}

// ==============变化可扩展的==================
// ali 接口
type Ali struct {

}
// 实现Recognition()抽象方法
func (a *Ali) Recognition(){
	fmt.Println("ali api 调用")
}

// bai 接口
type Bai struct {

}
// 实现Recognition()抽象方法
func (b *Bai) Recognition(){
	fmt.Println("baidu api 调用")
}



func main(){
	// 当然如果要彻底解决if else 的问题， 还需要工厂模式进行配合使用
	// 这里就不增加工厂模式，只是一个纯粹的strategy模式
	api := new(API)
	//使用阿里
	api.iapi = new(Ali)
	api.iapi.Recognition()
	//使用百度
	api.iapi = new(Bai)
	api.iapi.Recognition()
}

