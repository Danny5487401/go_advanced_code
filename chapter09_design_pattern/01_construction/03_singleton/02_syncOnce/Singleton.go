package main

/*
单例模式：
	在程序中我们只需要某个类实例化一次即可，保证一个类仅有一个实例，并提供一个获取实例的方法。

标准库案例源码分析：
	strings/replace.go
	线程安全且支持规则复用
	结构体定义
		type Replacer struct {
			once   sync.Once // 控制 r replacer 替换算法初始化
			r      replacer
			oldnew []string
		}
	func NewReplacer(oldnew ...string) *Replacer {
			****
		return &Replacer{oldnew: append([]string(nil), oldnew...)} //没有创建算法
	}
	当我们使用 strings.NewReplacer 创建 strings.Replacer 时，这里采用惰性算法，并没有在这时进行 build 解析替换规则并创建对应算法实例，
	func (r *Replacer) Replace(s string) string {
		r.once.Do(r.buildOnce) //初始化
		return r.r.Replace(s)
	}

	func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error) {
		r.once.Do(r.buildOnce) //初始化
		return r.r.WriteString(w, s)
	}
	而是在执行替换时( Replacer.Replace 和 Replacer.WriteString)进行的
	初始化算法
		func (r *Replacer) buildOnce() {
			r.r = r.build()
			r.oldnew = nil
		}

		func (b *Replacer) build() replacer {
			....
		}

*/

import (
	"fmt"
	"sync"
)

type Singleton struct{}

var (
	singleton *Singleton
	once      sync.Once
)

func GetSingleton() *Singleton {
	once.Do(func() {
		fmt.Println("Create Obj")
		singleton = new(Singleton)
	})
	return singleton
}
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			obj := GetSingleton()
			fmt.Printf("%p\n", obj) //打印其地址
			wg.Done()
		}()
	}
	wg.Wait()
	// 只打印一次 Create Obj
}
