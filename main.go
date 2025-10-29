// clWeather is a Golang command line tool for querying the weather.gov
// API for current weather information for a specified station/stations

// TODO- handle error better when station isnt available, printing a blank line to terminal
// TODO- use remaining colors or delete

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	url := "https://api.weather.gov/stations/" + s + "/observations/latest?require_qc=true"

	return url
}

// celsiusToFahrenheit converts celsius to fahrenheit
func celsiusToFahrenheit(c float64) float64 {
	value := (c * 9 / 5) + 32
	return value
}

// processData prepares JSON
func processData(d []byte) (datamodel.Response, error) {
	var station datamodel.Response
	err := json.Unmarshal(d, &station)
	if err != nil {
		return station, err
	}
	return station, nil
}

func requestObservation(stationID string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := makeURL(stationID)
	weatherResponse, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := io.ReadAll(weatherResponse.Body)
	if err != nil {
		log.Fatalf("Error reading data: %s\n", err)
	}

	sc := weatherResponse.StatusCode
	if sc >= 400 {
		fmt.Println(warn(stationID, ": Weather observation for this station is currently unavailable. Check spelling or try again later."))
		fmt.Println(fata("Server error, status code " + fmt.Sprint(sc)))
		return
	}

	defer func() {
		err := weatherResponse.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	ch <- string(responseData)
}

func checkRawMessage(s string) {
	if s == "" {
		fmt.Println(red("Raw METAR not curently available"))
		return
	}
	fmt.Println(green(s))

}

// presentResults is called to display the weather on command line
func presentResults(stations []string) {
	ch := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(len(stations))

	for _, s := range stations {
		go requestObservation(s, ch, &wg)
	}

	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		p, err := processData([]byte(result))
		if err != nil {
			log.Printf("error processing data: %s\n", err)
			return
		}

		fmt.Println(green(p.Properties.StationName))
		checkRawMessage(p.Properties.RawMessage)

		if p.Properties.Temperature.Value.Valid {
			t := p.Properties.Temperature.Value.Value
			f := celsiusToFahrenheit(t)
			s := fmt.Sprintf("%.2f", f)
			fmt.Println(teal("Temperature is " + s + "°F\n"))
		} else {
			fmt.Println(warn("Temperature is currently unavailable, please try again later.\n"))
		}
	}
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
	fmt.Printf("Searching station(s): %s\n\n", stations)
	presentResults(stations)

	defer func() {
		fmt.Println("Execution Time: ", time.Since(start))
	}()
}
