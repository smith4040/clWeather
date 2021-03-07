package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeURL(t *testing.T) {
	correctURL := "https://api.weather.gov/stations/kfwb/observations/latest"
	stationID := "kfwb"
	u := makeURL(stationID)
	assert.Equal(t, u, correctURL)
}

func TestRequestWeather(t *testing.T) {

}

func TestCelsiusToFahrenheit(t *testing.T) {
	got := celsiusToFahrenheit(0.00)
	want := 32.00
	t.Logf("Running test case: TestCelsiusToFahrenheit")
	assert.Equal(t, got, want)
}
