package _4_self_HashMap

import (
	"fmt"
	"testing"
)

func TestHashMap(t *testing.T) {

	// 新建一个哈希表
	hashMap := NewHashMap(16)

	// 放35个值
	for i := 0; i < 35; i++ {
		hashMap.Put(fmt.Sprintf("%d", i), fmt.Sprintf("v%d", i))
	}
	fmt.Println("cap:", hashMap.Capacity(), "len:", hashMap.Len())

	// 打印全部键值对
	hashMap.Range()

	key := "4"
	value, ok := hashMap.Get(key)
	if ok {
		fmt.Printf("get '%v'='%v'\n", key, value)
	} else {
		fmt.Printf("get %v not found\n", key)
	}

	// 删除键
	hashMap.Delete(key)
	fmt.Println("after delete cap:", hashMap.Capacity(), "len:", hashMap.Len())
	value, ok = hashMap.Get(key)
	if ok {
		fmt.Printf("get '%v'='%v'\n", key, value)
	} else {
		fmt.Printf("get %v not found\n", key)
	}

}
