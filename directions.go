package main

import (
	"html/template"
	"net/http"
	"os"
	"time"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"googlemaps.github.io/maps"
)

type Config struct {
	Name  string
	Value string
}

type Directions struct {
	Origin string
	Client *maps.Client
	Apikey string
	r      *http.Request
	Resp   string
	Leg    *maps.Leg
	Dir    *maps.Route
	Steps  []*Step
	Duration   time.Duration
	DurationInTraffic   time.Duration
	Distance maps.Distance
	Dcookie	*http.Cookie
	Ocookie *http.Cookie
}

type Step struct {
	Distance	string
	Duration   time.Duration
	Directions template.HTML
	Color	string
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
	d.r = r
	d.Apikey = d.GetApikey()
	ctx := appengine.NewContext(r)
	uc := urlfetch.Client(ctx)
	c, err := maps.NewClient(maps.WithAPIKey(d.Apikey), maps.WithHTTPClient(uc))
	d.Client = c
	if err != nil {
		d.Resp = err.Error()
	}
	return d
}

func NewStep(v *maps.Step) *Step {
	st := Step{Distance: v.Distance.HumanReadable, Duration: v.Duration,
		Directions: template.HTML(v.HTMLInstructions),
		Color: "none"}
	if (st.Duration/time.Second > 5*60) {
		st.Color = "red"
	}
	return &st
}

func (d *Directions) DirectionsNow() {
	t := time.Now()
	d.Directions(&t)
}

func (d *Directions) Directions(td *time.Time) {
	// really not sure where the cookie/session stuff fits best.
	// put it here for now
	// two cookies for the start and dest total.
	var origin, destination string

	origin = "1200 Crittenden Lane, Mountain View"
	destination = "90 Enterprise Way, Scotts Valley"
	cookie, err := d.r.Cookie("origin")
	if err == nil && cookie.Value != "" {
		origin = cookie.Value
	} else {
		cookie=  &http.Cookie{Name: "origin", Value: origin}
 	}
	d.Ocookie = cookie
	cookie, err = d.r.Cookie("destination")
	if err == nil && cookie.Value != "" {
		destination = cookie.Value
	}
		cookie=  &http.Cookie{Name: "destination", Value: destination}
	d.Dcookie = cookie

	r := &maps.DirectionsRequest{
		Mode:        maps.TravelModeDriving,
		Origin:      origin,
		Destination: destination,
		DepartureTime:  strconv.FormatInt(td.Unix(), 10),
	}
	ctx := appengine.NewContext(d.r)

	mkey := strconv.FormatInt(td.Unix(), 10)
	log.Infof(ctx, "memcache: " + mkey)
	if _, err := memcache.JSON.Get(ctx, mkey, &d.Dir); err == memcache.ErrCacheMiss {
		log.Infof(ctx, "item not in the cache")
		resp, _, err := d.Client.Directions(appengine.NewContext(d.r), r)
		d.Dir = &resp[0]

 		err = memcache.JSON.Set(ctx, &memcache.Item{Key: mkey, Object: d.Dir})
		if err != nil {
			log.Infof(ctx, err.Error())
		}
	} else if err != nil {
		log.Errorf(ctx, "error getting item: %v", err)
	}

	for _, v := range d.Dir.Legs[0].Steps {
		d.Steps = append(d.Steps, NewStep(v))
	}
	d.Leg = d.Dir.Legs[0]
	d.Distance = d.Leg.Distance
	d.Duration = d.Leg.Duration
	d.DurationInTraffic = d.Leg.DurationInTraffic
}
