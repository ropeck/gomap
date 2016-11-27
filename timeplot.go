package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"
	"strconv"
	"github.com/husobee/vestigo"
)

func drawday(td time.Time, r *http.Request) []interface{} {
	return drawday_base(td, r, false, false)
}

func drawday_base(td time.Time, r *http.Request, reverse bool, cache bool) []interface{} {
	d := NewDirections(r)
	var data []interface{}
	data = append(data, []string{"Time", "Leave", "Expected", "Delay"})
	t := td.Truncate(24*time.Hour).Add(24*time.Hour)
	for i := 0; i < 24; i++ {
		d.Directions(&t)
		data = append(data,[]interface{}{i*60, i*60,
			int(d.Duration.Seconds())/60+i*60,
			int(d.DurationInTraffic.Seconds())/60+i*60,})
		t = t.Add(time.Hour)
	}
	return data
}

// func drawdaylines(td time.Time) []interface{} {
//   midnight = td.replace(td.year,td.month,td.day,8,0,0,0)
//   memkey = 'lines:'+str(midnight)

//   daylist = []
//   data = {}
//   td = midnight
//   for i in range(7):
//     day = td.strftime("%a")
//     data[day] = drawday(td, reverse=True, cache=False)
//     daylist.append(day)
//     td = td + timedelta(days=1)

//   ret = [['Time'] + daylist]
//   try:
//     r = 0
//     for m in drawday(td)[1:]:
//       r = r + 1
//       row = [m[0]]
//       dn = 1
//       for day in daylist:
//         row.append(data[day][r][1])
//         dn = dn + 1
//         ret.append(row)
//   except runtime.DeadlineExceededError:
//     print "error: deadline exceeded in lookups"
//     # the lookups timed out, so the result is probably incomplete
//   return ret
// }

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
