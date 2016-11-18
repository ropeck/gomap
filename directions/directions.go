package directions
import (
	"net/http"
        "appengine"
        "appengine/datastore"
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
	for i, v := range res {
	  if v.Name == "APIKEY" {
	    c = v.Value
	  }
        }
	return c
}
