<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [http.Request源码分析](#httprequest%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [Request](#request)
  - [1.错误类型](#1%E9%94%99%E8%AF%AF%E7%B1%BB%E5%9E%8B)
  - [2.结构体定义](#2%E7%BB%93%E6%9E%84%E4%BD%93%E5%AE%9A%E4%B9%89)
  - [3. request请求头的一些字段的修改方法](#3-request%E8%AF%B7%E6%B1%82%E5%A4%B4%E7%9A%84%E4%B8%80%E4%BA%9B%E5%AD%97%E6%AE%B5%E7%9A%84%E4%BF%AE%E6%94%B9%E6%96%B9%E6%B3%95)
  - [4. 请求体的处理方法（针对post请求）：生成一个读取器(以及内部实现)](#4-%E8%AF%B7%E6%B1%82%E4%BD%93%E7%9A%84%E5%A4%84%E7%90%86%E6%96%B9%E6%B3%95%E9%92%88%E5%AF%B9post%E8%AF%B7%E6%B1%82%E7%94%9F%E6%88%90%E4%B8%80%E4%B8%AA%E8%AF%BB%E5%8F%96%E5%99%A8%E4%BB%A5%E5%8F%8A%E5%86%85%E9%83%A8%E5%AE%9E%E7%8E%B0)
  - [5. request 写入方法](#5-request-%E5%86%99%E5%85%A5%E6%96%B9%E6%B3%95)
  - [6. 根据读取器 读出一个请求对象(以及内部实现)方法自定义request请求](#6-%E6%A0%B9%E6%8D%AE%E8%AF%BB%E5%8F%96%E5%99%A8-%E8%AF%BB%E5%87%BA%E4%B8%80%E4%B8%AA%E8%AF%B7%E6%B1%82%E5%AF%B9%E8%B1%A1%E4%BB%A5%E5%8F%8A%E5%86%85%E9%83%A8%E5%AE%9E%E7%8E%B0%E6%96%B9%E6%B3%95%E8%87%AA%E5%AE%9A%E4%B9%89request%E8%AF%B7%E6%B1%82)
  - [7. 根据 读取器(*bufio.Reader)，读取一个请求实例方法（readRequest内部方法实现），以及读取器的相关方法（新建一个最大字节的读取器，以及读取器结构体对象的读取关闭方法](#7-%E6%A0%B9%E6%8D%AE-%E8%AF%BB%E5%8F%96%E5%99%A8bufioreader%E8%AF%BB%E5%8F%96%E4%B8%80%E4%B8%AA%E8%AF%B7%E6%B1%82%E5%AE%9E%E4%BE%8B%E6%96%B9%E6%B3%95readrequest%E5%86%85%E9%83%A8%E6%96%B9%E6%B3%95%E5%AE%9E%E7%8E%B0%E4%BB%A5%E5%8F%8A%E8%AF%BB%E5%8F%96%E5%99%A8%E7%9A%84%E7%9B%B8%E5%85%B3%E6%96%B9%E6%B3%95%E6%96%B0%E5%BB%BA%E4%B8%80%E4%B8%AA%E6%9C%80%E5%A4%A7%E5%AD%97%E8%8A%82%E7%9A%84%E8%AF%BB%E5%8F%96%E5%99%A8%E4%BB%A5%E5%8F%8A%E8%AF%BB%E5%8F%96%E5%99%A8%E7%BB%93%E6%9E%84%E4%BD%93%E5%AF%B9%E8%B1%A1%E7%9A%84%E8%AF%BB%E5%8F%96%E5%85%B3%E9%97%AD%E6%96%B9%E6%B3%95)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# http.Request源码分析

## Request
request 表示由服务器接收或由客户端发送的HTTP请求，例如客户端(client)在发送各种请求时，需要先新建一个请求对象，然后调用一些请求的方法开始自定义一些配置，服务端监听到该请求便会做出相应的应答

## 1.错误类型
```go
const (
	defaultMaxMemory = 32 << 20 // 32 MB			// 默认最大内存32 MB
)

var ErrMissingFile = errors.New("http: no such file")// 提供的文件字段名不在请求中或非文件字段时，FormFile返回

// 弃用的省略（并非http包中与协议错误相关的所有错误都属于ProtocolError类型。）
type ProtocolError struct {
	ErrorString string
}
func (pe *ProtocolError) Error() string { return pe.ErrorString }
var (
	// Pusher实现的Push方法返回ErrNotSupported，以指示HTTP/2push支持不可用
	ErrNotSupported = &ProtocolError{"feature not supported"}
	// 当请求的内容类型不包含“boundary”参数时，request.MultipartReader返回
	ErrMissingBoundary = &ProtocolError{"no multipart boundary param in Content-Type"}
	// Content-Type不是multipart/form-data时，request.MultipartReader返回
	ErrNotMultipart = &ProtocolError{"request Content-Type isn't multipart/form-data"}
)  
func badStringError(what, val string) error { return fmt.Errorf("%s %q", what, val) }

```
代码风格

- ErrMissingFile：这是请求对象的一个查找方法返回的错误，调用了errors包的New方法来返回一个错误类型的错误（即errors.New可以传string类型参数 来返回一个简单的字符串错误）
- ProtocolError ：ProtocolError表示HTTP协议错误结构体对象。该对象 实现了error接口(Error() string)，因此，所有该结构体实例都可以当成error返回传参
- badStringError：是一个根据两个 string参数返回一个error的函数，当必要时，可以调用此函数实现返回错误

## 2.结构体定义

```go
// Headers that Request.Write 处理自身应跳过
var reqWriteExcludeHeader = map[string]bool{
	"Host":              true, // 反正不在Header的map中
	"User-Agent":        true,
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Trailer":           true,
}

// Request 请求结构体：在客户机和服务器使用之间，字段语义略有不同。 除了以下字段的注释外，请参阅文档以了解Request.Write and RoundTripper
type Request struct {
	// 指定发送的HTTP请求的方法（GET、POST、PUT等），Go的HTTP客户端不支持使用CONNECT方法发送请求
	Method string

	// URL 指定被请求的URI（对于服务器请求）或要访问的URL（对于客户端请求，就是要连接的服务器地址）
	URL *url.URL

	// 入服务器请求的协议版本（主版本号.次版本号）
	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	// 请求头(为请求报文添加了一些附加信息)，不区分大小写，对于客户端请求，某些头（如内容长度和连接）会在需要时自动写入
	Header Header

	// 请求体，get请求没有请求体-Body字段（get请求的参数都在url里）
	Body io.ReadCloser

	// 获取请求体的副本
	GetBody func() (io.ReadCloser, error)

	// 请求Body的大小（字节数）
	ContentLength int64

	// 列出从最外层到最内层的传输编码。空列表表示“身份”编码，通常可以忽略；在发送和接收请求时，根据需要自动添加和删除分块编码。
	TransferEncoding []string

	// 在回复了此次请求后结束连接。对服务端来说就是回复了这个 request ，对客户端来说就是收到了 response
	// 且 对于服务端是Handlers 会自动调用Close, 对客户端，如果设置长连接(Transport.DisableKeepAlives=false)，就不会关闭。没设置就关闭
	Close bool

	// 对于服务器请求，Host指定查找URL的主机
	Host string

	// Form 包含解析的表单数据，包括URL字段的查询参数和补丁、POST或PUT表单数据, 此字段仅在调用ParseForm后可用。
	Form url.Values

	// PostForm 包含来自PATCH、POST或PUT body参数的解析表单数据, 此字段仅在调用ParseForm后可用。HTTP客户机忽略PostForm，而是使用Body。
	PostForm url.Values

	// MultipartForm 是经过解析的多部件 表单，包括文件上传. 此字段仅在解析表单即 调用ParseMultipartForm后可用。 HTTP客户机忽略MultipartForm，而是使用Body。
	MultipartForm *multipart.Form

	// Trailer 指定在 当body部分发送完成之后 发送的附加头
	// 对于服务器请求，尾部映射最初只包含尾部键，值为nil。（客户端声明它稍后将发送哪些trailers）当处理程序从主体读取时，它不能引用trailers
	// 从Body读取返回EOF后，可以再次读取Trailer并包含非nil值
	// 对于客户端请求，必须将拖车初始化为包含trailer键的map，以便稍后发送。
	// 很少有HTTP客户机、服务器或代理支持HTTP Trailer
	Trailer Header

	// RemoteAddr 允许HTTP服务器和其他软件记录发送请求的网络地址，通常用于日志记录。
	// 此字段不是由ReadRequest填写的，并且没有定义的格式。此包中的HTTP服务器在调用处理程序之前将RemoteAddr设置为“IP:port”地址。 HTTP客户端将忽略此字段。
	RemoteAddr string

	// RequestURI 是客户端发送到服务器的请求行（RFC 7230，第3.1.1节）的未修改的请求目标。通常应该改用URL字段。在HTTP客户端请求中设置此字段是错误的。
	RequestURI string

	// TLS 允许HTTP服务器和其他软件记录有关接收请求的TLS连接的信息。此字段不是由ReadRequest填写的。HTTP客户端将忽略此字段
	TLS *tls.ConnectionState

	// Cancel 是一个可选通道，它的关闭表示客户端请求应被视为已取消
	// 不推荐：改为使用NewRequestWithContext设置请求的上下文
	Cancel <-chan struct{}

	// Response是导致创建此请求的重定向响应。此字段仅在客户端重定向期间填充
	Response *Response

	// 请求的上下文。（修改时，通过使用WithContext复制整个请求来修改它）
	// 对于传出的客户机请求，上下文控制请求及其响应的整个生存期：获取连接、发送请求以及读取响应头和主体
	ctx context.Context
}

```
获取内部成员ctx(上下文)的方法，2个修改上下文方法(都用到了内部自定义clone函数)
```go
func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	// 1、首先，使用new(type), 返回一个 指向 Request类型的 零值 的指针
	r2 := new(Request)
	// 2、将*Request 赋值给 r2的指针（即将 指向请求数据的指针，赋值给新变量的指针。此时，新变量的指针也是指向原请求，即新变量r2是一个 数据和r相同的变量）
	*r2 = *r
	// 3、把r2的上下文修改成新上下文
	r2.ctx = ctx
	// 4、将 原请求的Url克隆进新请求的Url (cloneURL的本质就是上面指针修改方法, 内部还进行了user属性的克隆)
	r2.URL = cloneURL(r.URL)
	// 5、返回新的 请求
	return r2
}

func (r *Request) Clone(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	r2 := new(Request)
	*r2 = *r
	r2.ctx = ctx
	r2.URL = cloneURL(r.URL)

	// 对比上面方法，这里新增其他 属性的克隆：Header、Trailer、TransferEncoding、Form、PostForm、MultipartForm
	if r.Header != nil {
		r2.Header = r.Header.Clone()
	}
	if r.Trailer != nil {
		r2.Trailer = r.Trailer.Clone()
	}
	if s := r.TransferEncoding; s != nil {
		s2 := make([]string, len(s))
		copy(s2, s)
		r2.TransferEncoding = s2
	}
	r2.Form = cloneURLValues(r.Form)
	r2.PostForm = cloneURLValues(r.PostForm)
	r2.MultipartForm = cloneMultipartForm(r.MultipartForm)
	return r2
}

// ProtoAtLeast 报告请求中使用的HTTP协议是否至少 major.minor.（入服务器请求的协议版本）
func (r *Request) ProtoAtLeast(major, minor int) bool {
	return r.ProtoMajor > major ||
		r.ProtoMajor == major && r.ProtoMinor >= minor
}


```

- Context()：获取对象内部成员，可以使用一个公开的方法封装实现
- WithContext()： 更改请求的上下文方法一，传入新的上下文，返回修改后的请求–r的浅层副本，该方法 通过新建一个变量，使变量的指针指向原 请求，然后克隆请求的URL实现。此方法很少用。要使用上下文创建新请求，请使用NewRequestWithContext
- Clone()： 更改请求的上下文方法二，传入新的上下文，返回修改后的请求–r的深度副本， 通过除了克隆请求的URL，还有对请求结构里的请求头等一一克隆进一个新的请求结构体 实现，一般用此方法

## 3. request请求头的一些字段的修改方法
```go
// 如果在请求中发送，UserAgent将返回 发送请求的应用程序名称。
func (r *Request) UserAgent() string {
	return r.Header.Get("User-Agent")
}

// Cookies （所有cookie）解析并返回与请求一起发送的HTTP Cookies（其实就是调用请求头的解析cookie方法）
func (r *Request) Cookies() []*Cookie {
	return readCookies(r.Header, "")
}

// 当找不到Cookie时，请求的Cookie方法将返回ErrNoCookie
var ErrNoCookie = errors.New("http: named cookie not present")

// Cookie （单个cookie）返回请求中提供的指定Cookie，如果找不到则返回ErrNoCookie
// 调用请求头的过滤器解析cookie的方法，实现（入参name就是过滤器,即指定的cookie命名）
func (r *Request) Cookie(name string) (*Cookie, error) {
	for _, c := range readCookies(r.Header, name) {
		return c, nil
	}
	return nil, ErrNoCookie
}

// AddCookie将 cookie添加到请求中。
func (r *Request) AddCookie(c *Cookie) {
	// 1、清理入参cookie（只清理入参cookie的名称和值，不清理请求中已经存在的Cookie头。）
	s := fmt.Sprintf("%s=%s", sanitizeCookieName(c.Name), sanitizeCookieValue(c.Value))

	// 2、获取原请求头的cookie，然后设置为 ： c;s (不会附加多个Cookie头字段。 都写在同一行中，用；分隔)
	if c := r.Header.Get("Cookie"); c != "" {
		r.Header.Set("Cookie", c+"; "+s)
	} else {
		r.Header.Set("Cookie", s)
	}
}

// 如果在请求中发送，Referer返回引用URL（提供了Request的上下文信息的服务器，告诉服务器我是从哪个链接过来的，比如从我主页上链接到一个朋友那里， 他的服务器就能够从HTTP Referer中统计出每天有多少用户点击我主页上的链接访问 他的网站）
// Referer的拼写与请求本身一样错误，这是HTTP最早出现的一个错误。
func (r *Request) Referer() string {
	return r.Header.Get("Referer")
}

```
## 4. 请求体的处理方法（针对post请求）：生成一个读取器(以及内部实现)
```go
// multipartByReader 是一个sentinel(哨兵)值
// 它在Request.MultipartForm里的存在，表明了 请求体的解析 已经被 传递给了 一个MultipartReader实例，而不是还需要解析-ParseMultipartForm
var multipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),		// FileHeader描述一个multipart请求的（一个）文件记录的信息。
}
// 生成一个 多部件表单 读取器（内部multipartReader方法的封装）
// 实现过程： 只有在 请求体 字段MultipartForm 为nil时，才会调用生成 读取器
// 使用这个函数可 代替 请求体解析方法ParseMultipartForm() 将请求体作为流进行处理（即要么 解析处理，要么使用读取器处理）
func (r *Request) MultipartReader() (*multipart.Reader, error) {
	// 1、判断 如果 请求体 多部件表单字段，已经是一个 解析过的表单类型，说明已经调用了一次，不需要再处理
	if r.MultipartForm == multipartByReader {
		return nil, errors.New("http: MultipartReader called twice")
	}
	// 2、如果 请求体 多部件表单字段 不为空，说明 已经被解析了
	if r.MultipartForm != nil {
		return nil, errors.New("http: multipart handled by ParseMultipartForm")
	}
	// 3、只有当 该字段为空时， 就会先 将该字段设置为 多部件表单类型，然后调用 生成读取器方法来 返回
	r.MultipartForm = multipartByReader
	// 4、返回 一个读取器
	return r.multipartReader(true)
}

// 内部生成 读取器的过程：
// 主要是调用 multipart.NewReader()方法实现
func (r *Request) multipartReader(allowMixed bool) (*multipart.Reader, error) {
	// 1、获取、解析请求头里设置的 请求内容类型
	v := r.Header.Get("Content-Type")
	if v == "" {
		return nil, ErrNotMultipart
	}

	// 2、将设置的 类型 转换为小写后判断 是否是多部件类型 ，因为只有这个类型才需要生成 读取器
	d, params, err := mime.ParseMediaType(v)
	if err != nil || !(d == "multipart/form-data" || allowMixed && d == "multipart/mixed") {
		return nil, ErrNotMultipart
	}

	// 3、获取 内容类型的键为boundary的值
	boundary, ok := params["boundary"]
	if !ok {
		return nil, ErrMissingBoundary
	}

	// 4、（使用给定的MIMIE边缘值-boundary）创建一个新的 请求体 读取的多元读取器
	return multipart.NewReader(r.Body, boundary), nil
}

```

## 5. request 写入方法
```go
// 根据 写入器io.Writer 编写http请求，调用了内部write(),传参：不使用代理、空的附加头、不等待
// 它是头和正文; 此方法引用请求的字段：Host、URL、Method、Header、ContentLength、TransferEncoding、Body
func (r *Request) Write(w io.Writer) error {
	return r.write(w, false, nil, nil)
}

// WriteProxy与Write类似，但以HTTP代理所期望的形式写入请求。
func (r *Request) WriteProxy(w io.Writer) error {
	return r.write(w, true, nil, nil)
}

// 当请求中没有主机或URL时，通过Write返回errMissingHost.
var errMissingHost = errors.New("http: Request.Write on Request with no Host or URL set")

// 根据写入器 编写 请求，内部实现
// 除了写入器，可能会附加使用代理、附加请求头、是否等待参数
func (r *Request) write(w io.Writer, usingProxy bool, extraHeaders Header, waitForContinue func() bool) (err error) {
	// 1、根据请求的上下文 返回一个跟踪对象
	trace := httptrace.ContextClientTrace(r.Context())

	// 2、程序结束之后，调用 WroteRequest()函数来传入 任何错误进跟踪器
	if trace != nil && trace.WroteRequest != nil {
		defer func() {
			trace.WroteRequest(httptrace.WroteRequestInfo{
				Err: err,
			})
		}()
	}

	// 3、找到目标主机。首选Host:header，但如果没有给出，则使用请求URL中的Host。清理主机，以防主机中有意外的东西
	host := cleanHost(r.Host)
	if host == "" {
		if r.URL == nil {
			return errMissingHost
		}
		host = cleanHost(r.URL.Host)    // 这里判断了如果 请求的host为空，则使用 请求URL的host
	}

	// 4、host清理：HTTP客户机、代理或其他中介必须删除附加到传出URI的任何IPv6区域标识符。
 	host = removeZone(host)

	// 5、如果又代理，组合uri：uri = Scheme://host+r.URL.RequestURI()
	ruri := r.URL.RequestURI()
	if usingProxy && r.URL.Scheme != "" && r.URL.Opaque == "" {
		ruri = r.URL.Scheme + "://" + host + ruri
		// 6、如果没有代理， ruri = host
	} else if r.Method == "CONNECT" && r.URL.Path == "" {
		// CONNECT请求通常只提供主机和端口，而不是完整的URL
		ruri = host
		if r.URL.Opaque != "" {
			ruri = r.URL.Opaque
		}
	}
	if stringContainsCTLByte(ruri) {
		return errors.New("net/http: can't write control character in Request.URL")
	}
	
	var bw *bufio.Writer
	// 7、写入器转 字节写入器,再转缓冲写入器
	if _, ok := w.(io.ByteWriter); !ok {
		bw = bufio.NewWriter(w)
		w = bw
	}

	_, err = fmt.Fprintf(w, "%s %s HTTP/1.1\r\n", valueOrDefault(r.Method, "GET"), ruri)
	if err != nil {
		return err
	}

	// 8、标题行 写入追踪缓冲
	_, err = fmt.Fprintf(w, "Host: %s\r\n", host)
	if err != nil {
		return err
	}
	if trace != nil && trace.WroteHeaderField != nil {
		// WroteHeaderField在Transport  写入每个请求头之后被调用。在调用时，这些值可能已被缓冲，尚未写入网络
		trace.WroteHeaderField("Host", []string{host})
	}

	// 9、 使用defaultUserAgent，除非标头包含一个，否则该标头可能为空以不发送标头
	// 代理写入 追踪缓冲
	userAgent := defaultUserAgent
	if r.Header.has("User-Agent") {
		userAgent = r.Header.Get("User-Agent")
	}
	if userAgent != "" {
		_, err = fmt.Fprintf(w, "User-Agent: %s\r\n", userAgent)
		if err != nil {
			return err
		}
		if trace != nil && trace.WroteHeaderField != nil {
			trace.WroteHeaderField("User-Agent", []string{userAgent})
		}
	}

	// 进程主体，ContentLength，结束，尾部
	tw, err := newTransferWriter(r)		// 根据请求new一个新的 传输写入器
	if err != nil {
		return err
	}

	// 10、传输写入器，调用writeHeader()方法，把追踪器里的请求头写入传入的写入器里
	err = tw.writeHeader(w, trace)
	if err != nil {
		return err
	}

	// 11、调用请求头的writeSubset()方法，用wire格式写入Header。
	err = r.Header.writeSubset(w, reqWriteExcludeHeader, trace)
	if err != nil {
		return err
	}

	// 12、如果附加请求头存在，也写入
	if extraHeaders != nil {
		err = extraHeaders.write(w, trace)
		if err != nil {
			return err
		}
	}

	// 13、将"\r\n" 写入w
	_, err = io.WriteString(w, "\r\n")
	if err != nil {
		return err
	}

	if trace != nil && trace.WroteHeaders != nil {
		trace.WroteHeaders()
	}

	// 14、刷新并等待100
	if waitForContinue != nil {
		if bw, ok := w.(*bufio.Writer); ok {
			err = bw.Flush()
			if err != nil {
				return err
			}
		}
		if trace != nil && trace.Wait100Continue != nil {
			trace.Wait100Continue()
		}
		if !waitForContinue() {
			r.closeBody()
			return nil
		}
	}

	if bw, ok := w.(*bufio.Writer); ok && tw.FlushHeaders {
		if err := bw.Flush(); err != nil {
			return err
		}
	}

	// 15、写请求体
	err = tw.writeBody(w)
	if err != nil {
		if tw.bodyReadError == err {
			err = requestBodyReadError{err}
		}
		return err
	}

	if bw != nil {
		return bw.Flush()
	}
	return nil
}


```
## 6. 根据读取器 读出一个请求对象(以及内部实现)方法自定义request请求
新建一个请求的方法（使用后台空上下文包装NewRequestWithContext）
```go
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	return NewRequestWithContext(context.Background(), method, url, body)
}

// 使用给定的上下文 新建一个请求方法
// 如果提供的主体也是io.Closer, 返回的 请求体被设置为body 后就会被客户端 的 Do\Post\PostForm\Transport.RoundTrip关闭
// 要创建用于测试服务器处理程序的请求，请使用net/http/httptest包中的NewRequest函数、ReadRequest或手动更新请求字段
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	// 1、默认 Get请求
	if method == "" {
		method = "GET"
	}
	// 2、method、ctx 参数检查
	if !validMethod(method) {
		return nil, fmt.Errorf("net/http: invalid method %q", method)
	}
	if ctx == nil {
		return nil, errors.New("net/http: nil Context")
	}
	// 3、解析string的URL
	u, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}
	// 4、读取请求体
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	// 5、Url里的Host端口应该正常化
	u.Host = removeEmptyPort(u.Host)
	// 6、构建新请求
	req := &Request{
		ctx:        ctx,
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(Header),
		Body:       rc,
		Host:       u.Host,
	}
	// 7、body参数处理 (这里实际应用：jpush包内的 自己构建request)
	if body != nil {
		// 根据传入的不同类型的body,经过 读取body后，转成 io.ReadCloser 返回
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		default:
		}
		// 对于客户请求，请求.ContentLength0表示实际为0或未知。 要明确表示ContentLength为零的唯一方法是将Body设置为nil
		// 但是太多的代码依赖于NewRequest返回一个非空请求体，所以我们使用ReadCloser变量，并让http包也将sentinel变量显式地表示为零。
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = NoBody    // var NoBody = noBody{},空结构体
			req.GetBody = func() (io.ReadCloser, error) { return NoBody, nil }
		}
	}
	return req, nil
}

```

## 7. 根据 读取器(*bufio.Reader)，读取一个请求实例方法（readRequest内部方法实现），以及读取器的相关方法（新建一个最大字节的读取器，以及读取器结构体对象的读取关闭方法
```go
// ReadRequest 读取并解析来自Reader的传入请求
func ReadRequest(b *bufio.Reader) (*Request, error) {
	return readRequest(b, deleteHostHeader)
}

// readRequest的deleteHostHeader参数的常量
const (
	deleteHostHeader = true
	keepHostHeader   = false
)

// 从读取器中 解析出 请求*Request
func readRequest(b *bufio.Reader, deleteHostHeader bool) (req *Request, err error) {
	// 1、将读取器转成 *textproto.Reader
	tp := newTextprotoReader(b)
	req = new(Request)

	// 2、读取器第一行：获取索引.htmlHTTP/1.0协议
	var s string
	if s, err = tp.ReadLine(); err != nil {
		return nil, err
	}
	defer func() {
		putTextprotoReader(tp)
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	// 3、解析 第一行
	var ok bool
	req.Method, req.RequestURI, req.Proto, ok = parseRequestLine(s)
	if !ok {
		return nil, badStringError("malformed HTTP request", s)
	}
	// 4、检查解析结果
	if !validMethod(req.Method) {
		return nil, badStringError("invalid method", req.Method)
	}
	rawurl := req.RequestURI
	if req.ProtoMajor, req.ProtoMinor, ok = ParseHTTPVersion(req.Proto); !ok {
		return nil, badStringError("malformed HTTP version", req.Proto)
	}

	// 连接请求有两种不同的使用方式，都不使用完整的URL
	// 标准用法是通过HTTP代理来 贯穿HTTPS
	justAuthority := req.Method == "CONNECT" && !strings.HasPrefix(rawurl, "/")
	if justAuthority {
		rawurl = "http://" + rawurl
	}
	// 5、解析url
	if req.URL, err = url.ParseRequestURI(rawurl); err != nil {
		return nil, err
	}

	if justAuthority {
		// 把假“http://”去掉。
		req.URL.Scheme = ""
	}

	// 后续行：Key:value
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	req.Header = Header(mimeHeader)

	req.Host = req.URL.Host
	if req.Host == "" {
		req.Host = req.Header.get("Host")
	}
	if deleteHostHeader {
		delete(req.Header, "Host")
	}

	fixPragmaCacheControl(req.Header)

	req.Close = shouldClose(req.ProtoMajor, req.ProtoMinor, req.Header, false)

	err = readTransfer(req, b)
	if err != nil {
		return nil, err
	}

	if req.isH2Upgrade() {
		//因为它既不是块，也不是声明
		req.ContentLength = -1

		// 我们希望给处理程序一个劫持连接的机会，但是如果连接没有被劫持，我们需要阻止服务器进一步处理连接。设置为关闭以确保
		req.Close = true
	}
	return req, nil
}

```

最大字节读取器
```go
// 新建一个最大字节的读取器：类似于io.LimitReader但旨在限制传入请求主体的大小
func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser {
	return &maxBytesReader{w: w, r: r, n: n}
}

// 内部 最大字节 读取器结构体对象
type maxBytesReader struct {
	w   ResponseWriter
	r   io.ReadCloser // 底层读取器
	n   int64         // 剩余最大字节数
	err error         // 粘性错误
}

// 读取器 读取方法
func (l *maxBytesReader) Read(p []byte) (n int, err error) {
	if l.err != nil {
		return 0, l.err
	}
	if len(p) == 0 {
		return 0, nil
	}
	// 如果要求读取32KB字节，但只剩下5个字节，则无需读取32KB。6字节将回答我们是否达到极限或超过极限的问题
	if int64(len(p)) > l.n+1 {
		p = p[:l.n+1]
	}
	n, err = l.r.Read(p)

	if int64(n) <= l.n {
		l.n -= int64(n)
		l.err = err
		return n, err
	}

	n = int(l.n)
	l.n = 0

	// 服务器代码和客户端代码都使用maxBytesReader, 此“requestTooLarge”检查仅由服务器代码使用
	// 为了防止仅使用HTTP客户机代码（如cmd/go）的二进制文件也在HTTP服务器中链接，不要对服务器“*response”类型使用静态类型断言。改为检查此接口：
	type requestTooLarger interface {
		requestTooLarge()
	}
	if res, ok := l.w.(requestTooLarger); ok {
		res.requestTooLarge()
	}
	l.err = errors.New("http: request body too large")
	return n, l.err
}

// 读取器关闭方法
func (l *maxBytesReader) Close() error {
	return l.r.Close()
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

```