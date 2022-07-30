package main

import (
	"errors"
	"log"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {

	var singleSetCache singleflight.Group

	getAndSetCache := func(requestID int, cacheKey string) (string, error) {

		log.Printf("request %v start to get and set cache...", requestID)

		retChan := singleSetCache.DoChan(cacheKey, func() (ret interface{}, err error) {

			log.Printf("request %v is setting cache...", requestID)

			time.Sleep(3 * time.Second)

			log.Printf("request %v set cache success!", requestID)

			return "VALUE", nil

		})

		var ret singleflight.Result

		timeout := time.After(5 * time.Second)

		select { //加入了超时机制

		case <-timeout:

			log.Printf("time out!")

			return "", errors.New("time out")

		case ret = <-retChan: //从chan中取出结果

			return ret.Val.(string), ret.Err

		}

		return "", nil

	}

	cacheKey := "cacheKey"

	for i := 1; i < 10; i++ {

		go func(requestID int) {

			value, _ := getAndSetCache(requestID, cacheKey)

			log.Printf("request %v get value: %v", requestID, value)

		}(i)

	}

	time.Sleep(20 * time.Second)

}
