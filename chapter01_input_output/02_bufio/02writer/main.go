package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
	bufio:高效io读写
		buffer缓存
		io：input/output

	将io包下的Reader，Write对象进行包装，带缓存的包装，提高读写的效率

		func (b *Writer) Write(p []byte) (nn int, err error)
		func (b *Writer) WriteByte(c byte) error
		func (b *Writer) WriteRune(r rune) (size int, err error)
		func (b *Writer) WriteString(s string_test) (int, error)

*/

func main() {

	fileName := "chapter01_input_output/files/dannyBufioWrite.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	w1 := bufio.NewWriter(file)

	for i := 1; i <= 1000; i++ {

		w1.WriteString(fmt.Sprintf("%d:hello\n", i))
	}
	//刷新缓冲区
	w1.Flush()
}
