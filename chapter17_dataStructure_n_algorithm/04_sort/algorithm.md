# 算法分类
![](.sort_images/sort_category.png)
我们常见的排序算法可以分为两大类：

- 比较类排序：通过比较来决定元素间的相对次序，由于其时间复杂度不能突破O(nlogn)，因此也称为非线性时间比较类排序。
- 非比较类排序：不通过比较来决定元素间的相对次序，它可以突破基于比较排序的时间下界，以线性时间运行，因此也称为线性时间非比较类排序。

```css
排序算法	时间复杂度(平均)	时间复杂度(最坏)	时间复杂度(最优)	空间复杂度           稳定性

冒泡排序	O(𝑛2)	        O(𝑛2)	        O(𝑛)	        O(1)                稳定
快速排序	O(nlogn)        O(𝑛2)	        O(nlogn)        O(nlogn)～O(n)       不稳定
插入排序	O(𝑛2)	        O(𝑛2)	        O(𝑛)	        O(1)                稳定
希尔排序	O(nlogn)~O(𝑛2)	O(𝑛2)	        O(𝑛1.3)	        O(1)                不稳定
选择排序	O(𝑛2)	        O(𝑛2)	        O(𝑛2)	        O(1)                稳定
堆排序	O(nlogn)        O(nlogn)        O(nlogn)        O(1)                不稳定
归并排序	O(nlogn)        O(nlogn)        O(nlogn)        O(n)	            稳定

计数排序	O(n+k)          O(n+k)          O(n+k)          O(k)                稳定
桶排序	O(n+k)          O(𝑛2)           O(𝑛2)           O(n+k)                稳定
基数排序	O(n*k)]         O(n*k)          O(n*k)          O(n+k)              稳定
```
## 1.冒泡排序
![](.sort_images/bubble.gif)

步骤

1. 比较相邻的元素。如果第一个比第二个大，就交换他们两个。
2. 对每一对相邻元素作同样的工作，从开始第一对到结尾的最后一对。这步做完后，最后的元素会是最大的数。
3. 针对所有的元素重复以上的步骤，除了最后一个。
4. 持续每次对越来越少的元素重复上面的步骤，直到没有任何一对数字需要比较

## 2.快速排序
![](.sort_images/quick_sort.gif)

步骤

1. 从数列中挑出一个元素，称为 “基准”（pivot）;
2. 重新排序数列，所有元素比基准值小的摆放在基准前面，所有元素比基准值大的摆在基准的后面（相同的数可以到任一边）。在这个分区退出之后，该基准就处于数列的中间位置。这个称为分区（partition）操作；
3. 递归地（recursive）把小于基准值元素的子数列和大于基准值元素的子数列排序；


## Golang中 sort包内部实现了四种基本的排序算法

1. 插入排序insertionSort
2. 归并排序symMerge
3. 堆排序heapSort
4. 快速排序quickSort
```go

// 快速排序
func quickSort(data Interface, a, b, maxDepth int)

```


```go
// 插入排序
func insertionSort(data Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

```

```go
// 归并排序
func symMerge(data Interface, a, m, b int) {
	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[a] into data[m:b]
	// if data[a:m] only contains one element.
	if m-a == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] >= data[a] for m <= i < b.
		// Exit the search loop with i == b in case no such index exists.
		i := m
		j := b
		for i < j {
			h := int(uint(i+j) >> 1)
			if data.Less(h, a) {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[a] reaches the position before i.
		for k := a; k < i-1; k++ {
			data.Swap(k, k+1)
		}
		return
	}

	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[m] into data[a:m]
	// if data[m:b] only contains one element.
	if b-m == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] > data[m] for a <= i < m.
		// Exit the search loop with i == m in case no such index exists.
		i := a
		j := m
		for i < j {
			h := int(uint(i+j) >> 1)
			if !data.Less(m, h) {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[m] reaches the position i.
		for k := m; k > i; k-- {
			data.Swap(k, k-1)
		}
		return
	}

	mid := int(uint(a+b) >> 1)
	n := mid + m
	var start, r int
	if m > mid {
		start = n - b
		r = mid
	} else {
		start = a
		r = m
	}
	p := n - 1

	for start < r {
		c := int(uint(start+r) >> 1)
		if !data.Less(p-c, c) {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < m && m < end {
		rotate(data, start, m, end)
	}
	if a < start && start < mid {
		symMerge(data, a, start, mid)
	}
	if mid < end && end < b {
		symMerge(data, mid, end, b)
	}
}
```
```go
// 堆排序
func heapSort(data Interface, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}
```


sort包内置的四种排序方法是不公开的，只能被用于sort包内部使用。因此，对数据集合排序时， 不必考虑应当选择哪一种，只需要实现sort.Interface接口定义三个接口即可
```go
type Interface interface{
    Len() int //返回集合中的元素个数
    Less(i,j int) bool//i>j 返回索引i和元素是否比索引j的元素小
    Swap(i,j int)//交换i和j的值
}
```
这里其实隐含要求这个容器或数据集合是slice类型或Array类型。否则，没法按照索引号取值
逆序:sort包提供了Reverse()方法，允许将数据按Less()定义的排序方式逆序排序，而无需修改Less()代码。

Note：Go的sort包已经为基本数据类型都实现了sort功能，其函数名的最后一个字母是s，表示sort之意。比如：Ints, Float64s, Strings，等等。

