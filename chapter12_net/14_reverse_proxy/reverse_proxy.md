<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [reverse proxy](#reverse-proxy)
  - [反向代理使用场景](#%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第三方使用: k8s](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8-k8s)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# reverse proxy


Golang 中反向代理的实现主要使用了标准库的 net/http/httputil 包



## 反向代理使用场景
负载均衡（Load balancing）： 反向代理可以提供负载均衡解决方案，将传入的流量均匀地分布在不同的服务器之间，以防止单个服务器过载。

防止安全攻击： 由于真正的后端服务器永远不需要暴露公共 IP，所以 DDoS 等攻击只能针对反向代理进行， 这能确保在网络攻击中尽量多的保护你的资源，真正的后端服务器始终是安全的。

缓存： 假设你的实际服务器与用户所在的地区距离比较远，那么你可以在当地部署反向代理，它可以缓存网站内容并为当地用户提供服务。

SSL 加密： 由于与每个客户端的 SSL 通信会耗费大量的计算资源，因此可以使用反向代理处理所有与 SSL 相关的内容， 然后释放你真正服务器上的宝贵资源。


## 源码分析

```go
// go1.21.5/src/net/http/httputil/reverseproxy.go
type ReverseProxy struct {
	// Rewrite must be a function which modifies
	// the request into a new request to be sent
	// using Transport. Its response is then copied
	// back to the original client unmodified.
	// Rewrite must not access the provided ProxyRequest
	// or its contents after returning.
	//
	// The Forwarded, X-Forwarded, X-Forwarded-Host,
	// and X-Forwarded-Proto headers are removed from the
	// outbound request before Rewrite is called. See also
	// the ProxyRequest.SetXForwarded method.
	//
	// Unparsable query parameters are removed from the
	// outbound request before Rewrite is called.
	// The Rewrite function may copy the inbound URL's
	// RawQuery to the outbound URL to preserve the original
	// parameter string. Note that this can lead to security
	// issues if the proxy's interpretation of query parameters
	// does not match that of the downstream server.
	//
	// At most one of Rewrite or Director may be set.
	Rewrite func(*ProxyRequest)


	// 修改请求,默认设置 the X-Forwarded-For header
	Director func(*http.Request)

	// The transport used to perform proxy requests.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper

	// 刷新到客户端的刷新时间间隔
	// 流式请求下该参数会被忽略，所有反向代理请求将被立即刷新
	FlushInterval time.Duration


	ErrorLog *log.Logger


	BufferPool BufferPool

	// 用于修改响应结果及HTTP状态码，当返回结果error不为空时，会调用ErrorHandler
	ModifyResponse func(*http.Response) error

	// 用于处理后端和ModifyResponse返回的错误信息，默认将返回传递过来的错误信息，并返回HTTP 502se.
	ErrorHandler func(http.ResponseWriter, *http.Request, error)
}
```


初始化

```go
func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy {
	director := func(req *http.Request) {
		rewriteRequestURL(req, target)
	}
	return &ReverseProxy{Director: director}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	// 获取请求参数，例如请求的是/dir?id=123，那么rawQuery ：id=123
	targetQuery := target.RawQuery
	
	// 拼接 req 
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		// 使用"&"符号拼接请求参数
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}
```



处理请求

```go

func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transport := p.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	ctx := req.Context()
    // ..
	
	
    // 拷贝上游请求到下游请求
	outreq := req.Clone(ctx)
    // ..
    if p.Director != nil {
		// 修改请求（例如协议、参数、url等）
        p.Director(outreq)
        if outreq.Form != nil {
            outreq.URL.RawQuery = cleanQueryParams(outreq.URL.RawQuery)
        }
    }
	// 判断是否需要升级协议（Upgrade）
	reqUpType := upgradeType(outreq.Header)
	if !ascii.IsPrint(reqUpType) {
		p.getErrorHandler()(rw, req, fmt.Errorf("client tried to switch to invalid protocol %q", reqUpType))
		return
	}
	// 删除上游请求中的hop-by-hop Header，即不需要透传到下游的header
	removeHopByHopHeaders(outreq.Header)

	// Issue 21096: tell backend applications that care about trailer support
	// that we support trailers. (We do, but we don't go out of our way to
	// advertise that unless the incoming client request thought it was worth
	// mentioning.) Note that we look at req.Header, not outreq.Header, since
	// the latter has passed through removeHopByHopHeaders.
	if httpguts.HeaderValuesContainsToken(req.Header["Te"], "trailers") {
		outreq.Header.Set("Te", "trailers")
	}

	// After stripping all the hop-by-hop connection headers above, add back any
	// necessary for protocol upgrades, such as for websockets.
	if reqUpType != "" {
		outreq.Header.Set("Connection", "Upgrade")
		outreq.Header.Set("Upgrade", reqUpType)
	}

	if p.Rewrite != nil {
		// Strip client-provided forwarding headers.
		// The Rewrite func may use SetXForwarded to set new values
		// for these or copy the previous values from the inbound request.
		outreq.Header.Del("Forwarded")
		outreq.Header.Del("X-Forwarded-For")
		outreq.Header.Del("X-Forwarded-Host")
		outreq.Header.Del("X-Forwarded-Proto")

		// Remove unparsable query parameters from the outbound request.
		outreq.URL.RawQuery = cleanQueryParams(outreq.URL.RawQuery)

		pr := &ProxyRequest{
			In:  req,
			Out: outreq,
		}
		p.Rewrite(pr)
		outreq = pr.Out
	} else {
		if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			prior, ok := outreq.Header["X-Forwarded-For"]
			omit := ok && prior == nil // Issue 38079: nil now means don't populate the header
			if len(prior) > 0 {
				clientIP = strings.Join(prior, ", ") + ", " + clientIP
			}
			if !omit {
                // 设置X-Forward-For Header，追加当前节点IP
				outreq.Header.Set("X-Forwarded-For", clientIP)
			}
		}
	}

    // 。。

    
    // 向下游发起请求
	res, err := transport.RoundTrip(outreq)
	if err != nil {
		p.getErrorHandler()(rw, outreq, err)
		return
	}

	// 处理协议升级（httpcode 101）
	if res.StatusCode == http.StatusSwitchingProtocols {
		if !p.modifyResponse(rw, res, outreq) {
			return
		}
		p.handleUpgradeResponse(rw, outreq, res)
		return
	}
    // 删除不需要返回给上游的逐跳Header
	removeHopByHopHeaders(res.Header)

	// 修改响应体内容（如有需要）
	if !p.modifyResponse(rw, res, outreq) {
		return
	}
    // 拷贝下游响应头部到上游响应请求
	copyHeader(rw.Header(), res.Header)

    // ...
	
	// 返回HTTP状态码
	rw.WriteHeader(res.StatusCode)
    // 定时刷新内容到response
	err = p.copyResponse(rw, res.Body, p.flushInterval(res))
    // ..
	res.Body.Close() // close now, instead of defer, to populate res.Trailer

    // ..
}
```


## 第三方使用: k8s 

```shell
kubectl exec service1 -n ns 1 -- ls /root 
```

为进入目标pod的目标容器中执行命令（挂载标准输入和输出、标准错误的情景），kubectl exec访问kube-apiserver的connect接口（中间过程是通过http协议来握手，之后升级为spdy协议），
kube-apiserver把请求转发至对应节点的kubelet进程，而kubelet进程此时是一个反向代理，再把请求转发至cri shim程序. 



storage 注册
```go
func (c LegacyRESTStorageProvider) NewLegacyRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (LegacyRESTStorage, genericapiserver.APIGroupInfo, error) {
	// ..

	if resource := "pods"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = podStorage.Pod
		storage[resource+"/attach"] = podStorage.Attach
		storage[resource+"/status"] = podStorage.Status
		storage[resource+"/log"] = podStorage.Log
		storage[resource+"/exec"] = podStorage.Exec
		storage[resource+"/portforward"] = podStorage.PortForward
		storage[resource+"/proxy"] = podStorage.Proxy
		storage[resource+"/binding"] = podStorage.Binding
		if podStorage.Eviction != nil {
			storage[resource+"/eviction"] = podStorage.Eviction
		}
		storage[resource+"/ephemeralcontainers"] = podStorage.EphemeralContainers

	}
}
```

```go
// https://github.com/kubernetes/kubernetes/blob/4b8ec54d8e5ede787e2f94343d0723462f430127/staging/src/k8s.io/apiserver/pkg/endpoints/installer.go
func (a *APIInstaller) registerResourceHandlers(path string, storage rest.Storage, ws *restful.WebService) (*metav1.APIResource, *storageversion.ResourceInfo, error) {
    connecter, isConnecter := storage.(rest.Connecter)
    switch action.Verb {
    case "CONNECT":
	    for _, method := range connecter.ConnectMethods() {
			// 
			handler := metrics.InstrumentRouteFunc(action.Verb, group, version, resource, subresource, requestScope, metrics.APIServerComponent, deprecated, removedRelease, restfulConnectResource(connecter, reqScope, admit, path, isSubresource))
	    }
	}

}
```

```go
func restfulConnectResource(connecter rest.Connecter, scope handlers.RequestScope, admit admission.Interface, restPath string, isSubresource bool) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.ConnectResource(connecter, &scope, admit, restPath, isSubresource)(res.ResponseWriter, req.Request)
	}
}
```

```go

func ConnectResource(connecter rest.Connecter, scope *RequestScope, admit admission.Interface, restPath string, isSubresource bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		/*
			其他代码
		*/
		requestInfo, _ := request.RequestInfoFrom(ctx)
		metrics.RecordLongRunning(req, requestInfo, metrics.APIServerComponent, func() {
			// connecter对象的类型是ExecREST
			// 调用ExecREST结构体的Connect(...)方法来获得一个http handler
			handler, err := connecter.Connect(ctx, name, opts, &responder{scope: scope, req: req, w: w})
			if err != nil {
				scope.err(err, w, req)
				return
			}
			// 处理kubectl exec的请求
			handler.ServeHTTP(w, req)
		})
	}
}

```

```go
// https://github.com/kubernetes/kubernetes/blob/80060a502c3f86f00800fbeba7684a85f1ce5e17/pkg/registry/core/pod/rest/subresources.go

// Connect returns a handler for the pod exec proxy
func (r *ExecREST) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	execOpts, ok := opts.(*api.PodExecOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opts)
	}
	// 获取 node 地址 
	location, transport, err := pod.ExecLocation(ctx, r.Store, r.KubeletConn, name, execOpts)
	if err != nil {
		return nil, err
	}
	return newThrottledUpgradeAwareProxyHandler(location, transport, false, true, responder), nil
}


// kube-apiserver此时是一个反向代理，访问的是目标kubelet的接口，然后进行流拷贝
func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.MaxBytesPerSec = capabilities.Get().PerConnectionBandwidthLimitBytesPerSec
	return handler
}

```


```go
func (h *UpgradeAwareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // ...
	reverseProxyLocation := &url.URL{Scheme: h.Location.Scheme, Host: h.Location.Host}
	if h.AppendLocationPath {
		reverseProxyLocation.Path = h.Location.Path
	}

	proxy := httputil.NewSingleHostReverseProxy(reverseProxyLocation)
	proxy.Transport = h.Transport
	proxy.FlushInterval = h.FlushInterval
	proxy.ErrorLog = log.New(noSuppressPanicError{}, "", log.LstdFlags)
	if h.RejectForwardingRedirects {
		oldModifyResponse := proxy.ModifyResponse
		proxy.ModifyResponse = func(response *http.Response) error {
			code := response.StatusCode
			if code >= 300 && code <= 399 && len(response.Header.Get("Location")) > 0 {
				// close the original response
				response.Body.Close()
				msg := "the backend attempted to redirect this request, which is not permitted"
				// replace the response
				*response = http.Response{
					StatusCode:    http.StatusBadGateway,
					Status:        fmt.Sprintf("%d %s", response.StatusCode, http.StatusText(response.StatusCode)),
					Body:          io.NopCloser(strings.NewReader(msg)),
					ContentLength: int64(len(msg)),
				}
			} else {
				if oldModifyResponse != nil {
					if err := oldModifyResponse(response); err != nil {
						return err
					}
				}
			}
			return nil
		}
	}
	if h.Responder != nil {
		// if an optional error interceptor/responder was provided wire it
		// the custom responder might be used for providing a unified error reporting
		// or supporting retry mechanisms by not sending non-fatal errors to the clients
		proxy.ErrorHandler = h.Responder.Error
	}
	proxy.ServeHTTP(w, newReq)
}

```



## 参考

- [Golang 中的反向代理(ReverseProxy) 介绍与使用](https://www.cnblogs.com/FengZeng666/p/15643335.html)
- [kubernetes exec源码简析](https://blog.csdn.net/nangonghen/article/details/110411187)