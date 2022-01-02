package main

import "sync"

// 1.懒汉模式

type singleton1 struct{}

var ins1 *singleton1

func GetIns1() *singleton1 {
	if ins1 == nil {
		ins1 = &singleton1{}
	}
	return ins1
}

// 缺点：非协程安全。如果有多个线程同时调用了这个方法， 那么都会检测到instance为nil,就会创建多个对象。

// 2. 恶汉·模式

type singleton2 struct{}

var ins2 *singleton2 = &singleton2{}

func GetIns2() *singleton2 {
	return ins2
}

// 缺点：如果singleton创建初始化比较复杂耗时时，加载时间会延长。饿汉模式将在包加载的时候就创建单例对象，当程序中用不到该对象时，浪费了一部分空间
//和懒汉模式相比，更安全，但是会减慢程序启动速度，所以我们可以在进一步修改程序

// 3. 懒汉加锁

type singleton3 struct{}

var ins3 *singleton3
var mu sync.Mutex

func GetIns3() *singleton3 {
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

func GetIns4() *singleton4 {
	if ins4 == nil {
		mu4.Lock()
		defer mu4.Unlock()
		if ins4 == nil {
			ins4 = &singleton4{}
		}
	}
	return ins4
}

func main() {

}
