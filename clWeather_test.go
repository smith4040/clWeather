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

func TestPrintUsage(t *testing.T) {

}
