<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [URL(Uniform Resource Locator) 统一资源定位符](#urluniform-resource-locator-%E7%BB%9F%E4%B8%80%E8%B5%84%E6%BA%90%E5%AE%9A%E4%BD%8D%E7%AC%A6)
  - [背景](#%E8%83%8C%E6%99%AF)
  - [源码net/url分析](#%E6%BA%90%E7%A0%81neturl%E5%88%86%E6%9E%90)
    - [转义 QueryEscape](#%E8%BD%AC%E4%B9%89-queryescape)
    - [反转义 QueryUnescape](#%E5%8F%8D%E8%BD%AC%E4%B9%89-queryunescape)
    - [parse](#parse)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# URL(Uniform Resource Locator) 统一资源定位符

## 背景

URL 是连接用户与服务的媒介。本文将以 golang 1.11.1 版本标准库中的 net/url 进行分析，进而了解 URL 规范在实际应用中的实现。
golang 基本遵循 RFC3986 进行设计，因为一些兼容问题会作此许的修改
golang中为我们提供了url.Parse,专门用来解析协议的url，比如http、https,rtsp,rtmp 等协议.

为什么能够解析这些不同的协议呢，只要遵循下面的格式即可

```go
// 双斜杠地址
[scheme:][//[userinfo@]host][/]path[?query][#fragment]
 
// 非双斜杠地址
scheme:opaque[?query][#fragment]
```

我们可以称每个 [] 为一个组件。规范中对于每一个 [] 中都有一定的转义规则，转义的好处是可以更好得在互联网上进行传播，而不会丢失数据.

## 源码net/url分析

这个包对外主要提供了 URL 的解析 Parse，query 数据的转义与反转义 QueryEscape, QueryUnescape。我认为，转义和解析是这个包的主要功能。

### 转义 QueryEscape
RL 中的转义为把非安全的字符转义为包含一个百分号(%)(%)后面跟着两个表示字符 ASCII 码的十六进制数。
```go
// 转义 [?query] 组件中需要转义的字符
func QueryEscape(s string) string {
	return escape(s, encodeQueryComponent)
}

// 转义的基本方法，按 mode 转义不同字符，mode 有(encodePath, encodePathSegment, encodeHost,
// encodeZone, encodeUserPassword, encodeQueryComponent, encodeFragment)
func escape(s string, mode encoding) string {
	// 用于计算整个字符串需要占用多少空间，及判断是否需要转义
	spaceCount, hexCount := 0, 0

	// 先遍历一次所有字符，计算空格及转义字符有多少
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	// 没有需要转义的，直接返回字符串
	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	// 转义后原字符会用 "%AB" 表示，所以长度增加了 2 倍的转义字符数
	t := make([]byte, len(s)+2*hexCount)
	j := 0 // t 的索引，记录写入的位置
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent: // 可以看到转义对于 query 组件，会替换空格为 + 号
			t[j] = '+'
			j++
		case shouldEscape(c, mode): // 需要转义的字符
			// 可以看到转义算法很简单
			// 1. 添加一个百分号(%)
			// 2. 取这个字符的高 4 位对应的 16 进制
			// 3. 取这个字符的低 4 位对应的 16 进制
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default: // 不需要转义的
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// 根据 mode 类型判断 字符 c 是否需要转义，所有规则都在 RFC3986 中：https://tools.ietf.org/html/rfc3986
// 按规则判断是否转义，就不翻译了。
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost || mode == encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't user %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case encodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	if mode == encodeFragment {
		// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are
		// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not
		// need to be escaped. To minimize potential breakage, we apply two restrictions:
		// (1) we always escape sub-delims outside of the fragment, and (2) we always
		// escape single quote to avoid breaking callers that had previously assumed that
		// single quotes would be escaped. See issue #19917.
		switch c {
		case '!', '(', ')', '*':
			return false
		}
	}

	// Everything else must be escaped.
	return true
}
```

### 反转义 QueryUnescape
```go
// 此方法是 `QueryEscape` 的逆运算
// 转换每三个像 "%AB" 的字符为十六进制 0xAB.
// 当百分号(%)后跟没有跟着正确的十六进制则抛出异常
func QueryUnescape(s string) (string, error) {
	return unescape(s, encodeQueryComponent)
}

// 按 mode 类型来反转义字符串 s，一般按组件来调用这个方法。
func unescape(s string, mode encoding) (string, error) {
	// 计数，百分号(%)的个数，也就是有多少个转义的字符数
	n := 0
	hasPlus := false // 记录 query 组件中是否出现加(+)号
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			// 三种情况说明是不合法的 URL 转义
			// 1. 百分号(%)后不足 2 位。
			// 2. 百分号(%)后一位不是合法的十六进制字符
			// 3. 百分号(%)后二位不是合法的十六进制字符
			// 这里 ishex 其它就是判断当前字符是否在 "0123456789ABCDEF" 内
			// 用的字符直接比较，'0' <= c && c <= '9', 'a' <= c && c <= 'f'
			// 'A' <= c && c <= 'F'，我想这么比较应该会比字符串比较来得更快。
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return "", EscapeError(s)
			}
			// https://tools.ietf.org/html/rfc3986#page-21
			// 在 host 组件中 "%AB" 这种转义方式只能对于非 ASIIC 码中的字符
			// 不过在这个规范中 https://tools.ietf.org/html/rfc6874#section-2
			// 提及了 %25 允许在 host 组件的 IPv6作用域地址进行转义
			// unhex(s[i+1]) < 8 的意思是字符 s[i+1] 是 ASIIC 码中的字符
			// 我是这么理解的，unhex(s[i+1]) < 8 表示取值范围是 s[i+1] 的取值范围
			// 是 0-7， s[i+2] 取值为 0-15，正好是 8*16=128，十六进行的前 127 个为
			// ASIIC 码中的字符。
			if mode == encodeHost && unhex(s[i+1]) < 8 && s[i:i+3] != "%25" {
				return "", EscapeError(s[i : i+3])
			}
			if mode == encodeZone {
				// 这段没看懂，有朋友懂的还请告知，原注释如下
				// RFC 6874 says basically "anything goes" for zone identifiers
				// and that even non-ASCII can be redundantly escaped,
				// but it seems prudent to restrict %-escaped bytes here to those
				// that are valid host name bytes in their unescaped form.
				// That is, you can user escaping in the zone identifier but not
				// to introduce bytes you couldn't just write directly.
				// But Windows puts spaces here! Yay.
				v := unhex(s[i+1])<<4 | unhex(s[i+2])
				if s[i:i+3] != "%25" && v != ' ' && shouldEscape(v, encodeHost) {
					return "", EscapeError(s[i : i+3])
				}
			}
			i += 3
		case '+':
			hasPlus = mode == encodeQueryComponent
			i++
		default:
			if (mode == encodeHost || mode == encodeZone) && s[i] < 0x80 && shouldEscape(s[i], mode) {
				return "", InvalidHostError(s[i : i+1])
			}
			i++
		}
	}

	// 没有被转义，且 query 组件中不含加号(+)
	if n == 0 && !hasPlus {
		return s, nil
	}

	// 分配最小的空间
	t := make([]byte, len(s)-2*n)
	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			// 之前 escape 的逆运算
			t[j] = unhex(s[i+1])<<4 | unhex(s[i+2])
			j++
			i += 3
		case '+':
			if mode == encodeQueryComponent {
				t[j] = ' '
			} else {
				t[j] = '+'
			}
			j++
			i++
		default:
			t[j] = s[i]
			j++
			i++
		}
	}
	return string(t), nil
}
```

### parse
```go
// 解析 rawurl(URL 字符串，可以是相对路径，不包含 host 组件) 返回 URL 结构
func Parse(rawurl string) (*URL, error) {
	// 截取 # 之后的，获取 fragment 组件
	// 此 split 方法当第三个参数
	// 为 true 时，第二个值返回不包含 # 的两段；
	// 为 false 时， 第二个值返回包含 # 的两段；
	u, frag := split(rawurl, "#", true)
	url, err := parse(u, false)
	if err != nil {
		return nil, &Error{"parse", u, err}
	}
	if frag == "" {
		return url, nil
	}
	// 反转义 frag 成功则将解析出的 fragment 添加到 URL 结构体
	if url.Fragment, err = unescape(frag, encodeFragment); err != nil {
		return nil, &Error{"parse", rawurl, err}
	}
	return url, nil
}

// 解析 URL 基础方法，第二个参数 viaRequest 为 true 时只接受绝对地址，否则允许所有形式的相对url
// 此方法授受的 rawurl 不包含 fragment 组件
func parse(rawurl string, viaRequest bool) (*URL, error) {
	var rest string
	var err error

	if rawurl == "" && viaRequest {
		return nil, errors.New("empty url")
	}
	url := new(URL)

	if rawurl == "*" {
		url.Path = "*"
		return url, nil
	}

	// 获取 URL 中的 Scheme，例如："http:", "mailto:", "ftp:"。不能包含转义字符。
	if url.Scheme, rest, err = getscheme(rawurl); err != nil {
		return nil, err
	}
	// 这里可以看到协议是不区分大小写的
	url.Scheme = strings.ToLower(url.Scheme)

	// 以问号(?)结尾，且只有一个问号(?)。reset 重置，去除问号(?)
	if strings.HasSuffix(rest, "?") && strings.Count(rest, "?") == 1 {
		url.ForceQuery = true
		rest = rest[:len(rest)-1]
	} else {
		// 截取出 reset([//[userinfo@]host][/]path)和原始的 query 组件(url.RawQuery)的内容
		rest, url.RawQuery = split(rest, "?", true)
	}

	// reset 不是以 / 开头，按理说我们会截出像([//[userinfo@]host][/]path)，会以 // 开头
	if !strings.HasPrefix(rest, "/") {
		if url.Scheme != "" {
			// We consider rootless paths per RFC 3986 as opaque.
			url.Opaque = rest
			return url, nil
		}
		if viaRequest {
			return nil, errors.New("invalid URI for request")
		}

		// Avoid confusion with malformed schemes, like cache_object:foo/bar.
		// See golang.org/issue/16822.
		//
		// RFC 3986, §3.3:
		// In addition, a URI reference (Section 4.1) may be a relative-path reference,
		// in which case the first path segment cannot contain a colon (":") character.
		colon := strings.Index(rest, ":")
		slash := strings.Index(rest, "/")
		if colon >= 0 && (slash < 0 || colon < slash) {
			// First path segment has colon. Not allowed in relative URL.
			return nil, errors.New("first path segment in URL cannot contain colon")
		}
	}

	if (url.Scheme != "" || !viaRequest && !strings.HasPrefix(rest, "///")) && strings.HasPrefix(rest, "//") {
		var authority string

		// 下列对应关系
		// rest[2:] = [userinfo@]host][/]path
		// authority = [userinfo@]host
		// rest = [/]path
		authority, rest = split(rest[2:], "/", false)
		url.User, url.Host, err = parseAuthority(authority)
		if err != nil {
			return nil, err
		}
	}
	// 设置 URL 的 path 组件，如果 reset 被转义，则还会设置 URL 结构体的 RawPath 的值
	if err := url.setPath(rest); err != nil {
		return nil, err
	}
	return url, nil
}

// 设置 p 为 URL 的 path 组件，此方法将反转义 p
// p 被包含转义则 RawPath 为 p，Path 为反转义的 p
// 否则 RawPath 为空，Path 为 p
//  例:
// - setPath("/foo/bar")   Path="/foo/bar" , RawPath=""
// - setPath("/foo%2fbar") Path="/foo/bar" , RawPath="/foo%2fbar"
// p 包含不合法转义字符时，将抛出异常
func (u *URL) setPath(p string) error {
	path, err := unescape(p, encodePath)
	if err != nil {
		return err
	}
	u.Path = path
	if escp := escape(path, encodePath); p == escp {
		// 没有转义，原始的 path 既是空。
		u.RawPath = ""
	} else {
		u.RawPath = p
	}
	return nil
}
```