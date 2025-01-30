<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [netip](#netip)
  - [基本知识](#%E5%9F%BA%E6%9C%AC%E7%9F%A5%E8%AF%86)
    - [IPv6（Internet Protocol Version 6）](#ipv6internet-protocol-version-6)
  - [标准库 net.IP 的问题](#%E6%A0%87%E5%87%86%E5%BA%93-netip-%E7%9A%84%E9%97%AE%E9%A2%98)
  - [Addr 类型](#addr-%E7%B1%BB%E5%9E%8B)
  - [Prefix 类型](#prefix-%E7%B1%BB%E5%9E%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# netip 


Go1.18 新特性：引入新的 netIp 网络库.专门用于处理网络地址和相关操作。
它旨在提供一种更有效、更安全的方式来处理 IP 地址和网络前缀。
与传统的 net 包相比，net/netip 包在处理大量 IP 数据时表现出更优的性能，并且在 API 设计上更加简洁明了。


## 基本知识
### IPv6（Internet Protocol Version 6）
网络层协议的第二代标准协议，也被称为IPng（IP Next Generation），它所在的网络层提供了无连接的数据传输服务。
IPv6是IETF设计的一套规范，是IPv4的升级版本。它解决了目前IPv4存在的许多不足之处，IPv6和IPv4之间最显著的区别就是IP地址长度从原来的32位升级为128位。

```go
// /go1.23.0/src/net/ip.go
// IP address lengths (bytes).
const (
	IPv4len = 4
	IPv6len = 16
)

```
IPv6的128位IP地址书写格式有以下两种表示形式

- X:X:X:X:X:X:X:X

在这种形式中，128位的IPv6地址被分为8组


- X:X:X:X:X:X:d.d.d.d IPv4映射IPv6地址
```go
var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// 构造一个 IPv4 类型的 IP 地址
func IPv4(a, b, c, d byte) IP {
	p := make(IP, IPv6len)
	copy(p, v4InV6Prefix)
	p[12] = a
	p[13] = b
	p[14] = c
	p[15] = d
	return p
}
```

## 标准库 net.IP 的问题


Brad Fitzpatrick 对于标准库 net.IP 的问题: https://github.com/golang/go/issues/18804
```go
// go1.23.0/src/net/ip.go
type IP []byte
```

- 可变的。net.IP 的底层类型是 []byte，它的定义是：type IP []byte，这意味着你可以随意修改它。不可变数据结构更安全、更简单。
- 不可比较的。因为 Go 中 slice 类型是不可比较的，也就是说 net.IP 不支持 ==，也不能作为 map 的 key。
- 有两个 IP 地址类型，net.IP 表示基本的 IPv4 或 IPv6 地址，而 net.IPAddr 表示支持 zone scopes 的 IPv6。因为有两个类型，使用时就存在选择问题，到底使用哪个。



## Addr 类型
netip.Addr 类型可以表示 IPv4 或者 IPv6 地址，它是整个 netip 包的核心类型
```go
// /Users/python/go/go1.23.0/src/net/netip/netip.go
type Addr struct {
    // 它将 16 字节长的地址以大端序的方式存储在该 uint128 类型中，其中最高有效位在 uint128.hi 中、最低有效位在 uint128.lo 中
	addr uint128

	// Details about the address, wrapped up together and canonicalized.
	z unique.Handle[addrDetail]
}
```

```go
type uint128 struct {
	hi uint64
	lo uint64
}
```
例如，对于 IPv6 地址 0011:2233:4455:6677:8899:aabb:ccdd:eeff 按照如下方式保存：

- addr.hi = 0x0011223344556677
- addr.lo = 0x8899aabbccddeeff

使用两个 uint64 类型的字段而不是 [16]byte 来保存地址数据，使得 IP 地址的绝大多数操作都可以变成 64 位寄存器的算术/位运算，这样比字节操作运算更快。


对于 IPv4 地址，addr 字段保存的是 IPv4-mapped 格式的 IPv6 地址。此时为了快速判断一个 Addr 对象的 IP 地址类型，就需要依赖 z 字段了。
它使用了 unique 包来实现 addrDetail 类型值的 intern 化，用于表示该 IP 地址的细节信息，例如区分 IPv4/IPv6 地址、保存 IPv6 地址的 zone 信息。
使用 intern 技术来保存地址的详细信息可以节省内存，同时提高比较效率。

netip 包预定义了如下 Addr.z 类型的变量，分别表示 Addr 零值、IPv4 地址、IPv6 地址（不含 zone 信息）：
```go
var (
	z0    unique.Handle[addrDetail] //  Addr 零值
	z4    = unique.Make(addrDetail{}) // IPv4 地址
	z6noz = unique.Make(addrDetail{isV6: true}) // IPv6 地址
)
```


```go
// AddrFromSlice 从字节切片生成 Addr 对象，字节切片长度必须是 4 或 16，否则无效
func AddrFromSlice(slice []byte) (ip Addr, ok bool) {
	switch len(slice) {
	case 4: // ipv4
	    // 从 4 字节的 IPv4 地址数据生成 Addr 对象
		return AddrFrom4([4]byte(slice)), true
	case 16: // ipv6
		return AddrFrom16([16]byte(slice)), true
	}
	return Addr{}, false
}


func AddrFrom4(addr [4]byte) Addr {
	return Addr{
		addr: uint128{0, 0xffff00000000 | uint64(addr[0])<<24 | uint64(addr[1])<<16 | uint64(addr[2])<<8 | uint64(addr[3])},
		z:    z4,
	}
}
```


```go
// 解析字符串格式的 IP 地址，返回 Addr 对象。如果解析失败，则返回 Addr 零值以及错误
// 该函数可以正确解析形如 "192.0.2.1"、"2001:db8::68"、"fe80::1cc0:3e8c:119f:c2e1%ens18"（携带 zone 标识的 IPv6 地址）
// 对于 IPv4-mapped 格式的 IPv6 地址，例如 ::FFFF:192.168.1.1 会被认为是 IPv6 地址
func ParseAddr(s string) (Addr, error) {
	for i := 0; i < len(s); i++ { // 根据字符串中遇到的第一个 `.` 或 `:` 符号来判断待解析的字符串是 IPv4 还是 IPv6 地址，然后再根据对应的规则去解析地址中的各个字节
		switch s[i] {
		case '.':
			return parseIPv4(s)
		case ':':
			return parseIPv6(s)
		case '%':
			// Assume that this was trying to be an IPv6 address with
			// a zone specifier, but the address is missing.
			return Addr{}, parseAddrError{in: s, msg: "missing IPv6 address"}
		}
	}
	return Addr{}, parseAddrError{in: s, msg: "unable to parse IP"}
}
```

```go
// 返回该 IP 地址的位长，返回值可以是 0（Addr 零值）、32（IPv4 地址）或 128（IPv6 地址）
func (ip Addr) BitLen() int {
	switch ip.z {
	case z0:
		return 0
	case z4:
		return 32
	}
	return 128
}
```

判断某个 Addr 是否属于某类特殊地址：
```go
// 是否是链路本地地址
// 对于 IPv4，为 169.254.0.0/16
// 对于 IPv6，为 fe80::/10
func (ip Addr) IsLinkLocalUnicast() bool

// 是否是环回地址
// 对于 IPv4，为 127.0.0.0/8
// 对于 IPv6，为 ::1/128
func (ip Addr) IsLoopback() bool

// 是否是组播地址
// 对于 IPv4，为 224.0.0.0/4
// 对于 IPv6，为 0xFF00::/8
func (ip Addr) IsMulticast() bool

// 是否是接口本地多播地址
// https://datatracker.ietf.org/doc/html/rfc4291#section-2.7.1
// FF01::/16
func (ip Addr) IsInterfaceLocalMulticast() bool

// 是否是链路本地多播地址
// 对于 IPv4，为 224.0.0.0/24
// 对于 IPv6，为 0xFF02::/16
func (ip Addr) IsLinkLocalMulticast() bool

// 是否是全球单播地址
func (ip Addr) IsGlobalUnicast() bool

// 是否是私有地址
// 对于 IPv4，为 10.0.0.0/8、172.16.0.0/12 和 192.168.0.0/16
// 对于 IPv6，为 fc00::/7
func (ip Addr) IsPrivate() bool

// 是否是未指定地址
// 对于 IPv4，为 0.0.0.0
// 对于 IPv6，为 ::
func (ip Addr) IsUnspecified() bool
```

```go
// 用于对指定的 IP 地址进行位掩码运算，生成一个 Prefix 对象。此时 Prefix 对象中保存的 IP 地址就是一个掩码运算后的值（即一个网络地址，主机位全为 0）
func (ip Addr) Prefix(b int) (Prefix, error) {
	if b < 0 {
		return Prefix{}, errors.New("negative Prefix bits")
	}
	effectiveBits := b
	switch ip.z {
	case z0:
		return Prefix{}, nil
	case z4:
		if b > 32 {
			return Prefix{}, errors.New("prefix length " + itoa.Itoa(b) + " too large for IPv4")
		}
		effectiveBits += 96
	default:
		if b > 128 {
			return Prefix{}, errors.New("prefix length " + itoa.Itoa(b) + " too large for IPv6")
		}
	}
	ip.addr = ip.addr.and(mask6(effectiveBits))
	return PrefixFrom(ip, b), nil
}
```

## Prefix 类型

```go
type Prefix struct {
	ip Addr

	// 该网段的掩码长度 + 1
	// 因为 0 是有效的掩码长度，为了将其与 0 值本身进行区分，这里存储为掩码长度 + 1
	// 这样当 bitsPlusOne 为 0 时，代表是 Prefix 的零值
	// 当 bitsPlusOne 为 1 时，代表掩码长度为 0
	bitsPlusOne uint8
}
```

```go
// 基于指定的 IP 地址和掩码长度生成 Prefix，注意生成的 Prefix 只是简单地保存原始的 IP 地址，并不会对 IP 地址进行位掩码运算
func PrefixFrom(ip Addr, bits int) Prefix {
	var bitsPlusOne uint8
	if !ip.isZero() && bits >= 0 && bits <= ip.BitLen() {
		bitsPlusOne = uint8(bits) + 1
	}
	return Prefix{
		ip:          ip.withoutZone(),
		bitsPlusOne: bitsPlusOne,
	}
}

```

## 参考

- [IPv6 介绍](https://info.support.huawei.com/info-finder/encyclopedia/zh/IPv6.html)
- [Go1.18 新特性：引入新的 Netip 网络库](https://www.51cto.com/article/701196.html)
- [go 库学习之 netip](https://fuchencong.com/2024/10/27/go-library-netip/)