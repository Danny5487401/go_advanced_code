package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 方式一:
func IncX(d interface{}) int64 {
	v := reflect.ValueOf(d).Elem()
	f := v.FieldByName("X")
	x := f.Int()
	x++
	f.SetInt(x)
	return x
}

// 如果是 reflect.Type，可将其缓存，避免重复操作耗时。但 Value 显然不行，因为它和具体对象绑定，内部存储实例指针
// 方式二:
var offset uintptr = 0xFFFF // 避开offset=0的字段
func UnsafeIncX(d interface{}) int {
	if offset == 0xFFFF {
		t := reflect.TypeOf(d).Elem()
		x, _ := t.FieldByName("X")
		offset = x.Offset
	}
	p := (*[2]uintptr)(unsafe.Pointer(&d))
	px := (*int)(unsafe.Pointer(p[1] + offset))
	*px++
	return *px
}

// 如何设计缓存结构，这个 offset 变量自然不能用于实际开发?
// 用 map[Type]map[name]offset？显然不行。每次执行 reflect.TypeOf，这于性能优化不利。可除了 Type，还有什么可以作为 Key 使用？
// 要知道，接口由 itab 和 data 指针组成，相同类型（接口和实际类型组合）的 itab 指针相同，自然也可当作 key 来用
// 方式三
var cache = map[*uintptr]map[string]uintptr{}

func UnsafeCacheIncrX(d interface{}) int {
	itab := *(**uintptr)(unsafe.Pointer(&d))
	m, ok := cache[itab]
	if !ok {
		m = make(map[string]uintptr)
		cache[itab] = m
	}
	offset, ok := m["X"]
	if !ok {
		t := reflect.TypeOf(d).Elem()
		x, _ := t.FieldByName("X")
		offset = x.Offset
		m["X"] = offset
	}
	p := (*[2]uintptr)(unsafe.Pointer(&d))
	px := (*int)(unsafe.Pointer(p[1] + offset))
	*px++
	return *px
}

// 虽因引入 map 导致性能有所下降，但相比直接使用 reflect 还是提升很多

func main() {
	d := struct {
		X int
	}{100}
	fmt.Println(IncX(&d))
	fmt.Println(UnsafeIncX(&d))
	fmt.Println(UnsafeCacheIncrX(&d))
}
