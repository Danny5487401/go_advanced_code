package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 爬虫接口
type ICrawler interface {
	GetUserList() ([]*User, error)
}

// 具体实现
type MyCrawler struct {
	url string
}

func (c *MyCrawler) GetUserList() ([]*User, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userList []*User
	err = json.Unmarshal(data, &userList)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func GetAndPrintUsers(crawler ICrawler) {
	users, err := crawler.GetUserList()
	if err != nil {
		return
	}

	for _, u := range users {
		fmt.Println(u)
	}
}

type User struct {
	Name string
	Age  int
}
