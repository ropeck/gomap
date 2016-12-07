package main

// https://elithrar.github.io/article/testing-http-handlers-go/
import (
	"os"
	"testing"
	"time"

	"google.golang.org/appengine/aetest"
)

func TestDrawday(t *testing.T) {
	inst, _ := aetest.NewInstance(nil)
	r, _ := inst.NewRequest("GET", "/", nil)

	os.Setenv("TESTDATA_MODE", "READ") // read from canned API data

	td, _ := time.Parse("20060102150405", "20161206130000")
	data := drawday(td, r)
	// [24][4]int
	if data[0] != [4]int{0, 0, 40, 39} {
		t.Fatalf("data mismatched. got %v", data[0])
	}
}
