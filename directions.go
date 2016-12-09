package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
	"googlemaps.github.io/maps"
)

type Config struct {
	Name  string
	Value string
}

type Directions struct {
	Origin            string
	Destination       string
	Client            *maps.Client
	Apikey            string
	r                 *http.Request
	Resp              string
	Leg               *maps.Leg
	Dir               *maps.Route
	Steps             []*Step
	Duration          time.Duration
	DurationInTraffic time.Duration
	Distance          maps.Distance
	Dcookie           *http.Cookie
	Ocookie           *http.Cookie
}

type Step struct {
	Distance   string
	Duration   time.Duration
	Directions template.HTML
	Color      string
}

func (d *Directions) GetApikey() string {
	res := make([]Config, 10)
	ctx := appengine.NewContext(d.r)
	q := datastore.NewQuery("Config")
	_, _ = q.GetAll(ctx, &res)

	c := os.Getenv("APIKEY")
	for _, v := range res {
		if v.Name == "APIKEY" {
			c = v.Value
		}
	}
	return c
}

func NewDirections(r *http.Request) *Directions {
	var d = new(Directions)
	d.Origin = ""
	d.Destination = ""
	d.r = r
	d.Apikey = d.GetApikey()
	ctx := appengine.NewContext(r)
	uc := urlfetch.Client(ctx)
	c, err := maps.NewClient(maps.WithAPIKey(d.Apikey),
		maps.WithHTTPClient(uc))
	d.Client = c
	if err != nil {
		d.Resp = err.Error()
	}
	return d
}

func NewStep(v *maps.Step) *Step {
	st := Step{Distance: v.Distance.HumanReadable, Duration: v.Duration,
		Directions: template.HTML(v.HTMLInstructions),
		Color:      "none"}
	if st.Duration/time.Second > 5*60 {
		st.Color = "red"
	}
	return &st
}

func (d *Directions) DirectionsNow() {
	l, _ := time.LoadLocation("US/Pacific")
	t := time.Now().In(l)
	d.Directions(&t)
}

func (d *Directions) LookupDirections(tdd time.Time,
	origin string, destination string) []maps.Route {
	l, _ := time.LoadLocation("US/Pacific")
	dtime := strconv.FormatInt(tdd.In(l).Unix(), 10)
	ctx := appengine.NewContext(d.r)

	// check testdata canned API case. nil unless testing mode
	resp := testdata_read(tdd, origin, destination)
	if resp == nil {
		r := &maps.DirectionsRequest{
			Mode:          maps.TravelModeDriving,
			Origin:        origin,
			Destination:   destination,
			DepartureTime: dtime,
		}
		dir, _, err := d.Client.Directions(ctx, r)
		resp = dir
		if err != nil {
			log.Infof(ctx, err.Error())
			return nil
		}
		log.Infof(ctx, "new data %v", resp)
		testdata_save(resp, tdd, origin, destination)
	} else {
		log.Infof(ctx, "data read %v", resp)
	}
	log.Infof(ctx, "Lookup %v", resp)
	return resp
}

func hash_key(key string) string {
	data := []byte(key)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func timestamp(tdd time.Time) string {
	return tdd.Format("20060102150405")
}

func mkey(tdd time.Time, origin string, destination string) string {
	return timestamp(tdd) + ":" + origin + ":" + destination
}

func testdata_save(d []maps.Route, tdd time.Time, origin, destination string) {
	if os.Getenv("TESTDATA_MODE") == "SAVE" {
		b, _ := json.MarshalIndent(d, "", "\t")
		key := mkey(tdd, origin, destination)
		err := ioutil.WriteFile("testdata/"+timestamp(tdd)+"_"+hash_key(key), b, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func testdata_read(tdd time.Time, origin, destination string) []maps.Route {
	var resp []maps.Route = nil
	if os.Getenv("TESTDATA_MODE") == "READ" {
		key := mkey(tdd, origin, destination)
		_, err := os.Stat(key)
		if os.IsNotExist(err) {
			return resp
		}
		data, err := ioutil.ReadFile("testdata/" + timestamp(tdd) + ":" + hash_key(key))
		if err != nil {
			panic(err)
		}
		json.Unmarshal(data, &resp)
	}
	return resp
}

func (d *Directions) Directions(td *time.Time) {
	ctx := appengine.NewContext(d.r)
	// really not sure where the cookie/session stuff fits best.
	// put it here for now
	// two cookies for the start and dest total.
	var origin, destination string

	origin = "1200 Crittenden Lane, Mountain View"
	destination = "90 Enterprise Way, Scotts Valley"
	if d.Origin != "" {
		log.Infof(ctx, "unset %s %s", d.Origin, d.Destination)

		origin = d.Origin
		destination = d.Destination
	} else {
		cookie, err := d.r.Cookie("origin")
		if err == nil && cookie.Value != "" {
			origin = cookie.Value
		} else {
			cookie = &http.Cookie{Name: "origin", Value: origin}
		}
		d.Ocookie = cookie
		cookie, err = d.r.Cookie("destination")
		if err == nil && cookie.Value != "" {
			destination = cookie.Value
		}
		cookie = &http.Cookie{Name: "destination", Value: destination}
		d.Dcookie = cookie
	}
	d.Destination = destination
	d.Origin = origin

	// cache by intervals for better hit rate
	tdd := td.Truncate(30 * time.Minute)

	// for testing, all the time is in the past and should be used as is
	// maybe don't skip to next week if there's a cache hit?
	if tdd.Unix() < time.Now().Unix() {
		// look at next week for hints on past
		tdd = tdd.Add(time.Hour * 24 * 7)
	}

	if _, err := memcache.JSON.Get(ctx, mkey(tdd, origin, destination), &d.Dir); err == memcache.ErrCacheMiss {
		log.Infof(ctx, "item not in the cache: %s", mkey(tdd, origin, destination))

		resp := d.LookupDirections(tdd, origin, destination)
		d.Dir = &resp[0]
		forw := d.Dir.Legs[0]
		resp = d.LookupDirections(tdd, destination, origin)
		rev := resp[0].Legs[0]
		log.Infof(ctx, "rev %v", rev)
		if rev.DurationInTraffic > forw.DurationInTraffic {
			d.Dir = &resp[0]
		}
		d.Leg = d.Dir.Legs[0]

		err = memcache.JSON.Set(ctx,
			&memcache.Item{Key: mkey(tdd, origin, destination), Object: d.Dir})
		log.Infof(ctx, "cache update")
		if err != nil {
			log.Infof(ctx, err.Error())
		}
	} else if err != nil {
		log.Errorf(ctx, "error getting item: %v", err)
	}
	//	}
	log.Infof(ctx, "&v", d.Dir)
	for _, v := range d.Dir.Legs[0].Steps {
		d.Steps = append(d.Steps, NewStep(v))
	}
	d.Leg = d.Dir.Legs[0]
	d.Distance = d.Leg.Distance
	d.Duration = d.Leg.Duration
	d.DurationInTraffic = d.Leg.DurationInTraffic
}
