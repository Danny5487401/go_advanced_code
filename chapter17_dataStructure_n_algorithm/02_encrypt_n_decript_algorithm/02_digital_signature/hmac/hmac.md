# HMAC 哈希运算消息认证码(Hash-based Message Authentication Code)
它通过一个标准算法，在计算哈希的过程中，把key混入计算过程中。

和我们自定义的加salt算法不同，Hmac算法针对所有哈希算法都通用，无论是MD5还是SHA-1。
采用Hmac替代我们自己的salt算法，可以使程序算法更标准化，也更安全。

```go
//key随意设置 data 要加密数据
func Hmac(key, data string) string {
	hash:= hmac.New(md5.New, []byte(key)) // 创建对应的md5哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}
func HmacSha256(key, data string) string {
	hash:= hmac.New(sha256.New, []byte(key)) //创建对应的sha256哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

```