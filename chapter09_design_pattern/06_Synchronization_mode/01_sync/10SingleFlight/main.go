package main

import (
	"errors"
	"log"
	"sync"

	"golang.org/x/sync/singleflight"
)

/*
作用
	在一个服务中抑制对下游的多次重复请求。一个比较常见的使用场景是：我们在使用 Redis 对数据库中的数据进行缓存，发生缓存击穿时，大量的流量都会打到数据库上进而影响服务的尾延时
缓存击穿
	缓存在某个时间点过期的时候，恰好在这个时间点对这个Key有大量的并发请求过来，这些请求发现缓存过期一般都会从后端DB加载数据并回设到缓存，这个时候大并发的请求可能会瞬间把后端DB压垮。
*/

var g singleflight.Group

//获取数据
func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		//模拟从db中获取数据
		v, err, _ := g.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
			//set cache
		})
		if err != nil {
			log.Println(err)
			return "", err
		}
		//TODO: set cache
		data = v.(string)
	} else if err != nil {
		return "", err
	}
	return data, nil
}

var errorNotExist = errors.New("not exist")

//模拟从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {
	return "", errorNotExist
}

//模拟从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	return "data", nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(10)

	//模拟10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}
