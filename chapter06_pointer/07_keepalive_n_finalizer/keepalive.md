<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [runtime.KeepAlive](#runtimekeepalive)
  - [基本知识](#%E5%9F%BA%E6%9C%AC%E7%9F%A5%E8%AF%86)
    - [静态单赋值SSA（Static Single Assignment）](#%E9%9D%99%E6%80%81%E5%8D%95%E8%B5%8B%E5%80%BCssastatic-single-assignment)
      - [实现步骤](#%E5%AE%9E%E7%8E%B0%E6%AD%A5%E9%AA%A4)
      - [应用](#%E5%BA%94%E7%94%A8)
  - [实现](#%E5%AE%9E%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# runtime.KeepAlive


一个变量直至它的最后一次使用期间都是存活的，所以如果你在后面使用一个变量，那么 Go 必须让它一直存活到最后使用的时候.


## 基本知识

### 静态单赋值SSA（Static Single Assignment）


静态单赋值（Static Single Assignment, SSA）是一种中间表示形式，在编译器设计中广泛使用。
SSA 的核心思想是程序中的每个变量在其整个生命周期内只被赋值一次。通过这种方式，SSA 形式能够简化编译器的分析和优化过程


#### 实现步骤
将代码转换为 SSA 形式通常包括以下步骤：

- 构建控制流图（CFG）： 解析原始代码并构建其控制流图。
- 插入 Φ 函数： 在控制流图中找到所有合并节点，为每个变量插入必要的 Φ 函数。
- 重命名变量： 遍历控制流图，为每个变量赋予唯一的名称（即使一个变量在不同位置被重新定义也会有不同的名称）

转换前
```go
x = a + b;
if (c > 0) {
    x = x + 1;
} else {
    x = x - 1;
}
y = x * 2;
```

转换为 SSA 形式后的代码
```go
x1 = a + b;
if (c > 0) {
    x2 = x1 + 1;
} else {
    x3 = x1 - 1;
}
x4 = Φ(x2, x3);
y1 = x4 * 2;
```

#### 应用
SSA 形式在现代编译器中的应用非常广泛，特别是在以下方面：

1. 优化：
   - 常量传播：由于变量在 SSA 中只赋值一次，很容易跟踪常量值并进行替换。
   - 死代码消除：可以更容易识别和移除未使用的变量和代码段。
   - 循环不变代码外提：通过分析变量的唯一赋值点，可以确定哪些代码可以移出循环。
2. 并行化： SSA 形式为自动并行化和依赖关系分析提供了有利条件，因为数据依赖关系变得更加明确。
3. 编译器架构： 许多现代编译器（如 LLVM）采用 SSA 作为其中间表示形式，利用其简化的数据流分析和优化特性来生成高效的目标代码。



## 实现

```go
func KeepAlive(x any) {
	// Introduce a use of x that the compiler can't eliminate.
	// This makes sure x is alive on entry. We need x to be alive
	// on entry for "defer runtime.KeepAlive(x)"; see issue 21402.
	if cgoAlwaysFalse {
		println(x)
	}
}
```

实际实现

```go
// go1.24.3/src/cmd/compile/internal/ssagen/intrinsics.go

func initIntrinsics(cfg *intrinsicBuildConfig) {
	// ... 
    add := func(pkg, fn string, b intrinsicBuilder, archs ...*sys.Arch) {
           intrinsics.addForArchs(pkg, fn, b, archs...)
       }
    
	// ...
	add("runtime", "KeepAlive",
		func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
			data := s.newValue1(ssa.OpIData, s.f.Config.Types.BytePtr, args[0])
			s.vars[memVar] = s.newValue2(ssa.OpKeepAlive, types.TypeMem, data, s.mem())
			return nil
		},
		all...)
	
	// ...
}
```

当你的代码中使用了 runtime.KeepAlive()，Go 编译器会设置一个名为 OpKeepAlive 的静态单赋值(SSA)，然后剩余的编译就会知道将这个变量的存活期保证到使用了 runtime.KeepAlive() 的时刻.



## 参考
- [程序分析与优化 - 7 静态单赋值（SSA）](https://www.cnblogs.com/zhouronghua/p/16390138.html)
- [Go 语言中 runtime.KeepAlive() 方法的一些随笔](https://zhuanlan.zhihu.com/p/213744309)