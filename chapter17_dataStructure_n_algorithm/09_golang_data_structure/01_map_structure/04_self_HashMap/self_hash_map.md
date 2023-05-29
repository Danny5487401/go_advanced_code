<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [自我实现 hashmap](#%E8%87%AA%E6%88%91%E5%AE%9E%E7%8E%B0-hashmap)
  - [功能](#%E5%8A%9F%E8%83%BD)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 自我实现 hashmap

## 功能
- 初始化：新建一个 2^x 个长度的数组，一开始 x 较小。
- 添加键值：进行 hash(key) & (2^x-1)，定位到数组下标，查找数组下标对应的链表，如果链表有该键，更新其值，否则追加元素。
- 获取键值：进行 hash(key) & (2^x-1)，定位到数组下标，查找数组下标对应的链表，如果链表不存在该键，返回 false，否则返回该值以及 true。
- 删除键值：进行 hash(key) & (2^x-1)，定位到数组下标，查找数组下标对应的链表，如果链表不存在该键，直接返回，否则删除该键。
- 进行键值增删时如果数组容量太大或者太小，需要相应缩容或扩容。