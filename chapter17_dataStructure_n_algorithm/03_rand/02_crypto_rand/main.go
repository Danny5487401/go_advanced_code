package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

func main() {
	// 读取32个字节的数据，然后以base64编码的形式打印出来
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(b))

	// 你想要的是一个[0,n)的整数
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", n)
}
