package gomap

import (
	"net/http"
	"html/template"
)
var t, err = template.ParseFiles("base.html")

func hello(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "layout", "")
}

func init() {
  http.HandleFunc("/hello", hello)
}
