package main

import (
	"encoding/xml"
	"html/template"
	"net/http"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/unrolled/render"
)

// render demo
func main() {
	r := render.New(render.Options{
		Layout:        "layout",
		IsDevelopment: true,
		Funcs: []template.FuncMap{
			sprig.FuncMap(),
		},
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Welcome, visit sub pages now."))
	})
	ob := struct {
		XMLName xml.Name `json:"-" xml:"xmlname"`
		Name    string   `json:"name" xml:"name"`
		Items   []string `json:"items" xml:"items"`
	}{
		Name:  "danny",
		Items: []string{"demo", "app"},
	}
	mux.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, 200, ob)
	})
	mux.HandleFunc("/xml", func(w http.ResponseWriter, req *http.Request) {
		r.XML(w, 200, ob)
	})
	mux.HandleFunc("/html", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, 200, "index", ob)
	})
	http.ListenAndServe("127.0.0.1:3000", mux)
}
