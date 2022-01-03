package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// 1. 解析http
	Url := "https://root:123456@www.baidu.com:0000/api/v1/login?name=xiaoming&name=xiaoqing&age=24&age1=23#fffffff"
	//Parse函数解析Url为一个URL结构体，Url可以是绝对地址，也可以是相对地址
	//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	//	type URL struct {
	//    Scheme   string
	//    Opaque   string    // 编码后的不透明数据
	//    User     *Userinfo // 用户名和密码信息
	//    Host     string    // host或host:port
	//    Path     string
	//    RawQuery string // 编码后的查询字符串，没有'?'
	//    Fragment string // 引用的片段（文档位置），没有'#'
	//}
	u, err := url.Parse(Url)
	if err == nil {
		fmt.Printf("解析后的数据%+v/n", u)
	}
	//ParseRequestURI函数解析Url为一个URL结构体，本函数会假设Url是在一个HTTP请求里，
	//	因此会假设该参数是一个绝对URL或者绝对路径，并会假设该URL没有#fragment后缀
	u1, err := url.ParseRequestURI(Url)
	if err == nil {
		fmt.Println(u1)
	}

	//得到Scheme
	fmt.Println(u.Scheme)

	//User 包含认证信息 username and password
	user := u.User
	username := user.Username()
	password, pdSet := user.Password()
	fmt.Println(username, password, pdSet)

	////Host 包括主机名和端口信息，如过有端口，用strings.Split() 从Host中手动提取端口
	host := u.Host
	ho := strings.Split(host, ":")
	fmt.Println("主机名：", ho[0], "端口号：", ho[1])

	//获取path
	path := u.Path
	fmt.Println(path)

	//获取参数 将查询参数解析为一个map。 map以查询字符串为键，对应值字符串切片为值。
	fmt.Println(u.RawQuery)
	urlParam := u.RawQuery
	fmt.Println("urlParam:", urlParam)
	m, err := url.ParseQuery(urlParam)
	if err == nil {
		fmt.Println(m)
		for k, v := range m {
			fmt.Println(k, v)
		}
	}

	//得到查询片段信息
	fmt.Println(u.Fragment)

	////生成参数形如：name=xiaoming&name=xiaoqing&age=24&age1=23
	v := url.Values{}
	//添加一个k-v值
	v.Set("name", "xiaoming")
	v.Add("name", "xiaoqing")
	v.Set("Age", "23")
	fmt.Println(v)
	fmt.Println(v.Get("name"))
	v.Del("name")
	fmt.Println(v)
	//把map编码成name=xiaoming&name=xiaoqing&age=24&age1=23的形式
	v.Set("name", "xiaoming")
	v.Add("name", "xiaoqing")
	fmt.Println(v.Encode())

	http.HandleFunc("/go", gso)
	http.ListenAndServe(":8080", nil)
	//输入URL:  http://localhost:8080/go?a=111&b=456
}

func gso(w http.ResponseWriter, r *http.Request) {
	r1 := r.URL.String()
	r2 := r.RequestURI
	fmt.Println(r1)                         //go?a=111&b=456
	fmt.Println(r2)                         //go?a=111&b=456
	u, _ := url.Parse(r1)                   //将string解析成*URL格式
	fmt.Println(u)                          // go?a=111&b=456
	fmt.Println(u.RawQuery)                 //编码后的查询字符串，没有'?' a=111&b=456
	values, _ := url.ParseQuery(u.RawQuery) //返回Values类型的字典
	fmt.Println(values)                     // map[a:[111] b:[456]]
	fmt.Println(values["a"])                //[111]
	fmt.Println(values.Get("b"))            //456
}
