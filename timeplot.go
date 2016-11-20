package gomap

import (
	"github.com/ropeck/directions"
	"net/http"
	"html/template"
)
var t, err = template.ParseFiles("base.html", "index.html")

func hello(w http.ResponseWriter, r *http.Request) {
     d := directions.NewDirections(r)
     d.Directions()
     t.ExecuteTemplate(w, "layout", d)
}

func init() {
  http.HandleFunc("/hello", hello)
}
