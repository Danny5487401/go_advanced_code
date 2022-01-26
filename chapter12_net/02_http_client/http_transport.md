# http之transport源码详解
使用golang net/http库发送http请求，最后都是调用 transport的 RoundTrip方法中

go/go1.15.10/src/net/http/client.go
```go
//接口
type RoundTripper interface {
    RoundTrip(*Request) (*Response, error)
}
//实际调用
func send(ireq *Request, rt RoundTripper, deadline time.Time) (resp *Response, didTimeout func() bool, err error) {
	//--------
    resp, err = rt.RoundTrip(req)
    //--------
}
```
![](.http_images/roundTrip.png)
RoundTrip代表一个http事务，给一个请求返回一个响应。RoundTripper必须是并发安全的。RoundTripper接口的实现Transport结构体在源码包net/http/transport.go 中
```go
type Transport struct {
	idleMu     sync.Mutex
	wantIdle   bool                                // user has requested to close all idle conns  用户是否已关闭所有的空闲连接
	idleConn   map[connectMethodKey][]*persistConn // most recently used at end，保存从connectMethodKey(代表着不同的协议，不同的host，也就是不同的请求)到persistConn的映射
	/*
	idleConnCh 用来在并发http请求的时候在多个 goroutine 里面相互发送持久连接,也就是说，
	这些持久连接是可以重复利用的， 你的http请求用某个persistConn用完了，
	通过这个channel发送给其他http请求使用这个persistConn
	*/
	idleConnCh map[connectMethodKey]chan *persistConn
	idleLRU    connLRU
 
	reqMu       sync.Mutex
	reqCanceler map[*Request]func(error)   //请求取消器
 
	altMu    sync.Mutex   // guards changing altProto only
	altProto atomic.Value // of nil or map[string]RoundTripper, key is URI scheme  为空或者map[string]RoundTripper,key为URI  的scheme，用于自定义的协议及对应的处理请求的RoundTripper
 
	// Proxy specifies a function to return a proxy for a given
	// Request. If the function returns a non-nil error, the
	// request is aborted with the provided error.
	//
	// The proxy type is determined by the URL scheme. "http"
	// and "socks5" are supported. If the scheme is empty,
	// "http" is assumed.
	//
	// If Proxy is nil or returns a nil *URL, no proxy is used.
	Proxy func(*Request) (*url.URL, error)   //根据给定的Request返回一个代理，如果返回一个不为空的error，请求会终止
 
	// DialContext specifies the dial function for creating unencrypted TCP connections.
	// If DialContext is nil (and the deprecated Dial below is also nil),
	// then the transport dials using package net.
	/*
	DialContext用于指定创建未加密的TCP连接的dial功能，如果该函数为空，则使用net包下的dial函数
	*/
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
 
	// Dial specifies the dial function for creating unencrypted TCP connections.
	//
	// Deprecated: Use DialContext instead, which allows the transport
	// to cancel dials as soon as they are no longer needed.
	// If both are set, DialContext takes priority.
	/*
	Dial获取一个tcp连接，也就是net.Conn结构，然后就可以写入request，从而获取到response
	DialContext比Dial函数的优先级高
	*/
	Dial func(network, addr string) (net.Conn, error)
 
	// DialTLS specifies an optional dial function for creating
	// TLS connections for non-proxied HTTPS requests.
	//
	// If DialTLS is nil, Dial and TLSClientConfig are used.
	//
	// If DialTLS is set, the Dial hook is not used for HTTPS
	// requests and the TLSClientConfig and TLSHandshakeTimeout
	// are ignored. The returned net.Conn is assumed to already be
	// past the TLS handshake.
	/*
	DialTLS  为创建非代理的HTTPS请求的TLS连接提供一个可选的dial功能
	如果DialTLS为空，则使用Dial和TLSClientConfig
	如果设置了DialTLS，则HTTPS的请求不使用Dial的钩子，并且TLSClientConfig 和 TLSHandshakeTimeout会被忽略
	返回的net.Conn假设已经通过了TLS握手
	*/
	DialTLS func(network, addr string) (net.Conn, error)
 
	// TLSClientConfig specifies the TLS configuration to use with
	// tls.Client.
	// If nil, the default configuration is used.
	// If non-nil, HTTP/2 support may not be enabled by default.
	/*
      TLSClientConfig指定tls.Client使用的TLS配置信息
	如果为空，则使用默认配置
	如果不为空，默认情况下未启动HTTP/2支持
	*/
	TLSClientConfig *tls.Config
 
	// TLSHandshakeTimeout specifies the maximum amount of time waiting to
	// wait for a TLS handshake. Zero means no timeout.
	/*
	指定TLS握手的超时时间
	*/
	TLSHandshakeTimeout time.Duration
 
	// DisableKeepAlives, if true, prevents re-use of TCP connections
	// between different HTTP requests.
	DisableKeepAlives bool   //如果为true，则阻止在不同http请求之间重用TCP连接
 
	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression bool   //如果为true，则进制传输使用 Accept-Encoding: gzip
 
	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int   //指定最大的空闲连接数
 
	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int  //用于控制某一个主机的连接的最大空闲数
 
	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration   //指定空闲连接保持的最长时间，如果为0，则不受限制
 
	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	/*
	ResponseHeaderTimeout，如果非零，则指定在完全写入请求（包括其正文，如果有）之后等待服务器响应头的最长时间。
	此时间不包括读响应体的时间。
	*/
	ResponseHeaderTimeout time.Duration
 
	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	/*
   如果请求头是"Expect:100-continue",ExpectContinueTimeout  如果不为0，它表示等待服务器第一次响应头的最大时间
	零表示没有超时并导致正文立即发送，无需等待服务器批准。
	此时间不包括发送请求标头的时间。
	*/
	ExpectContinueTimeout time.Duration
 
	// TLSNextProto specifies how the Transport switches to an
	// alternate protocol (such as HTTP/2) after a TLS NPN/ALPN
	// protocol negotiation. If Transport dials an TLS connection
	// with a non-empty protocol name and TLSNextProto contains a
	// map entry for that key (such as "h2"), then the func is
	// called with the request's authority (such as "example.com"
	// or "example.com:1234") and the TLS connection. The function
	// must return a RoundTripper that then handles the request.
	// If TLSNextProto is not nil, HTTP/2 support is not enabled
	// automatically.
	/*
TLSNextProto指定在TLS NPN / ALPN协议协商之后传输如何切换到备用协议（例如HTTP / 2）。
	如果传输使用非空协议名称拨打TLS连接并且TLSNextProto包含该密钥的映射条目（例如“h2”），则使用请求的权限调用func（例如“example.com”或“example” .com：1234“）和TLS连接。
	该函数必须返回一个RoundTripper，然后处理该请求。 如果TLSNextProto不是nil，则不会自动启用HTTP / 2支持。
	*/
	TLSNextProto map[string]func(authority string, c *tls.Conn) RoundTripper
 
	// ProxyConnectHeader optionally specifies headers to send to
	// proxies during CONNECT requests.
	/*
	ProxyConnectHeader可选地指定在CONNECT请求期间发送给代理的标头。
	*/
	ProxyConnectHeader Header
 
	// MaxResponseHeaderBytes specifies a limit on how many
	// response bytes are allowed in the server's response
	// header.
	//
	// Zero means to use a default limit.
	/*
	指定服务器返回的响应头的最大字节数
	为0则使用默认的限制
	*/
	MaxResponseHeaderBytes int64
 
	// nextProtoOnce guards initialization of TLSNextProto and
	// h2transport (via onceSetNextProtoDefaults)
	//nextProtoOnce保护  TLSNextProto和 h2transport 的初始化
	nextProtoOnce sync.Once
	h2transport   *http2Transport // non-nil if http2 wired up，如果是http2连通，则不为nil
 
	// TODO: tunable on max per-host TCP dials in flight (Issue 13957)
}
```
RoundTrip方法会做两件事情：
- 调用 Transport 的 getConn 方法获取连接；
- 在获取到连接后，调用 persistConn 的 roundTrip 方法等待请求响应结果；
```go
//RoundTrip实现了RoundTripper接口
func (t *Transport) RoundTrip(req *Request) (*Response, error) {
	//初始化TLSNextProto  http2使用
	t.nextProtoOnce.Do(t.onceSetNextProtoDefaults)
	//获取请求的上下文
	ctx := req.Context()
	trace := httptrace.ContextClientTrace(ctx)
    // ...
 
	//如果该scheme有自定义的RoundTrip，则使用自定义的RoundTrip处理request，并返回response
	altProto, _ := t.altProto.Load().(map[string]RoundTripper)
	if altRT := altProto[scheme]; altRT != nil {
		if resp, err := altRT.RoundTrip(req); err != ErrSkipAltProtocol {
			return resp, err
		}
	}
    // ...
 
	for {
        select {
        case <-ctx.Done():
            req.closeBody()
            return nil, ctx.Err()
        default:
        }
		// treq gets modified by roundTrip, so we need to recreate for each retry.
		//  封装请求:初始化transportRequest,transportRequest是request的包装器
		treq := &transportRequest{Request: req, trace: trace}
		//根据用户的请求信息获取connectMethod  cm
		cm, err := t.connectMethodForRequest(treq)
		if err != nil {
			req.closeBody()
			return nil, err
		}
 

		//从缓存中获取一个连接，或者新建一个连接
		pconn, err := t.getConn(treq, cm)
		if err != nil {
			t.setReqCanceler(req, nil)
			req.closeBody()
			return nil, err
		}
 
		var resp *Response
		if pconn.alt != nil {
			// HTTP/2 path.
			t.setReqCanceler(req, nil) // not cancelable with CancelRequest
			resp, err = pconn.alt.RoundTrip(req)
		} else {
			resp, err = pconn.roundTrip(treq)
		}
		if err == nil {
			return resp, nil
		}
		// ...
	}
}
```


## 获取或则新建连接:
- 1. 调用 queueForIdleConn 获取空闲 connection；
- 2. 调用 queueForDial 等待创建新的 connection；
![](.http_transport_images/get_conn.png)
```go
func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (pc *persistConn, err error) {
    req := treq.Request
    trace := treq.trace
    ctx := req.Context()
    if trace != nil && trace.GetConn != nil {
        trace.GetConn(cm.addr())
    }   
    // 将请求封装成 wantConn 结构体
    w := &wantConn{
        cm:         cm,
        key:        cm.key(),
        ctx:        ctx,
        ready:      make(chan struct{}, 1),
        beforeDial: testHookPrePendingDial,
        afterDial:  testHookPostPendingDial,
    }
    defer func() {
        if err != nil {
            w.cancel(t, err)
        }
    }()

    // 获取空闲连接
    if delivered := t.queueForIdleConn(w); delivered {
        pc := w.pc
        ...
        t.setReqCanceler(treq.cancelKey, func(error) {})
        return pc, nil
    }

    // 创建连接
    t.queueForDial(w)

    select {
    // 获取到连接后进入该分支
    case <-w.ready:
        ...
        return w.pc, w.err
    ...
}
```

### 获取空闲连接 queueForIdleConn
![](.http_transport_images/get_queue_for_idle_conn.png)
成功获取 connection 分为如下几步：

1. 根据当前的请求的地址去空闲 connection 字典中查看存不存在空闲的 connection 列表；
2. 如果能获取到空闲的 connection 列表，那么获取到列表的最后一个 connection；
3. 返回

![](.http_transport_images/fail_to_get_queue_for_indle_conn.png)
当获取不到空闲 connection 时：

1. 根据当前的请求的地址去空闲 connection 字典中查看存不存在空闲的 connection 列表；
2. 不存在该请求的 connection 列表，那么将该 wantConn 加入到 等待获取空闲 connection 字典中；
```go
func (t *Transport) queueForIdleConn(w *wantConn) (delivered bool) {
    if t.DisableKeepAlives {
        return false
    }

    t.idleMu.Lock()
    defer t.idleMu.Unlock() 
    t.closeIdle = false

    if w == nil { 
        return false
    }

    // 计算空闲连接超时时间
    var oldTime time.Time
    if t.IdleConnTimeout > 0 {
        oldTime = time.Now().Add(-t.IdleConnTimeout)
    }
    // Look for most recently-used idle connection.
    // 找到key相同的 connection 列表
    if list, ok := t.idleConn[w.key]; ok {
        stop := false
        delivered := false
        for len(list) > 0 && !stop {
            // 找到connection列表最后一个
            pconn := list[len(list)-1] 
            // 检查这个 connection 是不是等待太久了
            tooOld := !oldTime.IsZero() && pconn.idleAt.Round(0).Before(oldTime)
            if tooOld { 
                go pconn.closeConnIfStillIdle()
            }
            // 该 connection 被标记为 broken 或 闲置太久 continue
            if pconn.isBroken() || tooOld { 
                list = list[:len(list)-1]
                continue
            }
            // 尝试将该 connection 写入到 w 中
            delivered = w.tryDeliver(pconn, nil)
            if delivered {
                // 操作成功，需要将 connection 从空闲列表中移除
                if pconn.alt != nil { 
                } else { 
                    t.idleLRU.remove(pconn)
                    list = list[:len(list)-1]
                }
            }
            stop = true
        }
        if len(list) > 0 {
            t.idleConn[w.key] = list
        } else {
            // 如果该 key 对应的空闲列表不存在，那么将该key从字典中移除
            delete(t.idleConn, w.key)
        }
        if stop {
            return delivered
        }
    } 
    // 如果找不到空闲的 connection
    if t.idleConnWait == nil {
        t.idleConnWait = make(map[connectMethodKey]wantConnQueue)
    }
  // 将该 wantConn 加入到 等待获取空闲 connection 字典中
    q := t.idleConnWait[w.key] 
    q.cleanFront()
    q.pushBack(w)
    t.idleConnWait[w.key] = q
    return false
}
```

### 建立连接 queueForDial
![](.http_transport_images/queue_for_dial.png)
尝试去建立连接，总共分为以下几个步骤
1. 在调用 queueForDial 方法的时候会校验 MaxConnsPerHost 是否未设置或已达上限；
   * 检验不通过则将当前的请求放入到 connsPerHostWait 等待字典中；
2. 如果校验通过那么会异步的调用 dialConnFor 方法创建连接；
3. dialConnFor 方法首先会调用 dialConn 方法创建 TCP 连接，然后启动两个异步线程来处理读写数据，然后调用 tryDeliver 将连接绑定到 wantConn 上面。

```go
func (t *Transport) queueForDial(w *wantConn) {
    w.beforeDial()
    // 小于零说明无限制，异步建立连接
    if t.MaxConnsPerHost <= 0 {
        go t.dialConnFor(w)
        return
    }

    t.connsPerHostMu.Lock()
    defer t.connsPerHostMu.Unlock()
    // 每个 host 建立的连接数没达到上限，异步建立连接
    if n := t.connsPerHost[w.key]; n < t.MaxConnsPerHost {
        if t.connsPerHost == nil {
            t.connsPerHost = make(map[connectMethodKey]int)
        }
        t.connsPerHost[w.key] = n + 1
        go t.dialConnFor(w)
        return
    }
    //每个 host 建立的连接数已达到上限，需要进入等待队列
    if t.connsPerHostWait == nil {
        t.connsPerHostWait = make(map[connectMethodKey]wantConnQueue)
    }
    q := t.connsPerHostWait[w.key]
    q.cleanFront()
    q.pushBack(w)
    t.connsPerHostWait[w.key] = q
}
```
```go
func (t *Transport) dialConnFor(w *wantConn) {
    defer w.afterDial()
    // 建立连接
    pc, err := t.dialConn(w.ctx, w.cm)
    // 连接绑定 wantConn
    delivered := w.tryDeliver(pc, err)
    // 建立连接成功，但是绑定 wantConn 失败
    // 那么将该连接放置到空闲连接字典或调用 等待获取空闲 connection 字典 中的元素执行
    if err == nil && (!delivered || pc.alt != nil) { 
        t.putOrCloseIdleConn(pc)
    }
    if err != nil {
        t.decConnsPerHost(w.key)
    }
}
```
```go
func (t *Transport) dialConn(ctx context.Context, cm connectMethod) (pconn *persistConn, err error) {
    // 创建连接结构体
    pconn = &persistConn{
        t:             t,
        cacheKey:      cm.key(),
        reqch:         make(chan requestAndChan, 1),
        writech:       make(chan writeRequest, 1),
        closech:       make(chan struct{}),
        writeErrCh:    make(chan error, 1),
        writeLoopDone: make(chan struct{}),
    }
    ...
    if cm.scheme() == "https" && t.hasCustomTLSDialer() {
        ...
    } else {
        // 建立 tcp 连接
        conn, err := t.dial(ctx, "tcp", cm.addr())
        if err != nil {
            return nil, wrapErr(err)
        }
        pconn.conn = conn 
    } 
    ...

    if s := pconn.tlsState; s != nil && s.NegotiatedProtocolIsMutual && s.NegotiatedProtocol != "" {
        if next, ok := t.TLSNextProto[s.NegotiatedProtocol]; ok {
            alt := next(cm.targetAddr, pconn.conn.(*tls.Conn))
            if e, ok := alt.(http2erringRoundTripper); ok {
                // pconn.conn was closed by next (http2configureTransport.upgradeFn).
                return nil, e.err
            }
            return &persistConn{t: t, cacheKey: pconn.cacheKey, alt: alt}, nil
        }
    }

    pconn.br = bufio.NewReaderSize(pconn, t.readBufferSize())
    pconn.bw = bufio.NewWriterSize(persistConnWriter{pconn}, t.writeBufferSize())
    //为每个连接异步处理读写数据
    go pconn.readLoop()
    go pconn.writeLoop()
    return pconn, nil
}
```