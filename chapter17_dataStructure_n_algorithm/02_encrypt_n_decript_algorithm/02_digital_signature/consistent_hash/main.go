package main

import (
	"fmt"
	"github.com/serialx/hashring"
)

func main() {
	memcacheServers := []string{
		"192.168.0.246:11212",
		"192.168.0.247:11212",
		"192.168.0.249:11212"}

	// 不带权重
	ring := hashring.New(memcacheServers)
	// 删除节点
	ring = ring.RemoveNode("192.168.0.246:11212")
	// 添加节点
	ring = ring.AddNode("192.168.0.250:11212")
	server, _ := ring.GetNode("my_key")
	fmt.Println(server)

}
