package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeURL(t *testing.T) {
	correctURL := "https://aviationweather.gov/api/data/metar?ids=kfwb&format=json&taf=true&hours=1.5"
	stationID := "kfwb"
	taf := "true"
	u := makeURL(stationID, taf)
	if correctURL != u {
		t.Errorf("got %v want %v", u, correctURL)
	}
}

func TestRequestObservation(t *testing.T) {
	// place holder demo code for testing requestWX
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}

		err := r.ParseForm()
		if err != nil {
			t.Error(err)
		}

		topic := r.Form.Get("topic")
		if topic != "meaningful-topic" {
			t.Errorf("Expected request to have ‘topic=meaningful-topic’, got: ‘%s’", topic)
		}
	}))
	defer ts.Close()
}

func TestCelsiusToFahrenheit(t *testing.T) {
	got := celsiusToFahrenheit(0.00)
	want := 32.00
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	g := celsiusToFahrenheit(15.00)
	w := 59.00
	if g != w {
		t.Errorf("got %v want %v", g, w)
	}

	gg := celsiusToFahrenheit(-10.00)
	ww := 14.00
	if gg != ww {
		t.Errorf("got %v want %v", gg, ww)
	}
}

// func TestProcessData(t *testing.T) {
// 	input := `{"properties":{"station":"ksgf"}}`
// 	got, err := processData([]byte(input))
// 	if err != nil {
// 		t.Error("Failure message: ", err)
// 	}
// 	want := "ksgf"

// 	if !reflect.DeepEqual(want, got.Properties.Station) {
// 		t.Fatal("Actual output doesn't match expected")
// 	}
// }
