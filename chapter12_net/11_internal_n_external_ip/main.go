package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	fmt.Println("---内网IP---")
	getInternal()
	fmt.Println("---外IP---")
	getExternal()

}

func getExternal() {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func getInternal() {
	//  InterfaceAddrs 返回系统的单播接口地址列表
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
}
