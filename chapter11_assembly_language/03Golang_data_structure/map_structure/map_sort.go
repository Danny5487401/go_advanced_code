package main

import (
	"fmt"
	"sort"
)

// 排序及传参
func main() {
	// 初始化
	// 排序：无序变有序
	mapScore := map[string]int{
		"Tom":   78,
		"Mary":  60,
		"Kevin": 90,
		"Danny": 100,
	}
	len := len(mapScore)
	names := make([]string, 0, len)
	for key, _ := range mapScore {
		names = append(names, key)
	}
	// 按字母升序
	sort.Strings(names)
	for _, name := range names {
		fmt.Println("Key:", name, "Value:", mapScore[name])
	}

	// 传递
	fmt.Println("传递前Tom的值", mapScore["Tom"])
	modify(mapScore)
	fmt.Println("传递后Tom的值", mapScore["Tom"])
	// 运行结果是可以看出key为"Tom"的值被修改了,说明map是引用类型
}

//修改
func modify(mapScore map[string]int) {
	mapScore["Tom"] = 66
}
