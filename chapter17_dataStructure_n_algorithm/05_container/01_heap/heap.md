<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [container/heap](#containerheap)
  - [heap 的应用](#heap-%E7%9A%84%E5%BA%94%E7%94%A8)
  - [完全二叉树的存储方式](#%E5%AE%8C%E5%85%A8%E4%BA%8C%E5%8F%89%E6%A0%91%E7%9A%84%E5%AD%98%E5%82%A8%E6%96%B9%E5%BC%8F)
  - [最大堆图解操作](#%E6%9C%80%E5%A4%A7%E5%A0%86%E5%9B%BE%E8%A7%A3%E6%93%8D%E4%BD%9C)
    - [1 添加元素: 从下到上堆化](#1-%E6%B7%BB%E5%8A%A0%E5%85%83%E7%B4%A0-%E4%BB%8E%E4%B8%8B%E5%88%B0%E4%B8%8A%E5%A0%86%E5%8C%96)
    - [2 拿出元素: 堆顶向下堆化](#2-%E6%8B%BF%E5%87%BA%E5%85%83%E7%B4%A0-%E5%A0%86%E9%A1%B6%E5%90%91%E4%B8%8B%E5%A0%86%E5%8C%96)
  - [container/heap 包](#containerheap-%E5%8C%85)
    - [添加元素 push](#%E6%B7%BB%E5%8A%A0%E5%85%83%E7%B4%A0-push)
    - [去除元素:Pop 操作](#%E5%8E%BB%E9%99%A4%E5%85%83%E7%B4%A0pop-%E6%93%8D%E4%BD%9C)
    - [额外方法](#%E9%A2%9D%E5%A4%96%E6%96%B9%E6%B3%95)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# container/heap

堆一般指二叉堆。是使用完全二叉树这种数据结构构建的一种实际应用。通过它的特性，分为最大堆和最小堆两种。
完全二叉树是指除叶子节点，所有层级是满节点，叶子节点从左向右排列填满。

![](.heap_images/heap_max_and_min.png)

- 最小堆就是在这颗二叉树中，任何一个节点的值比其所在子树的任意一个节点都要小。所以，根节点就是 heap 中最小的值.

- 最大堆就是在这颗二叉树中，任何一个节点的值都比起所在子树的任意一个节点值都要大。


首先要将所有的元素构建为一个完全二叉树。

在一个完全二叉树中，将数据重新按照堆的的特性排列，就可以将完全二叉树变成一个堆。这个过程叫做“堆化”。

在堆中，我们要删除一个元素一般从堆顶删除（可以取到最大值/最小值）。
删除之后，数据集就不能算作一个堆了，因为最顶层的元素没有了，数据集不符合完全二叉树的定义。这时，我们需要将堆的数据进行重新排列，也就是重新“堆化”。
同样的，在堆中新添加一个元素也需要重新做“堆化”的操作，来将数据集恢复到满足堆定义的状态。

## heap 的应用
- 定时器
- 优先级队列：比如kubernetes中的实现，FIFO-PriorityQueue
- heap 排序

## 完全二叉树的存储方式

对于二叉树来说，存储方式有2种，一种使用数组的形式来存储，一种使用链表的方式存储。

- 链表的方式相对浪费存储空间，因为要存储左右子树的指针，但扩缩容方便。
- 数组更加节省空间，更加方便定位节点，缺点则是扩缩容不便。



## 最大堆图解操作

![](.heap_images/full_x_tree.png)
这 index= 1开始, 0 没有存储数据，左节点 2 * i, 右节点 2*i +1, 父节点就是 i/2


一个最大堆，【插入】和【弹出】这两个能力，就需要做“堆化”，使得堆满足定义

### 1 添加元素: 从下到上堆化
![](.heap_images/heap_push.png)

```go
func (h *Heap) downToUpHeapify(pos int) {
    for pos / 2 > 0 && h.data[pos/2].Less(h.data[pos]) { // 如果存在父节点 & 值大于父节点
        h.swap(pos, pos/2) // 交换两个值的位置
        pos = pos /2 // 将操作节点变为父节点的位置
    }
}
```

### 2 拿出元素: 堆顶向下堆化

![](.heap_images/heap_pop.png)
```go
// 从上到下堆化
func (h *Heap) upToDownHeapify() {
    max := h.len
    i := 1
    pos := i
    for {
        if i * 2 <= max && h.data[i].Less(h.data[i*2]) { // 如果有左子树，且自己小于左子树
            pos = i*2 
        }

        if i *2 +1 <= max && h.data[pos].Less(h.data[i*2+1]) { // 如果有右子树，且自己小于右子树
            pos = i*2+1
        }
        if pos == i { // 如果位置没有变化，说明堆化结束
            break
        }

        h.swap(i, pos) // 交换当前位置和下一个位置的内容
        i = pos // 操作下一个位置
    }
}

```


## container/heap 包
Golang 的实现中，索引 0 是存储了数据的。 左节点 2 * i+1, 右节点 2*i +2 。可以实现最小堆和最大堆，通过Less实现。


下面以最小堆作为案例: 
```go
// go1.23.0/src/container/heap/example_intheap_test.go

// An IntHeap is a min-heap of ints.
type IntHeap []int
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
```

初始化
```go
func Init(h Interface) {
	// heapify
	n := h.Len()
	// i = 最后一个非叶子节点的 index； i >= 堆顶； index 自减
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}
```

```go
// 从上到下堆化
func down(h Interface, i0, n int) bool {
	i := i0 // 堆顶 index
	for {
		j1 := 2*i + 1  // 左孩子 index
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break // 堆化结束
		}
		j := j1 // 为了找出左右孩子中较小值的下标,初始化为左孩子
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) { // j2 = 右孩子；j2 小于堆长度 && 右孩子小于左孩子
			j = j2 // j = 2*i + 2 = 右孩子 
		}
		if !h.Less(j, i) { // 如果 堆顶i小于 j , 堆化结束
			break
		}
		h.Swap(i, j) // 交换堆顶元素和 j
		i = j // 切换到下一个操作 index
	}
    // 返回 元素是否有移动
    // 此处是一个特殊设计，用来判断向下堆化是否真的有操作
    // 当删除中间的元素时，如果向下堆化没有操作的话，就需要再做向上堆化
	return i > i0
}
```
从指定下标的节点开始，计算它左孩子的下标。如果左孩子下标没有超出底层数组的范围，那么找出左右孩子中较小值的下标（如果存在右孩子）。
如果父节点的值已经比较小的子节点值还小，那么无需继续调整。否则将较小的子节点与父节点进行交换，继续判断调整后的父节点是否满足最小堆条件





### 添加元素 push

```go
func Push(h Interface, x any) {
	h.Push(x) // 向数据集添加一个元素
	up(h, h.Len()-1) // 从下向上堆化
}

```

```go
// 从下向上堆化的内容
func up(h Interface, j int) { // h 表示堆，j 代表需要堆化的元素 index
	for {
		i := (j - 1) / 2 // 定义 j 的父 index
		if i == j || !h.Less(j, i) { // 如果两个元素相等 或者 父元素小于当前元素
			break
		}
		h.Swap(i, j) // 交换父元素和当前元素
		j = i // index 变为父元素的 index
	}
}
```
从指定下标处开始向上调整堆，首先找到这个节点的父节点，如果当前节点已经是根节点或者父节点的值已经比当前节点的值小，则直接返回，停止调整。
否则，将该节点与父节点交换位置，并从父节点开始继续向上调整，直至根节点


###  去除元素:Pop 操作

```go
func Pop(h Interface) any {
	n := h.Len() - 1
	h.Swap(0, n) // 交换堆顶和最后一个元素
	down(h, 0, n) // 从上到下堆化
	return h.Pop() // 弹出最后一个元素
}
```


### 额外方法

- Remove:从堆中删除并返回索引 i 处的元素。
```go
func Remove(h Interface, i int) any {
	n := h.Len() - 1 // Note: 表示 h.Len() - 1 表示堆长度
	if n != i { // 如果不是堆顶
		h.Swap(i, n) // 交换 需要删除的元素 和 最后一个元素
		if !down(h, i, n) { // 从上到下堆化
			up(h, i) // 如果没有成功，就从下向上堆化
		}
	}
	return h.Pop() // 如果待弹出的元素正好就是最后一个元素（只有一个元素的堆），那么直接弹出即可
}
```
- Fix: 在索引 i 处的元素更改其值后重新建立堆排序。
```go
func Fix(h Interface, i int) {
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}
```

## 参考

- [Golang Heap 源码剖析](https://www.cnblogs.com/reposkeeper-wx/p/golang-heap-yuan-ma-pou-xi.html)
- [数组中的第K个最大元素--leetcode 图解](https://leetcode.cn/problems/kth-largest-element-in-an-array/solutions/307351/shu-zu-zhong-de-di-kge-zui-da-yuan-su-by-leetcod-2/?envType=study-plan-v2&envId=top-100-liked)
- [Golang 源码阅读系列---container](https://juejin.cn/post/7414371017252913179)

