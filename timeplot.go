package main

import (
	"github.com/ropeck/directions"
	"net/http"
	"encoding/json"
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

func arrivedata(w http.ResponseWriter, r *http.Request) {
     d := directions.NewDirections(r)
     d.Directions()
     data := [][]string{{"Time", "Leave", "Expected", "Delay"}}
     // for h in drawday(td):
     //   data.append( stuff from h )
     b, _ := json.Marshal(data)
     w.Write(b)
}

func init() {
  http.HandleFunc("/arrivedata", arrivedata)
  http.HandleFunc("/arrive", arrive)
  http.HandleFunc("/", hello)
}
