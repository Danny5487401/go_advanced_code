# GC 垃圾回收
GC，全称 GarbageCollection，即垃圾回收，是一种自动内存管理的机制

当程序向操作系统申请的内存不再需要时，垃圾回收主动将其回收并供其他代码进行内存申请时候复用，或者将其归还给操作系统，
这种针对内存级别资源的自动回收过程，即为垃圾回收。而负责垃圾回收的程序组件，即为垃圾回收器.

![](.gc_images/gc_mutator_n_collector.png)

垃圾回收器的执行过程被划分为两个半独立的组件：

- 赋值器（Mutator）：这一名称本质上是在指代用户态的代码。因为对垃圾回收器而言，用户态的代码仅仅只是在修改对象之间的引用关系，也就是在对象图（对象之间引用关系的一个有向图）上进行操作。
    
- 回收器（Collector）：负责执行垃圾回收的代码.需要扫描内存中存活的堆对象，扫描完成后，未被扫描到的对象就是无法访问的堆上垃圾，需要将其占用内存回收掉。

## 基本概念

### 根对象
根对象在垃圾回收的术语中又叫做根集合，它是垃圾回收器在标记过程时最先检查的对象，包括
    
1. 全局变量：程序在编译期就能确定的那些存在于程序整个生命周期的变量。

2. 执行栈：每个 goroutine 都包含自己的执行栈，这些执行栈上包含栈上的变量及指向分配的堆内存区块的指针。

3. 寄存器：寄存器的值可能表示一个指针，参与计算的这些指针可能指向某些赋值器分配的堆内存区块。

### 垃圾

根对象不可达的对象

### 安全点Safe Point
safePoint 安全点顾名思义是指一些特定的位置，当线程运行到这些位置时，线程的一些状态可以被确定(the thread’s representation of it’s Java machine state is well described)，比如记录映射表OopMap的状态，从而确定GC Root的信息，使JVM可以安全的进行一些操作，比如开始GC


### 安全区域
指在一段代码片段之中，引用关系不会发生变化。这个区域中的任意地方开始GC都是安全的、我们也可以把安全区域看做是被扩展了的安全点。

## GC实现方式

所有的 GC 算法其存在形式可以归结为追踪（Tracing）和引用计数（Reference Counting）这两种形式的混合运用.

### 1. 追踪式 GC(可达性分析)
![](.gc_images/tracking_gc.png)

从根对象出发，根据对象之间的引用信息，一步步推进直到扫描完毕整个堆并确定需要保留的对象，从而回收所有可回收的对象。
Go、 Java、V8 对 JavaScript 的实现等均为追踪式 GC。

### 2. 引用计数式 GC
![](.gc_images/ref_gc.png)

每个对象自身包含一个被引用的计数器，当计数器归零时自动得到回收。因为此方法缺陷较多，在追求高性能时通常不被应用。Python、Objective-C 等均为引用计数式 GC。

#### 缺点
- 循环引用，内存泄漏

Python如何解决循环这个问题？

![](.gc_images/python_extend_GC.png)
1. python并没有解决这问题，只是有个双向链表去进行处理，将引用的对象ref减1，然后进行对比识别出来。


## 如何清理
![](.gc_images/GC_conclude_info.png) 

- 标记清除
- 标记复制
- 标记压缩



### 1. 标记清除
![](.gc_images/mark_gc.png)   
将不可达对象放回双向链表 

#### 优化方向

![](.gc_images/bit_mark.png)
- 给每个对象增加个属性标记，但是占用空间大，用**位图标记法**优化。
- 回收后有碎片，使用多个**空闲链表**进行减少碎片分布


### 2. 标记复制
![](.gc_images/mark_copy.png) 
从A区可用对象复制到B区，分配对象在B区，清理完A区，再拷贝回A区。

很明显占用空间大  

### 3. 标记压缩算法
![](.gc_images/mark_zip.png) 

可用对象在一端，可用空间在一端。

很明显整理的频率得高。


## GC优化方向
![](.gc_images/gc_optimazation.png)
- 增量式GC: 允许 collector 分多个小批次执行，每次造成的 mutator 停顿都很小，达到近似实时的效果
- 并发式GC : 利用CPU多核
- 分代GC: JAVA的JVM

### 增量收集器 
增量式（Incremental）的垃圾收集是减少程序最长暂停时间的一种方案，它可以将原本时间较长的暂停时间切分成多个更小的 GC 时间片，虽然从垃圾收集开始到结束的时间更长了，但是这也减少了应用程序暂停的最大时间：

### 并发收集器
并发（Concurrent）的垃圾收集不仅能够减少程序的最长暂停时间，还能减少整个垃圾收集阶段的时间，通过开启读写屏障、利用多核优势与用户程序并行执行，并发垃圾收集器确实能够减少垃圾收集对应用程序的影响

## Go 的 GC
Go 的 GC 目前使用的是无分代（对象没有代际之分）、不整理（回收过程中不对对象进行移动与整理）、并发（与用户代码并发执行）的三色标记清扫算法.

### 原因

1. 对象整理的优势是解决内存碎片问题以及“允许”使用顺序内存分配器。但 Go 运行时的分配算法基于 tcmalloc，基本上没有碎片问题。
并且顺序内存分配器在多线程的场景下并不适用。Go 使用的是基于 tcmalloc 的现代内存分配算法，对对象进行整理不会带来实质性的性能提升。

2. 分代 GC 依赖分代假设，即 GC 将主要的回收目标放在新创建的对象上（存活时间短，更倾向于被回收），而非频繁检查所有对象。
但 Go 的编译器会通过逃逸分析将大部分新生对象存储在栈上（栈直接被回收），只有那些需要长期存在的对象才会被分配到需要进行垃圾回收的堆中。
也就是说，分代 GC 回收的那些存活时间短的对象在 Go 中是直接被分配到栈上，当 goroutine 死亡后栈也会被直接回收，不需要 GC 的参与，进而分代假设并没有带来直接优势。
并且 Go 的垃圾回收器与用户代码并发执行，使得 STW 的时间与对象的代际、对象的 size 没有关系。
   
Go 团队更关注于如何更好地让 GC 与用户代码并发执行（使用适当的 CPU 来执行垃圾回收），而非减少停顿时间这一单一目标上。

### 三色标记法的流程如下 
三色标记法的关键是理解对象的三色抽象以及波面（wavefront）推进这两个概念
    
* 白色对象（可能死亡）：未被回收器访问到的对象。在回收开始阶段，所有对象均为白色，当回收结束后，白色对象均不可达。

* 灰色对象（波面）：已被回收器访问到的对象，但但该对象下的属性没有全被标记完,回收器需要对其中的一个或多个指针进行扫描，因为他们可能还指向白色对象。

* 黑色对象（确定存活）：已被回收器访问到的对象，其中所有字段都已被扫描，黑色对象中任何一个指针都不可能直接指向白色对象。


![参考动图过程](.gc_images/GC_dynamic.gif)

1. 所有对象最开始都是白色.
2. 从 root 开始找到所有可达对象，标记为灰色，放入待处理队列。
3. 遍历灰色对象队列，将其引用对象标记为灰色放入待处理队列，自身标记为黑色。
4. 循环步骤3直到灰色队列为空为止，此时所有引用对象都被标记为黑色，所有不可达的对象依然为白色，白色的就是需要进行回收的对象。

三色标记法相对于普通标记清扫，减少了 STW 时间. 这主要得益于标记过程是 "on-the-fly" 的，在标记过程中是不需要 STW 的，它与程序是并发执行的，这就大大缩短了 STW 的时间.


### GC时为什么要暂停用户线程？
- 首先，如果不暂停用户线程，就意味着期间会不断有垃圾产生，永远也清理不干净。
- 其次，用户线程的运行必然会导致对象的引用关系发生改变，这就会导致两种情况：漏标和错标。

#### 可能存在的问题
1. 多标-浮动垃圾问题

![](.gc_images/float_garbage.png)  

假设 E 已经被标记过了（变成灰色了），此时 D 和 E 断开了引用，按理来说对象 E/F/G 应该被回收的，但是因为 E 已经变为灰色了，其仍会被当作存活对象继续遍历下去。
最终的结果是：这部分对象仍会被标记为存活，即本轮 GC 不会回收这部分内存。

这部分本应该回收 但是没有回收到的内存，被称之为“浮动垃圾”。


解释方式二： 
![](.gc_images/ignore_mark.png)

假设GC已经在遍历对象B了，而此时用户线程执行了A.B=null的操作，切断了A到B的引用。

本来执行了A.B=null之后，B、D、E都可以被回收了，但是由于B已经变为灰色，它仍会被当做存活对象，继续遍历下去。

最终的结果就是本轮GC不会回收B、D、E，留到下次GC时回收，也算是浮动垃圾的一部分。

实际上，这个问题依然可以通过「写屏障」来解决，只要在A写B的时候加入写屏障，记录下B被切断的记录，重新标记时可以再把他们标为白色即可。




2. 漏标-悬挂指针问题

![](.gc_images/float_pointer.png)  

当 GC 线程已经遍历到 E 变成灰色，D变成黑色时，灰色 E 断开引用白色 G ，黑色 D 引用了白色 G。此时切回 GC 线程继续跑，因为 E 已经没有对 G 的引用了，所以不会将 G 放到灰色集合。尽管因为 D 重新引用了 G，但因为 D 已经是黑色了，不会再重新做遍历处理。

最终导致的结果是：G 会一直停留在白色集合中，最后被当作垃圾进行清除。这直接影响到了应用程序的正确性，是不可接受的，这也是 Go 需要在 GC 时解决的问题。


解释方式二： 假设GC线程已经遍历到B了，此时用户线程执行了以下操作：
![](.gc_images/wrong_mark.png)
```go
B.D=null;//B到D的引用被切断A.xx=D;//A到D的引用被建立
```

此时GC线程继续工作，由于B不再引用D了，尽管A又引用了D，但是因为A已经标记为黑色，GC不会再遍历A了，所以D会被标记为白色，最后被当做垃圾回收。
可以看到错标的结果比漏表严重的多，浮动垃圾可以下次GC清理，而把不该回收的对象回收掉，将会造成程序运行错误。

错标只有在满足下面两种情况下才会发生：
![](.gc_images/wrong_mark_situation.png)


### 如何解决上述**漏标**问题
满足三色不变性

![](.gc_images/three_colors.png)
- 强三色不变性： 黑对象不能直接引用白对象
- 弱三色不变性： 黑对象可以引用白对象，但是必须多一条灰色指向白色。

如何满足三色不变性？使用屏障技术（偏向硬件）

#### 内存屏障


内存屏障，是一种屏障指令，它能使CPU或编译器对在该屏障指令之前和之后发出的内存操作强制执行排序约束，在内存屏障前执行的操作一定会先于内存屏障后执行的操作。

垃圾收集中的屏障技术更像是一个钩子方法，它是在用户程序读取对象、创建新对象以及更新对象指针时执行的一段代码，根据操作类型的不同，我们可以将它们分成读屏障（Read barrier）和写屏障（Write barrier）两种，因为读屏障需要在读操作中加入代码片段，对用户程序的性能影响很大，所以编程语言往往都会采用写屏障保证三色不变性。


对于一个不需要对象拷贝的垃圾回收器来说， Read barrier（读屏障）代价是很高的，因为对于这类垃圾回收器来说是不需要保存读操作的版本指针问题。
相对来说 Write barrier（写屏障）代码更小，因为堆中的写操作远远小于堆中的读操作。

#### 写屏障
![](.gc_images/write_barrier.png)

##### Dijkstra 插入屏障--满足强三色

Go 1.7 之前使用的是 Dijkstra Write barrier（写屏障），使用的实现类似下面伪代码：
```go
writePointer(slot, ptr):
    shade(ptr)
    *slot = ptr
```
如果该对象是白色的话，shade(ptr)会将对象标记成灰色。这样可以保证强三色不变性，它会保证 ptr 指针指向的对象在赋值给 *slot 前不是白色。


Dijkstra 插入屏障的好处在于可以立刻开始并发标记，但由于产生了灰色赋值器，缺陷是需要标记终止阶段 STW 时进行重新扫描


![](.gc_images/write_barrier2.png)
在GC进行的过程中，应用程序新建了对象I，此时如果已经标记成黑的对象F引用了对象I，那么在本次GC执行过程中因为黑色对象不会再次扫描，所以如果I着色成白色的话，会被回收掉，这显然是不允许的。

写屏障主要做一件事情，修改原先的写逻辑，然后在对象新增的同时给它着色，并且着色为灰色。因此打开了写屏障可以保证了三色标记法在并发下安全正确地运行。那么有人就会问这些写屏障标记成灰色的对象什么时候回收呢？答案是后续的GC过程中回收，在新的GC过程中所有已存对象就又从白色开始逐步被标记啦。




##### Yuasa 删除写屏障--满足弱三色
Yuasa 在 1990 年的论文 Real-time garbage collection on general-purpose machines 中提出了删除写屏障，因为一旦该写屏障开始工作，它会保证开启写屏障时堆上所有对象的可达，所以也被称作快照垃圾收集（Snapshot GC）

其思想是当赋值器从灰色或白色对象中删除白色指针时，通过写屏障将这一行为通知给并发执行的回收器。

该算法会使用如下所示的写屏障保证增量或者并发执行垃圾收集时程序的正确性，伪代码实现如下：
```
writePointer(slot, ptr)
    shade(*slot)
    *slot = ptr
```


Yuasa 删除屏障的优势则在于不需要标记结束阶段的重新扫描，缺陷是依然会产生丢失的对象，需要在标记开始前对整个对象图进行快照

##### Hybrid write barrier 混合写屏障

Go 在 1.8 的时候使用 Hybrid write barrier（混合写屏障），结合了 Yuasa write barrier 和 Dijkstra write barrier ，实现的伪代码如下：
```go
writePointer(slot, ptr):
    shade(*slot)
    if current stack is grey:
        shade(ptr)
    *slot = ptr
```
混合写屏障的基本思想是：对正在被覆盖的对象进行着色，且如果当前栈未扫描完成， 则同样对指针进行着色。




## GC phase 垃圾收集阶段
![](.gc_images/gc_cycle.png)  
![](.gc_images/gc_cycle2.png)

GC 相关的代码在 runtime/mgc.go 文件下。通过注释介绍我们可以知道 GC 一共分为4个阶段：
```go
const (
	_GCoff             = iota // GC not running; sweeping in background, write barrier disabled
	_GCmark                   // GC marking roots and workbufs: allocate black, write barrier ENABLED
	_GCmarktermination        // GC mark termination: allocate black, P's help GC, write barrier ENABLED
)
```

1. sweep termination（标记准备阶段）
    
    - 会触发 STW ，所有的 P（处理器） 都会进入 safe-point（安全点）；
    
2. the mark phase（标记阶段）
    
    - 将 _GCoff GC 状态 改成 _GCmark，开启 Write Barrier （写入屏障）、mutator assists（协助线程），将根对象入队；
    - 恢复程序执行，mark workers（标记进程）和 mutator assists（协助线程）会开始并发标记内存中的对象。对于任何指针写入和新的指针值，都会被写屏障覆盖，而所有新创建的对象都会被直接标记成黑色；
    - GC 执行根节点的标记，这包括扫描所有的栈、全局对象以及不在堆中的运行时数据结构。扫描goroutine 栈绘导致 goroutine 停止，并对栈上找到的所有指针加置灰，然后继续执行 goroutine。
    - GC 在遍历灰色对象队列的时候，会将灰色对象变成黑色，并将该对象指向的对象置灰；
    - GC 会使用分布式终止算法（distributed termination algorithm）来检测何时不再有根标记作业或灰色对象，如果没有了 GC 会转为mark termination（标记终止）；
    
3. mark termination（标记终止）

   - STW，然后将 GC 阶段转为 _GCmarktermination,关闭 GC 工作线程以及 mutator assists（协助线程）；
   - 执行清理，如 flush mcache；
    
4. the sweep phase（清理阶段）

   - 将 GC 状态转变至 _GCoff，初始化清理状态并关闭 Write Barrier（写入屏障）；
   - 恢复程序执行，从此开始新创建的对象都是白色的；
   - 后台并发清理所有的内存管理单元


问题：
1. 生产速度（申请内存） 大于消费速度（GC）？

解决：在标记开始的时候，收集器会默认抢占 25% 的 CPU 性能，剩下的75%会分配给程序执行。
但是一旦收集器认为来不及进行标记任务了，就会改变这个 25% 的性能分配。这个时候收集器会抢占程序额外的 CPU，这部分被抢占 goroutine 有个名字叫 Mark Assist。
而且因为抢占 CPU的目的主要是 GC 来不及标记新增的内存，那么抢占正在分配内存的 goroutine 效果会更加好，所以分配内存速度越快的 goroutine 就会被抢占越多的资源。





## GC 触发条件

1. 主动触发，通过调用 runtime.GC() 来触发 GC，此调用阻塞式地等待当前 GC 运行完毕。

2. 被动触发，分为两种方式：

    * 使用系统 sysmon 监控，当超过两分钟没有产生任何 GC 时，强制触发 GC。
    
    * 使用步调（Pacing）算法，其核心思想是控制内存增长的比例


## 观察GC方式

1. 
```go
 go build main.go
 GODEBUG=gctrace=1 ./main
```
GODEBUG 变量可以控制运行时内的调试变量，参数以逗号分隔，格式为：name=val。本文着重点在 GC 的观察上，主要涉及 gctrace 参数，
我们通过设置 gctrace=1 后就可以使得垃圾收集器向标准错误流发出 GC 运行信息

2. 
```shell
go tool trace trace.out
```

3. debug.ReadGCStats

4. runtime.ReadMemStats

```shell
gc 18 @17.141s 0%: 0.21+4.8+0.007 ms clock, 1.7+0/0.45/4.8+0.062 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
gc 19 @18.151s 0%: 0.063+1.1+0.003 ms clock, 0.51+0/0.12/1.1+0.030 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
gc 20 @19.157s 0%: 0.12+3.5+0.008 ms clock, 0.98+0/3.6/0.094+0.064 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
```

涉及术语
- mark：标记阶段。
- markTermination：标记结束阶段。
- mutator assist：辅助 GC，是指在 GC 过程中 mutator 线程会并发运行，而 mutator assist 机制会协助 GC 做一部分的工作。
- heap_live：在 Go 的内存管理中，span 是内存页的基本单元，每页大小为 8kb，同时 Go 会根据对象的大小不同而分配不同页数的 span，而 heap_live 就代表着所有 span 的总大小。
- dedicated / fractional / idle：在标记阶段会分为三种不同的 mark worker 模式，分别是 dedicated、fractional 和 idle，它们代表着不同的专注程度，
       其中 dedicated 模式最专注，是完整的 GC 回收行为，fractional 只会干部分的 GC 行为，idle 最轻松。
       这里你只需要了解它是不同专注程度的 mark worker 就好了，详细介绍我们可以等后续的文章。

```shell
gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
```

含义
- 1 gc#：GC 执行次数的编号，每次叠加。
- 2 @#s：自程序启动后到当前的具体秒数。
- 3 #%：自程序启动以来在GC中花费的时间百分比。
- 4 #+...+#：GC 的标记工作共使用的 CPU 时间占总 CPU 时间的百分比。
- 5 #->#-># MB：分别表示 GC 启动时, GC 结束时, GC 活动时的堆大小.
- 6 #MB goal：下一次触发 GC 的内存占用阈值。
- 7 #P：当前使用的处理器 P 的数量

案例代码
```go

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
```
执行结果分析
```shell
gc 7 @0.140s 1%: 0.031+2.0+0.042 ms clock, 0.12+0.43/1.8/0.049+0.17 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
```


- gc 7：第 7 次 GC。
- @0.140s：当前是程序启动后的 0.140s。
- 1%：程序启动后到现在共花费 1% 的时间在 GC 上。
- 0.031+2.0+0.042 ms clock：
    0.031：表示单个 P 在 mark 阶段的 STW 时间。
    2.0：表示所有 P 的 mark concurrent（并发标记）所使用的时间。
    0.042：表示单个 P 的 markTermination 阶段的 STW 时间。

- 0.12+0.43/1.8/0.049+0.17 ms cpu：
    0.12：表示整个进程在 mark 阶段 STW 停顿的时间。
    0.43/1.8/0.049：0.43 表示 mutator assist 占用的时间，1.8 表示 dedicated + fractional 占用的时间，0.049 表示 idle 占用的时间。
    0.17ms：0.17 表示整个进程在 markTermination 阶段 STW 时间。

- 4->4->1 MB：
    4：表示开始 mark 阶段前的 heap_live 大小。
    4：表示开始 markTermination 阶段前的 heap_live 大小。
    1：表示被标记对象的大小。
- 5 MB goal：表示下一次触发 GC 回收的阈值是 5 MB。
- 4 P：本次 GC 一共涉及多少个 P。

- wall clock 是指开始执行到完成所经历的实际时间，包括其他程序和本程序所消耗的时间；cpu time 是指特定程序使用 CPU 的时间；他们存在以下关系：

    - wall clock < cpu time: 充分利用多核

    - wall clock ≈ cpu time: 未并行执行

    - wall clock > cpu time: 多核优势不明显



## 参考资料
1. [Go中内存分配源码实现](https://www.luozhiyun.com/archives/434) 
2. [Go语言设计](https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/)