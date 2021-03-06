package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeURL(t *testing.T) {
	correctURL := "https://api.weather.gov/stations/kfwb/observations/latest"
	stationID := "kfwb"
	u := makeURL(stationID)
	t.Logf("Running test case: TestMakeURL")
	assert.Equal(t, u, correctURL)
}

func TestRequestWeather(t *testing.T) {

}

func TestCelsiusToFahrenheit(t *testing.T){
	got := celsiusToFahrenheit(0)
	want := 32
	if got != want {
		t.Logf("got %q, want %q", got, want)
	}
}
