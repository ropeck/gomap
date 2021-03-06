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
	ctx := appengine.NewContext(r)
	data := [24][4]int{}

	l, _ := time.LoadLocation("US/Pacific")
	yy, mm, dd := td.In(l).Date()
	t := time.Date(yy, mm, dd, 0, 0, 0, 0, l)

	batch := 1
	ch := make(chan [4]int, batch)
	for i := 0; i < 24; i++ {
		go func(i int, t time.Time, ch chan [4]int) {
			d := NewDirections(r)
			d.Directions(&t)
			traf := int(d.DurationInTraffic.Seconds()) / 60
			a := [4]int{i * 60, i * 60,
				int(d.Duration.Seconds())/60 + i*60,
				traf + i*60}
			ch <- a
		}(i, t, ch)
		t = t.Add(time.Hour)
	}

	// read back all 24 hours and assign them to the slots
	for i := 0; i < 24; i++ {
		d := <-ch
		h := d[0] / 60
		log.Infof(ctx, fmt.Sprintf("hour %d %v", h, d))
		data[h] = d
	}
	return data
}

func drawdaylines(td time.Time, days string, r *http.Request) []interface{} {
	ctx := appengine.NewContext(r)
	daylist := []string{"Time"}
	data := make(map[time.Weekday]([24][4]int))
	td = previous_monday(td)
	max_days, _ := strconv.Atoi(days)
	for i := 0; i < max_days; i++ {
		day := td.Weekday()
		log.Infof(ctx, fmt.Sprintf("draw %s %s", td, day))
		data[day] = drawday(td, r)
		daylist = append(daylist, day.String())
		td = td.Add(time.Hour * 24)
	}
	ret := make([]interface{}, 0)
	ret = append(ret, daylist)

	for h := 0; h < 24; h++ {
		row := make([]interface{},max_days+1)
		for w := 0; w < max_days; w++ {
			d := data[time.Weekday(w)][h]
			row[w+1] = d[3] - d[2]
		}
		row[0] = fmt.Sprintf("%d:00", h)
		ret = append(ret, row)
		log.Infof(ctx, fmt.Sprintf("  %v", row))
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
	t, _ := template.ParseFiles("template/base.html",
		"template/index.html")
	t.ExecuteTemplate(w, "layout", d)
}

func arrive(w http.ResponseWriter, r *http.Request) {
	d := LocalNewDirections(w, r)
	t, _ := template.ParseFiles("template/base.html",
		"template/arrive.html")
	t.ExecuteTemplate(w, "layout", d)
}

func daily(w http.ResponseWriter, r *http.Request) {
	d := LocalNewDirections(w, r)
	t, _ := template.ParseFiles("template/base.html",
		"template/daily.html")
	t.ExecuteTemplate(w, "layout", d)
}

func to_hms(t time.Time) string {
	//	t = t.Truncate(time.Hour)
	return fmt.Sprintf("%d:%.02d", t.Hour(), t.Minute())
}

func previous_monday(t time.Time) time.Time {
	l, _ := time.LoadLocation("US/Pacific")
	wd := int(t.Weekday())
	yy, mm, dd := t.Add(time.Hour * -24 * time.Duration(wd)).Date()
	st := time.Date(yy, mm, dd, 0, 0, 0, 0, l)
	return st
}

func dailydata(w http.ResponseWriter, r *http.Request) {
	var data []interface{}

	data = append(data, []interface{}{"Time", "Delay"})

	tdarg := vestigo.Param(r, "date")
	i, _ := strconv.ParseInt(tdarg, 10, 64)

	// adjust start time to Midnight on Monday
	st := time.Unix(i/1000, 0)
	yy, mm, dd := st.Date()
	l, _ := time.LoadLocation("US/Pacific")
	st = time.Date(yy, mm, dd, 0, 0, 0, 0, l)
	for _, v := range drawday(st, r) {
		data = append(data, [2]interface{}{to_hms(st), v[3] - v[2]})
		st = st.Add(time.Hour)
	}
	b, _ := json.Marshal(data)
	w.Write(b)

}

func weekly(w http.ResponseWriter, r *http.Request) {
	d := LocalNewDirections(w, r)
	t, err := template.ParseFiles("template/base.html",
		"template/weekly.html")
	if err != nil {
		ctx := appengine.NewContext(r)
		log.Infof(ctx, err.Error())
		return
	}
	t.ExecuteTemplate(w, "layout", d)
}

func weeklydata(w http.ResponseWriter, r *http.Request) {
	tdarg := vestigo.Param(r, "date")
	days := vestigo.Param(r, "days")
	i, _ := strconv.ParseInt(tdarg, 10, 64)
	data := drawdaylines(time.Unix(i/1000, 0), days, r)
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

type WhenInfo struct {
	Distance          string
	Duration          string
	DurationInTraffic string
	Text              string
	Apikey            string
}

func whentogo(w http.ResponseWriter, r *http.Request) {
	wh := new(WhenInfo)

	d := LocalNewDirections(w, r)
	wh.Distance = d.Distance.HumanReadable
	wh.Duration = d.Duration.String()
	wh.DurationInTraffic = d.DurationInTraffic.String()
	wh.Text = ""
	wh.Apikey = d.Apikey
	l, _ := time.LoadLocation("US/Pacific")
	td := time.Now().In(l).Truncate(10 * time.Minute)
	for h := 0; h < 24; h++ {
		d.Directions(&td)
		td = td.Add(time.Minute * 10)
		delay := int((d.DurationInTraffic - d.Duration).Seconds() / 60)
		ar := td.Add(d.DurationInTraffic)
		wh.Text = wh.Text + fmt.Sprintf("%s %s (%d m %+d)\n",
			to_hms(td), to_hms(ar),
			int(d.DurationInTraffic.Seconds()/60), delay)
	}

	t, _ := template.ParseFiles("template/base.html",
		"template/whentogo.html")
	t.ExecuteTemplate(w, "layout", wh)
}

func init() {
	r := vestigo.NewRouter()
	r.Get("/whentogo", whentogo)
	r.Get("/dailydata/:date", dailydata)
	r.Get("/daily", daily)
	r.Get("/arrivedata/:date", arrivedata)
	r.Get("/arrive", arrive)
	r.Get("/weeklydata/:date/:days", weeklydata)
	r.Get("/weekly", weekly)
	r.Get("/", hello)
	http.Handle("/", r)
}
