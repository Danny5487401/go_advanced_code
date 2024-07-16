package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter10_function/02_advanced_function/03_go_zero_map_reduce/go-zero-folk/mr"
	"sync"

	"github.com/Danny5487401/go_advanced_code/chapter10_function/02_advanced_function/03_go_zero_map_reduce/product_srv"
	"github.com/Danny5487401/go_advanced_code/chapter10_function/02_advanced_function/03_go_zero_map_reduce/user_srv"
)

func main() {
	// 1. 未封装前
	//modifyBefore()

	// 2. 封装后
	modifyAfter()
}

func modifyAfter() {
	var userInfo *user_srv.User
	var productList []product_srv.Product
	// 并行执行，worker数目为函数数目
	func1 := func() (err error) {
		userInfo, err = user_srv.GetUser()
		return err
	}
	func2 := func() (err error) {
		productList, err = product_srv.GetProductList()
		return err
	}
	_ = mr.Finish(func1, func2)
	fmt.Printf("用户信息:%+v\n", userInfo)
	fmt.Printf("商品信息:%+v\n", productList)
}

func modifyBefore() {
	var wg sync.WaitGroup
	wg.Add(2)

	var userInfo *user_srv.User
	var productList []product_srv.Product

	go func() {
		defer wg.Done()
		userInfo, _ = user_srv.GetUser()
	}()

	go func() {
		defer wg.Done()
		productList, _ = product_srv.GetProductList()
	}()
	wg.Wait()
	fmt.Printf("用户信息:%+v\n", userInfo)
	fmt.Printf("商品信息:%+v\n", productList)
}
