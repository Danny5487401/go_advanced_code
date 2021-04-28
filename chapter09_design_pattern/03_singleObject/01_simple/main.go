package main

import "sync"

/*
单例模式：在程序的运行中只产生一个实例。

意图：保证一个类仅有一个实例，并提供一个访问它的全局访问点。

主要解决：一个全局使用的类频繁地创建与销毁。

何时使用：当您想控制实例数目，节省系统资源的时候。

如何解决：判断系统是否已经有这个单例，如果有则返回，如果没有则创建。

关键代码：构造函数是私有的。

使用场景：
1. 要求生产唯一序列号。
2. WEB 中的计数器，不用每次刷新都在数据库里加一次，用单例先缓存起来。
3. 创建的一个对象需要消耗的资源过多，比如 I/O 与数据库的连接等

优点：在内存里只有一个实例，减少了内存的开销，尤其是频繁的创建和销毁实例（比如管理学院首页页面缓存）。
		避免对资源的多重占用（比如写文件操作）。
缺点：没有接口，不能继承，与单一职责原则冲突，一个类应该只关心内部逻辑，而不关心外面怎么样来实例化。
 */

// 1.懒汉模式

type singleton1 struct{}
var ins1 *singleton1
func GetIns1() *singleton1{
	if ins1 == nil {
		ins1 = &singleton1{}
	}
	return ins1
}
// 缺点：非协程安全。如果有多个线程同时调用了这个方法， 那么都会检测到instance为nil,就会创建多个对象。


// 2. 恶汉·模式


type singleton2 struct{}
var ins2 *singleton2 = &singleton2{}
func GetIns2() *singleton2{
	return ins2
}
// 缺点：如果singleton创建初始化比较复杂耗时时，加载时间会延长。饿汉模式将在包加载的时候就创建单例对象，当程序中用不到该对象时，浪费了一部分空间
		//和懒汉模式相比，更安全，但是会减慢程序启动速度，所以我们可以在进一步修改程序


// 3. 懒汉加锁

type singleton3 struct{}
var ins3 *singleton3
var mu sync.Mutex
func GetIns3() *singleton3{
	mu.Lock()
	defer mu.Unlock()
	if ins3 == nil {
		ins3 = &singleton3{}
	}
	return ins3
}

// 缺点：虽然解决并发的问题，但每次加锁是要付出代价的

// 4. 双重锁

type singleton4 struct{}
var ins4 *singleton4
var mu4 sync.Mutex
func GetIns4() *singleton4{
	if ins4 == nil {
		mu4.Lock()
		defer mu4.Unlock()
		if ins4 == nil {
			ins4 = &singleton4{}
		}
	}
	return ins4
}

func main()  {

}
