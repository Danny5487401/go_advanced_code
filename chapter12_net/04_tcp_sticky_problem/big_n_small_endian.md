# 什么是大端、小端
![](.big_n_small_endian_images/biig_n_small_endian.png)

- 大端模式：高位字节排放在内存的低地址端，低位字节排放在内存的高地址端;

- 小端模式：低位字节排放在内存的低地址端，高位字节排放在内存的高地址端；

## 为什么需要大小端字节序转化

因为在涉及到网络传输、文件存储时，因为不同系统的大小端字节序不同，这是就需要大小端转化，才能保证读取到的数据是正确的。
在做arm和dsp通信的时候，就遇到个大小端转换的问题，因为arm是小端，dsp是大端，

```go
原理
func SwapEndianUin32(val uint32)  uint32{
 return (val & 0xff000000) >> 24 | (val & 0x00ff0000) >> 8 |
  (val & 0x0000ff00) << 8 | (val & 0x000000ff) <<24
}
```
标准库
```go
// use encoding/binary
// bigEndian littleEndian
func BigEndianAndLittleEndianByLibrary()  {
    var value uint32 = 10
    by := make([]byte,4)
    binary.BigEndian.PutUint32(by,value)
    fmt.Println("转换成大端后 ",by)
    fmt.Println("使用大端字节序输出结果：",binary.BigEndian.Uint32(by))
    little := binary.LittleEndian.Uint32(by)
    fmt.Println("大端字节序使用小端输出结果：",little)
}
// 结果：
//转换成大端后  [0 0 0 10]
//使用大端字节序输出结果： 10
//大端字节序使用小端输出结果： 167772160
```
grpc中对大端的应用

gRPC封装message时，在封装header时，特意指定了使用大端字节序，

```go
// msgHeader returns a 5-byte header for the message being transmitted and the
// payload, which is compData if non-nil or data otherwise.
func msgHeader(data, compData []byte) (hdr []byte, payload []byte) {
    hdr = make([]byte, headerLen)
    if compData != nil {
    hdr[0] = byte(compressionMade)
    data = compData
    } else {
    hdr[0] = byte(compressionNone)
    }
    
    // Write length of payload into buf
    binary.BigEndian.PutUint32(hdr[payloadLen:], uint32(len(data)))
    return hdr, data
}
```