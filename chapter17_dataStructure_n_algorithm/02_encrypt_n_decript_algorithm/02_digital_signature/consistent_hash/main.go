package main

import (
	"fmt"
	"github.com/serialx/hashring"
)

func main() {
	memcacheServers := []string{"192.168.0.246:11212",
		"192.168.0.247:11212",
		"192.168.0.249:11212"}

	ring := hashring.New(memcacheServers)
	ring = ring.RemoveNode("192.168.0.246:11212")
	ring = ring.AddNode("192.168.0.250:11212")
	server, _ := ring.GetNode("my_key")
	fmt.Println(server)

}
