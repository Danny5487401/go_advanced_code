# io pipe


```go
func Pipe() (*PipeReader, *PipeWriter)
```

io.Pipe会返回一个reader和writer,对 reader 读取（或写入writer）后，进程会被锁住，直到writer有新数据流进入或关闭（或reader把数据读走）

## 修改前
在向其他服务器发送json数据时，都需要先声明一个bytes缓存，然后通过json库把结构体中的内容 marshal 成字节流，再通过Post函数发送。

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  request := new(Person)
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&request)
  if err != nil {
     http.Error(w, err)
  }
  ......
})
```

我们不需要 ioutil.ReadAll 全部 body 再调用 Unmarshal, decoder 内置 buffer 流式解析即可。但是这个例子不完美，有很多问题
- json.NewDecoder 会一直读 r.Body, 未做长度限制
- 没有检查 Content-Type header, 只有 json 才允许 Decode
- 错误处理不够好，error 需要转换，不能直接返回 client

## 修改后
```go
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
    if r.Header.Get("Content-Type") != "" {
        value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            msg := "Content-Type header is not application/json"
            return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
        }
    }

    r.Body = http.MaxBytesReader(w, r.Body, 1048576)

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    err := dec.Decode(&dst)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        switch {
        case errors.As(err, &syntaxError):
            msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}
            ......
        }
    }

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
            msg := "Request body must only contain a single JSON object"
            return &malformedRequest{status: http.StatusBadRequest, msg: msg}
        }
    }
```


## 第三方应用：minio 

上面修改后只是一个 reader 的实现。在 minio 中，经常有 N 多个 io.Reader 或者 io.Writer 组合在一起，实现 io pipeline, 稍复杂一些


minio下载数据
```go
func (api objectAPIHandlers) getObjectHandler(ctx context.Context, objectAPI ObjectLayer, bucket, object string, w http.ResponseWriter, r *http.Request) {
    //......
    //调用后端具体实现
    gr, err := getObjectNInfo(ctx, bucket, object, rs, r.Header, readLock, opts)
    //......
    httpWriter := xioutil.WriteOnClose(w)
    if rs != nil || opts.PartNumber > 0 {
        statusCodeWritten = true
        w.WriteHeader(http.StatusPartialContent)
    }

    // Write object content to response body
    if _, err = xioutil.Copy(httpWriter, gr); err != nil {
    // ......
    }
    // ......
}
```

getObjectNInfo 调用后端具体实现，返回 GetObjectReader gr, 从 gr 中读取数据写回 http Writer

gr 实现有很多种，minio 支持 NAS，FS, EC 多种模式，可以从文件系统中读数据，可以从 remote http 中读取


## 参考资料
1. [io pipeline关于Minio源码](https://mp.weixin.qq.com/s/b_FBTSMZtAw0KqE_3mnMew)