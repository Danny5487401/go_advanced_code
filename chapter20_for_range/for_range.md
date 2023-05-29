<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [for range 源码分析](#for-range-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [Golang 1.5版本之前的 gcc 源码](#golang-15%E7%89%88%E6%9C%AC%E4%B9%8B%E5%89%8D%E7%9A%84-gcc-%E6%BA%90%E7%A0%81)
  - [go源码](#go%E6%BA%90%E7%A0%81)
  - [range对应的walkrange源码](#range%E5%AF%B9%E5%BA%94%E7%9A%84walkrange%E6%BA%90%E7%A0%81)
    - [walkrange函数中当数组切片的情况下](#walkrange%E5%87%BD%E6%95%B0%E4%B8%AD%E5%BD%93%E6%95%B0%E7%BB%84%E5%88%87%E7%89%87%E7%9A%84%E6%83%85%E5%86%B5%E4%B8%8B)
    - [walkrange函数中当节点是map的情况下](#walkrange%E5%87%BD%E6%95%B0%E4%B8%AD%E5%BD%93%E8%8A%82%E7%82%B9%E6%98%AFmap%E7%9A%84%E6%83%85%E5%86%B5%E4%B8%8B)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# for range 源码分析
语法糖（Syntactic sugar），也译为糖衣语法

    英国计算机科学家彼得·蘭丁发明的一个术语，指计算机语言中添加的某种语法，这种语法对语言的功能没有影响，但是更方便程序员使用。 
    语法糖让程序更加简洁，有更高的可读性。

range 是 Golang 语言定义的一种语法糖迭代器，1.5版本 Golang 引入自举编译器后 range 相关源码如下，
根据类型不同进行不同的处理，支持对切片和数组、map、通道、字符串类型的的迭代。编译器会对每一种 range 支持的类型做专门的 “语法糖还原”。

这里我们主要介绍数组切片和 map 的 for-range 迭代。字符串和通道的 range 迭代平时使用的不多，同时篇幅原因我们就不详细介绍了，
感兴趣可以自行查看 Golang 源码和参考文献中自举前 gcc 的源码

## Golang 1.5版本之前的 gcc 源码
```shell
//   for_temp := range
//   len_temp := len(for_temp)
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = for_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

## go源码
## range对应的walkrange源码
```go
//src/cmd/compile/internal/gc/range.go
// walkrange transforms various forms of ORANGE into
// simpler forms.  The result must be assigned back to n.
// Node n may also be modified in place, and may also be
// the returned node.
func walkrange(n *Node) *Node {
    //…………
    switch t.Etype {
        default:
            Fatalf("walkrange")

        case TARRAY, TSLICE:  
            //……  数组，切牌呢
        case TMAP:
           //  …… map
        case TCHAN:
            // …… 通道
        case TSRTING:
           //  ……字符串
    }
    ……
    n = walkstmt(n)

 lineno = lno
 return n
}
```
### walkrange函数中当数组切片的情况下
```go
case TARRAY, TSLICE:
    // order.stmt arranged for a copy of the array/slice variable if needed.
    // 创建副本
    ha := a

    hv1 := temp(types.Types[TINT])
    hn := temp(types.Types[TINT])

    init = append(init, nod(OAS, hv1, nil))
    init = append(init, nod(OAS, hn, nod(OLEN, ha, nil)))

    n.Left = nod(OLT, hv1, hn)
    n.Right = nod(OAS, hv1, nod(OADD, hv1, nodintconst(1)))

    // 只关心数据值
    // for range ha { body }
    if v1 == nil {
        break
    }

    // 只关心索引
    // for v1 := range ha { body }
    if v2 == nil {
        body = []*Node{nod(OAS, v1, hv1)}
        break
    }

    // 关心索引和数据的情况
    // for v1, v2 := range ha { body }
    if cheapComputableIndex(n.Type.Elem().Width) {
        // v1, v2 = hv1, ha[hv1]
        tmp := nod(OINDEX, ha, hv1)
        tmp.SetBounded(true)
        // Use OAS2 to correctly handle assignments
        // of the form "v1, a[v1] := range".
        a := nod(OAS2, nil, nil)
        a.List.Set2(v1, v2)
        a.Rlist.Set2(hv1, tmp)
        body = []*Node{a}
        break
}
```
array和slice分析特点：
1. 循环次数在循环开始前已经确定 
2. 循环的时候会创建每个元素的副本
3. 循环的时候短声明只会在开始时执行一次，后面都是重用
```go
func main() {
     slice := []int{0,1,2,3}
     m := make(map[int]*int)
     for key,val := range slice {
       m[key] = &val
       fmt.Println(key,&key)
       fmt.Println(val,&val)
     }
     for k,v := range m {
      fmt.Println(k,"->",*v)
     }
 }

```
    循环 index 和 value 在每次循环体中都会被重用，而不是新声明。
    for-range 循环里的短声明index,value :=相当于第一次是 := ，后面都是 =，所以变量地址是不变的，就相当于全局变量了。
    
    每次遍历会把被循环元素当前 key 和值赋给这两个全局变量，但是注意变量还是那个变量，地址不变，所以如果用的是地址的或者当前上下文环境值的话最后打印出来都是同一个值。
```go
/*
结果
    0 0xc0000b4008
    0 0xc0000b4010
    1 0xc0000b4008
    1 0xc0000b4010
    2 0xc0000b4008
    2 0xc0000b4010
    3 0xc0000b4008
    3 0xc0000b4010
    0 -> 3
    1 -> 3
    2 -> 3
    3 -> 3

结果分析

    key0、key1、key2、key3 其实都是短声明中的key变量，所以地址是一致的，val0、val1、val2、val3 其实都是短声明中的val变量，地址也一致

    最终遍历 map 进行输出时因为 map 赋值时用的是 val 的地址m[key] = &val,循环结束时 val 的值是3，所以最终输出时4个元素的值都是3。

    需要注意 map 的遍历输出结果 key 的顺序可能会不一致，比如2，0，1，3这样，那是因为 map 的遍历输出是无序的，后面会再说，但是对应的 value 的值都是3
 */


```
想要正确的map做法
```go
func main() {
     slice := []int{0,1,2,3}
     m := make(map[int]*int)
     for key,val := range slice {
       value := val    //增加临时变量，每次都是新声明的，地址也就不一样，也就能传过去正确的值
       m[key] = &value
       fmt.Println(key,&key)
       fmt.Println(val,&val)
     }
     for k,v := range m {
      fmt.Println(k,"->",*v)
     }
}
```
### walkrange函数中当节点是map的情况下
```go
	case TMAP:
		// 副本
		// order.stmt allocated the iterator for us.
		// we only use a once, so no copy needed.
		ha := a

		hit := prealloc[n]
		th := hit.Type
		n.Left = nil
		keysym := th.Field(0).Sym  // depends on layout of iterator struct.  See reflect.go:hiter
		elemsym := th.Field(1).Sym // ditto

		fn := syslook("mapiterinit")

		fn = substArgTypes(fn, t.Key(), t.Elem(), th)
		init = append(init, mkcall1(fn, nil, nil, typename(t), ha, nod(OADDR, hit, nil)))
		n.Left = nod(ONE, nodSym(ODOT, hit, keysym), nodnil())

		fn = syslook("mapiternext")
		fn = substArgTypes(fn, th)
		n.Right = mkcall1(fn, nil, nil, nod(OADDR, hit, nil))

		key := nodSym(ODOT, hit, keysym)
		key = nod(ODEREF, key, nil)
		if v1 == nil {
			body = nil
		} else if v2 == nil {
			body = []*Node{nod(OAS, v1, key)}
		} else {
			elem := nodSym(ODOT, hit, elemsym)
			elem = nod(ODEREF, elem, nil)
			a := nod(OAS2, nil, nil)
			a.List.Set2(v1, v2)
			a.Rlist.Set2(key, elem)
			body = []*Node{a}
		}
```


