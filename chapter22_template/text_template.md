# Text/template

实现数据驱动模板以生成文本输出，可以理解为一组文字按照特定格式动态嵌入另一组文字中

## 输入
模板的输入文本是任何格式的UTF-8编码文本。 {{ 和 }} 包裹的内容统称为 action，分为两种类型：
- 数据求值（data evaluations）
- 控制结构（control structures）

### 控制结构
循环操作：pipeline 的值必须是数组、切片、字典和通道中的一种，即可迭代类型的值，根据值的长度输出多个 T1
```go
{{ range pipeline }} T1 {{ end }}
*// 这个 else 比较有意思，如果 pipeline 的长度为 0 则输出 else 中的内容*
{{ range pipeline }} T1 {{ else }} T0 {{ end }}
*// 获取容器的下标*
{{ range $index, $value := pipeline }} T1 {{ end }}

// 案例
{{range .Field}}
  {{.ChildFieldOne}}  -- {{.ChildFieldTwo }}
{{ end }}
```
条件语句
```go
{{ if pipeline }} T1 {{ end }}
{{ if pipeline }} T1 {{ else }} T0 {{ end }}
{{ if pipeline }} T1 {{ else if pipeline }} T0 {{ end }}
*// 上面的语法其实是下面的简写*
{{ if pipeline }} T1 {{ else }}{{ if pipeline }} T0 { {end }}{{ end }}
{{ if pipeline }} T1 {{ else if pipeline }} T2 {{ else }} T0 {{ end }}

```

### 注释
```go
{{*/\* comment \*/*}}

```

### 裁剪空格
```go
// 裁剪 content 前后的空格
{{- content -}}

// 裁剪 content 前面的空格
{{- content }}

// 裁剪 content 后面的空格
{{ content -}}

```
### 管道函数
```css
用法1：

{{FuncName1}}

此标签将调用名称为“FuncName1”的模板函数（等同于执行“FuncName1()”，不传递任何参数）并输出其返回值。

用法2：

{{FuncName1 "参数值1" "参数值2"}}

此标签将调用“FuncName1("参数值1", "参数值2")”，并输出其返回值

用法3：

{{.Admpub|FuncName1}}

此标签将调用名称为“FuncName1”的模板函数（等同于执行“FuncName1(this.Admpub)”，将竖线“|”左边的“.Admpub”变量值作为函数参数传送）并输出其返回值。

```
### 文本输出
pipeline 代表的数据会产生与调用 fmt.Print 函数类似的输出，例如整数类型的 3 会转换成字符串 “3” 输出。
```go
{{ pipeline }}
```


## 模板函数
即：可以对某个字段使用函数操作。适用于稍微复杂的字段处理。

```go
type FuncMap map[string]interface{}
t = t.Funcs(template.FuncMap{"handleFieldName": HandleFunc})

```
内置模板函数
```go
var builtins = FuncMap{
    "and":      and,
    "call":     call,
    "html":     HTMLEscaper,
    "index":    index,
    "js":       JSEscaper,
    "len":      length,
    "not":      not,
    "or":       or,
    "print":    fmt.Sprint,
    "printf":   fmt.Sprintf,
    "println":  fmt.Sprintln,
    "urlquery": URLQueryEscaper,
}

```

## 第三方应用
- kratos应用：https://github.com/go-kratos/kratos/blob/main/cmd/kratos/internal/proto/add/template.go