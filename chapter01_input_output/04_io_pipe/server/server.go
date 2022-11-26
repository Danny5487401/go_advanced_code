package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	http.HandleFunc("/json", handleJson)

	http.ListenAndServe(":9999", nil)
}

func handleJson(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body := req.Body
		defer body.Close()
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			log.Println(err)
			resp.Write([]byte(err.Error()))
			return
		}
		j := map[string]interface{}{}
		if err := json.Unmarshal(bodyBytes, &j); err != nil {
			log.Println(err)
			resp.Write([]byte(err.Error()))
			return
		}
		resp.Write(bodyBytes)
	} else {
		resp.Write([]byte("请使用post方法!"))
	}
}
