package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	// 定律1：从接口值到反射对象  普通变量 -> 接口interface变量 -> 反射对象
	// 反射操作：通过反射，可以获取一个接口类型变量的 类型和数值
	Law1()

	// 定律2：反射对象->interface变量 -> 断言 普通变量
	Law2()

	// 定律3：当反射对象所存的值是可设置时，反射对象才可修改
	//Go 语言里的函数都是值传递，只要你传递的不是变量的指针，你在函数内部对变量的修改是不会影响到原始的变量的。
	//回到反射上来，当你使用 reflect.Typeof 和 reflect.Valueof 的时候，如果传递的不是接口变量的指针，
	//反射世界里的变量值始终将只是真实世界里的一个拷贝，你对该反射对象进行修改，并不能反映到真实世界里。
	// 可以使用Value.CanSet()检测是否可以通过此Value修改原始变量的值,可寻址的
	Law3()
}

func Law1() {
	// 1. 内置类型
	Typeof(123)

	// 2. 自定义类型
	Typeof(MyInt(456))

	// 3. 结构体类型
	coder := &Coder{Name: "danny"}
	typ := reflect.TypeOf(coder)
	val := reflect.ValueOf(coder)
	fmt.Println("kind of coder : ", typ.Kind()) // kind of coder :  ptr
	fmt.Println("type of coder : ", typ)        // type of coder :  *main.Coder
	fmt.Println("value of coder : ", val)       // value of coder :  danny

	typeofStringer := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()       // 获得接口类型
	fmt.Println("kind of Stringer : ", typeofStringer.Kind())           // kind of Stringer :  interface
	fmt.Println("type of Stringer : ", typeofStringer)                  // type of Stringer :  fmt.Stringer
	fmt.Println("PkgPath of Stringer : ", typeofStringer.PkgPath())     // PkgPath of Stringer :  fmt
	fmt.Printf("Method of Stringer : %#v \n", typeofStringer.Method(0)) // Method of Stringer : reflect.Method{Name:"String", PkgPath:"", Type:(*reflect.rtype)(0x104b7cba0), Func:reflect.Value{typ:(*reflect.rtype)(nil), ptr:(unsafe.Pointer)(nil), flag:0x0}, Index:0}

}

type MyInt int

//go:noinline
func (f MyInt) Ree() int {
	return int(f)
}

//go:noinline
func (f MyInt) String() string {
	return strconv.Itoa(f.Ree())
}

//go:noinline
func (f MyInt) print() {
	// 明明在源码中定义了print方法，为什么找不到该方法呢?
	// print方法是一个私有方法，不会被外部调用，但是main包范围内又没有调用者; Go编译器本着勤俭节约的原则，把print方法优化丢弃掉了
	fmt.Println("foo is " + f.String())
}

//go:noinline
func Typeof(i interface{}) {
	t := reflect.TypeOf(i)
	fmt.Println("值  ", i)
	fmt.Println("名称", t.Name())
	fmt.Println("类型", t.String())
	fmt.Println("方法")
	num := t.NumMethod()
	if num > 0 {
		for j := 0; j < num; j++ {
			fmt.Println("  ", t.Method(j).Name, t.Method(j).Type)
		}
	}
	fmt.Println("-------------------")
}

func Law2() {
	fmt.Println("-------------------")
	coder := &Coder{Name: "danny"}
	val := reflect.ValueOf(coder)
	c, ok := val.Interface().(*Coder)
	if ok {
		fmt.Println("类型断言成功,姓名是", c.Name) // danny
	} else {
		panic("type assert to *Coder error")
	}
}

func Law3() {
	fmt.Println("-------------------")
	z := 10 // 不是指针
	v2 := reflect.ValueOf(z)
	fmt.Println("settable:", v2.CanSet()) // settable: false

	p := reflect.ValueOf(&z)
	fmt.Println("settable:", p.CanSet()) // settable: false

	v3 := p.Elem()                        // 指针指向的元素
	fmt.Println("settable:", v3.CanSet()) // settable: true

	// Value.SetXXX()系列函数可设置Value中原始对象的值。 根据Value.Kind()的结果去获得变量的底层类型，然后选用该类别的Set函数。
	switch v3.Kind() {
	case reflect.Int:
		v3.SetInt(100) //设置值   源码会检测是否可设置mustBeAssignable,然后检测可导出
	default:
		fmt.Println("未知类型")
	}

	fmt.Printf("z新值是%d", z) // 100
}

type Coder struct {
	Name string
}

func (c *Coder) String() string {
	return c.Name
}
