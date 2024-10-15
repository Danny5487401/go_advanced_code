package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/Danny5487401/go_advanced_code/chapter11_assembly_language/02plan9/15_GoroutineId/GId_Package"
)

//获取goroutineId

/*
在任意的类型被定义之后，Go语言都会为该类型生成对应的类型信息。比如g结构体会生成一个type·runtime·g标识符表示g结构体的值类型信息，同时还有一个type·*runtime·g标识符表示指针类型的信息。
*/

func main() {
	// 方式一：汇编获取goroutineId
	GId := GId_Package.GetGoroutineId()
	fmt.Println("GoroutineId是", GetGoid(), GId)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// 汇编获取goroutineId
		GId = GId_Package.GetGoroutineId()
		fmt.Println("GoroutineId 是", GetGoid(), GId)
		wg.Done()
	}()
	wg.Wait()
}

// GetGoid 方式二：栈函数获取goroutineId
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

/*
stack 参考数据


goroutine 4 [running]:
go.uber.org/goleak/internal/stack.getStackBuffer(0x0)
	/Users/python/go/pkg/mod/go.uber.org/goleak@v1.1.12/internal/stack/stacks.go:124 +0x68
go.uber.org/goleak/internal/stack.getStacks(0x0)
	/Users/python/go/pkg/mod/go.uber.org/goleak@v1.1.12/internal/stack/stacks.go:73 +0x44
go.uber.org/goleak/internal/stack.Current()
	/Users/python/go/pkg/mod/go.uber.org/goleak@v1.1.12/internal/stack/stacks.go:118 +0x30
go.uber.org/goleak.Find({0x0, 0x0, 0x0})
	/Users/python/go/pkg/mod/go.uber.org/goleak@v1.1.12/leaks.go:55 +0x38
go.uber.org/goleak.VerifyNone({0x1009fbb48, 0x14000118340}, {0x0, 0x0, 0x0})
	/Users/python/go/pkg/mod/go.uber.org/goleak@v1.1.12/leaks.go:77 +0x40
github.com/Danny5487401/go_advanced_code/chapter02_goroutine/05_goroutine_leaks/02_avoid_leaks.TestGetDataWithGoleak(0x14000118340)
	/Users/python/Desktop/go_advanced_code/chapter02_goroutine/05_goroutine_leaks/02_avoid_leaks/goroutine_leak_test.go:12 +0x88
testing.tRunner(0x14000118340, 0x1009fad38)
	/Users/python/go/go1.18/src/testing/testing.go:1439 +0x178
*/
