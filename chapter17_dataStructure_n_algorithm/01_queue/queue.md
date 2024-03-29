<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [queue 队列](#queue-%E9%98%9F%E5%88%97)
  - [实现方式](#%E5%AE%9E%E7%8E%B0%E6%96%B9%E5%BC%8F)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# queue 队列

我们日常生活中，都需要将物品排列，或者安排事情的先后顺序。更通俗地讲，我们买东西时，人太多的情况下，我们要排队，排队也有先后顺序，有些人早了点来，排完队就离开了，有些人晚一点，才刚刚进去人群排队。


在计算机的世界里，会经常听见两种结构，栈（stack） 和 队列 (queue)。它们是一种收集数据的有序集合（Collection），只不过删除和访问数据的顺序不同。

- 栈：先进后出，先进队的数据最后才出来。在英文的意思里，stack 可以作为一叠的意思，这个排列是垂直的，你将一张纸放在另外一张纸上面，先放的纸肯定是最后才会被拿走，因为上面有一张纸挡住了它。
- 队列：先进先出，先进队的数据先出来。在英文的意思里，queue 和现实世界的排队意思一样，这个排列是水平的，先排先得

## 实现方式

我们可以用数据结构：链表（可连续或不连续的将数据与数据关联起来的结构），或 数组（连续的内存空间，按索引取值） 来实现 栈（stack） 和 队列 (queue)。

- 数组实现：能快速随机访问存储的元素，通过下标 index 访问，支持随机访问，查询速度快，但存在元素在数组空间中大量移动的操作，增删效率低。

- 链表实现：只支持顺序访问，在某些遍历操作中查询速度慢，但增删元素快。

