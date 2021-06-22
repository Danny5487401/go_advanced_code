package main

/*
注: sync.Pool的定位不是做类似连接池的东西，它的用途仅仅是增加对象重用的几率，减少gc的负担
背景：Go是自动垃圾回收的(garbage collector)，这大大减少了程序编程负担。但gc是一把双刃剑，带来了编程的方便但同时也增加了运行时开销，
	使用不当甚至会严重影响程序的性能。因此性能要求高的场景不能任意产生太多的垃圾（有gc但又不能完全依赖它挺恶心的），如何解决呢？
	那就是要重用对象了，我们可以简单的使用一个chan把这些可重用的对象缓存起来，但如果很多goroutine竞争一个chan性能肯定是问题

*/

import (
	"fmt"
	"sync"
)

func main() {
	p := &sync.Pool{
		New: func() interface{} {
			return 0
		},
	}

	a := p.Get().(int)
	p.Put(1)

	//runtime.GC() //注意注释前后的结果

	b := p.Get().(int)
	fmt.Println(a, b)

}

/*
上面我们可以看到pool创建的时候是不能指定大小的，所有sync.Pool的缓存对象数量是没有限制的（只受限于内存），
	因此使用sync.pool是没办法做到控制缓存对象数量的个数的。另外sync.pool缓存对象的期限是很诡异的，这是很多人错误理解的地方，

源码分析
	sync.Pool 的 init 函数
		func init() {
			runtime_registerPoolCleanup(poolCleanup)
		}
		func runtime_registerPoolCleanup(cleanup func())
	可以看到pool包在init的时候注册了一个poolCleanup函数，它会清除所有的pool里面的所有缓存的对象，该函数注册进去之后会在每次gc之前都会调用，
		因此sync.Pool缓存的期限只是两次gc之间这段时间

	正因为这样，我们是不可以使用sync.Pool去实现一个socket连接池的。

看sync_pool_structure.png流程
	一个goroutine固定在一个局部调度器P上，从当前 P 对应的 poolLocal 取值， 若取不到，则从对应的 shared 数组上取，若还是取不到；
	则尝试从其他 P 的 shared 中偷。 若偷不到，则调用 New 创建一个新的对象。池中所有临时对象在一次 GC 后会被全部清空。


*/
