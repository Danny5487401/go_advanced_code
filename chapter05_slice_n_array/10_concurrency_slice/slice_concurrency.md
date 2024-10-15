<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [并发访问 slice 如何做到优雅和安全](#%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE-slice-%E5%A6%82%E4%BD%95%E5%81%9A%E5%88%B0%E4%BC%98%E9%9B%85%E5%92%8C%E5%AE%89%E5%85%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# 并发访问 slice 如何做到优雅和安全
slice是对数组一个连续片段的引用，当slice长度增加的时候，可能底层的数组会被换掉。

当出在换底层数组之前，切片同时被多个goroutine拿到，并执行append操作。
那么很多goroutine的append结果会被覆盖，导致n个gouroutine append后，长度小于n
