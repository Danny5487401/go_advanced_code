<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/samber/lo](#githubcomsamberlo)
  - [案例](#%E6%A1%88%E4%BE%8B)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

#  github.com/samber/lo


samber/lo是一个用于实现函数式编程的Go语言库。该库相对于其他使用反射来实现的库来说，更加的快、同时还安全。
它提供了许多用于操作和转换数据的函数，以及一些实用的工具函数。

它提供了切片的许多辅助函数。例如：Filter、Slice、Fill、Map、FilterMap、FlatMap、GroupBy、PartitionBy等，还提供了类似Java中的try-catch机制的异常处理函数，例如：Try、TryWithErrorValue、TryCatch等



## 案例

filter 
```go
fmt.Println(lo.Filter[string]([]string{"hello", "good bye", "world", "fuck", "fuck who"}, func(s string, _ int) bool {
		return !strings.Contains(s, "fuck")
	}))
```
```go
func Filter[V any](collection []V, predicate func(item V, index int) bool) []V {
	result := make([]V, 0, len(collection))

	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}
```

