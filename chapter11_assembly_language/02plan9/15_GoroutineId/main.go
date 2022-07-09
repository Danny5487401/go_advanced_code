package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter11_assembly_language/02plan9/15_GoroutineId/GId_Package"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//获取goroutineId

/*
在任意的类型被定义之后，Go语言都会为该类型生成对应的类型信息。比如g结构体会生成一个type·runtime·g标识符表示g结构体的值类型信息，同时还有一个type·*runtime·g标识符表示指针类型的信息。
*/

func main() {
	// 方式一：汇编获取goroutineId
	GId := GId_Package.GetGroutineId()
	fmt.Println("GoroutineId是", GetGoid(), GId)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// 汇编获取goroutineId
		GId = GId_Package.GetGroutineId()
		fmt.Println("GoroutineId是", GetGoid(), GId)
		wg.Done()
	}()
	wg.Wait()
}

// 方式二：栈函数获取goroutineId
func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}
