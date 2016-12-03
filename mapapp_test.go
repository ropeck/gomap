package main

// https://elithrar.github.io/article/testing-http-handlers-go/
import (
	"fmt"
	"testing"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/log"
)

func TestDrawday(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	data := drawday(time.Now(), r) // [24][4]int
	log.Infof(appengine.NewContext(r), fmt.Sprint("%v", data))
}
