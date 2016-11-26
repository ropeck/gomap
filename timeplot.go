package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"
	"fmt"
)

func drawday(td time.Time, r *http.Request) [][]string {
	return drawday_base(td, r, false, false)
}

func drawday_base(td time.Time, r *http.Request, reverse bool, cache bool) [][]string {
	d := NewDirections(r)
	data := [][]string{{"Leave", "Expected", "Delay"}}
	//data = append(data, []string{"0", "1", "2"})
	// loop over hours of the day, collecting the directions result for
	// for each hour
	y,m,day := td.Add(time.Hour*24).Date()
	t := time.Date(y,m,day,0,0,0,0,td.Location())
	for i := 0; i < 24; i++ {
		d.Directions(&t)
		data = append(data, []string{t.String(),
			fmt.Sprintf("%v",d.Duration.Seconds()),
			fmt.Sprintf("%v",d.DurationInTraffic.Seconds())})
		t = t.Add(time.Hour)
	}
	return data
}

func LocalNewDirections(w http.ResponseWriter, r *http.Request) (*Directions) {
	d := NewDirections(r)
	d.DirectionsNow()
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
	//td := time.Now()
	data := drawday(time.Now(), r)
	b, _ := json.Marshal(data)
	w.Write(b)
}

func init() {
	http.HandleFunc("/arrivedata", arrivedata)
	http.HandleFunc("/arrive", arrive)
	http.HandleFunc("/", hello)
}
