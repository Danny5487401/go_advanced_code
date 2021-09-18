# 类型断言
##一。断言类型的语法
    x.(T)，这里x表示一个接口的类型，T表示一个类型（也可为接口类型）。
    一个类型断言检查一个接口对象x的动态类型是否和断言的类型T匹配。

## 二。分类
类型断言分两种情况：

    第一种，如果断言的类型T是一个具体类型，类型断言x.(T)就检查x的动态类型是否和T的类型相同。

		1。如果这个检查成功了，类型断言的结果是一个类型为T的对象，该对象的值为接口变量x的动态值。换句话说，具体类型的类型断言从它的操作对象中获得具体的值。
		2。如果检查失败，接下来这个操作会抛出panic，除非用两个变量来接收检查结果，如：f, ok := w.(*os.File)

	第二种，如果断言的类型T是一个接口类型，类型断言x.(T)检查x的动态类型是否满足T接口。

		1。如果这个检查成功，则检查结果的接口值的动态类型和动态值不变，但是该接口值的类型被转换为接口类型T。
        换句话说，对一个接口类型的类型断言改变了类型的表述方式，改变了可以获取的方法集合（通常更大），但是它保护了接口值内部的动态类型和值的部分。
		2。如果检查失败，接下来这个操作会抛出panic，除非用两个变量来接收检查结果，如：f, ok := w.(io.ReadWriter)
注意：

    如果断言的操作对象x是一个nil接口值，那么不论被断言的类型T是什么这个类型断言都会失败。
    我们几乎不需要对一个更少限制性的接口类型（更少的方法集合）做断言，因为它表现的就像赋值操作一样，除了对于nil接口值的情况。

    表达式是t,ok:=i.(T)，这个表达式也是可以断言一个接口对象（i）里不是nil，并且接口对象（i）存储的值的类型是 T，如果断言成功，就会返回其类型给t，并且此时 ok 的值 为true，表示断言成功。
    如果接口值的类型，并不是我们所断言的 T，就会断言失败，但和第一种表达式不同的是这个不会触发 panic，而是将 ok 的值设为false，表示断言失败，此时t为T的零值。所以推荐使用这种方式，可以保证代码的健壮性

## 三。反射reflect类型断言

从reflect.Value中获取接口interface的信息

    realValue := value.Interface().(已知的类型)
    可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic

    // Golang 对类型要求非常严格，类型一定要完全符合
    // 如下两个，一个是*float64，一个是float64，如果弄混，则会panic

	1. 从 Value 到实例
		该方法最通用，用来将 Value 转换为空接口，该空接口内部存放具体类型实例
		可以使用接口类型查询去还原为具体的类型
		//func (v Value) Interface() （i interface{})

## 四. 性能比较

    *空接口类型的类型断言代价并不高，与直接类型转换几乎没有性能差异
    *空接口类型使用type switch进行类型断言时，随着case的增多性能会直线下降
    *非空接口类型进行类型断言时，随着接口中方法的增多，性能会直线下降
    *直接进行方法调用要比非接口类型进行类型断言要高效很多
```go
var dst int64

// 空接口类型直接类型断言具体的类型
func Benchmark_efaceToType(b *testing.B) {
 b.Run("efaceToType", func(b *testing.B) {
  var ebread interface{} = int64(666)
  for i := 0; i < b.N; i++ {
   dst = ebread.(int64)
  }
 })
}

// 空接口类型使用TypeSwitch 只有部分类型
func Benchmark_efaceWithSwitchOnlyIntType(b *testing.B) {
 b.Run("efaceWithSwitchOnlyIntType", func(b *testing.B) {
  var ebread interface{} = 666
  for i := 0; i < b.N; i++ {
   OnlyInt(ebread)
  }
 })
}

// 空接口类型使用TypeSwitch 所有类型
func Benchmark_efaceWithSwitchAllType(b *testing.B) {
 b.Run("efaceWithSwitchAllType", func(b *testing.B) {
  var ebread interface{} = 666
  for i := 0; i < b.N; i++ {
   Any(ebread)
  }
 })
}

//直接使用类型转换
func Benchmark_TypeConversion(b *testing.B) {
 b.Run("typeConversion", func(b *testing.B) {
  var ebread int32 = 666

  for i := 0; i < b.N; i++ {
   dst = int64(ebread)
  }
 })
}

// 非空接口类型判断一个类型是否实现了该接口 两个方法
func Benchmark_ifaceToType(b *testing.B) {
 b.Run("ifaceToType", func(b *testing.B) {
  var iface Basic = &User{}
  for i := 0; i < b.N; i++ {
   iface.GetName()
   iface.SetName("1")
  }
 })
}

// 非空接口类型判断一个类型是否实现了该接口 12个方法
func Benchmark_ifaceToTypeWithMoreMethod(b *testing.B) {
 b.Run("ifaceToTypeWithMoreMethod", func(b *testing.B) {
  var iface MoreMethod = &More{}
  for i := 0; i < b.N; i++ {
   iface.Get()
   iface.Set()
   iface.One()
   iface.Two()
   iface.Three()
   iface.Four()
   iface.Five()
   iface.Six()
   iface.Seven()
   iface.Eight()
   iface.Nine()
   iface.Ten()
  }
 })
}

// 直接调用方法
func Benchmark_DirectlyUseMethod(b *testing.B) {
 b.Run("directlyUseMethod", func(b *testing.B) {
  m := &More{
   Name: "asong",
  }
  m.Get()
 })
}

```
