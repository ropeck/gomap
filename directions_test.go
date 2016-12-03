package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"testing"
	"time"

	"googlemaps.github.io/maps"

	"google.golang.org/appengine/aetest"
)

func TestNewDirections(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	r, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}
	//	c1 := appengine.NewContext(r)

	d := NewDirections(r)
	if d == nil {
		t.Errorf("no Directions returned")
	}
}

func TestGetApiKey(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	r, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}
	//	c1 := appengine.NewContext(r)

	savekey := os.Getenv("APIKEY")
	os.Setenv("APIKEY", "testing")
	d := NewDirections(r)

	if d.Apikey == "" {
		t.Errorf("APIKEY is missing")
	}
	os.Unsetenv("APIKEY")
	if savekey != "" {
		os.Setenv("APIKEY", savekey)
	}

}

func TestNewStep(t *testing.T) {
	var m = new(maps.Step)
	m.Duration = time.Hour
	m.HTMLInstructions = "directions"
	m.Duration = time.Second * 60
	m.Distance = maps.Distance{HumanReadable: "2.2 km", Meters: 2241}

	s := NewStep(m)

	correctResponse := &Step{
		Distance: "2.2 km", Duration: time.Minute,
		Directions: template.HTML("directions"),
		Color:      "none",
	}
	if actualResponse := s; !reflect.DeepEqual(actualResponse, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, actualResponse)
	}
}
func TestNewStepRed(t *testing.T) {
	var m = new(maps.Step)
	m.Duration = time.Hour
	m.HTMLInstructions = "directions"
	m.Duration = time.Second * 60 * 60
	m.Distance = maps.Distance{HumanReadable: "2.2 km", Meters: 2241}

	s := NewStep(m)

	correctResponse := &Step{
		Distance: "2.2 km", Duration: time.Hour,
		Directions: template.HTML("directions"),
		Color:      "red",
	}
	if actualResponse := s; !reflect.DeepEqual(actualResponse, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, actualResponse)
	}
}

func TestDirections(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	d := NewDirections(r)
	d.DirectionsNow()
	// assert some things about the directions
	// the content will change because of the API results
	// also, the origin / destination cookies are missing

	if &d.Distance == nil {
		t.Fatalf("Failed to set Distance")
	}
	b, _ := json.Marshal(d)
	fmt.Printf(string(b))

}
