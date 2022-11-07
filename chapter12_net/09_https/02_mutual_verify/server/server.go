package main

import (
	"crypto/tls"
	"fmt"
)

func main() {
	cert, err := tls.LoadX509KeyPair("chapter17_dataStructure_n_algorithm/06_certificate/02_x509/cert.pem", "chapter17_dataStructure_n_algorithm/06_certificate/02_x509/key.pem")
	if err != nil {
		fmt.Println("load x509 keypair error:", err)
		return
	}
	cfg := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAnyClientCert,
		InsecureSkipVerify: true,
	}
	listener, err := tls.Listen("tcp", ":18000", cfg)
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			return
		}
		fmt.Println("accept connection:", conn.RemoteAddr())
		go func() {
			for {
				// echo "request"
				var b = make([]byte, 16)
				_, err := conn.Read(b)
				if err != nil {
					fmt.Println("connection read error:", err)
					conn.Close()
					return
				}

				fmt.Println(string(b))
				_, err = conn.Write(b)
				if err != nil {
					fmt.Println("connection write error:", err)
					return
				}
			}
		}()
	}
}
