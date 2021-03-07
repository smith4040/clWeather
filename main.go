// clWeather is a Golang command line tool for querying the weather.gov
// API for current weather information for a specified station/stations

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	colour "github.com/fatih/color"
	flag "github.com/ogier/pflag"
	dataModel "github.com/smith4040/clWeather/datamodel"
)

var station string

// makeURL builds the URL for the endpoint to query
func makeURL(s string) string {
	const (
		apiURL          = "https://api.weather.gov/stations/"
		observationType = "/observations/latest"
	)

	url := apiURL + s + observationType
	return url
}

// requestWeather makes the GET request to weather service and prepares JSON
func requestWeather(stationID string) dataModel.Response {
	url := makeURL(stationID)

	weatherResponse, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error reading data: %s\n", err)
	}

	defer func() {
		err := weatherResponse.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	responseData, err := ioutil.ReadAll(weatherResponse.Body)
	if err != nil {
		log.Fatalf("Error reading data: %s\n", err)
	}

	var station dataModel.Response
	err = json.Unmarshal(responseData, &station)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	if weatherResponse.StatusCode >= 400 {
		colour.Red(stationID + ":" + " Weather observation for this station is currently unavailable. Check spelling and try again later.")
		log.Print("Response status code: ", weatherResponse.StatusCode)
		return station
	}
	return station
}

// presentResults is called to display the weather on command line
func presentResults(stations []string) {
	wg := sync.WaitGroup{}
	wg.Add(len(stations))

	for _, s := range stations {
		go func(s string) {
			result := requestWeather(s)

			if result.Properties.Temperature.Value.Valid {
				colour.HiGreen(result.Properties.RawMessage)
				t := result.Properties.Temperature.Value.Value
				u := fmt.Sprintf("%.2f", t)
				fmt.Println("Temperature is " + u + "°C")
				f := celsiusToFahrenheit(t)
				s := fmt.Sprintf("%.2f", f)
				colour.HiGreen("Temperature is " + s + "°F")
				fmt.Println("")
			} else {
				colour.HiGreen(result.Properties.RawMessage)
				colour.HiRed("Temperature is currently unavailable, please try again later.")
				fmt.Println("")
			}
			wg.Done()
		}(s)
	}
	wg.Wait()
	fmt.Println("All requests complete.")
}

// celsiusToFahrenheit converts celsius to fahrenheit
func celsiusToFahrenheit(c float64) float64 {
	value := (c * 9 / 5) + 32
	return value
}

// printUsage displays flags to the user if none are presented
func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.StringVarP(&station, "station", "s", "", "Station ID")
}

func main() {
	start := time.Now()

	flag.Parse()
	if flag.NFlag() == 0 {
		printUsage() // if user does not supply flags, print usage
	}

	stations := strings.Split(station, ",")
	fmt.Printf("Searching station(s): %s\n", stations)
	fmt.Println("")

	presentResults(stations)

	defer func() {
		fmt.Println("Execution Time: ", time.Since(start))
	}()
}
