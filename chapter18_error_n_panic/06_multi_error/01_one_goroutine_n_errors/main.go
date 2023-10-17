package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"

	"github.com/hashicorp/go-multierror"
)

var data = []byte(
	`a,b,c
	 foo
	 1,2,3
	 ,",
`)

// 单个 goroutine，多个错误

func main() {
	// 1. 最简单的实现，多次打印
	easyErr()

	// 2. HashiCorp 的 go-multierror实现：一次打印
	multiErrPkg()
}

func easyErr() {
	reader := csv.NewReader(bytes.NewBuffer(data))
	for {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf(err.Error())
		}
	}
}

func multiErrPkg() {
	var errs error
	reader := csv.NewReader(bytes.NewBuffer(data))
	for {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				break
			}
			errs = multierror.Append(errs, err)
		}
	}
	if errs != nil {
		log.Printf(errs.Error())
	}
}
