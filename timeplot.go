package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/husobee/vestigo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// use mutation function on the directions to pass options for cache and reverse

func drawday(td time.Time, r *http.Request) [24][4]int {
	return drawday_base(td, r, false, false)
}

func drawday_base(td time.Time, r *http.Request, reverse bool, cache bool) [24][4]int {
	data := [24][4]int{}

	t := td.Truncate(24 * time.Hour).Add(24 * time.Hour)

	ch := make(chan [4]int)
	for i := 0; i < 24; i++ {
		t := t.Add(time.Hour)
		go func(i int, t time.Time, ch chan [4]int) {
			d := NewDirections(r)
			d.Directions(&t)
			a := [4]int{i * 60, i * 60,
				int(d.Duration.Seconds())/60 + i*60,
				int(d.DurationInTraffic.Seconds())/60 + i*60}
			ch <- a
		}(i, t, ch)
	}

	// read back all 24 hours and assign them to the slots
	ctx := appengine.NewContext(r)
	for i := 0; i < 24; i++ {
		d := <-ch
		h := d[0] / 60
		log.Infof(ctx, fmt.Sprintf("hour %d %v", h, d))
		data[h] = d
	}
	return data
}

func drawdaylines(td time.Time, r *http.Request) []interface{} {
	ctx := appengine.NewContext(r)
	midnight := td.Truncate(time.Hour * 24)
	daylist := []string{"Time"}
	data := make(map[time.Weekday]([24][4]int))
	td = midnight
	for i := 0; i < 7; i++ {
		day := td.Weekday()
		data[day] = drawday(td, r) // this is where to pick out the data
		daylist = append(daylist, day.String())
		td = td.Add(time.Hour * 24)
	}
	ret := make([]interface{}, 0)
	ret = append(ret, daylist)

	for h := 0; h < 24; h++ {
		var row [8]int

		for w := 0; w < 7; w++ {

			//			log.Infof(ctx, fmt.Sprintf("drawdaylines %v", h))
			d := data[time.Weekday(w)][h]
			row[w+1] = d[3] - d[2]
			//			log.Infof(ctx, fmt.Sprintf(" %v", data[time.Weekday(w)]))
		}
		row[0] = h
		ret = append(ret, row)
	}
	log.Infof(ctx, fmt.Sprintf("travel %v", ret))

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
	var data []interface{}
	data = append(data, []string{"Time", "Leave", "Expected", "Delay"})

	tdarg := vestigo.Param(r, "date")
	i, _ := strconv.ParseInt(tdarg, 10, 64)

	for _, v := range drawday(time.Unix(i/1000, 0), r) {
		data = append(data, v)
	}
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
