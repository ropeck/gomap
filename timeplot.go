package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/kr/pretty"
)

func drawday(td time.Time, r *http.Request) [][]string {
	return drawday_base(td, r, false, false)
}

func drawday_base(td time.Time, r *http.Request, reverse bool, cache bool) [][]string {
	d := NewDirections(r)
	d.Directions()
	//data := [][]string{{"Leave", "Expected", "Delay"}}
	//data = append(data, []string{"0", "1", "2"})
	data := pretty.Sprint(d)
	return [][]string{{string(data)}}
}

func LocalNewDirections(w http.ResponseWriter, r *http.Request) (*Directions) {
	d := NewDirections(r)
	d.Directions()
	http.SetCookie(w, d.Ocookie)
	http.SetCookie(w, d.Dcookie)
	return d
}

func hello(w http.ResponseWriter, r *http.Request) {
        d := LocalNewDirections(w, r)
	t, _ := template.ParseFiles("base.html", "index.html")
	t.ExecuteTemplate(w, "layout", d)
}

func arrive(w http.ResponseWriter, r *http.Request) {
        d := LocalNewDirections(w, r)
	t, _ := template.ParseFiles("base.html", "arrive.html")
	t.ExecuteTemplate(w, "layout", d)
}

func arrivedata(w http.ResponseWriter, r *http.Request) {
        d := LocalNewDirections(w, r)
	//td := time.Now()
	//data := drawday(td, r)
	b, _ := json.Marshal(d.Resp)
	w.Write(b)
}

func init() {
	http.HandleFunc("/arrivedata", arrivedata)
	http.HandleFunc("/arrive", arrive)
	http.HandleFunc("/", hello)
}
