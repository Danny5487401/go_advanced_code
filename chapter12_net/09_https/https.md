# HTTPS (Secure Hypertext Transfer Protocol)安全超文本传输协议

## SAN(Subject Alternative Name)

SAN 是 SSL 标准 x509 中定义的一个扩展。使用了 SAN 字段的 SSL 证书，可以扩展此证书支持的域名，使得一个证书可以支持多个不同域名的解析.

### 应用
![](.https_images/san_hao123.png)
百度证书的扩展域名有这么多，其中还有了*.hao123.com


## 代码注意点

1. 问题：
> x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0 
 
原因：
> Since Go version 1.15, the deprecated, legacy behavior of treating the CommonName field on X.509 certificates as a hostname when no Subject Alternative Names (SAN) are present is now disabled by default.
> 
解决方式


> You may need the -addext flag.
```shell
openssl req -new -key certs/foo-bar.pem \
    -subj "/CN=foobar.mydomain.svc" \
    -addext "subjectAltName = DNS:foobar.mydomain.svc" \
    -out certs/foo-bar.csr \
    -config certs/foo-bar_config.txt
```




## 参考
1. [使用开启扩展SAN的证书](https://blog.csdn.net/m0_37322399/article/details/117308604?spm=1001.2101.3001.6650.2&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-2-117308604-blog-109230584.pc_relevant_default&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-2-117308604-blog-109230584.pc_relevant_default&utm_relevant_index=5)

