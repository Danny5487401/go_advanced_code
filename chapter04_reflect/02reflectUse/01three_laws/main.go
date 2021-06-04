package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 定律1：从接口值到反射对象  普通变量 -> 接口变量 -> 反射对象
	//反射操作：通过反射，可以获取一个接口类型变量的 类型和数值
	var x float64 = 3.4

	fmt.Println("type:", reflect.TypeOf(x))   //type: float64 获取某个变量的静态类型
	fmt.Println("value:", reflect.ValueOf(x)) //value: 3.4  获取某个变量的值

	fmt.Println("-------------------")
	//根据反射的值，来获取对应的类型和数值
	v := reflect.ValueOf(x)
	// reflect.Value.Kind() Kind获取变量值的底层类型（类别），注意不是类型，是Int、Float，还是Struct，还是Slice
	fmt.Println("kind is float64: ", v.Kind() == reflect.Float64)

	// 获取变量值的类型，效果等同于reflect.TypeOf
	fmt.Println("type : ", v.Type())

	fmt.Println("value : ", v.Float())

	// 自定义类型
	type MyInt int
	var y MyInt = 7
	v1 := reflect.ValueOf(y)                         // 接口变量存储了实际的值和值的类型
	fmt.Println("kind 底层类型 is float64: ", v1.Kind()) // kind 底层类型 is float64:  int
	// 获取变量值的类型，效果等同于reflect.TypeOf
	fmt.Println("type 定义类型 : ", v1.Type()) // type 定义类型 :  main.MyInt

	// 定律3：当反射对象所存的值是可设置时，反射对象才可修改
	// 可以使用Value.CanSet()检测是否可以通过此Value修改原始变量的值。
	z := 10
	v2 := reflect.ValueOf(z)
	fmt.Println("setable:", v2.CanSet()) // setable: false
	p := reflect.ValueOf(&z)
	fmt.Println("setable:", p.CanSet()) // setable: false
	v3 := p.Elem()
	fmt.Println("setable:", v3.CanSet()) // setable: true
	// Value.SetXXX()系列函数可设置Value中原始对象的值。 根据Value.Kind()的结果去获得变量的底层类型，然后选用该类别的Set函数。
	switch v3.Kind() {
	case reflect.Int:
		v3.SetInt(100) //设置值
	default:
		fmt.Println("未知类型")
	}

	fmt.Printf("z新值是%d", z) // 100
}

/*
	1. reflect.TypeOf： 直接给到了我们想要的type类型，如float64、int、各种pointer、struct 等等真实的类型
	2. reflect.ValueOf：直接给到了我们想要的具体的值，如1.2345这个具体数值，或者类似&{1 "Allen.Wu" 25} 这样的结构体struct的值
	3. 也就是说明反射可以将“接口类型变量”转换为“反射类型对象”，反射类型指的是reflect.Type和reflect.Value这两种
*/

// 反射三大定律
/*
第一条是最基本的：反射可以从接口值得到反射对象。

	反射是一种检测存储在 interface中的类型和值机制。这可以通过 TypeOf函数和 ValueOf函数得到。

第二条实际上和第一条是相反的机制，反射可以从反射对象获得接口值。

	它将 ValueOf的返回值通过 Interface()函数反向转变成 interface变量。

	前两条就是说 接口型变量和 反射类型对象可以相互转化，反射类型对象实际上就是指的前面说的 reflect.Type和 reflect.Value。
	func (v Value) Interface() (i interface{})  如果Value是结构体的非导出字段，调用该函数会导致panic。

   第三条不太好懂：如果需要操作一个反射变量，则其值必须可以修改。

	反射变量可设置的本质是它存储了原变量本身，这样对反射变量的操作，就会反映到原变量本身；反之，如果反射变量不能代表原变量，那么操作了反射变量，不会对原变量产生任何影响，这会给使用者带来疑惑。所以第二种情况在语言层面是不被允许的。
*/
