package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter10_function/02_advanced_function/04_mapReduce/mapReduce"
	"github.com/Danny5487401/go_advanced_code/chapter10_function/02_advanced_function/04_mapReduce/product_srv"
	"github.com/Danny5487401/go_advanced_code/chapter10_function/02_advanced_function/04_mapReduce/user_srv"
	"sync"
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
	// 并行执行
	_ = mapReduce.Finish(func() (err error) {
		userInfo, err = user_srv.GetUser()
		return err
	}, func() (err error) {
		productList, err = product_srv.GetProductList()
		return err
	})
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
