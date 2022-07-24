package main

import (
	"fmt"
	"log"
	"os"
)

//写出数据
func main() {

	//step1：打开文件
	fileName := "chapter01_input_output/files/dannyWrite.txt"
	//file,err := os.Open(fileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	//step3：关闭文件
	defer file.Close()

	//step2：写出数据
	//bs :=[]byte{65,66,67,68,69,70}//A,B,C,D,E,F
	bs := []byte{97, 98, 99, 100} //a,b,c,d
	//n,err := file.Write(bs)
	n, err := file.Write(bs[:2])
	fmt.Println(n)
	HandleErr(err)
	file.WriteString("\n")

	//直接写出字符串
	n, err = file.WriteString("HelloWorld")
	fmt.Println(n)
	HandleErr(err)

	file.WriteString("\n")
	n, err = file.Write([]byte("today"))
	fmt.Println(n)
	HandleErr(err)

}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
