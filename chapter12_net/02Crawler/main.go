package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	// \d是数字
	reQQEmail = `(\d+)@qq.com`
)

func GetEmail() {
	// 1.去网站拿数据
	resp, err := http.Get("https://tieba.baidu.com/p/6051076813?red_tag=1573533731")
	HandleError(err, "http.Get url")
	if resp != nil {
		defer resp.Body.Close()
	}
	/*
		解释：在这里，你同样需要检查 res 的值是否为 nil ，这是 http.Get 中的一个警告。\
		通常情况下，出错的时候，返回的内容应为空并且错误会被返回，可当你获得的是一个重定向 error 时， res 的值并不会为 nil ，但其又会将错误返回。
		上面的代码保证了无论如何 Body 都会被关闭，如果你没有打算使用其中的数据，那么你还需要丢弃已经接收的数据
	*/
	// 2.读取页面内容
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")

	// 字节转字符串
	pageStr := string(pageBytes)
	//fmt.Println(pageStr)

	// 3.过滤数据，过滤qq邮箱
	re := regexp.MustCompile(reQQEmail)
	// -1代表取全部
	results := re.FindAllStringSubmatch(pageStr, -1)
	//fmt.Println(results)
	// 遍历结果
	for i, result := range results {
		fmt.Println("第", i, "条：", "email:", result[0], "\t", "qq:", result[1])

	}

}

func main() {
	GetEmail()
}

// 处理异常
func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}
