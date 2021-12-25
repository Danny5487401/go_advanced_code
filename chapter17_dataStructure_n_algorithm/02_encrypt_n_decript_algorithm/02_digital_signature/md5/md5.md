# 数字签名
## hash 函数，有加密型和非加密型。
加密型的一般用于加密数据、数字摘要等，典型代表就是 md5、sha1、sha256、aes256 这种；
非加密型的一般就是查找。在 map 的应用场景中，用的是查找。选择 hash 函数主要考察的是两点：性能、碰撞概率
### MD5信息摘要算法(Message-DigestAlgorithm 5)
它可以把一个任意长度的字节数组转换成一个定长的整数，并且这种转换是不可逆的。对于任意长度的数据，转换后的MD5值长度是固定的，
而且MD5的转换操作很容易，只要原数据有一点点改动，转换后结果就会有很大的差异。正是由于MD5算法的这些特性，它经常用于对于一段信息产生信息摘要，以防止其被篡改。
其还广泛就于操作系统的登录过程中的安全验证，比如Unix操作系统的密码就是经过MD5加密后存储到文件系统中，当用户登录时输入密码后，
对用户输入的数据经过MD5加密后与原来存储的密文信息比对，如果相同说明密码正确，否则输入的密码就是错误的。

当前两个不同值产生相同的hash还是有可能的，只是这个可能性很小很小

MD5以512位为一个计算单位对数据进行分组，每一分组又被划分为16个32位的小组，经过一系列处理后，输出4个32位的小组，最后组成一个128位的哈希值。
对处理的数据进行512求余得到N和一个余数，如果余数不为448,填充1和若干个0直到448位为止，最后再加上一个64位用来保存数据的长度，这样经过预处理后，数据变成（N+1）x 512位。


go提供的hash接口
```go
type Hash interface {
    // 通过io.Writer接口的Write方法向hash中添加数据
    io.Writer
    
    // 返回添加b到当前的hash值后的新切片，不会改变底层的hash状态，这个方法就是返回计算后的hash值，只是它是字符切片
    Sum(b []byte) []byte
    
    // 重设hash为无数据输入的状态，就是清空hash之前写入的数据
    Reset()
    
    // 返回Sum会返回的切片的长度
    Size() int
    
    // 返回hash底层的块大小；Write方法可以接受任何大小的数据，
    // 但提供的数据是块大小的倍数时效率更高
    BlockSize() int
}
```

Hash包还有两个Hash接口：
```go
type Hash32 interface { // Hash32是一个被所有32位hash函数实现的公共接口。
    Hash
    Sum32() uint32
}
type Hash64 interface { // Hash64是一个被所有64位hash函数实现的公共接口。
    Hash
    Sum64() uint64
}
```
md5实现的Hash接口是16位的hash函数（它的Sum方法返回的字符切片长度为16位），Hash32和hash64是属于安全性更高的两个Hash函数，产生的hash值也更长。