<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [DNS（Domain Name System)](#dnsdomain-name-system)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
  - [DNS的记录类型](#dns%E7%9A%84%E8%AE%B0%E5%BD%95%E7%B1%BB%E5%9E%8B)
  - [工作原理](#%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86)
  - [net 源码分析](#net-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [基本方法](#%E5%9F%BA%E6%9C%AC%E6%96%B9%E6%B3%95)
    - [具体实现](#%E5%85%B7%E4%BD%93%E5%AE%9E%E7%8E%B0)
  - [应用 --> VictoriaMetrics 使用](#%E5%BA%94%E7%94%A8----victoriametrics-%E4%BD%BF%E7%94%A8)
  - [github.com/miekg/dns-->github.com/miekg/dns](#githubcommiekgdns--githubcommiekgdns)
    - [主要功能](#%E4%B8%BB%E8%A6%81%E5%8A%9F%E8%83%BD)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# DNS（Domain Name System)


主要作用是将主机名解析为IP地址的过程，完成了从域名到主机识别ip地址之间的转换；



## 基本概念


FQDN(Full Qualified Domain Name)：完整主机名，是由主机名和域名构成。例如http://www.baidu.com。当中，www就是web网站服务器的主机名，http://baidu.com 就是域名，主机名和域名之间用实心点号来表示


## DNS的记录类型


（1）A：地址记录（Address），返回域名指向的IP地址。

（2）NS：域名服务器记录（Name Server），返回保存下一级域名信息的服务器地址。该记录只能设置为域名，不能设置为IP地址。

（3）MX：邮件记录（Mail eXchange），返回接收电子邮件的服务器地址。

（4）CNAME：规范名称记录（Canonical Name），返回另一个域名，即当前查询的域名是另一个域名的跳转。

（5）PTR：逆向查询记录（Pointer Record），只用于从IP地址查询域名，详见下文。

 (6) SRV记录: SRV记录是一种用于服务发现和负载均衡的DNS记录类型。它包含了服务器的位置信息，如服务器的IP地址、端口号以及提供服务的服务器的优先级和权重等。通过查询SRV记录，客户端可以找到提供所需服务的服务器，并建立连接

```shell
# SRV记录
# 英文
_Service._Proto.Name TTL Class SRV Priority Weight Port Target

# 中文
_服务._协议.名称. TTL 类别 SRV 优先级 权重 端口 主机.

# 在Kubernetes里面，CoreDNS会为有名称的端口创建SRV记录，这些端口可以是svc或headless.svc的一部分。对每个命名端口，SRV记录了一个类似下列格式的记录：
# _port-name._port-protocol.my-svc.my-namespace.svc.cluster.local
```


 (7) TXT记录，一般用于某个主机名的标识和说明，通过设置TXT记录可以使别人更方便地联系到你。

 (8) AAAA记录: 是用于将域名解析到IPv6地址的一种DNS记录类型。

## 工作原理

DNS解析域名到IP要经过三个阶段：
1. 本地DNS缓存解析；
2. 本地DNS服务器解析，递归查询；
3. 根域及各级域名服务器解析，迭代查询
```shell
✗ dig www.baidu.com         

# 第一段是查询参数和统计。
; <<>> DiG 9.10.6 <<>> www.baidu.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 26841
;; flags: qr rd ra; QUERY: 1, ANSWER: 3, AUTHORITY: 5, ADDITIONAL: 9

# 第二段是查询内容。
;; QUESTION SECTION:
;www.baidu.com.                 IN      A

# 第三段是DNS服务器的答复
;; ANSWER SECTION:
www.baidu.com.          600     IN      CNAME   www.a.shifen.com.
www.a.shifen.com.       600     IN      A       183.2.172.177
www.a.shifen.com.       600     IN      A       183.2.172.17

# 第四段显示NS记录（Name Server的缩写），即哪些服务器负责管理a.shifen.com.的DNS记录。
;; AUTHORITY SECTION:
a.shifen.com.           569     IN      NS      ns5.a.shifen.com.
a.shifen.com.           569     IN      NS      ns4.a.shifen.com.
a.shifen.com.           569     IN      NS      ns1.a.shifen.com.
a.shifen.com.           569     IN      NS      ns2.a.shifen.com.
a.shifen.com.           569     IN      NS      ns3.a.shifen.com.

# 第五段是上面域名服务器的IP地址，这是随着前一段一起返回的。
;; ADDITIONAL SECTION:
ns4.a.shifen.com.       94      IN      A       14.215.177.229
ns4.a.shifen.com.       94      IN      A       111.20.4.28
ns5.a.shifen.com.       535     IN      A       180.76.76.95
ns1.a.shifen.com.       123     IN      A       110.242.68.42
ns2.a.shifen.com.       475     IN      A       220.181.33.32
ns3.a.shifen.com.       140     IN      A       36.155.132.12
ns3.a.shifen.com.       140     IN      A       153.3.238.162
ns5.a.shifen.com.       535     IN      AAAA    240e:bf:b801:1006:0:ff:b04f:346b
ns5.a.shifen.com.       535     IN      AAAA    240e:940:603:a:0:ff:b08d:239d

# 第六段是DNS服务器的一些传输信息
;; Query time: 20 msec
;; SERVER: 192.168.2.1#53(192.168.2.1)
;; WHEN: Wed Sep 24 19:56:39 CST 2025
;; MSG SIZE  rcvd: 351

```


## net 源码分析

解析结构体
```go
// go1.24.3/src/net/lookup.go
type Resolver struct {
	PreferGo bool // 是否使用Go内置解析器（绕过cgo）,相当于 GODEBUG=netdns=go,  

	// StrictErrors controls the behavior of temporary errors
	// (including timeout, socket errors, and SERVFAIL) when using
	// Go's built-in resolver. For a query composed of multiple
	// sub-queries (such as an A+AAAA address lookup, or walking the
	// DNS search list), this option causes such errors to abort the
	// whole query instead of returning a partial result. This is
	// not enabled by default because it may affect compatibility
	// with resolvers that process AAAA queries incorrectly.
	StrictErrors bool

	// Dial optionally specifies an alternate dialer for use by
	// Go's built-in DNS resolver to make TCP and UDP connections
	// to DNS services. The host in the address parameter will
	// always be a literal IP address and not a host name, and the
	// port in the address parameter will be a literal port number
	// and not a service name.
	// If the Conn returned is also a PacketConn, sent and received DNS
	// messages must adhere to RFC 1035 section 4.2.1, "UDP usage".
	// Otherwise, DNS messages transmitted over Conn must adhere
	// to RFC 7766 section 5, "Transport Protocol Selection".
	// If nil, the default dialer is used.
	Dial func(ctx context.Context, network, address string) (Conn, error)

	// lookupGroup merges LookupIPAddr calls together for lookups for the same
	// host. The lookupGroup key is the LookupIPAddr.host argument.
	// The return values are ([]IPAddr, error).
	lookupGroup singleflight.Group

}
```
PreferGo参数的适用场景是什么？
- 纯Go解析：当系统DNS配置不可靠（如容器环境缺少resolv.conf），或需要避免cgo依赖时，设置PreferGo: true使用Go内置解析器。
- 性能权衡：内置解析器可能比系统解析器稍慢，但兼容性更强。

```shell
export GODEBUG=netdns=go    # force pure Go resolver 纯go 方式
export GODEBUG=netdns=cgo   # force cgo resolver   cgo 方式
```


### 基本方法
基础地址解析

* LookupHost(host string) (addrs []string, err error)：解析域名对应的IP地址（A/AAAA记录）。
* func LookupIP(host string) ([]IP, error) ：按协议类型（“ip4”/“ip6”）解析IP地址。
* LookupAddr(ctx, addr)：反向解析IP地址到域名（PTR记录）。



资源记录查询

* LookupCNAME(ctx, host)：查询域名的CNAME别名记录。
* LookupMX(ctx, name)：查询邮件交换记录（同net.LookupMX，但支持上下文）。
* LookupNS(ctx, name)：查询名称服务器记录（同net.LookupNS，支持上下文）。
* LookupTXT(ctx, name)：查询TXT记录（如SPF/DKIM验证）


服务发现与高级解析

* LookupPort(ctx, network, service)：解析服务名称到端口（如"http"→80）。
* LookupSRV(ctx, service, proto, name)：查询SRV记录（用于微服务发现，如_http._tcp.example.com）。
* LookupNetIP(ctx, network, host)：返回netip.Addr（Go 1.20+新增，支持更灵活的IP处理）


### 具体实现
```go
func LookupHost(host string) (addrs []string, err error) {
	return DefaultResolver.LookupHost(context.Background(), host)
}

```

```go
func (r *Resolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
    // 校验
	
	// 解析
	return r.lookupHost(ctx, host)
}


func (r *Resolver) lookupHost(ctx context.Context, host string) (addrs []string, err error) {
	order, conf := systemConf().hostLookupOrder(r, host)
	if order == hostLookupCgo { // 查找顺序是 cgo
		return cgoLookupHost(ctx, host)
	}
	// 本地处理
	return r.goLookupHostOrder(ctx, host, order, conf)
}
```

cgo处理
```go
func cgoLookupHost(ctx context.Context, name string) (hosts []string, err error) {
	addrs, err := cgoLookupIP(ctx, "ip", name)
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		hosts = append(hosts, addr.String())
	}
	return hosts, nil
}


// 最终调用
func cgoLookupHostIP(network, name string) (addrs []IPAddr, err error) {
	var hints _C_struct_addrinfo
	*_C_ai_flags(&hints) = cgoAddrInfoFlags
	*_C_ai_socktype(&hints) = _C_SOCK_STREAM
	*_C_ai_family(&hints) = _C_AF_UNSPEC
	switch ipVersion(network) {
	case '4':
		*_C_ai_family(&hints) = _C_AF_INET
	case '6':
		*_C_ai_family(&hints) = _C_AF_INET6
	}

	h, err := syscall.BytePtrFromString(name)
	if err != nil {
		return nil, &DNSError{Err: err.Error(), Name: name}
	}
	var res *_C_struct_addrinfo
	// 调用到c标准库的 getaddrinfo
	gerrno, err := _C_getaddrinfo((*_C_char)(unsafe.Pointer(h)), nil, &hints, &res)
	if gerrno != 0 {
        // 戳五处理

	}
	defer _C_freeaddrinfo(res)

	for r := res; r != nil; r = *_C_ai_next(r) {
		// We only asked for SOCK_STREAM, but check anyhow.
		if *_C_ai_socktype(r) != _C_SOCK_STREAM {
			continue
		}
		switch *_C_ai_family(r) {
		case _C_AF_INET:
			sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(*_C_ai_addr(r)))
			addr := IPAddr{IP: copyIP(sa.Addr[:])}
			addrs = append(addrs, addr)
		case _C_AF_INET6:
			sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(*_C_ai_addr(r)))
			addr := IPAddr{IP: copyIP(sa.Addr[:]), Zone: zoneCache.name(int(sa.Scope_id))}
			addrs = append(addrs, addr)
		}
	}
	return addrs, nil
}

```


```go
func (r *Resolver) goLookupHostOrder(ctx context.Context, name string, order hostLookupOrder, conf *dnsConfig) (addrs []string, err error) {
	if order == hostLookupFilesDNS || order == hostLookupFiles {

		// 在本地/etc/hosts中查找
		addrs, _ = lookupStaticHost(name)
		if len(addrs) > 0 {
			return
		}

		if order == hostLookupFiles {
			return nil, newDNSError(errNoSuchHost, name, "")
		}
	}
	ips, _, err := r.goLookupIPCNAMEOrder(ctx, "ip", name, order, conf)
	if err != nil {
		return
	}
	addrs = make([]string, 0, len(ips))
	for _, ip := range ips {
		addrs = append(addrs, ip.String())
	}
	return
}
```


```go
func (r *Resolver) goLookupIPCNAMEOrder(ctx context.Context, network, name string, order hostLookupOrder, conf *dnsConfig) (addrs []IPAddr, cname dnsmessage.Name, err error) {
	if order == hostLookupFilesDNS || order == hostLookupFiles {
		var canonical string
		addrs, canonical = goLookupIPFiles(name)

		if len(addrs) > 0 {
			var err error
			cname, err = dnsmessage.NewName(canonical)
			if err != nil {
				return nil, dnsmessage.Name{}, err
			}
			return addrs, cname, nil
		}

		if order == hostLookupFiles {
			return nil, dnsmessage.Name{}, newDNSError(errNoSuchHost, name, "")
		}
	}

	if !isDomainName(name) {
		// See comment in func lookup above about use of errNoSuchHost.
		return nil, dnsmessage.Name{}, newDNSError(errNoSuchHost, name, "")
	}
	type result struct {
		p      dnsmessage.Parser
		server string
		error
	}
	
	if conf == nil {
		// 从 /etc/resolv.conf 获取配置
		conf = getSystemDNSConfig()
	}

	lane := make(chan result, 1)
	// 查找 A 和 AAAA 记录
	qtypes := []dnsmessage.Type{dnsmessage.TypeA, dnsmessage.TypeAAAA}
	if network == "CNAME" {
		qtypes = append(qtypes, dnsmessage.TypeCNAME)
	}
	switch ipVersion(network) {
	case '4':
		qtypes = []dnsmessage.Type{dnsmessage.TypeA}
	case '6':
		qtypes = []dnsmessage.Type{dnsmessage.TypeAAAA}
	}
	var queryFn func(fqdn string, qtype dnsmessage.Type)
	var responseFn func(fqdn string, qtype dnsmessage.Type) result
	if conf.singleRequest { // 单个请求
		queryFn = func(fqdn string, qtype dnsmessage.Type) {}
		responseFn = func(fqdn string, qtype dnsmessage.Type) result {
			dnsWaitGroup.Add(1)
			defer dnsWaitGroup.Done()
			p, server, err := r.tryOneName(ctx, conf, fqdn, qtype)
			return result{p, server, err}
		}
	} else {
		queryFn = func(fqdn string, qtype dnsmessage.Type) {
			dnsWaitGroup.Add(1)
			go func(qtype dnsmessage.Type) {
				p, server, err := r.tryOneName(ctx, conf, fqdn, qtype)
				lane <- result{p, server, err}
				dnsWaitGroup.Done()
			}(qtype)
		}
		responseFn = func(fqdn string, qtype dnsmessage.Type) result {
			return <-lane
		}
	}
	var lastErr error
	for _, fqdn := range conf.nameList(name) {
		for _, qtype := range qtypes {
			queryFn(fqdn, qtype)
		}
		hitStrictError := false
		for _, qtype := range qtypes {
			result := responseFn(fqdn, qtype)
			if result.error != nil {
				if nerr, ok := result.error.(Error); ok && nerr.Temporary() && r.strictErrors() {
					// This error will abort the nameList loop.
					hitStrictError = true
					lastErr = result.error
				} else if lastErr == nil || fqdn == name+"." {
					// Prefer error for original name.
					lastErr = result.error
				}
				continue
			}

			// Presotto says it's okay to assume that servers listed in
			// /etc/resolv.conf are recursive resolvers.
			//
			// We asked for recursion, so it should have included all the
			// answers we need in this one packet.
			//
			// Further, RFC 1034 section 4.3.1 says that "the recursive
			// response to a query will be... The answer to the query,
			// possibly preface by one or more CNAME RRs that specify
			// aliases encountered on the way to an answer."
			//
			// Therefore, we should be able to assume that we can ignore
			// CNAMEs and that the A and AAAA records we requested are
			// for the canonical name.

		loop:
			for {
				h, err := result.p.AnswerHeader()
				if err != nil && err != dnsmessage.ErrSectionDone {
					lastErr = &DNSError{
						Err:    errCannotUnmarshalDNSMessage.Error(),
						Name:   name,
						Server: result.server,
					}
				}
				if err != nil {
					break
				}
				switch h.Type {
				case dnsmessage.TypeA:
					a, err := result.p.AResource()
					if err != nil {
						lastErr = &DNSError{
							Err:    errCannotUnmarshalDNSMessage.Error(),
							Name:   name,
							Server: result.server,
						}
						break loop
					}
					addrs = append(addrs, IPAddr{IP: IP(a.A[:])})
					if cname.Length == 0 && h.Name.Length != 0 {
						cname = h.Name
					}

				case dnsmessage.TypeAAAA:
					aaaa, err := result.p.AAAAResource()
					if err != nil {
						lastErr = &DNSError{
							Err:    errCannotUnmarshalDNSMessage.Error(),
							Name:   name,
							Server: result.server,
						}
						break loop
					}
					addrs = append(addrs, IPAddr{IP: IP(aaaa.AAAA[:])})
					if cname.Length == 0 && h.Name.Length != 0 {
						cname = h.Name
					}

				case dnsmessage.TypeCNAME:
					c, err := result.p.CNAMEResource()
					if err != nil {
						lastErr = &DNSError{
							Err:    errCannotUnmarshalDNSMessage.Error(),
							Name:   name,
							Server: result.server,
						}
						break loop
					}
					if cname.Length == 0 && c.CNAME.Length > 0 {
						cname = c.CNAME
					}

				default:
					if err := result.p.SkipAnswer(); err != nil {
						lastErr = &DNSError{
							Err:    errCannotUnmarshalDNSMessage.Error(),
							Name:   name,
							Server: result.server,
						}
						break loop
					}
					continue
				}
			}
		}
		if hitStrictError {
			// If either family hit an error with StrictErrors enabled,
			// discard all addresses. This ensures that network flakiness
			// cannot turn a dualstack hostname IPv4/IPv6-only.
			addrs = nil
			break
		}
		if len(addrs) > 0 || network == "CNAME" && cname.Length > 0 {
			break
		}
	}
	if lastErr, ok := lastErr.(*DNSError); ok {
		// Show original name passed to lookup, not suffixed one.
		// In general we might have tried many suffixes; showing
		// just one is misleading. See also golang.org/issue/6324.
		lastErr.Name = name
	}
	sortByRFC6724(addrs)
	if len(addrs) == 0 && !(network == "CNAME" && cname.Length > 0) {
		if order == hostLookupDNSFiles {
			var canonical string
			addrs, canonical = goLookupIPFiles(name)
			if len(addrs) > 0 {
				var err error
				cname, err = dnsmessage.NewName(canonical)
				if err != nil {
					return nil, dnsmessage.Name{}, err
				}
				return addrs, cname, nil
			}
		}
		if lastErr != nil {
			return nil, dnsmessage.Name{}, lastErr
		}
	}
	return addrs, cname, nil
}

```


## 应用 --> VictoriaMetrics 使用

```go
// https://github.com/VictoriaMetrics/VictoriaMetrics/blob/41e413537167c6a90eb19f67cdab168387710ec0/lib/netutil/netutil.go
type resolver interface {
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
	LookupMX(ctx context.Context, name string) ([]*net.MX, error)
}

// Resolver is default DNS resolver.
var Resolver resolver

func init() {
	Resolver = &net.Resolver{
		PreferGo:     true, //  使用Go纯Go解析器（避免依赖系统cgo）
		StrictErrors: true,
	}
}

```


## github.com/miekg/dns-->github.com/miekg/dns
Go的 net 包提供了开箱即用的DNS查询功能，而第三方库如 miekg/dns 则赋予开发者构建自定义DNS服务器的能力。

github.com/miekg/dns 是一个用于 Go 语言的 DNS 库，它提供了丰富的功能来处理 DNS 查询、响应、解析和构建 DNS 消息。


### 主要功能
1. DNS 查询与响应：
支持各种类型的 DNS 查询（如 A、AAAA、MX、NS、TXT 等）。
支持发送和接收 DNS 消息。

2. DNS 消息构建与解析：
提供便捷的方法来构建 DNS 消息（如请求和响应）。
能够解析 DNS 消息并提取所需信息。

3. DNS 服务器：
可以用来构建自定义的 DNS 服务器。
支持处理多种类型的 DNS 请求。

4. 扩展性：
支持扩展和自定义，允许用户添加自己的功能和处理逻辑。

## 参考

- [深入解析Go语言net库：Resolver的全方位DNS解析实战指南](https://blog.csdn.net/tekin_cn/article/details/150571109)
- [go语言中强大的DNS库](https://www.cnblogs.com/guangdelw/p/18306424)
- [Go语言实现DNS解析与域名服务的实践与优化](https://juejin.cn/post/7526383867101642806)



