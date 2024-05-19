package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	str := "GO爱好者"
	fmt.Printf("%q \n", []rune(str)) // ['G' 'O' '爱' '好' '者']
	fmt.Printf("%x \n", []rune(str)) // [47 4f 7231 597d 8005]
	fmt.Printf("%x \n", []byte(str)) // 474fe788b1e5a5bde88085

	c, size := utf8.DecodeRune([]byte(str)) // 函数解码p开始位置的第一个utf-8编码的码值，返回该码值和编码的字节数。
	fmt.Printf("%c %d\n", c, size)

	c2, size2 := utf8.DecodeRuneInString(str) // 函数类似DecodeRune但输入参数是字符串。
	fmt.Printf("%c %d\n", c2, size2)

	count := utf8.RuneCountInString("王思聪")
	fmt.Println(count) //3

	single := '\u0015'
	fmt.Println(unicode.IsControl(single))
	single = '\ufe35'
	fmt.Println(unicode.IsControl(single))

	digit := '1'
	fmt.Printf("1 IsDigit 阿拉伯数字字符  %v \n", unicode.IsDigit(digit))
	fmt.Printf("1 IsNumber 是否数字字符，比如罗马数字 Ⅷ 也是数字字符 %v \n", unicode.IsNumber(digit))

	letter := 'Ⅷ'
	fmt.Printf("Ⅷ IsDigit 阿拉伯数字字符  %v \n", unicode.IsDigit(letter))
	fmt.Printf("Ⅷ IsNumber 是否数字字符，比如罗马数字 Ⅷ 也是数字字符 %v \n", unicode.IsNumber(letter))

	han := '你'
	fmt.Printf("你 是否汉字 %v \n", unicode.Is(unicode.Han, han))
	fmt.Printf("你 是否空格 %v \n", unicode.In(han, unicode.Gujarati, unicode.White_Space))

}

// 一个string类型的值会由若干个 Unicode 字符组成，每个 Unicode 字符都可以由一个rune类型的值来承载。
// 这些字符在底层都会被转换为 UTF-8 编码值，而这些 UTF-8 编码值又会以字节序列的形式表达和存储。
// 因此，一个string类型的值在底层就是一个能够表达若干个 UTF-8 编码值的字节序列。

// %q      单引号围绕的字符字面值，由Go语法安全地转义
// %x      十六进制表示
