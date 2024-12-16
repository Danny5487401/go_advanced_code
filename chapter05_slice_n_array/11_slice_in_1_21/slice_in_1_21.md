<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [切片泛型库](#%E5%88%87%E7%89%87%E6%B3%9B%E5%9E%8B%E5%BA%93)
  - [分类](#%E5%88%86%E7%B1%BB)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 切片泛型库


在 Go 1.21.0 版本中，引入了 切片泛型库，它提供了很多有用的函数，特别是在搜索、查找和排序等方面

## 分类

slices 库包含的函数可以分为以下类型：

- 搜索：通过二分查找算法查找指定元素。相关的函数有 BinarySearch 和 BinarySearchFunc
- 裁剪：删除切片中未使用的容量。相关的函数有 Clip
- 克隆：浅拷贝一个切片副本。相关的的函数有：Clone
- 压缩：将切片里连续的相同元素替换为一个元素。从而减少了切片的长度，相关的函数有：Compact 和 CompactFunc
- 大小比较：比较两个切片的大小。相关的函数有 Compare 和 CompareFunc
- 包含：判断切片是否包含指定元素。相关的函数有：Contains 和 ContainsFunc
- 删除：从切片中删除一个或多个元素。相关的函数有 Delete 和 DeleteFunc
- 等价比较：比较两个切片是否相等。相关的函数有：Equal 和 EqualFunc
- 扩容：增加切片的容量。相关的函数有：Grow
- 索引查找：查找指定元素在切片中的索引位置。相关的函数有：Index 和 IndexFunc
- 插入：往切片里插入一组值。相关的函数有：Insert
- 有序判断：判断切片是否按照升序排列。相关的函数有：IsSorted 和 IsSortedFunc
- 最大值：查找切片里的最大元素。相关的函数有：Max 和 MaxFunc
- 最小值：查找切片里的最小元素。相关的函数有：Min 和 MinFunc
- 替换：替换切片里的元素。相关的函数有：Replace
- 反转：反转切片的元素。相关的函数有：Reverse
- 排序：对切片里的元素进行升序排列。相关的函数有：Sort 和 SortFunc 以及 SortStableFunc




## 参考

- [Go Slices 切片泛型库](https://www.51cto.com/article/772183.html)