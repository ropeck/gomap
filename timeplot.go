package main

import (
	"github.com/ropeck/directions"
	"net/http"
	"html/template"
)

func hello(w http.ResponseWriter, r *http.Request) {
     t, _ := template.ParseFiles("base.html", "index.html")
     d := directions.NewDirections(r)
     d.Directions()
     t.ExecuteTemplate(w, "layout", d)
}

func arrive(w http.ResponseWriter, r *http.Request) {
     t, _ := template.ParseFiles("base.html", "arrive.html")
     d := directions.NewDirections(r)
     d.Directions()
     t.ExecuteTemplate(w, "layout", d)
}

func init() {
  http.HandleFunc("/arrive", arrive)
  http.HandleFunc("/", hello)
}
