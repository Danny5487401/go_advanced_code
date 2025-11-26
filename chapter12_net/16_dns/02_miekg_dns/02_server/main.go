package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
)

// handleA 处理A记录查询
func handleA(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	records := map[string]string{
		"example.com.": "93.184.216.34",
		"test.com.":    "192.0.2.1",
	}
	domain := r.Question[0].Name
	if ip, ok := records[domain]; ok && r.Question[0].Qtype == dns.TypeA {
		rr, _ := dns.NewRR(fmt.Sprintf("%s 3600 IN A %s", domain, ip))
		m.Answer = append(m.Answer, rr)
	} else {
		m.SetRcode(r, dns.RcodeNameError)
	}
	if err := w.WriteMsg(m); err != nil {
		log.Printf("写入响应失败: %v", err)
	}
}

func main() {
	dns.HandleFunc(".", handleA)
	server := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Println("DNS服务器启动在 :8053")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("启动DNS服务器失败: %v", err)
	}
}

/*
✗ nslookup -port=8053 example.com 127.0.0.1
Server:         127.0.0.1
Address:        127.0.0.1#8053

Name:   example.com
Address: 93.184.216.34

*/

/*
✗ dig @127.0.0.1 -p 8053 example.com

; <<>> DiG 9.10.6 <<>> @127.0.0.1 -p 8053 example.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 7017
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;example.com.                   IN      A

;; ANSWER SECTION:
example.com.            3600    IN      A       93.184.216.34

;; Query time: 0 msec
;; SERVER: 127.0.0.1#8053(127.0.0.1)
;; WHEN: Wed Nov 26 14:50:14 CST 2025
;; MSG SIZE  rcvd: 56

*/
