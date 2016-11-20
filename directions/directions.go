package directions
import (
	"os"
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

	c := os.Getenv("APIKEY")
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
  Apikey string
}

func NewDirections(r *http.Request) *Directions {
  var d = new(Directions)
  c, _ := maps.NewClient(maps.WithAPIKey(Apikey(r)))
  d.Client = c
  d.Apikey = Apikey(r)
  return d
}

func (d *Directions) Directions() {
  d.Origin = "something"
}

