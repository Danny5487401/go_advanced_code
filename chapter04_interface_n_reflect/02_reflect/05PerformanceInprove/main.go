package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 需求： 给结构体Test的内部元素X加一

type Test struct {
	Y string
	X int
}

// IncX 方式一:
func IncX(d interface{}) int64 {
	// 根据字段名获取元素
	v := reflect.ValueOf(d).Elem()
	f := v.FieldByName("X")

	// 转换成int64处理
	x := f.Int()

	// 加一操作
	x++

	// 设置值
	f.SetInt(x)
	return x
}

/*
分析：
	 如果是 reflect.Type，可将其缓存，避免重复操作耗时。
	但 reflect.Value 显然不行，因为它和具体对象绑定，内部存储实例指针
*/

// 方式二:
var offset uintptr = 0xFFFF // 避开offset=0的字段
func UnsafeIncX(d interface{}) int {
	if offset == 0xFFFF {
		t := reflect.TypeOf(d).Elem()
		x, _ := t.FieldByName("X")
		offset = x.Offset
	}
	// 用数组指针存取
	p := (*[2]uintptr)(unsafe.Pointer(&d))
	// 根据偏移获取第一个元素
	px := (*int)(unsafe.Pointer(p[1] + offset))
	*px++
	return *px
}

/*
分析：
	如何设计缓存结构，这个 offset 变量自然不能用于实际开发?
	用 map[Type]map[name]offset？显然不行。每次执行 reflect.TypeOf，这于性能优化不利。可除了 Type，还有什么可以作为 Key 使用？
	要知道，接口由 itab 和 data 指针组成，相同类型（接口和实际类型组合）的 itab 指针相同，自然也可当作 key 来用
*/

// 方式三
var cache = map[*uintptr]map[string]uintptr{}

func UnsafeCacheIncrX(d interface{}) int {
	itab := *(**uintptr)(unsafe.Pointer(&d))

	// 尝试获取结构体
	m, ok := cache[itab]
	if !ok {
		// 不存在，则加入
		m = make(map[string]uintptr)
		cache[itab] = m
	}

	// 尝试获取元素
	newOffset, ok := m["X"]
	if !ok {
		// 不存在，则加入元素
		t := reflect.TypeOf(d).Elem()
		x, _ := t.FieldByName("X")
		newOffset = x.Offset
		m["X"] = newOffset
	}

	p := (*[2]uintptr)(unsafe.Pointer(&d))
	px := (*int)(unsafe.Pointer(p[1] + newOffset))
	*px++
	return *px
}

// 虽因引入 map 导致性能有所下降，但相比直接使用 reflect 还是提升很多

func main() {
	d := Test{X: 100}
	fmt.Println("方式一：", IncX(&d))
	fmt.Println("方式二：", UnsafeIncX(&d))
	fmt.Println("方式三：", UnsafeCacheIncrX(&d))
}
