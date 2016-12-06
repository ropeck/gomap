package main

// https://elithrar.github.io/article/testing-http-handlers-go/
import (
	"testing"
	"time"

	"google.golang.org/appengine/aetest"
)

func TestDrawday(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)
	data := drawday(time.Now(), r) // [24][4]int
	if data[0] != [4]int{0, 0, 0, 0} {
		t.Fatalf("data mismatched. got %v", data[0])
	}
}
