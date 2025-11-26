package main

import (
	"fmt"

	"github.com/miekg/dns"
)

var (
	msg       = new(dns.Msg)
	dnsServer = "223.6.6.6:53"
	client    = new(dns.Client)
)

// 获取 A 记录
func ResolveARecord(domain string) error {
	// 创建 DNS 消息

	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	// 发送 DNS 查询
	response, _, err := client.Exchange(msg, dnsServer)
	if err != nil {
		fmt.Printf("Failed to query DNS for %s: %v\n", domain, err)
		return err
	}

	// 处理响应
	if response.Rcode != dns.RcodeSuccess {
		fmt.Printf("DNS query failed with Rcode %d\n", response.Rcode)
		return fmt.Errorf("DNS query failed with Rcode %d", response.Rcode)
	}
	// 判断是否来自权威服务器
	if response.Authoritative {
		fmt.Println("来自权威服务器")
	}
	// 判断是否递归查询
	if response.RecursionAvailable {
		fmt.Printf("%s,这是递归查询\n", domain)
	}

	// 打印 A 记录
	for _, ans := range response.Answer {
		if aRecord, ok := ans.(*dns.A); ok {
			fmt.Printf("A record for %s: %s\n", domain, aRecord.A.String())
		}
	}
	return nil
}

// 获取 AAAA 记录
func ResolveAAAARecord(domain string) error {
	// 创建 DNS 消息

	msg.SetQuestion(dns.Fqdn(domain), dns.TypeAAAA)

	// 发送 DNS 查询
	response, _, err := client.Exchange(msg, dnsServer)
	if err != nil {
		fmt.Printf("Failed to query DNS for %s: %v\n", domain, err)
		return err
	}

	// 处理响应
	if response.Rcode != dns.RcodeSuccess {
		fmt.Printf("DNS query failed with Rcode %d\n", response.Rcode)
		return fmt.Errorf("DNS query failed with Rcode %d", response.Rcode)
	}

	// 打印 AAAA 记录
	for _, ans := range response.Answer {
		if aRecord, ok := ans.(*dns.AAAA); ok {
			fmt.Printf("AAAA record for %s: %s\n", domain, aRecord.AAAA.String())
		}
	}
	return nil
}

// 获取 TXT 记录
func ResolveTXTRecord(domain string) error {
	// 创建 DNS 消息

	msg.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	// 发送 DNS 查询
	response, _, err := client.Exchange(msg, dnsServer)
	if err != nil {
		fmt.Printf("Failed to query DNS for %s: %v\n", domain, err)
		return err
	}

	// 处理响应
	if response.Rcode != dns.RcodeSuccess {
		fmt.Printf("DNS query failed with Rcode %d\n", response.Rcode)
		return fmt.Errorf("DNS query failed with Rcode %d", response.Rcode)
	}

	// 打印 TXT 记录
	for _, ans := range response.Answer {
		if aRecord, ok := ans.(*dns.TXT); ok {
			for _, txt := range aRecord.Txt {
				fmt.Printf("TXT record for %s: %s\n", domain, txt)
			}
		}
	}
	return nil
}

// 获取 NS 记录
func ResolveNSRecord(domain string) error {
	// 创建 DNS 消息

	msg.SetQuestion(dns.Fqdn(domain), dns.TypeNS)

	// 发送 DNS 查询
	response, _, err := client.Exchange(msg, dnsServer)
	if err != nil {
		fmt.Printf("Failed to query DNS for %s: %v\n", domain, err)
		return err
	}

	// 处理响应
	if response.Rcode != dns.RcodeSuccess {
		fmt.Printf("DNS query failed with Rcode %d\n", response.Rcode)
		return fmt.Errorf("DNS query failed with Rcode %d", response.Rcode)
	}

	// 打印 TXT 记录
	for _, ans := range response.Answer {
		if aRecord, ok := ans.(*dns.NS); ok {
			fmt.Printf("NS record for %s: %s\n", domain, aRecord.Ns)
		}
	}
	return nil
}

func main() {
	// 要查询的域名
	domain := "www.baidu.com"
	ResolveARecord(domain)
	ResolveAAAARecord(domain)
	ResolveTXTRecord(domain)
	ResolveNSRecord(domain)
}
