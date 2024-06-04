package main

import (
	"fmt"
	"reflect"
)

func main() {
	createInt()
	createFloat()
	createString()

	createArray()
	createSlice()
	createMap()
	createChan()
	createStruct()
	createPointer()
	createFunc()
}

func createInt() {
	val := reflect.New(reflect.TypeOf(0))
	val.Elem().SetInt(42)
	fmt.Println(val.Elem().Int()) // 输出：42
}

func createFloat() {
	val := reflect.New(reflect.TypeOf(0.0))
	val.Elem().SetFloat(3.14)
	fmt.Println(val.Elem().Float()) // 输出：3.1
}

func createString() {
	val := reflect.New(reflect.TypeOf(""))
	val.Elem().SetString("hello")
	fmt.Println(val.Elem().String()) // 输出：hello
}

func createArray() {
	typ := reflect.ArrayOf(3, reflect.TypeOf(0))
	val := reflect.New(typ)
	arr := val.Elem()
	arr.Index(0).SetInt(1)
	arr.Index(1).SetInt(2)
	arr.Index(2).SetInt(3)
	fmt.Println(arr.Interface()) // 输出：[1 2 3]
	arr1, ok := arr.Interface().([3]int)
	if !ok {
		fmt.Println("not a [3]int")
		return
	}

	fmt.Println(arr1) // [1 2 3]
}

func createSlice() {
	typ := reflect.SliceOf(reflect.TypeOf(0)) // 切片元素类型
	val := reflect.MakeSlice(typ, 3, 3)       // 动态创建切片实例
	val.Index(0).SetInt(1)
	val.Index(1).SetInt(2)
	val.Index(2).SetInt(3)
	fmt.Println(val.Interface()) // 输出：[1 2 3]

	sl, ok := val.Interface().([]int)
	if !ok {
		fmt.Println("sl is not a []int")
		return
	}
	fmt.Println(sl) // [1 2 3]
}

func createMap() {
	typ := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0))
	val := reflect.MakeMap(typ)
	key1 := reflect.ValueOf("one")
	value1 := reflect.ValueOf(1)
	key2 := reflect.ValueOf("two")
	value2 := reflect.ValueOf(2)
	val.SetMapIndex(key1, value1)
	val.SetMapIndex(key2, value2)
	fmt.Println(val.Interface()) // 输出：map[one:1 two:2]

	m, ok := val.Interface().(map[string]int)
	if !ok {
		fmt.Println("m is not a map[string]int")
		return
	}

	fmt.Println(m)
}

func createChan() {
	typ := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(0))
	val := reflect.MakeChan(typ, 0)
	go func() {
		val.Send(reflect.ValueOf(42))
	}()

	ch, ok := val.Interface().(chan int)
	if !ok {
		fmt.Println("ch is not a chan int")
		return
	}
	fmt.Println(<-ch) // 42
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() {
	fmt.Printf("Hello, my name is %s and I am %d years old\n", p.Name, p.Age)
}

func (p Person) SayHello(name string) {
	fmt.Printf("Hello, %s! My name is %s\n", name, p.Name)
}

func createStruct() {
	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(0),
		},
	})
	ptrVal := reflect.New(typ)
	val := ptrVal.Elem()
	val.FieldByName("Name").SetString("Danny")
	val.FieldByName("Age").SetInt(25)

	person := (*Person)(ptrVal.UnsafePointer())
	person.Greet()         // 输出：Hello, my name is Danny and I am 25 years old
	person.SayHello("Bob") // 输出：Hello, Bob! My name is Danny
}

func createPointer() {
	typ := reflect.PtrTo(reflect.TypeOf(Person{}))
	val := reflect.New(typ.Elem())
	val.Elem().FieldByName("Name").SetString("Alice")
	val.Elem().FieldByName("Age").SetInt(25)
	person := val.Interface().(*Person)
	fmt.Println(person.Name) // 输出：Alice
	fmt.Println(person.Age)  // 输出：25
}

func sum(args []reflect.Value) (res []reflect.Value) {
	a, b := args[0], args[1]
	if a.Kind() != b.Kind() {
		fmt.Println("format wrong")
		return nil
	}

	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []reflect.Value{reflect.ValueOf(a.Int() + b.Int())}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []reflect.Value{reflect.ValueOf(a.Uint() + b.Uint())}
	case reflect.Float32, reflect.Float64:
		return []reflect.Value{reflect.ValueOf(a.Float() + b.Float())}
	case reflect.String:
		return []reflect.Value{reflect.ValueOf(a.String() + b.String())}
	default:
		return []reflect.Value{}
	}
}

func createTwoSumFunc(fptr interface{}) {
	fn := reflect.ValueOf(fptr).Elem()

	v := reflect.MakeFunc(fn.Type(), sum)

	fn.Set(v)
}

func createFunc() {
	var intSum func(int, int) int64
	var floatSum func(float32, float32) float64
	var stringSum func(string, string) string

	createTwoSumFunc(&intSum)
	createTwoSumFunc(&floatSum)
	createTwoSumFunc(&stringSum)

	fmt.Println(intSum(1, 2))
	fmt.Println(floatSum(2.1, 3.5))
	fmt.Println(stringSum("Hello", "World"))
}
