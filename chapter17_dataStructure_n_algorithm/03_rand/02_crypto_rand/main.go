package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
)

func main() {
	// 需求：读取32个字节的数据，然后以base64编码的形式打印出来
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(b))

	// 需求：你想要的是一个[0,n)的整数
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", n)

	// 需求：想要生成区间[-m, n]的安全随机数
	fmt.Println(rangeRand(-1, 2))
}

// 指定区间的随机数
// 生成区间[-m, n]的安全随机数
func rangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}
