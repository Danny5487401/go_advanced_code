<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [实现无限缓存的channel](#%E5%AE%9E%E7%8E%B0%E6%97%A0%E9%99%90%E7%BC%93%E5%AD%98%E7%9A%84channel)
  - [特点](#%E7%89%B9%E7%82%B9)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 实现无限缓存的channel


Go语言的channel有两种类型，一种是无缓存的channel，一种是有缓存的buffer，这两种类型的channel大家都比较熟悉了，但是对于有缓存的channel,它的缓存长度在创建channel的时候就已经确定了，中间不能扩缩容，这导致在一些场景下使用有问题，或者说不太适合特定的场景。

## 特点
- 不会阻塞write。 它总是能处理write的数据，或者放入到待读取的channel中，或者放入到缓存中
- 无数据时read会被阻塞。当没有可读的数据时，从channel中读取的goroutine会被阻塞
- 读写都是通过channel操作。 内部的缓存不会暴露出来
- 能够查询当前待读取的数据数量。因为缓存中可能也有待处理的数据，所以需要返回len(buffer)+len(chan)
- 关闭channel后，还未读取的channel还是能够被读取，读取完之后才能发现channel已经完毕。这和正常的channel的逻辑是一样的，这种情况叫"drain"未读的数据