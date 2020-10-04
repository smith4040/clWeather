package main

import "testing"

func TestMakeURL(t *testing.T) {
	correctURL := "https://api.weather.gov/stations/kfwb/observations/latest"
	stationID := "kfwb"
	u := makeURL(stationID)
	if u != correctURL {
		t.Errorf("URL was incorrect, got: %s, want: %s.", u, correctURL)
	}
}
