package map_gc_trace

import (
	"runtime"
	"time"
)

/*
GODEBUG 变量可以控制运行时内的调试变量，参数以逗号分隔，格式为：name=val。本文着重点在 GC 的观察上，主要涉及 gctrace 参数，
我们通过设置 gctrace=1 后就可以使得垃圾收集器向标准错误流发出 GC 运行信息

gc 18 @17.141s 0%: 0.21+4.8+0.007 ms clock, 1.7+0/0.45/4.8+0.062 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
gc 19 @18.151s 0%: 0.063+1.1+0.003 ms clock, 0.51+0/0.12/1.1+0.030 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
gc 20 @19.157s 0%: 0.12+3.5+0.008 ms clock, 0.98+0/3.6/0.094+0.064 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
涉及术语
	mark：标记阶段。
	markTermination：标记结束阶段。
	mutator assist：辅助 GC，是指在 GC 过程中 mutator 线程会并发运行，而 mutator assist 机制会协助 GC 做一部分的工作。
	heap_live：在 Go 的内存管理中，span 是内存页的基本单元，每页大小为 8kb，同时 Go 会根据对象的大小不同而分配不同页数的 span，而 heap_live 就代表着所有 span 的总大小。
	dedicated / fractional / idle：在标记阶段会分为三种不同的 mark worker 模式，分别是 dedicated、fractional 和 idle，它们代表着不同的专注程度，
		其中 dedicated 模式最专注，是完整的 GC 回收行为，fractional 只会干部分的 GC 行为，idle 最轻松。
		这里你只需要了解它是不同专注程度的 mark worker 就好了，详细介绍我们可以等后续的文章。

gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
含义
	gc#：GC 执行次数的编号，每次叠加。
	@#s：自程序启动后到当前的具体秒数。
	#%：自程序启动以来在GC中花费的时间百分比。
	#+...+#：GC 的标记工作共使用的 CPU 时间占总 CPU 时间的百分比。
	#->#-># MB：分别表示 GC 启动时, GC 结束时, GC 活动时的堆大小.
	#MB goal：下一次触发 GC 的内存占用阈值。
	#P：当前使用的处理器 P 的数量
案例
gc 7 @0.140s 1%: 0.031+2.0+0.042 ms clock, 0.12+0.43/1.8/0.049+0.17 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
gc 7：第 7 次 GC。
@0.140s：当前是程序启动后的 0.140s。
1%：程序启动后到现在共花费 1% 的时间在 GC 上。
0.031+2.0+0.042 ms clock：
	0.031：表示单个 P 在 mark 阶段的 STW 时间。
	2.0：表示所有 P 的 mark concurrent（并发标记）所使用的时间。
	0.042：表示单个 P 的 markTermination 阶段的 STW 时间。

0.12+0.43/1.8/0.049+0.17 ms cpu：
	0.12：表示整个进程在 mark 阶段 STW 停顿的时间。
	0.43/1.8/0.049：0.43 表示 mutator assist 占用的时间，1.8 表示 dedicated + fractional 占用的时间，0.049 表示 idle 占用的时间。
	0.17ms：0.17 表示整个进程在 markTermination 阶段 STW 时间。

4->4->1 MB：
	4：表示开始 mark 阶段前的 heap_live 大小。
	4：表示开始 markTermination 阶段前的 heap_live 大小。
	1：表示被标记对象的大小。
5 MB goal：表示下一次触发 GC 回收的阈值是 5 MB。
4 P：本次 GC 一共涉及多少个 P。
*/
const capacity = 50000

var d interface{}

func main() {
	//value值方式
	//d = Value()

	// value指针方式
	d = pointer()
	for i := 0; i < 20; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}

}

func Value() interface{} {
	m := make(map[int]int, capacity)
	for i := 0; i < capacity; i++ {
		m[i] = i
	}
	return m
}

func pointer() interface{} {
	m := make(map[int]*int, capacity)
	for i := 0; i < capacity; i++ {
		v := i
		m[i] = &v
	}
	return m
}
