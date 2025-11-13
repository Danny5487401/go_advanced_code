package main

import (
	"fmt"
	"reflect"
)

// CanAddr 方法和 CanSet 方法不一样的地方在于：对于一些结构体内的私有字段，我们可以获取它的地址，但是不能设置它。

type FooStruct struct {
	A string
	b int
}

type Foo interface {
	Name() string
}

func (f FooStruct) Name() string {
	return f.A
}

type FooPointer struct {
	A string
}

func (f *FooPointer) Name() string {
	return f.A
}

func main() {
	{
		// canAddr vs canSet
		a := FooStruct{}
		val := reflect.ValueOf(&a)
		fmt.Println(val.Elem().FieldByName("b").CanSet())  // false
		fmt.Println(val.Elem().FieldByName("b").CanAddr()) // true
	}

	{
		/*  canSet
		元素本身是不能设置，

		元素的指针是肯定不能设置的，因为它是指针。

		元素的指针指向的元素是空的时候也不能设置

		元素的指针指向的元素(Elem)是可以设置的
		*/
		var n = 33
		var pn = &n
		var ppn = &pn
		var pn2 *int = nil

		// 指针是不能set的，指针指向的元素可以set
		fmt.Println("reflect.ValueOf(n).CanSet():", reflect.ValueOf(n).CanSet())     // reflect.ValueOf(n).CanSet(): false
		fmt.Println("reflect.ValueOf(pn).CanSet():", reflect.ValueOf(pn).CanSet())   // reflect.ValueOf(pn).CanSet(): false
		fmt.Println("reflect.ValueOf(pn2).CanSet():", reflect.ValueOf(pn2).CanSet()) // reflect.ValueOf(pn2).CanSet(): false
		fmt.Println("reflect.ValueOf(ppn).CanSet():", reflect.ValueOf(ppn).CanSet()) // reflect.ValueOf(ppn).CanSet(): false

		// 指针是不能set的，指针指向的元素可以set
		//fmt.Println("reflect.ValueOf(n).Elem().CanSet():", reflect.ValueOf(n).Elem().CanSet()) // 值没有elem
		fmt.Println("reflect.ValueOf(pn).Elem().CanSet():", reflect.ValueOf(pn).Elem().CanSet())   // reflect.ValueOf(pn).Elem().CanSet(): true
		fmt.Println("reflect.ValueOf(pn2).Elem().CanSet():", reflect.ValueOf(pn2).Elem().CanSet()) // reflect.ValueOf(pn2).Elem().CanSet(): false
		fmt.Println("reflect.ValueOf(ppn).Elem().CanSet():", reflect.ValueOf(ppn).Elem().CanSet()) // reflect.ValueOf(ppn).Elem().CanSet(): true

		// 指针是不能set的，指针指向的元素可以set
		fmt.Println("reflect.ValueOf(n).CanAddr():", reflect.ValueOf(n).CanAddr())     // reflect.ValueOf(n).CanAddr(): false
		fmt.Println("reflect.ValueOf(pn).CanAddr():", reflect.ValueOf(pn).CanAddr())   // reflect.ValueOf(pn).CanAddr(): false
		fmt.Println("reflect.ValueOf(pn2).CanAddr():", reflect.ValueOf(pn2).CanAddr()) // reflect.ValueOf(pn2).CanAddr(): false
		fmt.Println("reflect.ValueOf(ppn).CanAddr():", reflect.ValueOf(ppn).CanAddr()) // reflect.ValueOf(ppn).CanAddr(): false

		// 指针是不能set的，指针指向的元素可以set
		// fmt.Println("reflect.ValueOf(n).Elem().CanAddr():", reflect.ValueOf(n).Elem().CanAddr())// 值没有elem
		fmt.Println("reflect.ValueOf(pn).Elem().CanAddr():", reflect.ValueOf(pn).Elem().CanAddr())   // reflect.ValueOf(pn).Elem().CanAddr(): true
		fmt.Println("reflect.ValueOf(pn2).Elem().CanAddr():", reflect.ValueOf(pn2).Elem().CanAddr()) // reflect.ValueOf(pn2).Elem().CanAddr(): false
		fmt.Println("reflect.ValueOf(ppn).Elem().CanAddr():", reflect.ValueOf(ppn).Elem().CanAddr()) // reflect.ValueOf(ppn).Elem().CanAddr(): true

	}

	{
		// 对于复杂的 slice， map， struct， pointer 等方法
		{
			// slice
			a := []int{1, 2, 3}
			val := reflect.ValueOf(&a)
			val.Elem().SetLen(2)
			val.Elem().Index(0).SetInt(4)
			fmt.Println(a) // [4,2]
		}
		{
			// map
			a := map[int]string{
				1: "foo1",
				2: "foo2",
			}
			val := reflect.ValueOf(&a)
			key3 := reflect.ValueOf(3)
			val3 := reflect.ValueOf("foo3")
			val.Elem().SetMapIndex(key3, val3)
			fmt.Println(val) // &map[1:foo1 2:foo2 3:foo3]
		}
		{
			// struct
			a := FooStruct{}
			val := reflect.ValueOf(&a)
			val.Elem().FieldByName("A").SetString("foo2")
			fmt.Println(a) // {foo2}
		}
		{
			// pointer
			a := &FooPointer{}
			val := reflect.ValueOf(a)
			val.Elem().FieldByName("A").SetString("foo2")
			fmt.Println(a) //&{foo2}
		}
	}
}
