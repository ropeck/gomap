package directions
import (
	"net/http"
        "appengine"
        "appengine/datastore"
	"googlemaps.github.io/maps"
)

type Config struct {
	Name string
	Value string
  }

func Apikey(r *http.Request) string {
	res := make([]Config, 10)
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("Config")
	_, _ = q.GetAll(ctx, &res)

	var c string
	for _, v := range res {
	  if v.Name == "APIKEY" {
	    c = v.Value
	  }
        }
	return c
}

type Directions struct {
	Origin string
  Client *maps.Client
}

func NewDirections(r *http.Request) *Directions {
  var d = new(Directions)
  c, _ := maps.NewClient(maps.WithAPIKey(Apikey(r)))
  d.Client = c
  return d
}

func (d *Directions) Directions() {
  d.Origin = "something"
}

