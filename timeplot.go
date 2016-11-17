package gomap

import (
	"html/template"
)

func hello(w http.ResponseWriter, r *http.Request) {
  t, _ = template.ParseFiles("templates/base.html", "templates/index.html")
  t.ExecuteTemplate(w, "layout", "")
}

