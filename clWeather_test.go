package main

import "testing"

func TestRequestWeather(t *testing.T) {
	stationID := "kfwb"
	url := "https://api.weather.gov/stations/" + (stationID) + "/observations/latest"
	correctURL := "https://api.weather.gov/stations/kfwb/observations/latest"
	if url != correctURL {
		t.Errorf("URL was incorrect, got: %s, want: %s.", url, correctURL)
	}
}
