# strconv包

## strconv/itoa （int或int64转字符串）

对外提供的接口,Itoa函数具体的实现是由FormatInt()函数实现
```go
func Itoa(i int) string {
	return FormatInt(int64(i), 10)
}

```

【扩展阅读 a的典故】这是C语言遗留下的典故。C语言中没有string类型而是用字符数组(array)表示字符串，所以Itoa对很多C系的程序员很好理解。
```go
// 快速模式
const fastSmalls = true

// 小整数最大值
const nSmalls = 100

```

FormatInt是真正整数转字符串的入口函数.可以分为快速模式和一般模式.
```go
func FormatInt(i int64, base int) string {
	// i属于小整数(0 <= i <= nSmalls),通过快速模式转换
	// 快速模式仅支持10进制转换
	if fastSmalls && 0 <= i && i < nSmalls && base == 10 {
		return small(int(i))
	}

   // 一般模式
	_, s := formatBits(nil, uint64(i), base, i < 0, false)
	return s
}
```

快速模式(small)
```go
// 转换范围
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

const smallsString = "00010203040506070809" +
    "10111213141516171819" +
    "20212223242526272829" +
    "30313233343536373839" +
    "40414243444546474849" +
    "50515253545556575859" +
    "60616263646566676869" +
    "70717273747576777879" +
    "80818283848586878889" +
    "90919293949596979899"

// small returns the string for an i with 0 <= i < nSmalls.
func small(i int) string {
	if i < 10 {
		return digits[i : i+1]
	}
	return smallsString[i*2 : i*2+2]
}

```

一般模式(formatBits)
```go
func formatBits(dst []byte, u uint64, base int, neg, append_ bool) (d []byte, s string) {
	// neg == true, 负数
	// _append == true, 字符串追加 

	// 将运算后结果,写入a数组
	var a [64 + 1]byte

	// 数组长度,然后对该值进行递减
	i := len(a)

	// 1. 10进制,使用 / 和 % 进行运算.
	// 2. Base是2幂,使用shift 和 masks运算.
	// 3. 其它情况,从低位开始通过/运算
	if base == 10 {
		// 32位平台
		if host32bit {
			// 每次运算1e9,分多次执行.
			// 从低位开始运算
			for u >= 1e9 {
				q := u / 1e9
				us := uint(u - q*1e9)			// 低位待运算整数
				for j := 4; j > 0; j-- {
					is := us % 100 * 2      	// %100被用来取后两位整数
					us /= 100
					
					// 每次运算2位
					i -= 2
					a[i+1] = smallsString[is+1] 
					a[i+0] = smallsString[is+0]
				}

				
				// 运算剩下的一位(1e9)
				i--
				a[i] = smallsString[us*2+1]

				// Loop 1e9
				u = q
			}
		}

		// 整数>=2位,按2位进行运算
		us := uint(u)
		for us >= 100 {
			is := us % 100 * 2
			us /= 100
			i -= 2
			a[i+1] = smallsString[is+1]
			a[i+0] = smallsString[is+0]
		}

		// 不满足2位时,1位进行运算
		is := us * 2
		i--
		a[i] = smallsString[is+1]
		if us >= 10 {
			i--
			a[i] = smallsString[is]
		}

	} else if isPowerOfTwo(base) {
		// 1. 对于Base是2的幂,可以通过shift和masks进行运算,提升效率.
		// 2. Base位移操作<=36,即最大shift为1<<5(32).
		// 3. 通过使用&-ind(7),告诉编译器位移应小于8;小于寄存器宽度,代码优化。
		shift := uint(bits.TrailingZeros(uint(base))) & 7
		b := uint64(base)
		m := uint(base) - 1 // == 1<<shift - 1
		for u >= b {
			i--
			a[i] = digits[uint(u)&m]
			u >>= shift
		}

		// u < base
		i--
		a[i] = digits[uint(u)]

	} else {
		// 其它进制情况
		b := uint64(base)

		// 整数大于base,从低位开始运算，每次运算1位
		for u >= b {
			i--
			q := u / b
			a[i] = digits[uint(u-q*b)]
			u = q
		}

		// u < base
		i--
		a[i] = digits[uint(u)]
	}

	// 负数,增加'-'
	if neg {
		i--
		a[i] = '-'
	}

	// 和另一个字符串进行追加
	if append_ {
		d = append(dst, a[i:]...)
		return
	}

	s = string(a[i:])
	return
}
```


## Atoi 字符串转整形
Atoi函数是对外提供的接口函数.具体实现由ParseInt函数实现。
```go
func Atoi(s string) (int, error) {
	const fnAtoi = "Atoi"

	// 字符串长度
	sLen := len(s)

	// int类型小整数通过快速模式运算
	if intSize == 32 && (0 < sLen && sLen < 10) ||
		intSize == 64 && (0 < sLen && sLen < 19) {
	
		s0 := s

		// 如果有符号表示正负,从第二位开始运算
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
		}

		// 字符串转整形运算
		n := 0
		for _, ch := range []byte(s) {
			ch -= '0'
			n = n*10 + int(ch)   // 按10进制运算
		}

		// 如果负数,转换完成之后,加上"-"号
		if s0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// int类型大整数,通过普通模式运算
	i64, err := ParseInt(s, 10, 0)
	if nerr, ok := err.(*NumError); ok {
		nerr.Func = fnAtoi
	}
	return int(i64), err
}
```

```go
// ParseInt参数:
// s - 待转换字符串
// base - 转换进制
// bitsize - 参数结果必须正确的整数类型,eg: int(0), int8(8), int16(16), int32(32), int64(64)
func ParseInt(s string, base int, bitSize int) (i int64, err error) {
	const fnParseInt = "ParseInt"


	// Pick off leading sign.
	s0 := s

	// 默认非负数
	neg := false

	// 如果要转换的整数类型,存在符号标示,从字符串第二位开始运算
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		// 负数
		neg = true
		s = s[1:]
	}

	// 处理无符号整数
	un, err = ParseUint(s, base, bitSize)
	if err != nil && err.(*NumError).Err != ErrRange {
		err.(*NumError).Func = fnParseInt
		err.(*NumError).Num = s0
		return 0, err
	}

	// 转成int64类型
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}
```