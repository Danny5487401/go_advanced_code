package main

// sync pool 使用，比如fmt.Printf就有用到
//fmt.Printf->Fprintf(os.Stdout, format, a...)-->newPrinter()-->ppFree sync.pool实例

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
)

/* 应用场景
1.当我们必须重用共享的和长期存在的对象（例如，数据库连接）时。
2.用于优化内存分配。
用法：
1. Get() interface{} 用来从并发池中取出元素。
2. Put(interface{}) 将一个对象加入并发池。
3. New func() interface{}

 */
func NewConnection(num int) *Connection {
	return &Connection{id: num}
}

type Connection struct {
	id int
}

func main()  {

	pool := &sync.Pool{}

	// put 归还 内存对象
	pool.Put(NewConnection(1))
	pool.Put(NewConnection(2))
	pool.Put(NewConnection(3))

	// get 获取 内存对象
	connection := pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)

	connection = pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)

	connection = pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)

	/*
	考虑一个写入缓冲区并将结果持久保存到文件中的函数示例。使用sync.Pool，我们可以通过在不同的函数调用之间重用同一对象来重用为缓冲区分配的空间
	 */

}
/*结果
1
3
2
 */

//Note: 是Get()方法会从并发池中随机取出对象，无法保证以固定的顺序获取并发池中存储的对象。

/*
步骤：第一步是检索先前分配的缓冲区（如果是第一个调用，则创建一个缓冲区，但这是抽象的）。然后，defer操作是将缓冲区放回sync.Pool中
 */
func writeFile(pool *sync.Pool, filename string) error {
	buf := pool.Get().(*bytes.Buffer)

	defer pool.Put(buf)

	// Reset 缓存区，不然会连接上次调用时保存在缓存区里的字符串foo
	// 编程foo foo 以此类推
	buf.Reset()

	//写的内容
	buf.WriteString("foo")

	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
