package main

import (
	"fmt"
	"net/netip"
)

func main() {
	// Addr 是表示单个IP地址的基础类型
	// IPv4 地址
	addrV4, ok := netip.AddrFromSlice([]byte{192, 0, 2, 1})
	if !ok {
		fmt.Println("无效的IPv4地址")
	}
	fmt.Println(addrV4.Is4())

	// IPv6 地址
	addrV6, ok := netip.AddrFromSlice([]byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	if !ok {
		fmt.Println("无效的IPv6地址：")
	}
	fmt.Println(addrV6.Is6())

	// Prefix 类型代表一个网络段，它由一个 Addr 类型的 IP 地址和一个表示网络掩码长度的整数组成
	// 创建一个 IPv4 前缀
	prefixV4 := netip.PrefixFrom(addrV4, 24)
	isValid := prefixV4.IsValid()
	if isValid {
		fmt.Println("这是一个有效的前缀：", prefixV4)
	} else {
		fmt.Println("这是一个无效的前缀")
	}

	// 创建一个 IPv6 前缀
	prefixV6 := netip.PrefixFrom(addrV6, 64)
	fmt.Println(prefixV6.String())

	// 解析IP地址
	addr, err := netip.ParseAddr("192.0.2.10")
	if err != nil {
		fmt.Println("地址解析失败：", err)
	} else {
		fmt.Println("解析的地址是：", addr)
	}

	// 检查IP地址是否属于某个前缀
	if prefixV4.Contains(addr) {
		fmt.Println(addr, "属于", prefixV4)
	} else {
		fmt.Println(addr, "不属于", prefixV4)
	}

}
