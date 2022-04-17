package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

// 生成消息认证码
// plainText 明文
// key 密钥
// 返回 消息认证码
func GenerateMAC(plainText, key []byte) []byte {
	hash := hmac.New(sha1.New, key)
	hash.Write(plainText)
	hashText := hash.Sum(nil)
	return hashText
}

// 消息认证
// plainText 明文
// key 密钥
// hashText 消息认证码
// 返回 是否是原消息
func VerifyMAC(plainText, key, hashText []byte) bool {
	hash := hmac.New(sha1.New, key)
	hash.Write(plainText)
	return hmac.Equal(hashText, hash.Sum(nil))
}

func main() {
	plainText := []byte("消息")
	key := []byte("私钥")
	hashText := GenerateMAC(plainText, key)
	ok := VerifyMAC(plainText, key, hashText)
	if ok {
		fmt.Printf("%s 是原消息\n", plainText)
	}
	fakeText := []byte("假消息")
	ok = VerifyMAC(plainText, key, fakeText)
	if !ok {
		fmt.Printf("%s 是假消息\n", fakeText)
	}
}
