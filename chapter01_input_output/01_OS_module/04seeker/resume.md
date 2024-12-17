<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [断点续传](#%E6%96%AD%E7%82%B9%E7%BB%AD%E4%BC%A0)
  - [文件 io seeker 接口: io的偏移操作封装](#%E6%96%87%E4%BB%B6-io-seeker-%E6%8E%A5%E5%8F%A3-io%E7%9A%84%E5%81%8F%E7%A7%BB%E6%93%8D%E4%BD%9C%E5%B0%81%E8%A3%85)
  - [网络断点续传原理](#%E7%BD%91%E7%BB%9C%E6%96%AD%E7%82%B9%E7%BB%AD%E4%BC%A0%E5%8E%9F%E7%90%86)
    - [range](#range)
    - [Content-Range](#content-range)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 断点续传


## 文件 io seeker 接口: io的偏移操作封装
在文件中移动文件指针位置的方法。它可以用来定位文件中的特定位置，以便读取或写入数据

使用
```go
file, _ := os.OpenFile("chapter01_input_output/files/danny.txt", os.O_RDWR, 0)
file.Seek(3, io.SeekCurrent)
```

```go
// go1.23.0/src/io/io.go

// Seek whence values.
const (
	SeekStart   = 0 // 从文件开头开始计算偏移量
	SeekCurrent = 1 // 从当前位置开始计算偏移量。
	SeekEnd     = 2 // 从文件末尾开始计算偏移量。
)
```
```go
// go1.23.0/src/os/file.go
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	// 校验文件
	if err := f.checkValid("seek"); err != nil {
		return 0, err
	}
	
	r, e := f.seek(offset, whence)
	if e == nil && f.dirinfo.Load() != nil && r != 0 {
		e = syscall.EISDIR
	}
	if e != nil {
		return 0, f.wrapErr("seek", e)
	}
	return r, nil
}

```




## 网络断点续传原理

HTTP1.1 协议（RFC2616）开始支持获取文件的部分内容，这为并行下载以及断点续传提供了技术支持：通过在 Header里两个参数Range和Content-Range实现

客户端发请求时对应的是 Range，服务器端响应时对应的是 Content-Range。



### range 
Range是一个请求头，它告知了服务器返回文件的哪一部分。在一个 Range 中，可以一次性请求多个部分，服务器会以 multipart 文件的形式将其返回。

如果服务器返回的是范围响应，则需要使用 206 Partial Content 状态码。

如果所请求的范围不合法，那么服务器会返回 416 Range Not Satisfiable 状态码，表示客户端错误。

同时，服务器允许忽略Range，从而返回整个文件，此时状态码仍然是200


### Content-Range
```shell
$ curl --location --head 'https://download.jetbrains.com/go/goland-2020.2.2.exe'
HTTP/1.1 200 Connection established

HTTP/2 302
content-type: text/html
content-length: 138
location: https://download-cdn.jetbrains.com/go/goland-2020.2.2.exe
date: Tue, 17 Dec 2024 02:08:17 GMT
server: nginx
strict-transport-security: max-age=31536000; includeSubdomains;
x-frame-options: DENY
x-content-type-options: nosniff
x-xss-protection: 1; mode=block;
x-geocode: SG
x-cache: Miss from cloudfront
via: 1.1 869c20a0b6637fa4614a52064a4bf808.cloudfront.net (CloudFront)
x-amz-cf-pop: SIN2-P2
alt-svc: h3=":443"; ma=86400
x-amz-cf-id: cBvxluBrxCWGZZ5uM5Qhcv5kYKzN88n6WuEx4N24a7LOJqm2s7A7AQ==

HTTP/1.1 200 Connection established

HTTP/2 200
content-type: binary/octet-stream
content-length: 338589968
date: Tue, 17 Dec 2024 02:08:19 GMT
last-modified: Wed, 12 May 2021 19:41:30 GMT
etag: "1312fd0956b8cd529df1100d5e01837f-41"
x-amz-version-id: RMqUwQm39.p77fpawMulHTL_4YCJWjeH
accept-ranges: bytes
server: AmazonS3
x-cache: Miss from cloudfront
via: 1.1 5659c4bfa12ab1d4105fc650d6eb1624.cloudfront.net (CloudFront)
x-amz-cf-pop: SIN2-P3
alt-svc: h3=":443"; ma=86400
x-amz-cf-id: 1KKrT0i1vLyoBju9KCL2vikPf0-8Xa6U-tR9-3Jgxjg5ZSXvuQzIfg==
```
Accept-Ranges: bytes 表示界定范围的单位是 bytes，这里 Content-Length 也是很有用的信息，因为它提供了要检索的图片的完整大小！


如果站点返回的Header中不包括Accept-Ranges，那么它有可能不支持范围请求。一些站点会明确将其值设置为 “none”，以此来表明不支持