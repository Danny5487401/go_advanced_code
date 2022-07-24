# 数组结构：
一组相同元素类型的集合,数组是一块固定大小的连续的内存空间

## 源码：cmd/compile/internal/types.Array
```go
type Array struct {
	Elem  *Type // element type
	Bound int64 // number of elements; <0 if unknown yet
}
```

## 初始化
1. 一种是定义指定大小的数值，第二种则是通过go的语法糖[…]{1,2,3}
2. 第二种通过[…]初始化的方式,会在编译期间调用cmd/compile/internal/gc.typecheckcomplit自动计算出数组的长度。

语法常量:cmd/compile/internal/gc/syntax.go
```go
type Op uint8

// Node ops.
const (
    OXXX Op = iota
    OTARRAY  // []int, [8]int, [N]int or [...]int
        // misc
    ODDD        // func f(args ...int) or f(l...) or var a = [...]int{0, 1, 2}
)
```
类型检查
```go
func typecheckcomplit(n *Node) (res *Node) {
    .....
    if n.Right.Op == OTARRAY && n.Right.Left != nil && n.Right.Left.Op == ODDD {
        n.Right.Right = typecheck(n.Right.Right, ctxType)
        if n.Right.Right.Type == nil {
            n.Type = nil
            return n
        }
        //元素的类型
        elemType := n.Right.Right.Type

        //计算出元素的个数
        length := typecheckarraylit(elemType, -1, n.List.Slice(), "array literal")

        n.Op = OARRAYLIT
        //调用NewArray完成初始化
        n.Type = types.NewArray(elemType, length)
        n.Right = nil
        return n
    }
    .....
}
```
上面两种初始化过程,在go编译期间，最终都会调用NewArray来完成，并且一开始就确定该数组是会被分配在堆上还是在栈上
```go

// 数组的初始化源代码在cmd/compile/internal/types.NewArray
func NewArray(elem *Type, bound int64) *Type {
	if bound < 0 {
		Fatalf("NewArray: invalid bound %v", bound)
	}
	t := New(TARRAY)
	t.Extra = &Array{Elem: elem, Bound: bound}
	t.SetNotInHeap(elem.NotInHeap())
	return t
}
```



