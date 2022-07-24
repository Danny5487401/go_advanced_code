package main

// 通用形状接口
type shape interface {
	getType() string
	accept(visitor)
}
