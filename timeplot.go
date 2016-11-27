package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"
	"strconv"
	"github.com/husobee/vestigo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func drawday(td time.Time, r *http.Request) []interface{} {
	return drawday_base(td, r, false, false)
}

func drawday_base(td time.Time, r *http.Request, reverse bool, cache bool) []interface{} {
	d := NewDirections(r)
	var data []interface{}
	data = append(data, []string{"Time", "Leave", "Expected", "Delay"})
	//data = append(data, []string{"0", "1", "2"})
	// loop over hours of the day, collecting the directions result for
	// for each hour
	ctx := appengine.NewContext(d.r)
	t := td.Truncate(24*time.Hour)
	for i := 0; i < 24; i++ {
		log.Infof(ctx, t.String())
		d.Directions(&t)
		data = append(data,[]interface{}{i*60, i*60,
			int(d.Duration.Seconds())/60+i*60,
			int(d.DurationInTraffic.Seconds())/60+i*60,})
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
	r.Get("/", hello)
	http.Handle("/",r)
}
