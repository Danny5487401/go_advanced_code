<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [fanIn-fanOut模式](#fanin-fanout%E6%A8%A1%E5%BC%8F)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# fanIn-fanOut模式

以汽车组装为例，汽车生产线上有个阶段是给小汽车装4个轮子，可以把这个阶段任务交给4个人同时去做，这4个人把轮子都装完后，再把汽车移动到生产线下一个阶段。

这个过程中，就有任务的分发，和任务结果的收集。其中任务分发是FAN-OUT，任务收集是FAN-IN

FAN-OUT模式：
    多个goroutine从同一个通道读取数据，直到该通道关闭。OUT是一种张开的模式，所以又被称为扇出，可以用来分发任务。

FAN-IN模式：
    1个goroutine从多个通道读取数据，直到这些通道关闭。IN是一种收敛的模式，所以又被称为扇入，用来收集处理的结果。