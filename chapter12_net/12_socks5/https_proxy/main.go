package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	urli := url.URL{}
	urlproxy, _ := urli.Parse("http://127.0.0.1:7890")
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
	}
	if resp, err := c.Get("https://www.google.com"); err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
	}
}
