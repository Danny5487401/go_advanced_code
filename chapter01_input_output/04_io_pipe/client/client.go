package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	//modifyBefore()
	modifyAfter()
}

// 修改后
func modifyAfter() {
	cli := http.Client{}

	msg := struct {
		Name, Addr string
		Price      float64
	}{
		Name:  "hello",
		Addr:  "beijing",
		Price: 123.56,
	}
	r, w := io.Pipe()
	go func() {
		// 注意这边的逻辑！！
		//  创建goroutine给writer
		defer func() {
			time.Sleep(time.Second * 2)
			log.Println("encode完成")
			// 只有这里关闭了，Post方法才会返回
			w.Close()
		}()
		log.Println("管道准备输出")
		// 只有Post开始读取数据，这里才开始encode，并传输
		err := json.NewEncoder(w).Encode(msg)
		log.Println("管道输出数据完毕")
		if err != nil {
			log.Fatalln("encode json failed:", err)
		}
	}()
	log.Println("开始从管道读取数据")
	resp, err := cli.Post("http://localhost:9999/json", "application/json", r)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("POST传输完成")

	body := resp.Body
	defer body.Close()

	if bodyBytes, err := ioutil.ReadAll(body); err == nil {
		log.Println("response:", string(bodyBytes))
	} else {
		log.Fatalln(err)
	}
}

// 修改前
func modifyBefore() {
	cli := http.Client{}

	msg := struct {
		Name, Addr string
		Price      float64
	}{
		Name:  "hello",
		Addr:  "beijing",
		Price: 123.56,
	}
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(msg)
	resp, err := cli.Post("http://localhost:9999/json", "application/json", buf)

	if err != nil {
		log.Fatalln(err)
	}

	body := resp.Body
	defer body.Close()

	if bodyBytes, err := io.ReadAll(body); err == nil {
		log.Println("response:", string(bodyBytes))
	} else {
		log.Fatalln(err)
	}
}
