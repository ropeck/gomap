package main

import (
	"github.com/ropeck/directions"
	"net/http"
	"encoding/json"
	"html/template"
	"time"
)

func drawday(td time.Time) ([][]string) {
  return drawday_base(td, false, false)
}

func drawday_base(td time.Time, reverse bool, cache bool) ([][]string){
  data := [][]string{{"time.Time", "Leave", "Expected", "Delay"}}
  data = append(data, []string{"0","1","2"})
  return data
}


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
     td := time.Now()
     data := drawday(td)
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
