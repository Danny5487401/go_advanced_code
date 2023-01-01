package main

import (
	"fmt"
	"github.com/OneOfOne/xxhash"
)

// 将一个键进行Hash
func XXHash(key []byte) uint64 {
	h := xxhash.New64()
	h.Write(key)
	return h.Sum64()
}

func main() {
	keys := []string{"hi", "my", "friend", "I", "love", "you", "my", "apple"}
	var length uint64 = 8
	var k int64 = 3
	for _, key := range keys {
		value := XXHash([]byte(key))
		fmt.Printf("xxhash('%s')=%d\n", key, value)
		// xxhash(key) % 8 = 0，1，2，3，4，5，6，7
		fmt.Printf("xxhash('%s')=%v\n", key, value%length)
		// 恒等式 hash % 2^k = hash & (2^k-1)，表示截断二进制的位数，保留后面的 k 位
		fmt.Printf("xxhash('%s')=%v\n", key, value&(1<<k-1))
	}
}
