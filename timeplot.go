package gomap

import (
	"directions"
	"net/http"
	"html/template"
)
var t, err = template.ParseFiles("base.html", "index.html")

func hello(w http.ResponseWriter, r *http.Request) {
     d := directions.NewDirections(r)
     t.ExecuteTemplate(w, "layout", d)
}

func init() {
  http.HandleFunc("/hello", hello)
}
