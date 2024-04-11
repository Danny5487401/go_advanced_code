<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [html/template](#htmltemplate)
  - [html/template包上下文感知](#htmltemplate%E5%8C%85%E4%B8%8A%E4%B8%8B%E6%96%87%E6%84%9F%E7%9F%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



# html/template


## html/template包上下文感知

上下文感知具体指的是根据所处环境css、js、html、url的path、url的query，自动进行不同格式的转义。

上下文感知的自动转义能让程序更加安全，比如防止XSS攻击(例如在表单中输入带有<script>...</script>的内容并提交，会使得用户提交的这部分script被执行)

```go
func process(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("test.html")
    content := `I asked: <i>"What's up?"</i>`
    t.Execute(w, content)
}

```
上面content是Execute的第二个参数，它的内容是包含了特殊符号的字符串

```html
// test.html
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
        <title>Go Web Programming</title>
    </head>
    <body>
        <div>{{ . }}</div>
        <div><a href="/{{ . }}">Path</a></div>
        <div><a href="/?q={{ . }}">Query</a></div>
        <div><a onclick="f('{{ . }}')">Onclick</a></div>
    </body>
</html>
```

上面test.html中有4个不同的环境，分别是html环境、url的path环境、url的query环境以及js环境。虽然对象都是{{.}}，但解析执行后的值是不一样的。


使用curl获取源代码
```html
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>Go Web Programming</title>
</head>

<body>
    <div>I asked: &lt;i&gt;&#34;What&#39;s up?&#34;&lt;/i&gt;</div>
    <div>
        <a href="/I%20asked:%20%3ci%3e%22What%27s%20up?%22%3c/i%3e">
            Path
        </a>
    </div>
    <div>
        <a href="/?q=I%20asked%3a%20%3ci%3e%22What%27s%20up%3f%22%3c%2fi%3e">
            Query
        </a>
    </div>
    <div>
        <a onclick="f('I asked: \x3ci\x3e\x22What\x27s up?\x22\x3c\/i\x3e')">
            Onclick
        </a>
    </div>
</body>

</html>
```