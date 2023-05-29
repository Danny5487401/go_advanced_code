<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Golang的并发核心思路](#golang%E7%9A%84%E5%B9%B6%E5%8F%91%E6%A0%B8%E5%BF%83%E6%80%9D%E8%B7%AF)
  - [需求：](#%E9%9C%80%E6%B1%82)
    - [方式一](#%E6%96%B9%E5%BC%8F%E4%B8%80)
    - [方式二：](#%E6%96%B9%E5%BC%8F%E4%BA%8C)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Golang的并发核心思路
Golang并发核心思路是关注数据流动。数据流动的过程交给channel，数据处理的每个环节都交给goroutine，把这些流程画起来，有始有终形成一条线，那就能构成流水线模型。

## 需求：
计算一个整数切片中元素的平方值并把它打印出来。
### 方式一
非并发的方式是使用for遍历整个切片，然后计算平方，打印结果。

### 方式二：
使用流水线模型实现这个简单的功能，从流水线的角度，可以分为3个阶段：

1. 遍历切片，这是生产者。
2. 计算平方值。
3. 打印结果，这是消费者。