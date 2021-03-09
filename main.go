// clWeather is a Golang command line tool for querying the weather.gov
// API for current weather information for a specified station/stations

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	datamodel "github.com/smith4040/clWeather/datamodel"
)

// makeURL builds the URL for the endpoint to query
func makeURL(s string) string {
	url := fmt.Sprintf("https://api.weather.gov/stations/%s/observations/latest", s)

	return url
}

// celsiusToFahrenheit converts celsius to fahrenheit
func celsiusToFahrenheit(c float64) float64 {
	value := (c * 9 / 5) + 32
	return value
}

// processData prepares JSON
func processData(d []byte) datamodel.Response {
	var station datamodel.Response
	err := json.Unmarshal(d, &station)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return station
}

func requestObservation(stationID string) (datamodel.Response, error) {
	url := makeURL(stationID)
	weatherResponse, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
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

	p := processData(responseData)
	sc := weatherResponse.StatusCode
	if sc >= 400 {
		colour.Red(stationID + ":" + " Weather observation for this station is currently unavailable. Check spelling or try again later.")
		return p, errors.New("Server error, status code " + fmt.Sprint(sc))
	}
	return p, nil
}

// presentResults is called to display the weather on command line
func presentResults(stations []string) {
	wg := sync.WaitGroup{}
	wg.Add(len(stations))

	for _, s := range stations {
		go func(s string) {
			result, err := requestObservation(s)
			if err != nil {
				log.Fatal(err)
			}

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

// printUsage displays flags to the user if none are presented
func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	start := time.Now()
	var station string

	flag.StringVar(&station, "s", "", "Station ID")
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
