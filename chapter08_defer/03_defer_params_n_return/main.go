package main

import "fmt"

/*
defer 是 Go 语言提供的一种用于注册延迟调用的机制，每一次 defer 都会把函数压入栈中，当前函数返回前再把延迟函数取出并执行.
defer 语句定义时，对外部变量的引用是有两种方式:
1. 函数参数  :   作为函数参数，则在 defer 定义时就把值传递给 defer，并被缓存起来；
2. 作为闭包引用 :  作为闭包引用的话，则会在 defer 函数真正调用时根据整个上下文确定当前的值
*/

/* 注意这两个函数
1.
func increaseA() int {
    var i int
    defer func() {
        i++
    }()
    return i
}

2.
func increaseB() (r int) {
    defer func() {
        r++
    }()
    return r
}

*/

/* 理解return xxx： 这条语句并不是一个原子指令，经过编译之后，变成了三条指令：
1. 返回值 = xxx
2. 调用 defer 函数
3. 空的 return

1,3 步才是 return 语句真正的命令，第 2 步是 defer 定义的语句，这里就有可能会操作返回值

*/

//// 三个例子
//1. 闭包引用
func f1() (r int) {
	defer func() {
		r++
	}()
	return 0
}

/*
// 第一题拆解过程
func f1() (r int) {

    // 1.赋值
    r = 0

    // 2.闭包引用，返回值被修改
    defer func() {
        r++
    }()

    // 3.空的 return
    return
}

*/

//2. 闭包引用
func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

/*
// 第二题拆解过程：

func f2() (r int) {
    t := 5
    // 1.赋值
    r = t

    // 2.闭包引用，但是没有修改返回值 r
    defer func() {
        t = t + 5
    }()

    // 3.空的 return
    return
}
*/

//3. 函数参数,负责参数
func f3() (r int) {

	defer func(r int) {
		r = r + 5
	}(r)
	// 在有具名返回值的函数中（这里具名返回值为 1）,执行 return 的时候实际上已经将 r 的值重新赋值为 1
	return 1
}

/* 分析

// 第三题拆解过程：
func f3() (r int) {

    // 1.赋值
    r = 1

    // 2.r 作为函数参数，不会修改要返回的那个 r 值
    defer func(r int) {
        r = r + 5
    }(r)

    // 3.空的 return
    return
}

r 是作为函数参数使用，是一份复制，defer 语句里面的 r 和 外面的 r 其实是两个变量，里面变量的改变不会影响外层变量 r，所以不是返回 6 ，而是返回 1。
*/

//4. 闭包引用,具名返回
func f4() (r int) {
	r = 1
	defer func() {
		r = r + 5
	}()
	// 在有具名返回值的函数中（这里具名返回值为 1）,执行 return 的时候实际上已经将 r 的值重新赋值为 2
	return 2
}

/*
func f4() (r int) {
	r = 1
	r = 2
    // 2.闭包引用，修改返回值 r
    defer func() {
        r = r + 5
    }()
    return
}

*/

func main() {
	fmt.Println("第一个函数结果", f1()) // 第一个函数结果 1
	fmt.Println("第二个函数结果", f2()) // 第二个函数结果 5
	fmt.Println("第三个函数结果", f3()) // 第三个函数结果 1
	fmt.Println("第四个函数结果", f4()) // 第四个函数结果 7

}

//Note:
/*
func increaseA() int {
    var i int
    defer func() {
        i++
    }()
    return i
}
函数 increaseA() 是匿名返回值，返回局部变量，同时 defer 函数也会操作这个局部变量。
	对于匿名返回值来说，可以假定有一个变量存储返回值，比如假定返回值变量为 anony，上面的返回语句可以拆分成以下过程：
annoy = i
i++
return

结论：由于 i 是整型，会将值拷贝给 anony，所以 defer 语句中修改 i 值，对函数返回值不造成影响，所以返回 0 。

*/
