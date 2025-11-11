<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/google/go-cmp](#githubcomgooglego-cmp)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/google/go-cmp

reflect.DeepEqual 的替代品

## 特点
- reflect.DeepEqual不够灵活，无法提供选项实现我们想要的行为，例如允许浮点数误差，对test 不友好 .
- 支持自定义比较函数：你可以编写自定义比较函数，以处理特定类型的值的比较。这允许你在比较复杂的数据结构时定义自己的比较逻辑。
- 不会比较未导出字段（即字段名首字母小写的字段）。遇到未导出字段，cmp.Equal() 直接panic.
- 自定义比较选项：你可以使用 cmp.Options 结构来自定义比较的行为。这包括忽略特定字段、指定自定义比较函数、配置忽略类型的选项等。这使得你可以精确控制比较的方式
- 友好的错误报告：当比较失败时，cmp 生成清晰和有用的错误报告，帮助你理解为什么两个值不相等。这有助于快速识别和修复问题。



## 源码分析

重要结构体
```go
type state struct {
	// These fields represent the "comparison state".
	// Calling statelessCompare must not result in observable changes to these.
	result    diff.Result // The current result of comparison
	curPath   Path        // The current path in the value tree
	curPtrs   pointerPath // The current set of visited pointers
	reporters []reporter  // Optional reporters

	// recChecker checks for infinite cycles applying the same set of
	// transformers upon the output of itself.
	recChecker recChecker

	// dynChecker triggers pseudo-random checks for option correctness.
	// It is safe for statelessCompare to mutate this value.
	dynChecker dynChecker

	// These fields, once set by processOption, will not change.
	exporters []exporter // List of exporters for structs with unexported fields
	opts      Options    // List of all fundamental and filter options
}
```

选项 option 接口用来配置 equal 和 diff：基础option 有 [Ignore], [Transformer], and [Comparer]
```go
type Option interface {
	// filter applies all filters and returns the option that remains.
	// Each option may only read s.curPath and call s.callTTBFunc.
	//
	// An Options is returned only if multiple comparers or transformers
	// can apply simultaneously and will only contain values of those types
	// or sub-Options containing values of those types.
	filter(s *state, t reflect.Type, vx, vy reflect.Value) applicableOption
}

// applicableOption represents the following types:
//
//	Fundamental: ignore | validator | *comparer | *transformer
//	Grouping:    Options
type applicableOption interface {
	Option

	// apply executes the option, which may mutate s or panic.
	apply(s *state, vx, vy reflect.Value)
}
```

案例 Transformer option

```go
type transformer struct {
	core
	name string
	typ  reflect.Type  // T
	fnc  reflect.Value // func(T) R
}

//  \p{L}：表示Unicode字母
//  \p{N}：表示Unicode数字
const identRx = `[_\p{L}][_\p{L}\p{N}]*`

var identsRx = regexp.MustCompile(`^` + identRx + `(\.` + identRx + `)*$`)

func Transformer(name string, f interface{}) Option {
	v := reflect.ValueOf(f)
    // 判断函数的格式是否符合 trFunc  // func(T) R
	if !function.IsType(v.Type(), function.Transformer) || v.IsNil() {
		panic(fmt.Sprintf("invalid transformer function: %T", f))
	}
	if name == "" {
		name = function.NameOf(v)
		if !identsRx.MatchString(name) {
			name = "λ" // Lambda-symbol as placeholder name
		}
	} else if !identsRx.MatchString(name) {
		// 转换器名字要求判断
		panic(fmt.Sprintf("invalid name: %q", name))
	}
	tr := &transformer{name: name, fnc: reflect.ValueOf(f)}
	if ti := v.Type().In(0); ti.Kind() != reflect.Interface || ti.NumMethod() > 0 {
		tr.typ = ti
	}
	return tr
}
```





相等判断入口
```go
// github.com/google/go-cmp@v0.6.0/cmp/compare.go
func Equal(x, y interface{}, opts ...Option) bool {
	// state 初始化
	s := newState(opts)
	// 对比
	s.compareAny(rootStep(x, y))
	// 判断对比结果的NumDiff是否为0
	return s.result.Equal()
}

```
具体根据类型判断
```go
func (s *state) compareAny(step PathStep) {
	// Update the path stack.
	s.curPath.push(step)
	defer s.curPath.pop()
	for _, r := range s.reporters {
		r.PushStep(step)
		defer r.PopStep()
	}
	s.recChecker.Check(s.curPath)

	// Cycle-detection for slice elements (see NOTE in compareSlice).
	t := step.Type()
	vx, vy := step.Values()
	if si, ok := step.(SliceIndex); ok && si.isSlice && vx.IsValid() && vy.IsValid() {
		px, py := vx.Addr(), vy.Addr()
		if eq, visited := s.curPtrs.Push(px, py); visited {
			s.report(eq, reportByCycle)
			return
		}
		defer s.curPtrs.Pop(px, py)
	}

	// Rule 1: Check whether an option applies on this node in the value tree.
	if s.tryOptions(t, vx, vy) {
		return
	}

	// Rule 2: Check whether the type has a valid Equal method.
	// 调用自定义 Equal 方法
	if s.tryMethod(t, vx, vy) {
		return
	}

	// Rule 3: Compare based on the underlying kind.
	// 根据 kind 类型金习惯判断
	switch t.Kind() {
	case reflect.Bool:
		s.report(vx.Bool() == vy.Bool(), 0)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.report(vx.Int() == vy.Int(), 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s.report(vx.Uint() == vy.Uint(), 0)
	case reflect.Float32, reflect.Float64:
		s.report(vx.Float() == vy.Float(), 0)
	case reflect.Complex64, reflect.Complex128:
		s.report(vx.Complex() == vy.Complex(), 0)
	case reflect.String:
		s.report(vx.String() == vy.String(), 0)
	case reflect.Chan, reflect.UnsafePointer:
		s.report(vx.Pointer() == vy.Pointer(), 0)
	case reflect.Func:
		s.report(vx.IsNil() && vy.IsNil(), 0)
	case reflect.Struct:
		s.compareStruct(t, vx, vy)
	case reflect.Slice, reflect.Array:
		s.compareSlice(t, vx, vy)
	case reflect.Map:
		s.compareMap(t, vx, vy)
	case reflect.Ptr:
		s.comparePtr(t, vx, vy)
	case reflect.Interface:
		s.compareInterface(t, vx, vy)
	default:
		panic(fmt.Sprintf("%v kind not handled", t.Kind()))
	}
}
```

这里判断结构体作为案例

```go
func (s *state) compareStruct(t reflect.Type, vx, vy reflect.Value) {
	var addr bool
	var vax, vay reflect.Value // Addressable versions of vx and vy

	var mayForce, mayForceInit bool
	step := StructField{&structField{}}
	for i := 0; i < t.NumField(); i++ {
		step.typ = t.Field(i).Type
		step.vx = vx.Field(i)
		step.vy = vy.Field(i)
		step.name = t.Field(i).Name
		step.idx = i
		// 判断是否是未导出子段
		step.unexported = !isExported(step.name)
		if step.unexported {
			if step.name == "_" {
				continue
			}
			// Defer checking of unexported fields until later to give an
			// Ignore a chance to ignore the field.
			if !vax.IsValid() || !vay.IsValid() {
				// For retrieveUnexportedField to work, the parent struct must
				// be addressable. Create a new copy of the values if
				// necessary to make them addressable.
				addr = vx.CanAddr() || vy.CanAddr()
				vax = makeAddressable(vx)
				vay = makeAddressable(vy)
			}
			if !mayForceInit {
				for _, xf := range s.exporters {
					mayForce = mayForce || xf(t)
				}
				mayForceInit = true
			}
			step.mayForce = mayForce
			step.paddr = addr
			step.pvx = vax
			step.pvy = vay
			step.field = t.Field(i)
		}
		// 递归调用
		s.compareAny(step)
	}
}

```





## 参考

- [每日一库 go-cmp](https://darjun.github.io/2020/03/20/godailylib/go-cmp/)