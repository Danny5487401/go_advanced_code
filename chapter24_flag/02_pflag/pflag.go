package main

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"
)

// 支持长选项、短选项、默认值和使用文本，并将标志的值存储在指针中,
// 使用默认的CommandLine
var cliName = flag.StringP("name", "n", "nick", "Input Your Name")
var cliAge = flag.IntP("age", "a", 22, "Input Your Age")
var cliGender = flag.StringP("gender", "g", "male", "Input Your Gender")
var cliOK = flag.BoolP("ok", "o", false, "Input Are You OK")
var cliDes = flag.StringP("des-detail", "d", "", "Input Description")
var cliOldFlag = flag.StringP("badflag", "b", "just for test", "Input badflag")

func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return flag.NormalizedName(name)
}

func main() {
	// 设置标准化参数名称的函数
	// 如果我们创建了名称为 --des-detail 的参数，但是用户却在传参时写成了 --des_detail 或 --des.detail 会怎么样？
	// 默认情况下程序会报错退出，但是我们可以通过 pflag 提供的 SetNormalizeFunc 功能轻松的解决这个问题
	flag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)

	// 为 age 参数设置 NoOptDefVal 默认值，通过简便的方式为参数设置默认值之外的值
	flag.Lookup("age").NoOptDefVal = "25"

	// 把 badflag 参数标记为即将废弃的，请用户使用 des-detail 参数
	flag.CommandLine.MarkDeprecated("badflag", "please use --des-detail instead")

	// 把 badflag 参数的 shorthand 简写形式标记为即将废弃的，请用户使用 des-detail 的 shorthand 参数
	flag.CommandLine.MarkShorthandDeprecated("badflag", "please use -d instead")

	// 在帮助文档中隐藏参数 gender
	flag.CommandLine.MarkHidden("gender")

	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()

	name, _ := flag.CommandLine.GetString("name")
	fmt.Println("name=", name)
	fmt.Println("age=", *cliAge)
	fmt.Println("gender=", *cliGender)
	fmt.Println("ok=", *cliOK)
	fmt.Println("des=", *cliDes)

	fmt.Printf("argument number is: %v\n", flag.NArg())
	fmt.Printf("argument list is: %v\n", flag.Args())
	fmt.Printf("the first argument is: %v\n", flag.Arg(0))
}
