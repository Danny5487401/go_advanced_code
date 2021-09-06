package GId_Package

import (
	"runtime"
	"strings"
	"unsafe"
)

var offsetDictMap = map[string]int64{
	"go1.15": 152,
	"go1.10": 152,
	"go1.9":  152,
	"go1.8":  192,
}

// 获取不同go版本的偏移量
var g_goid_offset2 = func() int64 {
	goversion := runtime.Version()
	for key, off := range offsetDictMap {
		if goversion == key || strings.HasPrefix(goversion, key) {
			return off
		}
	}
	panic("unsupported go version:" + goversion)
}

const g_goid_offset = 152

// 通过枚举获取
func getg() unsafe.Pointer
func GetGroutineId() int64 {
	g := getg()
	p := (*int64)(unsafe.Pointer(uintptr(g) + g_goid_offset))
	return *p
}
