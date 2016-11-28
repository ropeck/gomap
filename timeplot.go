package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/husobee/vestigo"
)

func drawday(td time.Time, r *http.Request) [][]int {
	return drawday_base(td, r, false, false)
}

func drawday_base(td time.Time, r *http.Request, reverse bool, cache bool) [][]int {
	d := NewDirections(r)
	var data [][]int
//	data = append(data, []string{"Time", "Leave", "Expected", "Delay"})
	t := td.Truncate(24 * time.Hour).Add(24 * time.Hour)
	for i := 0; i < 24; i++ {
		d.Directions(&t)
		data = append(data, []int{i * 60, i * 60,
			int(d.Duration.Seconds())/60 + i*60,
			int(d.DurationInTraffic.Seconds())/60 + i*60})
		t = t.Add(time.Hour)
	}
	return data
}

func drawdaylines(td time.Time, r *http.Request) []interface{} {
	midnight := td.Truncate(time.Hour * 24)
	daylist := []string{"Time"}
	data := make(map[time.Weekday]([][]int))
	td = midnight
	for i := 0; i < 7; i++ {
		day := td.Weekday()
		data[day] = drawday(td, r)  // this is where to pick out the data
		
		daylist = append(daylist, day.String())
		td = td.Add(time.Hour * 24)
	}
	ret := make([]interface{}, 25)
	ret = append(ret, daylist)

	ctx := appengine.NewContext(r)
	
	for h := 0; h < 24; h++ {
		row := make([]int, 8)

		for w := 0; w < 7; w++ {
			log.Infof(ctx, fmt.Sprintf("drawdaylines %v", time.Weekday(w)))
			log.Infof(ctx, fmt.Sprintf("drawdaylines %v", data[time.Weekday(w)][h+1]))
			stat := make([]int, 4)
			stat[0] = h
			for i, v := range data[time.Weekday(w)][h] {
				log.Infof(ctx, fmt.Sprintf(" %v", v))
				stat[i] = 1
			}
			row[w+1] = stat[3]
		}
		ret = append(ret, row)
	}
	return ret
}

func LocalNewDirections(w http.ResponseWriter, r *http.Request) *Directions {
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

func travel(w http.ResponseWriter, r *http.Request) {
	d := LocalNewDirections(w, r)
	t, _ := template.ParseFiles("base.html", "travel.html")
	t.ExecuteTemplate(w, "layout", d)
}

func traveldata(w http.ResponseWriter, r *http.Request) {
	tdarg := vestigo.Param(r, "date")
	i, _ := strconv.ParseInt(tdarg, 10, 64)
	data := drawdaylines(time.Unix(i/1000, 0), r)
	b, _ := json.Marshal(data)
	w.Write(b)
}

func arrivedata(w http.ResponseWriter, r *http.Request) {
	tdarg := vestigo.Param(r, "date")
	i, _ := strconv.ParseInt(tdarg, 10, 64)
	data := drawday(time.Unix(i/1000, 0), r)
	b, _ := json.Marshal(data)
	w.Write(b)
}

func init() {
	r := vestigo.NewRouter()
	r.Get("/arrivedata/:date", arrivedata)
	r.Get("/arrivedata", arrivedata)
	r.Get("/arrive", arrive)
	r.Get("/traveldata/:date", traveldata)
	r.Get("/travel", travel)
	r.Get("/", hello)
	http.Handle("/", r)
}
