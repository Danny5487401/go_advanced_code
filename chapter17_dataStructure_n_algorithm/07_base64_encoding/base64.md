<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [base64](#base64)
  - [编/解码原理](#%E7%BC%96%E8%A7%A3%E7%A0%81%E5%8E%9F%E7%90%86)
  - [具体Encoding类型-在Go源码](#%E5%85%B7%E4%BD%93encoding%E7%B1%BB%E5%9E%8B-%E5%9C%A8go%E6%BA%90%E7%A0%81)
    - [StdEncoding](#stdencoding)
    - [URLEncoding](#urlencoding)
    - [RawStdEncoding](#rawstdencoding)
    - [RawURLEncoding](#rawurlencoding)
    - [编码](#%E7%BC%96%E7%A0%81)
    - [解码](#%E8%A7%A3%E7%A0%81)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# base64
Base64不是一种加密方式,而是一种编码方式.是一种任意二进制到文本字符串的编码方法，常用于在URL、Cookie、网页中传输少量二进制数据。

Base64是由64个字符的字母表定义的基数为64的编码/解码方案，可以将二进制数据转换为字符传输，是网络上最常见的用于传输8Bit字节码的编码方式之一。注意：采用Base64编码具有不可读性，需要解码后才能阅读。

首先使用Base64编码需要一个含有64个字符的表，这个表由大小写字母、数字、+和/组成。

目前Base64被广泛应用于计算机的各个领域，由于不同场景下对特殊字符的处理（+，/）不同，因此又根据应用场景又出现了Base64的各种改进的“变种”。
因此在使用时，必须先确认使用的是哪种encoding类型，才能正确编/解码。

## 编/解码原理

把每3个8Bit的字节转换为4个6Bit的字节（3 x 8 = 4 x 6 = 24），然后把6Bit再添两位高位0，组成四个8Bit的字节，也就是说，转换后的字符串理论上将要比原来的长1/3。
不足3个的byte，根据具体规则决定是否填充。

6bit意味着总共有2^6即64种情况，与base64的字符表可以一一对应。

解码则是一个反向的过程，则需要将每4个byte的数据根据base64转换表，转换为3个byte的数据。

## 具体Encoding类型-在Go源码

### StdEncoding
适用环境：标准环境

根据RFC 4648标准定义实现，包含特殊字符'+'、'/'，不足的部分采用'='填充，根据规则，最多有2个'='。
```go
const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// StdEncoding is the standard base64 encoding, as defined in
// RFC 4648.
var StdEncoding = NewEncoding(encodeStd)

```

### URLEncoding

适用环境：url传输

因为URL编码器会把标准Base64中的'/'和'+'字符变为形如"%XX"的形式，而这些"%"号在存入数据库时还需要再进行转换，因此采用'-'、'_'代替'/'、'+'，不足的部分采用'='填充，根据规则，最多有2个'='。

```go
const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
// URLEncoding is the alternate base64 encoding defined in RFC 4648.
// It is typically used in URLs and file names.
var URLEncoding = NewEncoding(encodeURL)

```

### RawStdEncoding
除了不填充'='外，与StdEncoding一致。
```go
// RawStdEncoding is the standard raw, unpadded base64 encoding,
// as defined in RFC 4648 section 3.2.
// This is the same as StdEncoding but omits padding characters.
var RawStdEncoding = StdEncoding.WithPadding(NoPadding)
```

### RawURLEncoding
除了不填充'='外，与URLEncoding一致
```go
// RawURLEncoding is the unpadded alternate base64 encoding defined in RFC 4648.
// It is typically used in URLs and file names.
// This is the same as URLEncoding but omits padding characters.
var RawURLEncoding = URLEncoding.WithPadding(NoPadding)

```

### 编码
```go
func (enc *Encoding) Encode(dst, src []byte) {
    if len(src) == 0 {
        return
    }
    // enc is a pointer receiver, so the use of enc.encode within the hot
    // loop below means a nil check at every operation. Lift that nil check
    // outside of the loop to speed up the encoder.
    _ = enc.encode

    di, si := 0, 0
    n := (len(src) / 3) * 3
    for si < n {
        // Convert 3x 8bit source bytes into 4 bytes
        // 将3x8 bit转换为4bytes
    	// 选取前3位byte，分别左移16、8、0位，然后进行逻辑或，得到的结果的前8位、中8位和最后8位分别对应原始的3个byte数据，如此组成新的24 bit数据val。
        val := uint(src[si+0])<<16 | uint(src[si+1])<<8 | uint(src[si+2])//通过移位运算实现，前8位为src[si+0]，中间8位为src[si+1]，最后8位为src[si+2]

        // 由val数据分别右移18、12、6、0，可以得到前6、12、18、24位数据，所有目前的数据均在后低6位中。
        // 为获取低6位，与0x3F（00111111）进行逻辑与&操作即可。
        dst[di+0] = enc.encode[val>>18&0x3F]
        dst[di+1] = enc.encode[val>>12&0x3F]
        dst[di+2] = enc.encode[val>>6&0x3F]
        dst[di+3] = enc.encode[val&0x3F]

        si += 3
        di += 4
    }

    remain := len(src) - si
    if remain == 0 {
        return
    }
    // Add the remaining small block
    val := uint(src[si+0]) << 16
    if remain == 2 {
        val |= uint(src[si+1]) << 8
    }

    dst[di+0] = enc.encode[val>>18&0x3F]
    dst[di+1] = enc.encode[val>>12&0x3F]

    // 每3个byte进行相关的4byte转换，当有剩余的byte不足3个时，此时如果需要填充，缺几个byte则补几个'='。
    switch remain {//填充
    case 2:
        dst[di+2] = enc.encode[val>>6&0x3F]
        if enc.padChar != NoPadding {
            dst[di+3] = byte(enc.padChar)
        }
    case 1:
        if enc.padChar != NoPadding {
            dst[di+2] = byte(enc.padChar)
            dst[di+3] = byte(enc.padChar)
        }
    }
}

```

### 解码

```go
func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
    if len(src) == 0 {
        return 0, nil
    }

    // Lift the nil check outside of the loop. enc.decodeMap is directly
    // used later in this function, to let the compiler know that the
    // receiver can't be nil.
    _ = enc.decodeMap

    si := 0
    for strconv.IntSize >= 64 && len(src)-si >= 8 && len(dst)-n >= 8 {
        if dn, ok := assemble64(//是否有效的base64字符
            // 8转6处理过程
            enc.decodeMap[src[si+0]],
            enc.decodeMap[src[si+1]],
            enc.decodeMap[src[si+2]],
            enc.decodeMap[src[si+3]],
            enc.decodeMap[src[si+4]],
            enc.decodeMap[src[si+5]],
            enc.decodeMap[src[si+6]],
            enc.decodeMap[src[si+7]],
        ); ok {
        	// 如果成功，将获取的前6个byte存入dst中
            binary.BigEndian.PutUint64(dst[n:], dn)
            n += 6
            si += 8
        } else {
            var ninc int
            // 包含无效字符的处理
            si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
            n += ninc
            if err != nil {
                return n, err
            }
        }
    }

    for len(src)-si >= 4 && len(dst)-n >= 4 {
    	// 4转3与8转6完全原理一致，只是使用uint32转换
        if dn, ok := assemble32(
            enc.decodeMap[src[si+0]],
            enc.decodeMap[src[si+1]],
            enc.decodeMap[src[si+2]],
            enc.decodeMap[src[si+3]],
        ); ok {
            binary.BigEndian.PutUint32(dst[n:], dn)
            n += 3
            si += 4
        } else {
            var ninc int
            si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
            n += ninc
            if err != nil {
                return n, err
            }
        }
    }

    for si < len(src) {
        var ninc int
        si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
        n += ninc
        if err != nil {
            return n, err
        }
    }
    return n, err
}

```