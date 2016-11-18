package gomap

import (
	"directions"
	"net/http"
	"html/template"
)
var t, err = template.ParseFiles("base.html", "index.html")

func hello(w http.ResponseWriter, r *http.Request) {
	var m = make(map[string]string)
	m["apikey"] = directions.Apikey(r)
	t.ExecuteTemplate(w, "layout", m)
}

func init() {
  http.HandleFunc("/hello", hello)
}
