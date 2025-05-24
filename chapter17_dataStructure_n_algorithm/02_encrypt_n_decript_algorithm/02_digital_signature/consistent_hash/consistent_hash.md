<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [一致性hash](#%E4%B8%80%E8%87%B4%E6%80%A7hash)
  - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [算法特点](#%E7%AE%97%E6%B3%95%E7%89%B9%E7%82%B9)
    - [均衡性(Balance)](#%E5%9D%87%E8%A1%A1%E6%80%A7balance)
    - [单调性(Monotonicity)](#%E5%8D%95%E8%B0%83%E6%80%A7monotonicity)
    - [分散性(Spread)](#%E5%88%86%E6%95%A3%E6%80%A7spread)
  - [库](#%E5%BA%93)
    - [github.com/serialx/hashring -->buildkit 使用](#githubcomserialxhashring---buildkit-%E4%BD%BF%E7%94%A8)
    - [go-zero实现](#go-zero%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 一致性hash 
哈希算法是一种从任何一种数据中创建小的数字“指纹”的方法，解决如何将数据映射到固定槽问题。而一致性哈希算法是为了解决当槽数目变化的时候如何将映射结果的变化降到最小

## 使用场景

- 分布式缓存。可以在 redis cluster 这种存储系统上构建一个 cache proxy，自由控制路由。而这个路由规则就可以使用一致性hash算法
- 服务发现
- 分布式调度任务

## 算法特点
### 均衡性(Balance)
均衡性主要指,通过算法分配, 集群中各节点应该要尽可能均衡.

### 单调性(Monotonicity)
单调性主要指当集群发生变化时, 已经分配到老节点的key, 尽可能的任然分配到之前节点,以防止大量数据迁移, 这里一般的hash取模就很难满足这点,而一致性hash算法能够将发生迁移的key数量控制在较低的水平

### 分散性(Spread)
分散性主要针对同一个key, 当在不同客户端操作时,可能存在客户端获取到的缓存集群的数量不一致,从而导致将key映射到不同节点的问题,这会引发数据的不一致性.好的hash算法应该要尽可能避免分散性.

## 库

### github.com/serialx/hashring -->buildkit 使用
使用 kentama 是一致性哈希的一种实现，主要思想如下：

- 将server和key同时映射到环形连续统（0~2^32)
```go
// /Users/python/go/pkg/mod/github.com/serialx/hashring@v0.0.0-20200727003509-22c0c7ab6b1b/hashring.go
// 将服务的key按该hash算法计算,得到在服务在一致性hash环上的位置.
func (h *HashRing) generateCircle() {
	totalWeight := 0
	for _, node := range h.nodes {
		if weight, ok := h.weights[node]; ok {
			totalWeight += weight
		} else {
			totalWeight += 1
			h.weights[node] = 1
		}
	}

	for _, node := range h.nodes {
		weight := h.weights[node]

		for j := 0; j < weight; j++ {
			nodeKey := node + "-" + strconv.FormatInt(int64(j), 10)
			key := h.hashFunc([]byte(nodeKey))
			h.ring[key] = node
			h.sortedKeys = append(h.sortedKeys, key)
		}
	}

	sort.Sort(HashKeyOrder(h.sortedKeys))
}
```
- 为了将key->server，找到第一个比key的映射值大的server的映射值，则key就映射到这台server上，如果没找到，则映射至第一台server
```go
func (h *HashRing) GetNode(stringKey string) (node string, ok bool) {
	pos, ok := h.GetNodePos(stringKey)
	if !ok {
		return "", false
	}
	return h.ring[h.sortedKeys[pos]], true
}

func (h *HashRing) GetNodePos(stringKey string) (pos int, ok bool) {
	if len(h.ring) == 0 {
		return 0, false
	}

	key := h.GenKey(stringKey)

	nodes := h.sortedKeys
	pos = sort.Search(len(nodes), func(i int) bool { return key.Less(nodes[i]) })

	if pos == len(nodes) {
		// Wrap the search, should return First node
		return 0, true
	} else {
		return pos, true
	}
}
```
- 为了平衡性，可以添加一层虚拟节点到物理节点的映射，将key首先映射到虚拟节点，然后再映射到物理节点



Note: 这里使用的是 md5 算法

### go-zero实现
```go
// github.com/zeromicro/go-zero@v1.3.5/core/hash/consistenthash.go

type ConsistentHash struct {
  hashFunc Func       // hash 函数
  replicas int       // 虚拟节点放大因子
  keys     []uint64     // 存储虚拟节点hash
  ring     map[uint64][]interface{}     // 虚拟节点与实际node的对应关系
  nodes    map[string]lang.PlaceholderType // 实际节点存储【便于快速查找，所以使用map】
  lock     sync.RWMutex
}
```