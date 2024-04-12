<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [错误码 自动化生成方式](#%E9%94%99%E8%AF%AF%E7%A0%81-%E8%87%AA%E5%8A%A8%E5%8C%96%E7%94%9F%E6%88%90%E6%96%B9%E5%BC%8F)
  - [安装stringer](#%E5%AE%89%E8%A3%85stringer)
  - [使用方式](#%E4%BD%BF%E7%94%A8%E6%96%B9%E5%BC%8F)
    - [注意点](#%E6%B3%A8%E6%84%8F%E7%82%B9)
  - [好处](#%E5%A5%BD%E5%A4%84)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [工具 github 源码](#%E5%B7%A5%E5%85%B7-github-%E6%BA%90%E7%A0%81)
    - [工具生成的代码](#%E5%B7%A5%E5%85%B7%E7%94%9F%E6%88%90%E7%9A%84%E4%BB%A3%E7%A0%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 错误码 自动化生成方式
go generate + stringer

stringer命令旨在自动创建满足 fmt.Stringer 的方法。 它为指定类型生成String()并将其描述为字符串。常可用于​定义错误码时同时生成错误信息等场景。

## 安装stringer
```shell script
go install golang.org/x/tools/cmd/stringer
```


## 使用方式
```go
//go:generate stringer -type ErrCode

//头部添加generate标签
// 执行go generate，会在同一个目录下生成一个文件errcode_string.go。文件名格式是类型名小写_string.go

// 也可以通过-output选项指定输出文件名，例如下面就是指定输出文件名为code_string.go
//go:generate stringer -type ErrCode -output code_string.go

//但是我们更希望的是能返回后面的注释作为错误描述。这就需要使用stringer的-linecomment选项
//go:generate stringer -type ErrCode -linecomment -output code_string.go
```

解析

* -type ：指定stringer命令作用的类型名

* -output选项：指定输出文件

* -linecomment 选项：返回后面的注释作为错误描述



### 注意点
1. go:generate前面只能使用//注释，注释必须在行首，前面不能有空格且//与go:generate之间不能有空格！！！
2. go:generate可以在任何 Go 源文件中，最好在类型定义的地方


## 好处
生成的代码做了一些优化，减少了字符串对象的数量

## 源码分析
### [工具 github 源码](https://github.com/golang/tools/blob/627959a8e32af98dce9ff3e65ef6f491d3dcb9f6/cmd/stringer/stringer.go)




### 工具生成的代码
1. 使用硬编码值生成String()函数效率更高：
```go
var (
	_ErrCode_index_2 = [...]uint8{0, 15, 36, 54}
	_ErrCode_index_3 = [...]uint8{0, 12, 33}
)

func (i ErrCode) String() string {
	switch {
	case i == 200:
		return _ErrCode_name_0
	case i == 1000000:
		return _ErrCode_name_1
	case 3000000 <= i && i <= 3000002:
		i -= 3000000
		return _ErrCode_name_2[_ErrCode_index_2[i]:_ErrCode_index_2[i+1]]
	case 4000000 <= i && i <= 4000001:
		i -= 4000000
		return _ErrCode_name_3[_ErrCode_index_3[i]:_ErrCode_index_3[i+1]]
	default:
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

```

基准更快，使用 map 的速度要慢得多，因为它必须进行函数调用，并且存储桶中的查找不像访问切片索引那样简单。


2. 自检
```go
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OK-200]
	_ = x[ServerCommonError-1000000]
	_ = x[ServerInvalidParams-1000000]
	_ = x[ServerTimout-1000000]
	_ = x[TicketNotExit-3000000]
	_ = x[TicketStatusNotOK-3000001]
	_ = x[TicketUpdateFail-3000002]
	_ = x[BookNotFoundError-4000000]
	_ = x[BookHasBeenBorrowedError-4000001]
}

```


## 参考




