# 错误码生成方式
go generate + stringer
# 安装stringer
```shell script

 go get golang.org/x/tools/cmd/stringer
```
# 使用方式
```go
//go:generate stringer -type ErrCode

//头部添加generate标签
// 执行go generate，会在同一个目录下生成一个文件errcode_string.go。文件名格式是类型名小写_string.go

// 也可以通过-output选项指定输出文件名，例如下面就是指定输出文件名为code_string.go
//go:generate stringer -type ErrCode -output code_string.go

//但是我们更希望的是能返回后面的注释作为错误描述。这就需要使用stringer的-linecomment选项
//go:generate stringer -type ErrCode -linecomment -output code_string.go
```

# 好处
生成的代码做了一些优化，减少了字符串对象的数量