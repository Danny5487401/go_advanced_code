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
RoundTrip方法
```go
//RoundTrip实现了RoundTripper接口
func (t *Transport) RoundTrip(req *Request) (*Response, error) {
	//初始化TLSNextProto  http2使用
	t.nextProtoOnce.Do(t.onceSetNextProtoDefaults)
	//获取请求的上下文
	ctx := req.Context()
	trace := httptrace.ContextClientTrace(ctx)
 
	//错误处理
	if req.URL == nil {
		req.closeBody()
		return nil, errors.New("http: nil Request.URL")
	}
	if req.Header == nil {
		req.closeBody()
		return nil, errors.New("http: nil Request.Header")
	}
	scheme := req.URL.Scheme
	isHTTP := scheme == "http" || scheme == "https"
	//如果是http或https请求，对Header中的数据进行校验
	if isHTTP {
		for k, vv := range req.Header {
			if !httplex.ValidHeaderFieldName(k) {
				return nil, fmt.Errorf("net/http: invalid header field name %q", k)
			}
			for _, v := range vv {
				if !httplex.ValidHeaderFieldValue(v) {
					return nil, fmt.Errorf("net/http: invalid header field value %q for key %v", v, k)
				}
			}
		}
	}
 
	//如果该scheme有自定义的RoundTrip，则使用自定义的RoundTrip处理request，并返回response
	altProto, _ := t.altProto.Load().(map[string]RoundTripper)
	if altRT := altProto[scheme]; altRT != nil {
		if resp, err := altRT.RoundTrip(req); err != ErrSkipAltProtocol {
			return resp, err
		}
	}
 
	//如果不是http请求，则关闭并退出
	if !isHTTP {
		req.closeBody()
		return nil, &badStringError{"unsupported protocol scheme", scheme}
	}
 
	//对请求的Method进行校验
	if req.Method != "" && !validMethod(req.Method) {
		return nil, fmt.Errorf("net/http: invalid method %q", req.Method)
	}
 
	//请求的host为空，则返回
	if req.URL.Host == "" {
		req.closeBody()
		return nil, errors.New("http: no Host in request URL")
	}
 
	for {
		// treq gets modified by roundTrip, so we need to recreate for each retry.
		//初始化transportRequest,transportRequest是request的包装器
		treq := &transportRequest{Request: req, trace: trace}
		//根据用户的请求信息获取connectMethod  cm
		cm, err := t.connectMethodForRequest(treq)
		if err != nil {
			req.closeBody()
			return nil, err
		}
 
		// Get the cached or newly-created connection to either the
		// host (for http or https), the http proxy, or the http proxy
		// pre-CONNECTed to https server. In any case, we'll be ready
		// to send it requests.
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
		if !pconn.shouldRetryRequest(req, err) {
			// Issue 16465: return underlying net.Conn.Read error from peek,
			// as we've historically done.
			if e, ok := err.(transportReadFromServerError); ok {
				err = e.err
			}
			return nil, err
		}
		testHookRoundTripRetried()
 
		// Rewind the body if we're able to.  (HTTP/2 does this itself so we only
		// need to do it for HTTP/1.1 connections.)
		if req.GetBody != nil && pconn.alt == nil {
			newReq := *req
			var err error
			newReq.Body, err = req.GetBody()
			if err != nil {
				return nil, err
			}
			req = &newReq
		}
	}
}
```

获取或则新建连接:
- 1. 调用 queueForIdleConn 获取空闲 connection；
- 2. 调用 queueForDial 等待创建新的 connection；
```go
 func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (*persistConn, error) {
	req := treq.Request
	trace := treq.trace
	ctx := req.Context()
 
	//GetConn是钩子函数在获取连接前调用
	if trace != nil && trace.GetConn != nil {
		trace.GetConn(cm.addr())
	}
	//如果可以获取到空闲的连接
	if pc, idleSince := t.getIdleConn(cm); pc != nil {
		if trace != nil && trace.GotConn != nil {  //GotConn是钩子函数，成功获取连接后调用
			trace.GotConn(pc.gotIdleConnTrace(idleSince))
		}
		// set request canceler to some non-nil function so we
		// can detect whether it was cleared between now and when
		// we enter roundTrip
		/*
		将请求的canceler设置为某些非零函数，以便我们可以检测它是否在现在和我们进入roundTrip之间被清除
		*/
		t.setReqCanceler(req, func(error) {})
		return pc, nil
	}
 
	type dialRes struct {
		pc  *persistConn
		err error
	}
	dialc := make(chan dialRes)
 
	// Copy these hooks so we don't race on the postPendingDial in
	// the goroutine we launch. Issue 11136.
	testHookPrePendingDial := testHookPrePendingDial
	testHookPostPendingDial := testHookPostPendingDial
 
	//该内部函数handlePendingDial的主要作用是，新开启一个协程，当新建连接完成后但没有被使用，将其放到连接池（缓存）中或将其关闭
	handlePendingDial := func() {
		testHookPrePendingDial()
		go func() {
			if v := <-dialc; v.err == nil {
				t.putOrCloseIdleConn(v.pc)
			}
			testHookPostPendingDial()
		}()
	}
 
	cancelc := make(chan error, 1)
	t.setReqCanceler(req, func(err error) { cancelc <- err })
 
	go func() {//开启一个协程新建一个连接
		pc, err := t.dialConn(ctx, cm)
		dialc <- dialRes{pc, err}
	}()
 
	idleConnCh := t.getIdleConnCh(cm)
	select {
	case v := <-dialc: //获取新建的连接
		// Our dial finished.
		if v.pc != nil { //如果新建的连接不为nil，则返回新建的连接
			if trace != nil && trace.GotConn != nil && v.pc.alt == nil {
				trace.GotConn(httptrace.GotConnInfo{Conn: v.pc.conn})
			}
			return v.pc, nil
		}
		// Our dial failed. See why to return a nicer error
		// value.
		select {
		case <-req.Cancel:
			// It was an error due to cancelation, so prioritize that
			// error value. (Issue 16049)
			return nil, errRequestCanceledConn
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case err := <-cancelc:
			if err == errRequestCanceled {
				err = errRequestCanceledConn
			}
			return nil, err
		default:
			// It wasn't an error due to cancelation, so
			// return the original error message:
			return nil, v.err
		}
	case pc := <-idleConnCh:  //如果在新建连接的过程中，有空闲的连接，则返回该空闲的连接
		// Another request finished first and its net.Conn
		// became available before our dial. Or somebody
		// else's dial that they didn't use.
		// But our dial is still going, so give it away
		// when it finishes:
		//如果在dial连接的时候，有空闲的连接，但是这个时候我们仍然正在新建连接，所以当它新建完成后将其放到连接池或丢弃
		handlePendingDial()
		if trace != nil && trace.GotConn != nil {
			trace.GotConn(httptrace.GotConnInfo{Conn: pc.conn, Reused: pc.isReused()})
		}
		return pc, nil
	case <-req.Cancel:
		handlePendingDial()
		return nil, errRequestCanceledConn
	case <-req.Context().Done():
		handlePendingDial()
		return nil, req.Context().Err()
	case err := <-cancelc:
		handlePendingDial()
		if err == errRequestCanceled {
			err = errRequestCanceledConn
		}
		return nil, err
	}
}
```